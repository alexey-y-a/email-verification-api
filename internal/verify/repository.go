package verify

import (
	"database/sql"
	"email-verification-api/internal/model"
	"time"
)

type VerificationRepository struct {
	DB *sql.DB
}

func NewVerificationRepository(db *sql.DB) *VerificationRepository {
	return &VerificationRepository{DB: db}
}

func (r *VerificationRepository) Create(verification *model.Verification) error {
	query := `
	INSERT INTO verifications (email, hash, created_at, expires_at, verified)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Exec(query,
		verification.Email,
		verification.Hash,
		verification.CreatedAt,
		verification.ExpiresAt,
		verification.Verified,
	)
	return err
}

func (r *VerificationRepository) FindByHash(hash string) (*model.Verification, error) {
	query := `SELECT id, email, hash, created_at, expires_at, verified FROM verifications WHERE hash = $1`
	row := r.DB.QueryRow(query, hash)

	var v model.Verification
	err := row.Scan(&v.ID, &v.Email, &v.Hash, &v.CreatedAt, &v.ExpiresAt, &v.Verified)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if v.ExpiresAt.Before(time.Now()) {
		return nil, nil
	}

	return &v, nil
}

func (r *VerificationRepository) MarkAsVerified(hash string) error {
	query := `UPDATE verifications SET verified = true WHERE hash = $1`
	_, err := r.DB.Exec(query, hash)
	return err
}

func (r *VerificationRepository) DeleteByHash(hash string) error {
	query := `DELETE FROM verifications WHERE hash = $1`
	_, err := r.DB.Exec(query, hash)
	return err
}
