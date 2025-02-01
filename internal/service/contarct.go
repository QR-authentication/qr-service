package service

type DBRepo interface {
	StoreToken(token, uuid, ip string) error
	GetTokenStatus(token string) (string, error)
	UpdateTokenStatusToExpired(token string) error
	UpdateTokenStatusToScanned(token string) error
}
