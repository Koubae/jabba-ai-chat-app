package auth

type AccessToken struct {
	ApplicationId string
	UserId        int64
	Username      string
	Issuer        string
	Role          string
}
