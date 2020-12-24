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
	ErrPaymentNotFound = errors.New("payment not found")
	ErrFavoriteNotFound = errors.New("favorit not found")
)

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites	  []*types.Favorite
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

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
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

	return account, nil
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymetByID(paymentID)
	if err != nil {
		return err
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == payment.AccountID {
			account = acc
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount 
	payment.Amount = 0
	
	return nil

}

func (s *Service) FindPaymetByID(paymentID string) (*types.Payment, error) {
	for _, pay := range s.payments {
		if pay.ID == paymentID {
			return pay, nil
		}
	}

	return nil, ErrPaymentNotFound

}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	var payment *types.Payment
	for _, pay := range s.payments{
		if pay.ID == paymentID {
			payment = pay
			break
		}
	}

	payment, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	return payment, nil

}

func (s *Service) findPaymentAndAccountByPaymentID(paymentID string) (*types.Payment, *types.Account, error) {
	payment, err := s.FindPaymetByID(paymentID)
	if err != nil {
		return nil, nil, err
	}

	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return nil, nil, err
	}

	return payment, account, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	targetPayment, targetAccount, err := s.findPaymentAndAccountByPaymentID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite := &types.Favorite{
		ID:        uuid.New().String(),
		AccountID: targetAccount.ID,
		Name:      name,
		Amount:    targetPayment.Amount,
		Category:  targetPayment.Category,
	}

	s.favorites = append(s.favorites, favorite)

	return favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	var favorite *types.Favorite
	for _, fav := range s.favorites {
		if fav.ID == favoriteID {
			favorite = fav
			break
		}
	}

	if favorite == nil {
		return nil, ErrFavoriteNotFound
	}

	var payment *types.Payment
	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category) 
	if err != nil {
		return nil, err
	}

	return payment, nil
}