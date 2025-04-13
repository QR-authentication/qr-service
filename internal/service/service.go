package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/skip2/go-qrcode"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	qrproto "github.com/QR-authentication/qr-proto/qr-proto"

	"github.com/QR-authentication/qr-service/internal/config"
	"github.com/QR-authentication/qr-service/internal/model"
)

type Service struct {
	qrproto.UnimplementedQRServiceServer
	repository DBRepo
	signingKey []byte
}

func New(DBRepo DBRepo, signingKey string) *Service {
	return &Service{
		repository: DBRepo,
		signingKey: []byte(signingKey),
	}
}

func (s *Service) CreateQR(ctx context.Context, _ *emptypb.Empty) (*qrproto.CreateQROut, error) {
	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "failed to find UUID")
	}

	claims := model.QRClaims{
		UUID:   uuid,
		Random: generateRandomString(32),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.signingKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to sign token: %v", err)
	}

	if err = s.repository.StoreToken(ctx, tokenString, uuid); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store token in repository: %v", err)
	}

	qrImg, err := qrcode.Encode(tokenString, qrcode.Medium, 256)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate QR code: %v", err)
	}

	return &qrproto.CreateQROut{
		Qr: base64.StdEncoding.EncodeToString(qrImg),
	}, nil
}

func (s *Service) VerifyQR(ctx context.Context, in *qrproto.VerifyQRIn) (*qrproto.VerifyQROut, error) {
	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return &qrproto.VerifyQROut{AccessGranted: false}, status.Error(codes.Unauthenticated, "failed to find UUID")
	}

	latestAction, err := s.repository.GetLatestAction(ctx, uuid)
	if err != nil {
		return &qrproto.VerifyQROut{AccessGranted: false}, status.Errorf(codes.Internal, "failed to get user latest action: %v", err)
	}

	if latestAction == in.Action {
		return &qrproto.VerifyQROut{AccessGranted: false}, nil
	}

	token := s.parseAndValidateToken(in.Token)
	claims, ok := token.Claims.(*model.QRClaims)
	if !ok {
		return &qrproto.VerifyQROut{AccessGranted: false}, status.Error(codes.InvalidArgument, "invalid token claims")
	}

	isScanned, err := s.repository.TokenStatusIsScanned(ctx, in.Token)
	if err != nil {
		return &qrproto.VerifyQROut{AccessGranted: false}, status.Errorf(codes.Internal, "failed to get token status: %v", err)
	}
	if isScanned {
		return &qrproto.VerifyQROut{AccessGranted: false}, nil
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		if err = s.repository.UpdateTokenStatusToExpired(ctx, in.Token); err != nil {
			return &qrproto.VerifyQROut{AccessGranted: false}, status.Errorf(codes.Internal, "failed to update expired token: %v", err)
		}
		return &qrproto.VerifyQROut{AccessGranted: false}, nil
	}

	if err = s.repository.UpdateTokenStatusToScanned(ctx, in.Action, in.Token); err != nil {
		return &qrproto.VerifyQROut{AccessGranted: false}, status.Errorf(codes.Internal, "failed to update token status: %v", err)
	}

	return &qrproto.VerifyQROut{AccessGranted: true}, nil
}

func (s *Service) parseAndValidateToken(tokenString string) *jwt.Token {
	token, _ := jwt.ParseWithClaims(tokenString, &model.QRClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.signingKey, nil
	})

	return token
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("failed to generate random string: %v", err))
	}

	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}
