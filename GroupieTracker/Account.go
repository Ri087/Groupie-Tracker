package GroupieTracker

import (
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type AccountStruct struct {
	AllId            map[string]string        // Mail TO Id
	AllToken         map[string]string        // Token TO Id
	AllHashtag       map[string]string        // EntireName TO Id
	AllTabToken      map[string][]string      // Id TO []Token
	AllMail          map[string]string        // Id TO Mail
	AllPassword      map[string][]byte        // Id TO Password
	AllAccount       map[string]ProfilAccount // Id TO UserInformations
	Creation         CreationAccount          // Verification create account
	Connection       ConnectionAccount        // Verification connection account
	ProfilParameters EveryProfilParameters    // Profil
}

func AccountStructureInit() *AccountStruct {
	AccStruct := &AccountStruct{}
	GetAllIdJson(AccStruct)
	GetAllMailJson(AccStruct)
	GetAllPasswordJson(AccStruct)
	GetAllTokenJson(AccStruct)
	GetAllTabTokenJson(AccStruct)
	GetAllAccount(AccStruct)
	GetAllHashtag(AccStruct)
	CreationAccountReset(AccStruct)
	ConnectionAccountReset(AccStruct)
	ProfilAccountReset(AccStruct)
	ProfilVisitReset(AccStruct)
	FriendsParametersReset(AccStruct)
	NameParametersReset(AccStruct)
	ArtistsProfilReset(AccStruct)
	return AccStruct
}

func GetAllIdJson(AccStruct *AccountStruct) {
	save, err := ioutil.ReadFile("./GroupieTracker/Account/AllId.json")
	if err != nil {
		ioutil.WriteFile("./GroupieTracker/Account/AllId.json", []byte("{}"), 0644)
		save, _ = ioutil.ReadFile("./GroupieTracker/Account/AllId.json")
	}
	json.Unmarshal(save, &AccStruct.AllId)
}

func SaveAllId(AccStruct *AccountStruct) {
	AllId, _ := json.Marshal(AccStruct.AllId)
	ioutil.WriteFile("./GroupieTracker/Account/AllId.json", AllId, 0644)
}

func GetAllMailJson(AccStruct *AccountStruct) {
	save, err := ioutil.ReadFile("./GroupieTracker/Account/AllMail.json")
	if err != nil {
		ioutil.WriteFile("./GroupieTracker/Account/AllMail.json", []byte("{}"), 0644)
		save, _ = ioutil.ReadFile("./GroupieTracker/Account/AllMail.json")
	}
	json.Unmarshal(save, &AccStruct.AllMail)
}

func SaveAllMail(AccStruct *AccountStruct) {
	AllMail, _ := json.Marshal(AccStruct.AllMail)
	ioutil.WriteFile("./GroupieTracker/Account/AllMail.json", AllMail, 0644)
}

func GetAllPasswordJson(AccStruct *AccountStruct) {
	save, err := ioutil.ReadFile("./GroupieTracker/Account/AllPassword.json")
	if err != nil {
		ioutil.WriteFile("./GroupieTracker/Account/AllPassword.json", []byte("{}"), 0644)
		save, _ = ioutil.ReadFile("./GroupieTracker/Account/AllPassword.json")
	}
	json.Unmarshal(save, &AccStruct.AllPassword)
}

func SaveAllPassword(AccStruct *AccountStruct) {
	AllPassword, _ := json.Marshal(AccStruct.AllPassword)
	ioutil.WriteFile("./GroupieTracker/Account/AllPassword.json", AllPassword, 0644)
}

func GetAllTokenJson(AccStruct *AccountStruct) {
	save, err := ioutil.ReadFile("./GroupieTracker/Account/AllToken.json")
	if err != nil {
		ioutil.WriteFile("./GroupieTracker/Account/AllToken.json", []byte("{}"), 0644)
		save, _ = ioutil.ReadFile("./GroupieTracker/Account/AllToken.json")
	}
	json.Unmarshal(save, &AccStruct.AllToken)
}

func SaveAllTabToken(AccStruct *AccountStruct) {
	AllTabToken, _ := json.Marshal(AccStruct.AllTabToken)
	ioutil.WriteFile("./GroupieTracker/Account/AllTabToken.json", AllTabToken, 0644)
}

func GetAllTabTokenJson(AccStruct *AccountStruct) {
	save, err := ioutil.ReadFile("./GroupieTracker/Account/AllTabToken.json")
	if err != nil {
		ioutil.WriteFile("./GroupieTracker/Account/AllTabToken.json", []byte("{}"), 0644)
		save, _ = ioutil.ReadFile("./GroupieTracker/Account/AllTabToken.json")
	}
	json.Unmarshal(save, &AccStruct.AllTabToken)
}

func SaveAllToken(AccStruct *AccountStruct) {
	AllToken, _ := json.Marshal(AccStruct.AllToken)
	ioutil.WriteFile("./GroupieTracker/Account/AllToken.json", AllToken, 0644)
}

func GetAllAccount(AccStruct *AccountStruct) {
	save, err := ioutil.ReadFile("./GroupieTracker/Account/AllAccount.json")
	if err != nil {
		ioutil.WriteFile("./GroupieTracker/Account/AllAccount.json", []byte("{}"), 0644)
		save, _ = ioutil.ReadFile("./GroupieTracker/Account/AllAccount.json")
	}
	json.Unmarshal(save, &AccStruct.AllAccount)
}

func SaveAllAccount(AccStruct *AccountStruct) {
	AllAccount, _ := json.Marshal(AccStruct.AllAccount)
	ioutil.WriteFile("./GroupieTracker/Account/AllAccount.json", AllAccount, 0644)
}

func GetAllHashtag(AccStruct *AccountStruct) {
	save, err := ioutil.ReadFile("./GroupieTracker/Account/AllHashtag.json")
	if err != nil {
		ioutil.WriteFile("./GroupieTracker/Account/AllHashtag.json", []byte("{}"), 0644)
		save, _ = ioutil.ReadFile("./GroupieTracker/Account/AllHashtag.json")
	}
	json.Unmarshal(save, &AccStruct.AllHashtag)
}

func SaveAllHashtag(AccStruct *AccountStruct) {
	AllHashtag, _ := json.Marshal(AccStruct.AllHashtag)
	ioutil.WriteFile("./GroupieTracker/Account/AllHashtag.json", AllHashtag, 0644)
}

func LoadUserByToken(Token string, AccStruct *AccountStruct) {
	AccStruct.ProfilParameters.Profil = AccStruct.AllAccount[AccStruct.AllToken[Token]]
}

//

// Account creation

type CreationAccount struct {
	UserInformation CreationUserInformation
	Verification    CreationVerification
}

type CreationUserInformation struct {
	Mail                 string
	Password             string
	PasswordConfirmation string
	PasswordEncrypted    []byte
}

type CreationVerification struct {
	Mail          bool
	Password      bool
	PasswordCheck bool
	AlreadyExist  bool
}

func CreationAccountReset(AccStruct *AccountStruct) {
	CreationUserInformationReset(AccStruct)
	CreationVerificationReset(AccStruct)
}

func CreationUserInformationReset(AccStruct *AccountStruct) {
	AccStruct.Creation.UserInformation.Mail = ""
	AccStruct.Creation.UserInformation.Password = ""
	AccStruct.Creation.UserInformation.PasswordConfirmation = ""
	AccStruct.Creation.UserInformation.PasswordEncrypted = []byte{}
}

func CreationVerificationReset(AccStruct *AccountStruct) {
	AccStruct.Creation.Verification.Mail = false
	AccStruct.Creation.Verification.Password = false
	AccStruct.Creation.Verification.PasswordCheck = false
	AccStruct.Creation.Verification.AlreadyExist = false
}

func CreationUserInformationFill(mail, pwd, pwdc string, AccStruct *AccountStruct) {
	AccStruct.Creation.UserInformation.Mail = mail
	AccStruct.Creation.UserInformation.Password = pwd
	AccStruct.Creation.UserInformation.PasswordConfirmation = pwdc
	AccStruct.Creation.UserInformation.PasswordEncrypted = Cryptage(pwd)
}

func Cryptage(pwd string) []byte {
	crypt := sha256.New()
	crypt.Write([]byte(pwd))
	return crypt.Sum(nil)
}

func CreationVerificationEntryUser(AccStruct *AccountStruct) bool {
	var a bool
	var dot bool
	for i, c := range AccStruct.Creation.UserInformation.Mail {
		if i == 0 && string(c) == "@" {
			AccStruct.Creation.Verification.Mail = true
			return false
		}
		if string(c) == "@" {
			a = true
		}
		if a && string(c) == "." && i != len(AccStruct.Creation.UserInformation.Mail)-1 {
			dot = true
		}
	}
	if !dot {
		AccStruct.Creation.Verification.Mail = true
		return false
	}

	if len(AccStruct.Creation.UserInformation.Password) < 6 {
		AccStruct.Creation.Verification.Password = true
		return false
	}
	if AccStruct.Creation.UserInformation.Password != AccStruct.Creation.UserInformation.PasswordConfirmation {
		AccStruct.Creation.Verification.PasswordCheck = true
		return false
	}
	if AccStruct.AllId[AccStruct.Creation.UserInformation.Mail] != "" {
		AccStruct.Creation.Verification.AlreadyExist = true
		return false
	}
	return true
}

func AccountCreation(AccStruct *AccountStruct) {
	rand.Seed(time.Now().UnixNano())
	TabNumbers := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	var Id string
	for Id == "" || AccStruct.AllMail[Id] != "" {
		Id = ""
		for i := 0; i < 16; i++ {
			Id += TabNumbers[rand.Intn(len(TabNumbers))]
		}
	}
	AccStruct.AllId[AccStruct.Creation.UserInformation.Mail] = Id
	SaveAllId(AccStruct)
	AccStruct.AllMail[Id] = AccStruct.Creation.UserInformation.Mail
	SaveAllMail(AccStruct)
	AccStruct.AllPassword[Id] = AccStruct.Creation.UserInformation.PasswordEncrypted
	SaveAllPassword(AccStruct)
	AccStruct.AllAccount[Id] = AccStruct.ProfilParameters.Profil
	SaveAllAccount(AccStruct)
}

//

// Connection account

type ConnectionAccount struct {
	UserInformation ConnectionUserInformation
	Verification    ConnectionVerification
}

type ConnectionUserInformation struct {
	Mail              string
	PasswordEncrypted []byte
}

type ConnectionVerification struct {
	Mail     bool
	Password bool
}

func ConnectionAccountReset(AccStruct *AccountStruct) {
	ConnectionUserInformationReset(AccStruct)
	ConnectionVerificationReset(AccStruct)
}

func ConnectionUserInformationReset(AccStruct *AccountStruct) {
	AccStruct.Connection.UserInformation.Mail = ""
	AccStruct.Connection.UserInformation.PasswordEncrypted = []byte{}
}

func ConnectionVerificationReset(AccStruct *AccountStruct) {
	AccStruct.Connection.Verification.Mail = false
	AccStruct.Connection.Verification.Password = false
}

func ConnectionUserInformationFill(mail, pwd string, AccStruct *AccountStruct) {
	AccStruct.Connection.UserInformation.Mail = mail
	AccStruct.Connection.UserInformation.PasswordEncrypted = Cryptage(pwd)
}

func ConnectionVerificationEntryUser(AccStruct *AccountStruct) bool {
	if AccStruct.AllId[AccStruct.Connection.UserInformation.Mail] == "" {
		AccStruct.Connection.Verification.Mail = true
		return false
	}
	if ConnectionVerifyPasswordNotSame(AccStruct.Connection.UserInformation.PasswordEncrypted, AccStruct.AllPassword[AccStruct.AllId[AccStruct.Connection.UserInformation.Mail]]) {
		AccStruct.Connection.Verification.Password = true
		return false
	}
	return true
}

func ConnectionVerifyPasswordNotSame(pwd []byte, AccPwd []byte) bool {
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

//

func SetAuthentificationToken(w http.ResponseWriter, AccStruct *AccountStruct) {
	rand.Seed(time.Now().UnixNano())
	TabForToken := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	var Token string
	for Token == "" || AccStruct.AllToken[Token] != "" {
		Token = ""
		for i := 0; i < 33; i++ {
			Token += TabForToken[rand.Intn(len(TabForToken))]
		}
	}
	if AccStruct.AllId[AccStruct.Creation.UserInformation.Mail] == "" {
		AccStruct.AllToken[Token] = AccStruct.AllId[AccStruct.Connection.UserInformation.Mail]
		AccStruct.AllTabToken[AccStruct.AllId[AccStruct.Connection.UserInformation.Mail]] = append(AccStruct.AllTabToken[AccStruct.AllId[AccStruct.Connection.UserInformation.Mail]], Token)
	} else {
		AccStruct.AllToken[Token] = AccStruct.AllId[AccStruct.Creation.UserInformation.Mail]
		AccStruct.AllTabToken[AccStruct.AllId[AccStruct.Creation.UserInformation.Mail]] = append(AccStruct.AllTabToken[AccStruct.AllId[AccStruct.Creation.UserInformation.Mail]], Token)
	}
	SaveAllToken(AccStruct)
	SaveAllTabToken(AccStruct)
	http.SetCookie(w, &http.Cookie{Name: "TOKEN", Value: Token})
}

//

// Profil Account

type EveryProfilParameters struct {
	Profil  ProfilAccount
	Friends FriendsParameters
	Name    NameParameters
	Visit   ProfilVisit
	Artists ArtistsProfil
}

type ProfilAccount struct {
	User       UserInformations
	Parameters ShowParameters
}

type UserInformations struct {
	PP           string
	Banner       string
	Name         string
	Numbers      string
	EntireName   string
	Friends      []string
	ArtistsLiked map[string]bool
}

func ProfilUserInformationFill(name, hashtag string, AccStruct *AccountStruct) {
	AccStruct.ProfilParameters.Profil.User.Name = name
	AccStruct.ProfilParameters.Profil.User.Numbers = strings.ToUpper(hashtag)
	AccStruct.ProfilParameters.Profil.User.EntireName = name + "#" + strings.ToUpper(hashtag)
}

type ShowParameters struct {
	PublicAccount      bool
	PublicFriends      bool
	OnlyFriendsFriends bool
}

func ProfilAccountReset(AccStruct *AccountStruct) {
	AccStruct.ProfilParameters.Profil.User.PP = ""
	AccStruct.ProfilParameters.Profil.User.Banner = ""
	AccStruct.ProfilParameters.Profil.User.Name = ""
	AccStruct.ProfilParameters.Profil.User.Numbers = ""
	AccStruct.ProfilParameters.Profil.User.EntireName = ""
	AccStruct.ProfilParameters.Profil.User.Friends = []string{}
	AccStruct.ProfilParameters.Profil.User.ArtistsLiked = make(map[string]bool)
	AccStruct.ProfilParameters.Profil.Parameters.PublicAccount = false
	AccStruct.ProfilParameters.Profil.Parameters.PublicFriends = false
	AccStruct.ProfilParameters.Profil.Parameters.OnlyFriendsFriends = false
}

func ProfilSettings(value string, AccStruct *AccountStruct) {
	if value == "PublicProfil" {
		AccStruct.ProfilParameters.Profil.Parameters.PublicAccount = true
	} else if value == "FriendsProfil" {
		AccStruct.ProfilParameters.Profil.Parameters.PublicAccount = false
	} else if value == "PublicFriends" {
		AccStruct.ProfilParameters.Profil.Parameters.PublicFriends = true
		AccStruct.ProfilParameters.Profil.Parameters.OnlyFriendsFriends = false
	} else if value == "Friends" {
		AccStruct.ProfilParameters.Profil.Parameters.PublicFriends = false
		AccStruct.ProfilParameters.Profil.Parameters.OnlyFriendsFriends = true
	} else {
		AccStruct.ProfilParameters.Profil.Parameters.PublicFriends = false
		AccStruct.ProfilParameters.Profil.Parameters.OnlyFriendsFriends = false
	}
}

type ProfilVisit struct {
	Exist   bool
	Profil  bool
	Friends bool
}

func ProfilVisitReset(AccStruct *AccountStruct) {
	AccStruct.ProfilParameters.Visit.Exist = false
	AccStruct.ProfilParameters.Visit.Profil = false
	AccStruct.ProfilParameters.Visit.Friends = false
}

func ProfilAuthorizeVisit(Token string, AccStruct *AccountStruct) {
	IsFriend := false
	if Token != "" {
		for _, i := range AccStruct.ProfilParameters.Profil.User.Friends {
			if i == Token {
				IsFriend = true
				break
			}
		}
	}
	if AccStruct.ProfilParameters.Profil.Parameters.PublicAccount || IsFriend {
		AccStruct.ProfilParameters.Visit.Profil = true
	}
	if AccStruct.ProfilParameters.Profil.Parameters.PublicFriends {
		AccStruct.ProfilParameters.Visit.Friends = true
	} else if AccStruct.ProfilParameters.Profil.Parameters.OnlyFriendsFriends && IsFriend {
		AccStruct.ProfilParameters.Visit.Friends = true
	}
}

type FriendsParameters struct {
	Friends     []FriendsInformations
	NotExist    bool
	AlreadyHave bool
	ThatsU      bool
}

type FriendsInformations struct {
	Profil UserInformations
	Id     string
}

func FriendsParametersReset(AccStruct *AccountStruct) {
	AccStruct.ProfilParameters.Friends.Friends = []FriendsInformations{}
	AccStruct.ProfilParameters.Friends.NotExist = false
	AccStruct.ProfilParameters.Friends.AlreadyHave = false
	AccStruct.ProfilParameters.Friends.ThatsU = false
}

type NameParameters struct {
	AlreadyExist  bool
	NameIsOk      bool
	Hashtag       bool
	LengthNameInf bool
	LengthNameSup bool
	LengthNumbers bool
}

func NameParametersReset(AccStruct *AccountStruct) {
	AccStruct.ProfilParameters.Name.AlreadyExist = false
	AccStruct.ProfilParameters.Name.NameIsOk = false
	AccStruct.ProfilParameters.Name.Hashtag = false
	AccStruct.ProfilParameters.Name.LengthNameInf = false
	AccStruct.ProfilParameters.Name.LengthNameSup = false
	AccStruct.ProfilParameters.Name.LengthNumbers = false
}

func NameVerificationEntryUser(AccStruct *AccountStruct) bool {
	if len(AccStruct.ProfilParameters.Profil.User.Name) < 3 {
		AccStruct.ProfilParameters.Name.LengthNameInf = true
		return false
	}
	if len(AccStruct.ProfilParameters.Profil.User.Name) > 16 {
		AccStruct.ProfilParameters.Name.LengthNameSup = true
		return false
	}
	if len(AccStruct.ProfilParameters.Profil.User.Numbers) != 4 {
		AccStruct.ProfilParameters.Name.LengthNumbers = true
		return false
	}
	if AccStruct.AllHashtag[AccStruct.ProfilParameters.Profil.User.EntireName] != "" {
		AccStruct.ProfilParameters.Name.AlreadyExist = true
		return false
	}
	return true
}

type ArtistsProfil struct {
	Artists     []ApiArtiste
	ArtistCheck bool
}

func ArtistsProfilReset(AccStruct *AccountStruct) {
	AccStruct.ProfilParameters.Artists.Artists = []ApiArtiste{}
	AccStruct.ProfilParameters.Artists.ArtistCheck = false
}

func ArtistsProfilFill(AccStruct *AccountStruct) {
	for Artist := range AccStruct.ProfilParameters.Profil.User.ArtistsLiked {
		ArtistTemp := ApiArtiste{}
		json.Unmarshal(GetReadAll("https://groupietrackers.herokuapp.com/api/artists/"+Artist), &ArtistTemp)
		AccStruct.ProfilParameters.Artists.Artists = append(AccStruct.ProfilParameters.Artists.Artists, ArtistTemp)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(AccStruct.ProfilParameters.Artists.Artists), func(i, j int) {
		AccStruct.ProfilParameters.Artists.Artists[i], AccStruct.ProfilParameters.Artists.Artists[j] = AccStruct.ProfilParameters.Artists.Artists[j], AccStruct.ProfilParameters.Artists.Artists[i]
	})
}
