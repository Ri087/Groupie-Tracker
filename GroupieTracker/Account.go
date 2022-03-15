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
	EveryUserInfos map[int]InfosUser
	EveryId        map[string]int
	User           InfosUser
	Friend         UserFriend
	PseudoCheck    PseudoCheck
	FriendCheck    FriendCheck
	AuthorizeVisit AuthorizeVisit
	IdUsers        int
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

type InfosUser struct {
	Public bool
	Pseudo string
	Mail   string
	// SelectArtistes map[int]Artists
	Friends             map[int]UserFriend
	ShowFriendsToPublic bool
	ShowFriendsToFriend bool
}

type UserFriend struct {
	Pseudo string
	Mail   string
}

type PseudoCheck struct {
	PseudoNotOk bool
	WrongPseudo bool
}

type FriendCheck struct {
	WrongFriend bool
	ThatsU      bool
}

type AuthorizeVisit struct {
	User        InfosUser
	ShowFriends bool
	Authorize   bool
	Existant    bool
}

func AccStructureInit() *AccountStruct {
	AccStruct := &AccountStruct{}
	GoodCreationReset(AccStruct)
	AccStruct.EveryAccount = TakeEveryAccount()
	AccStruct.AuthToken = TakeEveryToken()
	AccStruct.EveryUserInfos = TakeEveryInfosUser()
	AccStruct.EveryId = TakeEveryId()
	AccStruct.IdUsers = len(AccStruct.EveryAccount)
	AccStruct.PseudoCheck = PseudoCheck{false, false}
	AccStruct.FriendCheck = FriendCheck{false, false}
	return AccStruct
}

func VisitAuthorizeReset(AccStruct *AccountStruct) {
	AccStruct.AuthorizeVisit.ShowFriends = false
	AccStruct.AuthorizeVisit.Authorize = false
	AccStruct.AuthorizeVisit.Existant = false
	AccStruct.AuthorizeVisit.User = InfosUser{}
}

func PseudoAndFriendReset(AccStruct *AccountStruct) {
	AccStruct.PseudoCheck.PseudoNotOk = false
	AccStruct.PseudoCheck.WrongPseudo = false
	AccStruct.FriendCheck.ThatsU = false
	AccStruct.FriendCheck.WrongFriend = false
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
	AccStruct.IdUsers += 1
	AccStruct.EveryUserInfos[AccStruct.IdUsers] = InfosUser{false, "", mail, map[int]UserFriend{}, false, false}
	c, _ := json.Marshal(AccStruct.EveryUserInfos)
	ioutil.WriteFile("./GroupieTracker/Account/InfoUsers.json", c, 0644)
	AccStruct.EveryId[mail] = AccStruct.IdUsers
	d, _ := json.Marshal(AccStruct.EveryId)
	ioutil.WriteFile("./GroupieTracker/Account/IdByMail.json", d, 0644)
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

func TakeEveryInfosUser() map[int]InfosUser {
	save, err := ioutil.ReadFile("./GroupieTracker/Account/InfoUsers.json")
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
		os.Exit(1)
	}
	var EveryInfosUser map[int]InfosUser
	err = json.Unmarshal(save, &EveryInfosUser)
	if err != nil {
		log.Fatalf("Error with the file: %s", err)
		os.Exit(1)
	}
	return EveryInfosUser
}

func TakeEveryId() map[string]int {
	save, err := ioutil.ReadFile("./GroupieTracker/Account/IdByMail.json")
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
		os.Exit(1)
	}
	var EveryId map[string]int
	err = json.Unmarshal(save, &EveryId)
	if err != nil {
		log.Fatalf("Error with the file: %s", err)
		os.Exit(1)
	}
	return EveryId
}

func AuthentificationToken(mail string, AccStruct *AccountStruct, w http.ResponseWriter) {
	TabForToken := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	Token := ""
	for i := 0; i < 33; i++ {
		Token += TabForToken[rand.Intn(len(TabForToken))]
	}
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

func DeleteToken(value string, AccStruct *AccountStruct) {
	delete(AccStruct.AuthToken, value)
	b, _ := json.Marshal(AccStruct.AuthToken)
	ioutil.WriteFile("./GroupieTracker/Account/Token.json", b, 0644)
}

func GetUserInfos(value string, AccStruct *AccountStruct) {
	GetUserById(AccStruct.EveryId[AccStruct.AuthToken[value]], AccStruct)
}

func GetUserById(id int, AccStruct *AccountStruct) {
	AccStruct.User = AccStruct.EveryUserInfos[id]
}

func SavePseudo(value, pseudo string, AccStruct *AccountStruct) {
	GetUserInfos(value, AccStruct)
	if AccStruct.User.Pseudo != "" {
		return
	}
	AccStruct.User.Pseudo = pseudo
	id := AccStruct.EveryId[AccStruct.User.Mail]
	AccStruct.EveryUserInfos[id] = AccStruct.User
	c, _ := json.Marshal(AccStruct.EveryUserInfos)
	ioutil.WriteFile("./GroupieTracker/Account/InfoUsers.json", c, 0644)
}

func GetFriendById(id int, AccStruct *AccountStruct) {
	user := AccStruct.EveryUserInfos[id]
	AccStruct.Friend.Mail = user.Mail
	AccStruct.Friend.Pseudo = user.Pseudo
}

func AddFriend(id int, AccStruct *AccountStruct) {
	AccStruct.User.Friends[id] = AccStruct.Friend
	c, _ := json.Marshal(AccStruct.EveryUserInfos)
	ioutil.WriteFile("./GroupieTracker/Account/InfoUsers.json", c, 0644)
}

func VisitProfil(AccStruct *AccountStruct) {
	if AccStruct.AuthorizeVisit.User.Public {
		AccStruct.AuthorizeVisit.Authorize = true
	} else if id := AccStruct.EveryId[AccStruct.User.Mail]; AccStruct.User.Mail != "" && AccStruct.AuthorizeVisit.User.Friends[id].Mail == AccStruct.User.Mail {
		AccStruct.AuthorizeVisit.Authorize = true
	}
}

func ShowedFriends(AccStruct *AccountStruct) {
	if AccStruct.AuthorizeVisit.User.ShowFriendsToPublic {
		AccStruct.AuthorizeVisit.ShowFriends = true
	} else if id := AccStruct.EveryId[AccStruct.User.Mail]; AccStruct.User.Mail != "" && AccStruct.AuthorizeVisit.User.ShowFriendsToFriend && AccStruct.AuthorizeVisit.User.Friends[id].Mail == AccStruct.User.Mail {
		AccStruct.AuthorizeVisit.ShowFriends = true
	}
}
