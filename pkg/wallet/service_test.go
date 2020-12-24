package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/AlisherGulomzoda/wallet/pkg/types"
	"github.com/google/uuid"
)

type testAccount struct {
	phone types.Phone
	balance types.Money
	payments []struct {
		amount types.Money
		category types.PaymrntCategoty
	}
}

func (t *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := t.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, err
	}

	err = t.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, err
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = t.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, err
		}
	}

	return account, payments, nil
	
}

var defaultTestAccount = testAccount {
	phone: "+992000000001",
	balance: 10_000_00,
	payments: []struct {
		amount types.Money
		category types.PaymrntCategoty
	} {
		{amount: 1_000_00, category: "Auto"},
	},
}

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (t *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	account, err := t.RegisterAccount(phone)
	if err != nil {
		return nil, err
	}

	err = t.Deposit(account.ID, balance)
	if err != nil {
		return nil, err
	}

	return account, nil

}

func TestService_FindAccountByID_success(t *testing.T) {
	s := &Service{}
	account, err := s.RegisterAccount("+992935811031")
	if err != nil {
		t.Error(err)
		return
	}

	got, err := s.FindAccountByID(account.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(got, account) {
		t.Errorf("Result invaled: got - %v, account - %v", got, account)
	}
}

func TestService_FindAccountByID_fail(t *testing.T) {
	s := &Service{}
	_, err := s.RegisterAccount("+992935811031")
	if err != nil {
		t.Error(err)
		return
	}

	newID := int64(uuid.New().ID())
	_, err = s.FindAccountByID(newID)
	if err == nil {
		t.Error(err)
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	got, err := s.FindPaymetByID(payment.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(got, payment) {
		t.Errorf("wrong payment returned: got - %v, payment - %v", got, payment)
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	s := newTestService()
 	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FindPaymetByID(uuid.New().String())
	if err == nil {
		t.Error(err)
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("must return ErrPaymentNotFound: %v", err)
	}
}

func TestService_Reject_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Error(err)
		return
	}

	savedPayment, err := s.FindPaymetByID(payment.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if savedPayment.Status != types.PaymentStatusFail {
		t.Error(savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Error(err)
		return
	}

	if savedAccount.Balance != defaultTestAccount.balance {
		t.Error(err)
		return
	}

}

func TestService_Reject_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.Reject(uuid.New().String())
	if err != ErrPaymentNotFound {
		t.Error(err)
		return
	}

}

func TestService_Repeat_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	savedPayment, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		t.Error(err)
		return
	}

	savedRepeat, err := s.Repeat(payment.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(savedPayment.Amount, savedRepeat.Amount) {
		t.Errorf("inviled result: payment - %v, repeat - %v", savedPayment, savedRepeat)
		return
	}
}

func TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	savedFavorite, err := s.FavoritePayment(payment.ID, "auto")
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(payment.Amount, savedFavorite.Amount) {
		t.Errorf("inviled result: payAmount = %v, saveFavor = %v", payment.Amount, savedFavorite.Amount)
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	savedFavorite, err := s.FavoritePayment(payment.ID, "auto")
	if err != nil {
		t.Error(err)
		return
	}

	savedPayFav, err := s.PayFromFavorite(savedFavorite.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(savedPayFav.Amount, payment.Amount) {
		t.Errorf("Inviled resulte: sPF - %v, PA - %v", savedPayFav, payment)
		return
	}
}

func TestService_PayFromFavorite_fail(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	_, err = s.FavoritePayment(payment.ID, "auto")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.PayFromFavorite(uuid.New().String())
	if err != ErrFavoriteNotFound {
		t.Error(err)
		return
	}

}

func TestService_PayFromFavorite_rules(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return 
	}

	payment, err := s.FavoritePayment(uuid.New().String(), "megafon")
	if err != ErrPaymentNotFound {
		t.Error(err)
		return
	}

	fmt.Println(payment)
}
