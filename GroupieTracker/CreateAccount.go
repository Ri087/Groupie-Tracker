package GroupieTracker

import (
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"io/ioutil"
	"log"
)

type Account struct {
	Name     string
	Password []byte
	Mail     string
}

type CheckCreation struct {
	Name  bool
	Pwd   bool
	Pwdc  bool
	Mail  bool
	Exist bool
}

func CheckGoodCreation(name, pwd, pwdc, mail string, CC *CheckCreation, Acc *Account) bool {
	if len(name) < 3 {
		CC.Name = true
		return true
	}
	if len(pwd) < 6 {
		CC.Pwd = true
		return true
	}
	if pwd != pwdc {
		CC.Pwdc = true
		return true
	}
	var a bool
	var dot bool
	for i, c := range mail {
		if i == 0 && string(c) == "@" {
			CC.Mail = true
			return true
		}
		if string(c) == "@" {
			a = true
		}
		if a && string(c) == "." && i != len(mail)-1 {
			dot = true
		}
	}
	if !dot {
		CC.Mail = true
		return true
	}
	return IfExist(name, pwd, mail, CC, Acc)
}

func IfExist(name, pwd, mail string, CC *CheckCreation, Acc *Account) bool {
	ID := base32.StdEncoding.EncodeToString(Cryptage(mail))
	files, err := ioutil.ReadDir("./GroupieTracker/Account/Login/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.Name() == ID+".json" {
			CC.Exist = true
			return true
		}
	}
	CreateAccount(name, pwd, mail, ID, Acc)
	return false
}

func CreateAccount(name, pwd, mail, id string, Acc *Account) {
	Acc.Name = name
	Acc.Password = Cryptage(pwd)
	Acc.Mail = mail
	b, _ := json.Marshal(Acc)
	ioutil.WriteFile("./GroupieTracker/Account/Login/"+id+".json", b, 0644)
}

func Cryptage(pwd string) []byte {
	crypt := sha256.New()
	crypt.Write([]byte(pwd))
	return crypt.Sum(nil)
}
