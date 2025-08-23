package verify

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"time"

	"email-verification-api/config"

	"github.com/jordan-wright/email"
)

type VerificationService struct {
	Config     *config.Config
	Repository *VerificationRepository
}

func NewVerificationService(cfg *config.Config, db *sql.DB) *VerificationService {
	repo := NewVerificationRepository(db)
	return &VerificationService{Config: cfg, Repository: repo}
}

func (s *VerificationService) GenerateHash(email string) string {
	h := sha256.New()
	h.Write([]byte(email + time.Now().String() + "secret_salt"))
	return fmt.Sprintf("%x", h.Sum(nil))[:10]
}

func (s *VerificationService) SendVerificationEmail(to, hash string) error {
	e := email.NewEmail()
	e.From = s.Config.SMTPEmail
	e.To = []string{to}
	e.Subject = "Подтвердите ваш email"
	verifyURL := fmt.Sprintf("http://localhost:%s/verify/%s", s.Config.AppPort, hash)
	e.HTML = []byte(fmt.Sprintf(`
		<h1>Привет!</h1>
		<p>Нажмите на ссылку ниже, чтобы подтвердить email:</p>
		<a href="%s">%s</a>
	`, verifyURL, verifyURL))

	smtpAddr := s.Config.SMTPHost + ":" + s.Config.SMTPPort
	return e.Send(smtpAddr, email.PlainAuth("", s.Config.SMTPEmail, s.Config.SMTPPassword, s.Config.SMTPHost))
}
