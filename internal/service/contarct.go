package service

type DBRepo interface {
	StoreToken(token, uuid string) error
	TokenStatusIsScanned(token string) (bool, error)
	UpdateTokenStatusToExpired(token string) error
	UpdateTokenStatusToScanned(token string) error
}
