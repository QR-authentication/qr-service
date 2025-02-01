package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	qrproto "github.com/QR-authentication/qr-proto/qr-proto"
	"github.com/dgrijalva/jwt-go"
	"github.com/skip2/go-qrcode"
)

type Service struct {
	qrproto.UnimplementedQRServiceServer
	repository DBRepo
	signingKey string
}

func New(DBRepo DBRepo, signingKey string) *Service {
	return &Service{
		repository: DBRepo,
		signingKey: signingKey,
	}
}

func (s *Service) CreateQR(ctx context.Context, in *qrproto.CreateQRIn) (*qrproto.CreateQROut, error) {
	claims := jwt.MapClaims{
		//"exp":     time.Now().Add(time.Second * 30).Unix(),
		"exp":    time.Now().Add(time.Hour * 30).Unix(),
		"random": generateRandomString(32),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.signingKey))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	if err = s.repository.StoreToken(tokenString, in.Uuid, in.Ip); err != nil {
		return nil, fmt.Errorf("failed to store token: %w", err)
	}

	qrImg, err := qrcode.Encode(tokenString, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to create qr image: %w", err)
	}

	qrBase64 := base64.StdEncoding.EncodeToString(qrImg)

	return &qrproto.CreateQROut{QR: qrBase64}, nil
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomBytes := make([]byte, length)
	_, _ = rand.Read(randomBytes)

	for i := range randomBytes {
		randomBytes[i] = charset[randomBytes[i]%byte(len(charset))]
	}

	return string(randomBytes)
}

func (s *Service) VerifyQR(ctx context.Context, in *qrproto.VerifyQRIn) (*qrproto.VerifyQROut, error) {
	token, err := s.parseAndValidateToken(in.Token)
	if err != nil {
		return &qrproto.VerifyQROut{AccessGranted: false}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &qrproto.VerifyQROut{AccessGranted: false}, fmt.Errorf("failed to get token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return &qrproto.VerifyQROut{AccessGranted: false}, fmt.Errorf("failed to missing expiration claim")
	}

	if time.Now().Unix() > int64(exp) {
		log.Println(time.Now().Unix(), int64(exp))
		err = s.repository.UpdateTokenStatusToExpired(in.Token)
		if err != nil {
			return &qrproto.VerifyQROut{AccessGranted: false}, err
		}
	}

	tokenStatus, err := s.repository.GetTokenStatus(in.Token)
	if err != nil {
		return &qrproto.VerifyQROut{AccessGranted: false}, err
	}

	if tokenStatus != "pending" {
		return &qrproto.VerifyQROut{AccessGranted: false}, fmt.Errorf("token is not valid for access (status: %s)", tokenStatus)
	}

	if err = s.repository.UpdateTokenStatusToScanned(in.Token); err != nil {
		return &qrproto.VerifyQROut{AccessGranted: false}, err
	}

	return &qrproto.VerifyQROut{AccessGranted: true}, nil
}

func (s *Service) parseAndValidateToken(tokenString string) (*jwt.Token, error) {
	fmt.Println("Parsing token:", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		fmt.Println("Token Header:", token.Header)

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.signingKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return token, nil
}
