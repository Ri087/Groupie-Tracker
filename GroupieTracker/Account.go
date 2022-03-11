package GroupieTracker

import (
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type AccountStruct struct {
	GoodCreation   GoodCreation
	GoodConnection GoodConnection
	EveryAccount   map[string][]byte
	AuthToken      map[string]string
}

type GoodCreation struct {
	Mail  bool
	Pwd   bool
	Pwdc  bool
	Exist bool
}

type GoodConnection struct {
	Mail bool
	Pwd  bool
}

func AccStructureInit() *AccountStruct {
	AccStruct := &AccountStruct{}
	GoodCreationReset(AccStruct)
	AccStruct.EveryAccount = TakeEveryAccount()
	AccStruct.AuthToken = TakeEveryToken()
	return AccStruct
}

func VerifEntryUser(mail, pwd, pwdc string, AccStruct *AccountStruct) bool {
	var a bool
	var dot bool
	for i, c := range mail {
		if i == 0 && string(c) == "@" {
			AccStruct.GoodCreation.Mail = true
			return false
		}
		if string(c) == "@" {
			a = true
		}
		if a && string(c) == "." && i != len(mail)-1 {
			dot = true
		}
	}
	if !dot {
		AccStruct.GoodCreation.Mail = true
		return false
	}

	if len(pwd) < 6 {
		AccStruct.GoodCreation.Pwd = true
		return false
	}
	if pwd != pwdc {
		AccStruct.GoodCreation.Pwdc = true
		return false
	}
	AccStruct.EveryAccount = TakeEveryAccount()
	if len(AccStruct.EveryAccount[mail]) != 0 {
		AccStruct.GoodCreation.Exist = true
		return false
	}
	return true
}

func CreateAccount(mail, pwd string, AccStruct *AccountStruct) {
	AccStruct.EveryAccount[mail] = Cryptage(pwd)
	b, _ := json.Marshal(AccStruct.EveryAccount)
	ioutil.WriteFile("./GroupieTracker/Account/Acc.json", b, 0644)
}

func Cryptage(pwd string) []byte {
	crypt := sha256.New()
	crypt.Write([]byte(pwd))
	return crypt.Sum(nil)
}

func TakeEveryToken() map[string]string {
	save, err := ioutil.ReadFile("./GroupieTracker/Account/Token.json")
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
		os.Exit(1)
	}
	var EveryToken map[string]string
	err = json.Unmarshal(save, &EveryToken)
	if err != nil {
		log.Fatalf("Error with the file: %s", err)
		os.Exit(1)
	}
	return EveryToken
}

func TakeEveryAccount() map[string][]byte {
	save, err := ioutil.ReadFile("./GroupieTracker/Account/Acc.json")
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
		os.Exit(1)
	}
	var EveryAcc map[string][]byte
	err = json.Unmarshal(save, &EveryAcc)
	if err != nil {
		log.Fatalf("Error with the file: %s", err)
		os.Exit(1)
	}
	return EveryAcc
}

func AuthentificationToken(mail string, AccStruct *AccountStruct, w http.ResponseWriter) {
	TabForToken := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	Token := "FirstOne"
	rand.Seed(time.Now().UnixNano())
	for AccStruct.AuthToken[Token] != "" {
		Token = ""
		for i := 0; i < 33; i++ {
			Token += TabForToken[rand.Intn(len(TabForToken))]
		}
	}
	AccStruct.AuthToken[Token] = mail
	b, _ := json.Marshal(AccStruct.AuthToken)
	ioutil.WriteFile("./GroupieTracker/Account/Token.json", b, 0644)
	http.SetCookie(w, &http.Cookie{Name: "AUTHENTIFICATION_TOKEN", Value: Token})
}

func GoodCreationReset(AccStruct *AccountStruct) {
	AccStruct.GoodCreation.Mail = false
	AccStruct.GoodCreation.Pwd = false
	AccStruct.GoodCreation.Pwdc = false
	AccStruct.GoodCreation.Exist = false
}

func VerifConnectionUser(mail, pwd string, AccStruct *AccountStruct) bool {
	if len(AccStruct.EveryAccount[mail]) == 0 {
		AccStruct.GoodConnection.Mail = true
		return false
	}
	if VerifNotSamePwd(Cryptage(pwd), AccStruct.EveryAccount[mail]) {
		AccStruct.GoodConnection.Pwd = true
		return false
	}
	return true
}

func VerifNotSamePwd(pwd []byte, AccPwd []byte) bool {
	if len(pwd) != len(AccPwd) {
		return true
	}
	for k, i := range pwd {
		if i != AccPwd[k] {
			return true
		}
	}
	return false
}

func GoodConnectionReset(AccStruct *AccountStruct) {
	AccStruct.GoodConnection.Mail = false
	AccStruct.GoodConnection.Pwd = false
}
