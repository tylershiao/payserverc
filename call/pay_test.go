package call

import (
	"encoding/json"
	"errors"
	"fmt"
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
		Status: common.StatusROrderCreate, Amount: 100, Fee: 2}
}

func expectCurrentOrder(mock sqlmock.Sqlmock, order *common.OrderPay) {
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `order_pay` WHERE `order_pay`.`id` = ? AND `order_pay`.`deleted_at` IS NULL ORDER BY `order_pay`.`id` LIMIT ? FOR UPDATE")).
		WithArgs(order.ID, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "mid", "order_id", "m_order_id", "amount", "fee", "status", "is_dropped", "realtime_settle", "callback_url"}).
		AddRow(order.ID, order.MID, order.OrderID, order.MOrderID, order.Amount, order.Fee, order.Status, order.IsDropped, order.RealTimeSettle, order.CallbackURL))
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
			order.IsDropped = true
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
	// The failure callback loaded Create before an earlier success committed Finish.
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

func callbackNotifications(t *testing.T) (string, <-chan common.PayCallbackReq) {
	t.Helper()
	notifications := make(chan common.PayCallbackReq, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req common.PayCallbackReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode merchant notification: %v", err)
		}
		notifications <- req
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)
	return server.URL, notifications
}

func assertNoCallbackNotification(t *testing.T, notifications <-chan common.PayCallbackReq) {
	t.Helper()
	select {
	case req := <-notifications:
		t.Fatalf("unexpected merchant notification: %+v", req)
	case <-time.After(50 * time.Millisecond):
	}
}

func expectDropUpdate(mock sqlmock.Sqlmock, order *common.OrderPay) *sqlmock.ExpectedExec {
	// UpdateColumn must only change the drop flag, without changing status or timestamps.
	return mock.ExpectExec(regexp.QuoteMeta("UPDATE `order_pay` SET `is_dropped`=? WHERE id = ? AND `order_pay`.`deleted_at` IS NULL")).
		WithArgs(true, order.ID)
}

func TestPayCallbackUpstreamDropProbabilityBoundaries(t *testing.T) {
	for _, value := range []int{0, 4, 5, 999} {
		t.Run(fmt.Sprintf("draw_%d", value), func(t *testing.T) {
			mock := callbackDB(t)
			order := callbackOrder()
			var notifications <-chan common.PayCallbackReq
			order.CallbackURL, notifications = callbackNotifications(t)
			drop := value < 5
			mock.ExpectBegin()
			expectCurrentOrder(mock, order)
			if drop {
				expectDropUpdate(mock, order).WillReturnResult(sqlmock.NewResult(0, 1))
			} else {
				expectOrderUpdate(mock, order).WillReturnResult(sqlmock.NewResult(0, 1))
				expectMerchant(mock)
				expectBalance(mock, false).WillReturnResult(sqlmock.NewResult(0, 1))
				expectFund(mock, false).WillReturnResult(sqlmock.NewResult(1, 1))
			}
			mock.ExpectCommit()
			draws := 0
			processed, err := payCallback(order, true, func(n int) int {
				draws++
				if n != 1000 {
					t.Fatalf("expected 1000 possible outcomes, got %d", n)
				}
				return value
			})
			if err != nil || processed != !drop || draws != 1 {
				t.Fatalf("unexpected callback result: processed=%t err=%v draws=%d", processed, err, draws)
			}
			if drop {
				assertNoCallbackNotification(t, notifications)
				return
			}
			select {
			case req := <-notifications:
				if req.OrderStatus != common.StatusROrderFinish || req.OrderAmount != 10000 {
					t.Fatalf("invalid normal success notification: %+v", req)
				}
			case <-time.After(time.Second):
				t.Fatal("missing normal success notification")
			}
		})
	}
}

func TestPayCallbackIgnoresPersistedDroppedSuccess(t *testing.T) {
	mock := callbackDB(t)
	stale := callbackOrder()
	var notifications <-chan common.PayCallbackReq
	stale.CallbackURL, notifications = callbackNotifications(t)
	current := *stale
	current.IsDropped = true
	mock.ExpectBegin()
	expectCurrentOrder(mock, &current)
	mock.ExpectCommit()
	processed, err := payCallback(stale, true, func(int) int {
		t.Fatal("already dropped success must not draw again")
		return 999
	})
	if err != nil || processed {
		t.Fatalf("unexpected dropped callback result: processed=%t err=%v", processed, err)
	}
	// A query that loaded the order before it was dropped must not settle it either.
	mock.ExpectBegin()
	expectCurrentOrder(mock, &current)
	mock.ExpectCommit()
	if err := PayCallback(stale, true); err != nil {
		t.Fatal(err)
	}
	assertNoCallbackNotification(t, notifications)
}

func TestPayCallbackFinishedDuplicateDoesNotDraw(t *testing.T) {
	mock := callbackDB(t)
	stale := callbackOrder()
	current := *stale
	current.Status = common.StatusROrderFinish
	mock.ExpectBegin()
	expectCurrentOrder(mock, &current)
	mock.ExpectCommit()
	processed, err := payCallback(stale, true, func(int) int {
		t.Fatal("completed order must not be randomly dropped by a duplicate callback")
		return 0
	})
	if err != nil || processed {
		t.Fatalf("unexpected duplicate result: processed=%t err=%v", processed, err)
	}
}

func TestPayCallbackDroppedFailureKeepsOriginalBehavior(t *testing.T) {
	mock := callbackDB(t)
	order := callbackOrder()
	order.IsDropped = true
	var notifications <-chan common.PayCallbackReq
	order.CallbackURL, notifications = callbackNotifications(t)
	mock.ExpectBegin()
	expectCurrentOrder(mock, order)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `order_pay` SET `api_order_id`=?,`callback_time`=?,`status`=?,`updated_at`=? WHERE (id = ? and status <> ?) AND `order_pay`.`deleted_at` IS NULL")).
		WithArgs(order.APIOrderID, sqlmock.AnyArg(), common.StatusROrderFail, sqlmock.AnyArg(), order.ID, common.StatusROrderFail).
		WillReturnResult(sqlmock.NewResult(0, 1))
	expectMerchant(mock)
	mock.ExpectCommit()
	processed, err := payCallback(order, false, func(int) int {
		t.Fatal("failure callback must not draw for a drop")
		return 0
	})
	if err != nil || !processed {
		t.Fatalf("unexpected failure callback result: processed=%t err=%v", processed, err)
	}
	select {
	case req := <-notifications:
		if req.OrderStatus != common.StatusROrderFail {
			t.Fatalf("invalid failure notification: %+v", req)
		}
	case <-time.After(time.Second):
		t.Fatal("missing original failure notification")
	}
}

func TestPayCallbackDropPersistenceFailureDoesNotSettleOrNotify(t *testing.T) {
	for _, stage := range []string{"mark", "commit"} {
		t.Run(stage, func(t *testing.T) {
			mock := callbackDB(t)
			order := callbackOrder()
			var notifications <-chan common.PayCallbackReq
			order.CallbackURL, notifications = callbackNotifications(t)
			failure := errors.New("injected drop " + stage + " failure")
			mock.ExpectBegin()
			expectCurrentOrder(mock, order)
			mark := expectDropUpdate(mock, order)
			if stage == "mark" {
				mark.WillReturnError(failure)
				mock.ExpectRollback()
			} else {
				mark.WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit().WillReturnError(failure)
			}
			processed, err := payCallback(order, true, func(int) int { return 0 })
			if !errors.Is(err, failure) || processed {
				t.Fatalf("unexpected persistence failure result: processed=%t err=%v", processed, err)
			}
			assertNoCallbackNotification(t, notifications)
		})
	}
}
