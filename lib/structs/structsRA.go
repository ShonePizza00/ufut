package structsUFUT

type TokenResponse struct {
	JWT string `json:"jwt"`
	RT  string `json:"rt"`
}

type UserGeneral struct {
	Login      string `json:"login"`
	PasswdHash string `json:"passwdHash"`
	UserID     string `json:"userID"`
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
	Passwd    string `json:"password"`
	NewPasswd string `json:"newPassword"`
}

type UserUpdatePasswdHash struct {
	Login         string `json:"login"`
	PasswdHash    string `json:"passwdHash"`
	NewPasswdHash string `json:"newPasswdHash"`
}

type JWTUpdate struct {
	OldRT      string `json:"oldRT"`
	NewRT      string `json:"newRT"`
	TimeOffset int64  `json:"timeOffset"`
}
