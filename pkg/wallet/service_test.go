package wallet

import (
	// "github.com/google/uuid"
	"fmt"
	"reflect"
	"testing"

	"github.com/AlisherGulomzoda/wallet/pkg/types"
)

type testService struct {
	*Service
}

func newTestServicw() *testService {
	return &testService{Service : &Service{}}
}

func (s *testService) addAccountWithBalancce(phone types.Phone, balance types.Money) (*types.Account, error) {
	account, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, fmt.Errorf("can't regiter account, error = %v", err)
	}

	err = s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("can't deposit account, %v", err)
	}

	return account, nil

}

func TestService_Register(t *testing.T) {
	svc := Service{}
	_, err := svc.RegisterAccount("+992000000000")
	if err != nil {
		t.Error(err)
	}

	_, err = svc.RegisterAccount("+992000000000")
	if err != ErrPhoneAlreadyRegitered {
		t.Error(err)
	}
}

func TestService_Deposit(t *testing.T) {
	svc := Service{}
	err := svc.Deposit(1, 0)
	if err != ErrAmountMustGreateZero {
		t.Error(err)
	}

	err = svc.Deposit(1, 1)
	if err != ErrAccountNotFound {
		t.Error(err)
	}

	account, err := svc.RegisterAccount("+992000000010")
	if err != nil {
		t.Error(err)
	}

	err = svc.Deposit(account.ID, 1)
	if err != nil {
		t.Error(err)
	}
}

func TestService_Pay(t *testing.T) {
	svc := Service{}
	_, err := svc.Pay(1, 0, "auto")
	if err != ErrAmountMustGreateZero {
		t.Error(err)
	}

	_, err = svc.Pay(1, 1, "auto")
	if err != ErrAccountNotFound {
		t.Error(err)
	}

	account, err := svc.RegisterAccount("+992000000000")
	if err != nil {
		t.Error(err)
	}

	_, err = svc.Pay(account.ID, 1, "auto")
	if err != ErrBalanceNotAmount {
		t.Error(err)
	}

	account.Balance = 100

	_, err = svc.Pay(account.ID, 100, "auto")
	if err != nil {
		t.Error(err)
	}
}

func TestService_FindbyAccountById_success(t *testing.T) {
	svc := Service{}
	account, err := svc.RegisterAccount("+992000000000")
	if err != nil {
		t.Error(err)
	}

	_, err = svc.FindAccountByID(account.ID)
	if err != nil {
		t.Error(err)
	}
}

func TestService_FindByAccountByID_notFound(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+992000000000")
	_, err := svc.FindAccountByID(2)
	// тут даст false, так как err (уже имеет что то внутри)
	if err != ErrAccountNotFound {
		t.Error(err)
	}
}

func TestFindPaymentByID_success(t *testing.T) {
	svc := &Service{}

	phone := types.Phone("+992000000000")

	account, err := svc.RegisterAccount(phone)
	if err != nil {
		t.Error(err)
		return
	}

	err = svc.Deposit(account.ID, 1000)
	if err != nil {
		t.Error(err)
		return
	}

	pay, err := svc.Pay(account.ID, 500, "auto")
	if err != nil {
		t.Error(err)
		return
	}

	got, err := svc.FindPaymetByID(pay.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(got, pay) {
		t.Error(err)
		return
	}
}

func TestService_Repeat_success(t *testing.T) {
	s := newTestServicw()
	account, err := s.addAccountWithBalancce("+992935811034", 10_000_00)
	if err != nil {
		t.Error(err)
	}

	payment, err := s.Pay(account.ID, 1_000_00, "auto")
	if err != nil {
		t.Errorf("FindPaymentByID(): can't create payment, error = %v", err)
	}

	payment, err = s.Repeat(payment.ID)
	if err != nil {
		t.Error(err)
	}
}

// func TestService_FavoritePayment_success(t *testing.T) {
// 	s := newTestServicw()
// 	account, err := s.addAccountWithBalancce("+992935811036", 10_000_00)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	payment, err := s.Pay(account.ID, 1_00, "card")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	favorite, err := s.FavoritePayment(payment.ID, "car")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	if !reflect.DeepEqual(payment.Amount, favorite.Amount) {
// 		t.Errorf("Inviled resulte: payment.Account %v, favorite.Amount %v", payment.Amount, favorite.Amount)
// 	}

// }