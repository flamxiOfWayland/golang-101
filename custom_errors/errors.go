package customerrors

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrorNoBitcoin = errors.New("no bitcoin wallet found")
	ErrorNoMoney   = fmt.Errorf("no money")
)

func GetAllTheBitcoins() (string, error) {
	stat, err := os.Stat("/bitcoin/wallet")
	if err != nil {
		return "", errors.Join(ErrorNoBitcoin, err)
	}
	return stat.Name(), nil
}

func AreWeRich() error {
	defer func() {
		fmt.Println("we don't give a fuck either way")
	}()

	fmt.Println("we are rich if we have a lot of bitcoins")
	wallet, err := GetAllTheBitcoins()
	if err != nil {
		return errors.Join(ErrorNoMoney, err)
	}

	fmt.Printf("we are rich, our wallet: %s\n", wallet)

	return nil
}
