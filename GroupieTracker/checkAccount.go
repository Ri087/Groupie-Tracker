package GroupieTracker

import (
	"encoding/base32"
	"io/ioutil"
	"log"
)

func CheckAccount(name, pwd, mail string, CC *CheckCreation, Acc *Account) bool {
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
