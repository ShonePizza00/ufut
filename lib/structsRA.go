package structsUFUT

type TokenResponse struct {
	Token string `json:"token"`
}

type UserGeneral struct {
	Login      string `json:"login"`
	PasswdHash string `json:"passwdHash"`
	Token      string `json:"token"`
}

type UserData struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNum"`
}

type UpdatePasswdRequest struct {
	Login  string `json:"login"`
	Passwd string `json:"password"`
}

type UserUpdatePasswd struct {
	Login     string `json:"login"`
	Passwd    string `json:"passwd"`
	NewPasswd string `json:"newPasswd"`
}

type UserUpdatePasswdHash struct {
	Login         string `json:"login"`
	PasswdHash    string `json:"passwdHash"`
	NewPasswdHash string `json:"newPasswdHash"`
	NewToken      string `json:"newToken"`
}
