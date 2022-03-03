package GroupieTracker

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Account struct {
	Name     string
	Password []byte
	Mail     string
}

func CreateAccount(name, pwd, mail, id string, Acc *Account) {
	Acc.Name = name
	Acc.Password = Cryptage(pwd)
	Acc.Mail = mail
	b, _ := json.Marshal(Acc)
	fmt.Println(id)
	ioutil.WriteFile("./GroupieTracker/Account/Login/"+id+".json", b, 0644)
}

func Cryptage(pwd string) []byte {
	crypt := sha256.New()
	crypt.Write([]byte(pwd))
	return crypt.Sum(nil)
}
