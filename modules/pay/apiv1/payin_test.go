package apiv1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	"sync/atomic"
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

func TestPayInQueryReturnsStoredStatusWithoutSideEffects(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var notifications, queries atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notifications.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	previousChannel := allpay.All.TKPay
	allpay.All.TKPay = func(b *base.Base) {
		queries.Add(1)
	}
	defer func() { allpay.All.TKPay = previousChannel }()
	for status := common.StatusROrderWaitting; status < common.StatusROrderAll; status++ {
		for _, dropped := range []bool{false, true} {
			t.Run(fmt.Sprintf("status_%d/dropped_%t", status, dropped), func(t *testing.T) {
				mock := queryDB(t)
				// Only the merchant-scoped read is allowed: transactions, writes and extra reads fail the test.
				mock.ExpectQuery(regexp.QuoteMeta(orderLookupSQL)).WithArgs("merchant-7", 12, 1).
					WillReturnRows(queryOrderRows(status, dropped, server.URL))
				code, data := queryRequest(t, 12, "")
				if code != app.CodeOK || data == nil || data.OrderStatus != status || data.OrderAmount != 10000 ||
					data.PlatformOrderID != "platform-7" || data.MerchantOrderID != "merchant-7" || data.MerchantID != 12 {
					t.Fatalf("unexpected query response: code=%d data=%+v", code, data)
				}
				if queries.Load() != 0 || notifications.Load() != 0 {
					t.Fatalf("query produced side effects: upstream queries=%d merchant notifications=%d", queries.Load(), notifications.Load())
				}
			})
		}
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

func TestPayInQueryReportsDatabaseReadFailure(t *testing.T) {
	mock := queryDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(orderLookupSQL)).WithArgs("merchant-7", 12, 1).
		WillReturnError(errors.New("injected database read failure"))
	if code, data := queryRequest(t, 12, ""); code != app.CodeInnerError || data != nil {
		t.Fatalf("unexpected database error response: code=%d data=%+v", code, data)
	}
}
