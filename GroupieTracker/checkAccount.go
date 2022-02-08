package GroupieTracker

import (
	"crypto/sha256"
	"io/ioutil"
	"log"
	"strconv"
)

func CheckAccount(name, pwd, mail string, CC *CheckCreation, Acc *Account) bool {
	ID := IDMail(mail)
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

func IDMail(mail string) string {
	var ID string
	for _, i := range sha256.Sum256([]byte(mail)) {
		ID += strconv.Itoa(int(i))
	}
	return ID
}
