package main

import (
	"fmt"
	"github.com/AlisherGulomzoda/wallet/pkg/wallet"
)


func main() {
	svc := &wallet.Service{}

	account, err := svc.RegisterAccount("+992935811031")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = svc.Deposit(account.ID, 10)
	if err != nil {
		switch err {
		case wallet.ErrAccountNotFound:
			fmt.Println(wallet.ErrAccountNotFound)
		case wallet.ErrPhoneAlreadyRegitered:
			fmt.Println(wallet.ErrPhoneAlreadyRegitered)
		case wallet.ErrAmountMustGreateZero:
			fmt.Println(wallet.ErrAmountMustGreateZero)
		}
		return
	}

	fmt.Println(account.Balance)

}