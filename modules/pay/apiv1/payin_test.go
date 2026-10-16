package apiv1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"server/call"
	"server/common"
	"server/db"
	mdb "server/db/mysql"
	"server/modules/pay/allpay"
	"server/modules/pay/app"
	"server/modules/pay/base"
	"server/modules/pay/values"
	"server/sign"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func queryDB(t *testing.T) sqlmock.Sqlmock {
	t.Helper()
	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	orm, err := gorm.Open(mysql.New(mysql.Config{Conn: conn, SkipInitializeWithVersion: true}), &gorm.Config{
		DisableAutomaticPing: true, Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	previousDB, previousMerchants := db.Mysql(), call.ConfigMerchats
	db.SetDB(mdb.NewMysqlClient(orm))
	call.ConfigMerchats = []*common.Merchats{{MID: 12, Key: "test-key", Open: true}, {MID: 13, Key: "other-key", Open: true}}
	t.Cleanup(func() {
		db.SetDB(previousDB)
		call.ConfigMerchats = previousMerchants
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
		conn.Close()
	})
	return mock
}

const orderLookupSQL = "SELECT * FROM `order_pay` WHERE (`order_pay`.`m_order_id` = ? AND `order_pay`.`mid` = ?) AND `order_pay`.`deleted_at` IS NULL ORDER BY `order_pay`.`id` LIMIT ?"
const orderReloadSQL = "SELECT * FROM `order_pay` WHERE `order_pay`.`id` = ? AND `order_pay`.`deleted_at` IS NULL ORDER BY `order_pay`.`id` LIMIT ?"

func queryOrderRows(status int, dropped bool, callbackURL string) *sqlmock.Rows {
	oldTime := time.Date(2026, 1, 1, 1, 0, 0, 0, time.UTC)
	return sqlmock.NewRows([]string{"id", "order_id", "m_order_id", "api_order_id", "mid", "amount", "fee", "status", "is_dropped", "pay_channel", "callback_url", "created_at", "updated_at"}).
		AddRow(7, "platform-7", "merchant-7", "upstream-7", 12, 100, 2, status, dropped, values.TKPay, callbackURL, oldTime, oldTime)
}

func queryRequest(t *testing.T, mid int, signature string) (int, *PayInQueryResp) {
	t.Helper()
	req := &PayInQueryReq{MerchantID: mid, MerchantOrderID: "merchant-7", Signature: signature}
	if signature == "" {
		key := "test-key"
		if mid == 13 {
			key = "other-key"
		}
		req.Signature = sign.APISignMD5(req, key, "Sign", "signature")
	}
	body, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/payin/query", bytes.NewReader(body))
	context.Request.Header.Set("Content-Type", "application/json")
	PayInQuery(context)
	var response struct {
		Code int             `json:"code"`
		Data *PayInQueryResp `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	return response.Code, response.Data
}

func expectQueryPayment(mock sqlmock.Sqlmock, status int, dropped bool, callbackURL string) {
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(orderReloadSQL+" FOR UPDATE")).WithArgs(7, 1).
		WillReturnRows(queryOrderRows(status, dropped, callbackURL))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `order_pay` SET `api_order_id`=?,`callback_time`=?,`status`=?,`updated_at`=? WHERE (id = ? and status <> ?) AND `order_pay`.`deleted_at` IS NULL")).
		WithArgs("upstream-7", sqlmock.AnyArg(), common.StatusROrderFinish, sqlmock.AnyArg(), 7, common.StatusROrderFinish).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT \\* FROM `merchats` WHERE mid = \\? ORDER BY `merchats`.`id` LIMIT \\? FOR UPDATE").
		WithArgs(12, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "mid", "balance_settle"}).AddRow(3, 12, 300))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `merchats` SET `balance_settle`=balance_settle + 98.00,`updated_at`=? WHERE id = ?")).
		WithArgs(sqlmock.AnyArg(), 3).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO `fund_details`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func TestPayInQueryRoutesAndReturnsRefreshedStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		name        string
		status      int
		dropped     bool
		queryStatus int
		queries     int
		finish      bool
	}{
		{"already_successful", common.StatusROrderFinish, true, base.QueryStatusUnknown, 0, false},
		{"dropped_pending", common.StatusROrderCreate, true, base.QueryStatusUnknown, 0, true},
		{"dropped_failed", common.StatusROrderFail, true, base.QueryStatusUnknown, 0, true},
		{"upstream_successful", common.StatusROrderCreate, false, base.QueryStatusOrderSuccess, 1, true},
		{"upstream_pending", common.StatusROrderCreate, false, base.QueryStatusUnknown, 1, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mock := queryDB(t)
			notifications := make(chan struct{}, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				notifications <- struct{}{}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			queries := 0
			previousChannel := allpay.All.TKPay
			allpay.All.TKPay = func(b *base.Base) {
				b.Sub = &struct{ base.Sub }{}
				b.SubReq = func() error {
					queries++
					if b.QueryPayReq.OrderID != "platform-7" || b.QueryPayReq.APIOrderID != "upstream-7" {
						t.Errorf("incorrect upstream query: %+v", b.QueryPayReq)
					}
					b.QueryPayResp.Status = tc.queryStatus
					return nil
				}
			}
			defer func() { allpay.All.TKPay = previousChannel }()
			mock.ExpectQuery(regexp.QuoteMeta(orderLookupSQL)).WithArgs("merchant-7", 12, 1).
				WillReturnRows(queryOrderRows(tc.status, tc.dropped, server.URL))
			wantStatus := tc.status
			if tc.finish {
				expectQueryPayment(mock, tc.status, tc.dropped, server.URL)
				wantStatus = common.StatusROrderFinish
			}
			if tc.status != common.StatusROrderFinish {
				// Exact primary-key SQL prevents stale status/timestamp predicates from hiding the updated order.
				mock.ExpectQuery(regexp.QuoteMeta(orderReloadSQL)).WithArgs(7, 1).
					WillReturnRows(queryOrderRows(wantStatus, tc.dropped, server.URL))
			}
			code, data := queryRequest(t, 12, "")
			if code != app.CodeOK || data == nil || data.OrderStatus != wantStatus || data.OrderAmount != 10000 || data.PlatformOrderID != "platform-7" {
				t.Fatalf("unexpected query response: code=%d data=%+v", code, data)
			}
			if queries != tc.queries {
				t.Fatalf("upstream queries=%d, want %d", queries, tc.queries)
			}
			if tc.finish {
				select {
				case <-notifications:
				case <-time.After(time.Second):
					t.Fatal("missing merchant notification")
				}
			}
		})
	}
}

func TestPayInQueryRejectsInvalidSignatureAndScopesMerchant(t *testing.T) {
	t.Run("invalid_signature", func(t *testing.T) {
		queryDB(t)
		if code, data := queryRequest(t, 12, "invalid"); code != app.CodeSign || data != nil {
			t.Fatalf("unexpected authentication response: code=%d data=%+v", code, data)
		}
	})
	t.Run("other_merchant", func(t *testing.T) {
		mock := queryDB(t)
		// Merchant 13 cannot look up merchant 12's order, even with a valid signature.
		mock.ExpectQuery(regexp.QuoteMeta(orderLookupSQL)).WithArgs("merchant-7", 13, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		if code, data := queryRequest(t, 13, ""); code != app.CodeOrderNotExist || data != nil {
			t.Fatalf("unexpected merchant scope response: code=%d data=%+v", code, data)
		}
	})
}

func TestPayInQueryDoesNotReportSuccessWhenCompletionFails(t *testing.T) {
	mock := queryDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(orderLookupSQL)).WithArgs("merchant-7", 12, 1).
		WillReturnRows(queryOrderRows(common.StatusROrderCreate, true, ""))
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(orderReloadSQL+" FOR UPDATE")).WithArgs(7, 1).
		WillReturnRows(queryOrderRows(common.StatusROrderCreate, true, ""))
	mock.ExpectExec("UPDATE `order_pay`").WillReturnError(errors.New("injected order update failure"))
	mock.ExpectRollback()
	if code, data := queryRequest(t, 12, ""); code != app.CodeInnerError || data != nil {
		t.Fatalf("reported unsuccessful completion as success: code=%d data=%+v", code, data)
	}
}
