package main

import (
	"fmt"
	"io"
	"log"
	"os"

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

	data := byte(0b0100_0111)
	fmt.Println(int(data))
	fmt.Println(string(data))

	fmt.Println(account.Balance)

	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}

	log.Print(wd)

	file, err := os.Open("data/readme.txt")
	if err != nil {
		log.Print(err)
		return
	}

	defer func () {		
		if err := file.Close(); err != nil {
			log.Print(err)
		}
	}()

	log.Printf("#%v", file)

	content := make([]byte, 0)
	buf := make([]byte, 4096)

	for{
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}

		if err != nil {
			log.Print(err)
			return
		}	
		content = append(content, buf[:read]...)
	}

	dat := string(content)
	log.Print(dat)

	newFile, err := os.Create("data/message.txt")
	if err != nil {
		log.Print(err)
		return
	}
	defer func (){
		if cerr := newFile.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	_, err = newFile.Write([]byte("Hello from GO!"))
	if err != nil {
		log.Print(err)
		return
	}

}
