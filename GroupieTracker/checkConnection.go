package GroupieTracker

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type CheckCo struct {
	Mail bool
	Pwd  bool
}

func CheckConnection(mail, pwd string, CC *CheckCo, Acc *Account) bool {
	ID := string(base64.StdEncoding.EncodeToString(Cryptage(mail)))
	files, err := ioutil.ReadDir("./GroupieTracker/Account/Login/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.Name() == ID+".json" {
			LoginAcc(ID, Acc)
			if !bytes.Equal(Acc.Password, Cryptage(pwd)) {
				CC.Pwd = true
				return false
			}
			return true
		}
	}
	CC.Mail = true
	return false
}

func LoginAcc(ID string, Acc *Account) {
	save, err := ioutil.ReadFile("./GroupieTracker/Account/Login/" + ID + ".json")
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
		os.Exit(1)
	}
	err = json.Unmarshal(save, Acc)
	if err != nil {
		log.Fatalf("Error with the file: %s", err)
		os.Exit(1)
	}
}
