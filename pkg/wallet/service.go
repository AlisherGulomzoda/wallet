package wallet

import (
	"errors"

	"github.com/AlisherGulomzoda/wallet/pkg/types"
	"github.com/google/uuid"
)

var (
	ErrPhoneAlreadyRegitered = errors.New("phone already regitered")
	ErrAmountMustGreateZero = errors.New("amount must be greated then zero")
	ErrAccountNotFound = errors.New("account not found")
	ErrBalanceNotAmount = errors.New("balance little then amount")
)

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneAlreadyRegitered
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil

}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustGreateZero
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound 
	}

	account.Balance += amount

	return nil

}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymrntCategoty) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustGreateZero
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrBalanceNotAmount
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID: paymentID,
		AccountID: account.ID,
		Amount: amount,
		Category: category,
		Status: types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)

	return payment, nil

}
