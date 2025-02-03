package postgres

import (
	"fmt"
	"log"

	"github.com/QR-authentication/qr-service/internal/config"
	guuid "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
)

type Repository struct {
	connection *sqlx.DB
}

func New(cfg *config.Config) *Repository {
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	conn, err := sqlx.Connect("postgres", conStr)
	if err != nil {
		log.Fatal("error connect: ", err)
	}

	return &Repository{
		connection: conn,
	}
}

func (r *Repository) Close() {
	_ = r.connection.Close()
}

func (r *Repository) StoreToken(token, uuid, ip string) error {
	query := `INSERT INTO tokens (id, token, uuid, ip_address) VALUES ($1, $2, $3, $4)`
	_, err := r.connection.Exec(query, guuid.New(), token, uuid, ip)
	if err != nil {
		return fmt.Errorf("failed to insert token: %w", err)
	}

	return nil
}

func (r *Repository) GetTokenStatus(token string) (string, error) {
	var tokenStatus string
	log.Println(token)

	query := `SELECT status FROM tokens WHERE token = $1`

	err := r.connection.Get(&tokenStatus, query, token)
	if err != nil {
		return "", fmt.Errorf("failed to get token status: %w", err)
	}

	return tokenStatus, nil
}

func (r *Repository) UpdateTokenStatusToExpired(token string) error {
	query := `UPDATE tokens SET status = 'expired' WHERE token = $1`

	_, err := r.connection.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to update token status to expired: %w", err)
	}

	return nil
}

func (r *Repository) UpdateTokenStatusToScanned(token string) error {
	query := `UPDATE tokens SET status = 'scanned', scanned_at = NOW() WHERE token = $1`

	_, err := r.connection.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to update token status to scanned: %w", err)
	}

	return nil
}
