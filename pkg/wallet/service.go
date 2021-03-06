package wallet

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/AlisherGulomzoda/wallet/pkg/types"
	"github.com/google/uuid"
)

var (
	ErrPhoneAlreadyRegitered = errors.New("phone already regitered")
	ErrAmountMustGreateZero  = errors.New("amount must be greated then zero")
	ErrAccountNotFound       = errors.New("account not found")
	ErrBalanceNotAmount      = errors.New("balance little then amount")
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrFavoriteNotFound      = errors.New("favorit not found")
	ErrFileNotFound			 = errors.New("file not found")
)

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
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

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategoty) (*types.Payment, error) {
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
		ID:        paymentID,
		AccountID: account.ID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
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
	for _, pay := range s.payments {
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

func (s *Service) ExportToFile(path string) error {
	wd, err := os.Getwd()
	if err != nil {
		log.Print(err)
	}

	log.Print(wd)

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	result := ""
	for _, account := range s.accounts {
		result += strconv.FormatInt(int64(account.ID), 10) + ";"
		result += string(account.Phone) + ";"
		result += strconv.FormatInt(int64(account.Balance), 10) + "|"
	}

	err = ioutil.WriteFile(path, []byte(result), 0666)
	if err != nil {
		return err
	}

	return nil

}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	content := make([]byte, 0)
	buf := make([]byte, 4096)

	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}

		if err != nil {
			log.Print(err)
			return err
		}

		content = append(content, buf[:read]...)
	}

	data := string(content)
	splitSlice := strings.Split(data, "|")

	for _, split := range splitSlice {
		if split != "" {
			datas := strings.Split(split, ";")

			id, err := strconv.Atoi(datas[0])
			if err != nil {
				log.Println(err)
				return err
			}

			balance, err := strconv.Atoi(datas[2])
			if err != nil {
				log.Println(err)
				return err
			}

			newAccount := &types.Account{
				ID:      int64(id),
				Phone:   types.Phone(datas[1]),
				Balance: types.Money(balance),
			}

			s.accounts = append(s.accounts, newAccount)
		}
	}

	return nil

}

func (s *Service) Export(dir string) error {
	if s.accounts != nil {
		result := ""
		for _, account := range s.accounts {
			result += strconv.Itoa(int(account.ID)) + ";"
			result += string(account.Phone) + ";"
			result += strconv.Itoa(int(account.Balance)) + "\n"
		}

		err := actionByFile(dir+"/accounts.dump", result)
		if err != nil {
			return err
		}
	}

	if s.payments != nil {
		result := ""
		for _, payment := range s.payments {
			result += payment.ID + ";"
			result += strconv.Itoa(int(payment.AccountID)) + ";"
			result += strconv.Itoa(int(payment.Amount)) + ";"
			result += string(payment.Category) + ";"
			result += string(payment.Status) + "\n"
		}

		err := actionByFile(dir+"/payments.dump", result)
		if err != nil {
			return err
		}
	}

	if s.favorites != nil {
		result := ""
		for _, favorite := range s.favorites {
			result += favorite.ID + ";"
			result += strconv.Itoa(int(favorite.AccountID)) + ";"
			result += favorite.Name + ";"
			result += strconv.Itoa(int(favorite.Amount)) + ";"
			result += string(favorite.Category) + "\n"
		}

		err := actionByFile(dir+"/favorites.dump", result)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) actionByAccounts(path string) error {
	byteData, err := ioutil.ReadFile(path)
	if err == nil {
		datas := string(byteData)
		splits := strings.Split(datas, "\n")

		for _, split := range splits {
			if len(split) == 0 {
				break
			}

			data := strings.Split(split, ";")

			id, err := strconv.Atoi(data[0])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			phone := types.Phone(data[1])

			balance, err := strconv.Atoi(data[2])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			account, err := s.FindAccountByID(int64(id))
			if err != nil {
				acc, err := s.RegisterAccount(phone)
				if err != nil {
					log.Println("err from register account")
					return err
				}

				acc.Balance = types.Money(balance)
			} else {
				account.Phone = phone
				account.Balance = types.Money(balance)
			}
		}
	} else {
		log.Println(ErrFileNotFound.Error())
	}

	return nil
}

func (s *Service) actionByPayments(path string) error {
	byteData, err := ioutil.ReadFile(path)
	if err == nil {
		datas := string(byteData)
		splits := strings.Split(datas, "\n")

		for _, split := range splits {
			if len(split) == 0 {
				break
			}

			data := strings.Split(split, ";")
			id := data[0]

			accountID, err := strconv.Atoi(data[1])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			amount, err := strconv.Atoi(data[2])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			category := types.PaymentCategoty(data[3])

			status := types.PaymentStatus(data[4])

			payment, err := s.FindPaymetByID(id)
			if err != nil {
				newPayment := &types.Payment{
					ID:        id,
					AccountID: int64(accountID),
					Amount:    types.Money(amount),
					Category:  types.PaymentCategoty(category),
					Status:    types.PaymentStatus(status),
				}

				s.payments = append(s.payments, newPayment)
			} else {
				payment.AccountID = int64(accountID)
				payment.Amount = types.Money(amount)
				payment.Category = category
				payment.Status = status
			}
		}
	} else {
		log.Println(ErrFileNotFound.Error())
	}

	return nil
}

func (s *Service) actionByFavorites(path string) error {
	byteData, err := ioutil.ReadFile(path)
	if err == nil {
		datas := string(byteData)
		splits := strings.Split(datas, "\n")

		for _, split := range splits {
			if len(split) == 0 {
				break
			}

			data := strings.Split(split, ";")
			id := data[0]

			accountID, err := strconv.Atoi(data[1])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			name := data[2]

			amount, err := strconv.Atoi(data[3])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			category := types.PaymentCategoty(data[4])

			favorite, err := s.FindFavoriteByID(id)
			if err != nil {
				newFavorite := &types.Favorite{
					ID:        id,
					AccountID: int64(accountID),
					Name:      name,
					Amount:    types.Money(amount),
					Category:  types.PaymentCategoty(category),
				}

				s.favorites = append(s.favorites, newFavorite)
			} else {
				favorite.AccountID = int64(accountID)
				favorite.Name = name
				favorite.Amount = types.Money(amount)
				favorite.Category = category
			}
		}
	} else {
		log.Println(ErrFileNotFound.Error())
	}

	return nil
}

func (s *Service) FindFavoriteByID(id string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == id {
			return favorite, nil
		}
	}

	return nil, ErrFavoriteNotFound
}

func actionByFile(path, data string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = file.WriteString(data)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
	if len(payments) == 0 {
		log.Print(ErrPaymentNotFound)
		return nil
	}

	if len(payments) <= records {
		result := ""
		for _, payment := range payments {
			result += payment.ID + ";"
			result += strconv.Itoa(int(payment.AccountID)) + ";"
			result += strconv.Itoa(int(payment.Amount)) + ";"
			result += string(payment.Category) + ";"
			result += string(payment.Status) + "\n"
		}

		err := actionByFile(dir+"/payments.dump", result)
		if err != nil {
			return err
		}

		return nil
	}

	result := ""
	k := 1
	for i, payment := range payments {
		result += payment.ID + ";"
		result += strconv.Itoa(int(payment.AccountID)) + ";"
		result += strconv.Itoa(int(payment.Amount)) + ";"
		result += string(payment.Category) + ";"
		result += string(payment.Status) + "\n"

		if (i+1)%records == 0 {
			err := actionByFile(dir+"/payments"+strconv.Itoa(k)+".dump", result)
			if err != nil {
				return err
			}
			k++
			result = ""
		}
	}

	if result != "" {
		err := actionByFile(dir+"/payments"+strconv.Itoa(k)+".dump", result)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) SumPayments(goroutines int) types.Money {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	var summ types.Money = 0
	if goroutines == 0 || goroutines == 1 {
		wg.Add(1)
		go func(payments []*types.Payment) {
			defer wg.Done()
			for _, payment := range payments {
				summ += payment.Amount
			}
		}(s.payments)
	} else {
		from := 0
		count := len(s.payments) / goroutines
		for i := 1; i <= goroutines; i++ {
			wg.Add(1)
			last := len(s.payments) - i*count
			if i == goroutines {
				last = 0
			}
			to := len(s.payments) - last
			go func(payments []*types.Payment) {
				defer wg.Done()
				s := types.Money(0)
				for _, payment := range payments {
					s += payment.Amount
				}
				mu.Lock()
				defer mu.Unlock()
				summ += s
			}(s.payments[from:to])
			from += count
		}
	}

	wg.Wait()

	return summ
}

func (s *Service) FilterPayments(accountID int64, goroutines int) ([]types.Payment, error) {
	filteredPayments := []types.Payment{}
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	if goroutines == 0 || goroutines == 1 {
		wg.Add(1)
		go func(payments []*types.Payment) {
			defer wg.Done()
			for _, payment := range payments {
				if payment.AccountID == accountID {
					filteredPayments = append(filteredPayments, types.Payment{
						ID:        payment.ID,
						AccountID: payment.AccountID,
						Amount:    payment.Amount,
						Category:  payment.Category,
						Status:    payment.Status,
					})
				}
			}
		}(s.payments)
	} else {
		from := 0
		count := len(s.payments) / goroutines
		for i := 1; i <= goroutines; i++ {
			wg.Add(1)
			last := len(s.payments) - i*count
			if i == goroutines {
				last = 0
			}
			to := len(s.payments) - last
			go func(payments []*types.Payment) {
				defer wg.Done()
				separetePayments := []types.Payment{}
				for _, payment := range payments {
					if payment.AccountID == accountID {
						separetePayments = append(separetePayments, types.Payment{
							ID:        payment.ID,
							AccountID: payment.AccountID,
							Amount:    payment.Amount,
							Category:  payment.Category,
							Status:    payment.Status,
						})
					}
				}
				mu.Lock()
				defer mu.Unlock()
				filteredPayments = append(filteredPayments, separetePayments...)
			}(s.payments[from:to])
			from += count
		}
	}

	wg.Wait()

	if len(filteredPayments) == 0 {
		return nil, ErrAccountNotFound
	}

	return filteredPayments, nil
}

func (s *Service) FilterPaymentsByFn(filter func(payment types.Payment) bool, goroutines int) ([]types.Payment, error) {
	filteredPayments := []types.Payment{}
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	if goroutines == 0 || goroutines == 1 {
		wg.Add(1)
		go func(payments []*types.Payment) {
			defer wg.Done()
			for _, payment := range payments {
				if filter(*payment) {
					filteredPayments = append(filteredPayments, *payment)
				}
			}
		}(s.payments)
	} else {
		from := 0
		count := len(s.payments) / goroutines
		for i := 1; i <= goroutines; i++ {
			wg.Add(1)
			last := len(s.payments) - i*count
			if i == goroutines {
				last = 0
			}
			to := len(s.payments) - last
			go func(payments []*types.Payment) {
				defer wg.Done()
				separetePayments := []types.Payment{}
				for _, payment := range payments {
					if filter(*payment) {
						separetePayments = append(separetePayments, *payment)
					}
				}
				mu.Lock()
				defer mu.Unlock()
				filteredPayments = append(filteredPayments, separetePayments...)
			}(s.payments[from:to])
			from += count
		}
	}

	wg.Wait()

	if len(filteredPayments) == 0 {
		return nil, ErrAccountNotFound
	}

	return filteredPayments, nil
}

func (s *Service) SumPaymentsWithProgress() <- chan types.Progress {
	size := 100_0000

	amountOfMoney := make([]types.Money, 0)
	for _, pay := range s.payments {
		amountOfMoney = append(amountOfMoney, pay.Amount)
	}

	wg := sync.WaitGroup{}
	goroutines := (len(amountOfMoney) + 1) / size
	ch := make(chan types.Progress)
	if goroutines <= 0 {
		goroutines = 1
	}
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(ch chan<- types.Progress, amountOfMoney []types.Money, part int) {
			sum := 0
			defer wg.Done()
			for _, val := range amountOfMoney {
				sum += int(val)

			}
			ch <- types.Progress{
				Result: types.Money(sum),
			}
		}(ch, amountOfMoney, i)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	return ch
}
