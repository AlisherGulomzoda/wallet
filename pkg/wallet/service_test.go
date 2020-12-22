package wallet

import (
	"testing"
)

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