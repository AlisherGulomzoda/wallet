package wallet

import (
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

func TestService_Reject_success(t *testing.T) {
	s := &Service{}
	account, err := s.RegisterAccount("+992935811032")
	if err != nil {
		t.Error(err)
		return
	}

	err = s.Deposit(account.ID, 10_000_00)
	if err != nil {
		t.Error(err)
		return
	}

	payment, err := s.Pay(account.ID, 10_00, "Coffee")
	if err != nil {
		t.Error(err)
		return
	}

	err = s.Reject(payment.ID)
	if err != nil {
		t.Error(err)
		return
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


