package GroupieTracker

import (
	"encoding/json"
	"io/ioutil"
)

type Account struct {
	Name     string
	Password string
	Mail     string
}

func CreateAccount(name, pwd, mail, id string, Acc *Account) {
	Acc.Name = name
	Acc.Password = pwd
	Acc.Mail = mail
	b, _ := json.Marshal(Acc)
	ioutil.WriteFile("./GroupieTracker/Account/Login/"+id+".json", b, 0644)
}
