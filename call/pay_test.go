package call

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"server/common"
	"server/db"
	mdb "server/db/mysql"
	"server/sign"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func callbackDB(t *testing.T) sqlmock.Sqlmock {
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
	previousDB, previousMerchants := db.Mysql(), ConfigMerchats
	db.SetDB(mdb.NewMysqlClient(orm))
	ConfigMerchats = []*common.Merchats{{MID: 12, Key: "test-key"}}
	t.Cleanup(func() {
		db.SetDB(previousDB)
		ConfigMerchats = previousMerchants
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
		conn.Close()
	})
	return mock
}

func callbackOrder() *common.OrderPay {
	return &common.OrderPay{ID: 7, MID: 12, OrderID: "platform-7", MOrderID: "merchant-7", APIOrderID: "upstream-7",
		Status: common.StatusROrderCreate, Amount: 100, Fee: 2, IsDropped: true}
}

func expectCurrentOrder(mock sqlmock.Sqlmock, order *common.OrderPay) {
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `order_pay` WHERE `order_pay`.`id` = ? AND `order_pay`.`deleted_at` IS NULL ORDER BY `order_pay`.`id` LIMIT ? FOR UPDATE")).
		WithArgs(order.ID, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "mid", "order_id", "m_order_id", "amount", "fee", "status", "realtime_settle", "callback_url"}).
		AddRow(order.ID, order.MID, order.OrderID, order.MOrderID, order.Amount, order.Fee, order.Status, order.RealTimeSettle, order.CallbackURL))
}

func expectOrderUpdate(mock sqlmock.Sqlmock, order *common.OrderPay) *sqlmock.ExpectedExec {
	return mock.ExpectExec(regexp.QuoteMeta("UPDATE `order_pay` SET `api_order_id`=?,`callback_time`=?,`status`=?,`updated_at`=? WHERE (id = ? and status <> ?) AND `order_pay`.`deleted_at` IS NULL")).
		WithArgs(order.APIOrderID, sqlmock.AnyArg(), common.StatusROrderFinish, sqlmock.AnyArg(), order.ID, common.StatusROrderFinish)
}

func expectMerchant(mock sqlmock.Sqlmock) {
	mock.ExpectQuery("SELECT \\* FROM `merchats` WHERE mid = \\? ORDER BY `merchats`.`id` LIMIT \\? FOR UPDATE").
		WithArgs(12, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "mid", "balance", "balance_settle"}).AddRow(3, 12, 200, 300))
}

func expectBalance(mock sqlmock.Sqlmock, realTime bool) *sqlmock.ExpectedExec {
	col := "balance_settle"
	if realTime {
		col = "balance"
	}
	return mock.ExpectExec(regexp.QuoteMeta("UPDATE `merchats` SET `"+col+"`="+col+" + 98.00,`updated_at`=? WHERE id = ?")).WithArgs(sqlmock.AnyArg(), 3)
}

func expectFund(mock sqlmock.Sqlmock, realTime bool) *sqlmock.ExpectedExec {
	balType, balance := common.FundDetailsBalanceTypeSettle, 398
	if realTime {
		balType, balance = common.FundDetailsBalanceTypeAvailable, 298
	}
	return mock.ExpectExec("INSERT INTO `fund_details`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil,
		12, "platform-7", common.FundDetailsTypeRecharge, common.GetFundDetailsTypeName(common.FundDetailsTypeRecharge),
		float64(98), balType, common.GetFundDetailsBalanceTypeName(balType), float64(balance))
}

func TestPayCallbackCreditsOnceAndNotifiesAfterCommit(t *testing.T) {
	for _, realTime := range []bool{false, true} {
		t.Run(map[bool]string{false: "settlement", true: "available"}[realTime], func(t *testing.T) {
			mock := callbackDB(t)
			order := callbackOrder()
			order.RealTimeSettle = realTime
			type notification struct {
				request common.PayCallbackReq
				err     error
			}
			notifications := make(chan notification, 2)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req common.PayCallbackReq
				err := json.NewDecoder(r.Body).Decode(&req)
				if err == nil {
					err = mock.ExpectationsWereMet()
				}
				notifications <- notification{request: req, err: err}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			order.CallbackURL = server.URL
			persisted := *order
			mock.ExpectBegin()
			expectCurrentOrder(mock, &persisted)
			expectOrderUpdate(mock, order).WillReturnResult(sqlmock.NewResult(0, 1))
			expectMerchant(mock)
			expectBalance(mock, realTime).WillReturnResult(sqlmock.NewResult(0, 1))
			expectFund(mock, realTime).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			// Accounting and notification must use the locked database order, even with stale input fields.
			order.Amount, order.Fee, order.MID = 1000, 100, 13
			order.RealTimeSettle = !realTime
			if err := PayCallback(order, true); err != nil {
				t.Fatal(err)
			}
			select {
			case got := <-notifications:
				if got.err != nil {
					t.Fatalf("notification preceded transaction commit: %v", got.err)
				}
				if got.request.OrderStatus != common.StatusROrderFinish || got.request.OrderAmount != 10000 {
					t.Fatalf("invalid notification: %+v", got.request)
				}
				if got.request.Signature != sign.SignMD5(&got.request, "test-key", "Sign", "signature") {
					t.Fatal("invalid callback signature")
				}
			case <-time.After(time.Second):
				t.Fatal("missing merchant notification")
			}
			// A second query can still carry the original pending status; the database condition wins.
			persisted.Status = common.StatusROrderFinish
			mock.ExpectBegin()
			expectCurrentOrder(mock, &persisted)
			mock.ExpectCommit()
			if err := PayCallback(order, true); err != nil {
				t.Fatal(err)
			}
			select {
			case <-notifications:
				t.Fatal("duplicate merchant notification")
			case <-time.After(50 * time.Millisecond):
			}
		})
	}
}

func TestPayCallbackRollsBackAndDoesNotNotify(t *testing.T) {
	for _, stage := range []string{"order", "balance", "fund", "commit"} {
		t.Run(stage, func(t *testing.T) {
			mock := callbackDB(t)
			order := callbackOrder()
			notifications := make(chan struct{}, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				notifications <- struct{}{}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			order.CallbackURL = server.URL
			failure := errors.New("injected " + stage + " failure")
			mock.ExpectBegin()
			expectCurrentOrder(mock, order)
			update := expectOrderUpdate(mock, order)
			if stage == "order" {
				update.WillReturnError(failure)
			} else {
				update.WillReturnResult(sqlmock.NewResult(0, 1))
				expectMerchant(mock)
				balance := expectBalance(mock, false)
				if stage == "balance" {
					balance.WillReturnError(failure)
				} else {
					balance.WillReturnResult(sqlmock.NewResult(0, 1))
					fund := expectFund(mock, false)
					if stage == "fund" {
						fund.WillReturnError(failure)
					} else {
						fund.WillReturnResult(sqlmock.NewResult(1, 1))
					}
				}
			}
			if stage == "commit" {
				mock.ExpectCommit().WillReturnError(failure)
			} else {
				mock.ExpectRollback()
			}
			if err := PayCallback(order, true); !errors.Is(err, failure) {
				t.Fatalf("expected transaction error %v, got %v", failure, err)
			}
			select {
			case <-notifications:
				t.Fatal("notified merchant about uncommitted payment")
			case <-time.After(50 * time.Millisecond):
			}
		})
	}
}

func TestPayCallbackIgnoresStaleFailureAfterSuccess(t *testing.T) {
	mock := callbackDB(t)
	stale := callbackOrder()
	current := *stale
	current.Status = common.StatusROrderFinish
	// The failure callback loaded Create before a dropped-order query committed Finish.
	mock.ExpectBegin()
	expectCurrentOrder(mock, &current)
	mock.ExpectCommit()
	if err := PayCallback(stale, false); err != nil {
		t.Fatal(err)
	}
	// Another merchant query must still observe Finish and must not credit the balance again.
	mock.ExpectBegin()
	expectCurrentOrder(mock, &current)
	mock.ExpectCommit()
	if err := PayCallback(stale, true); err != nil {
		t.Fatal(err)
	}
}
