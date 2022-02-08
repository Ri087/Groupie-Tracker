package GroupieTracker

type CheckCreation struct {
	Name  bool
	Pwd   bool
	Pwdc  bool
	Mail  bool
	Exist bool
}

func CheckCrea(name, pwd, pwdc, mail string, CC *CheckCreation, Acc *Account) bool {
	if len(name) < 3 || len(name) > 18 {
		CC.Name = true
		return true
	}
	if len(pwd) < 6 || len(pwd) > 32 {
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
	return CheckAccount(name, pwd, mail, CC, Acc)
}
