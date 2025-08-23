package verify

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"email-verification-api/db"
	"email-verification-api/internal/model"

	_ "github.com/lib/pq"
)

var testService *VerificationService

func TestMain(m *testing.M) {
	dbConn, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=secret dbname=email_verification sslmode=disable")
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatal("Не удалось пинговать БД:", err)
	}

	if err := db.Migrate(dbConn); err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	testService = NewVerificationService(nil, dbConn)

	os.Exit(m.Run())
}

func TestGenerateHash(t *testing.T) {
	email := "test@example.com"
	hash1 := testService.GenerateHash(email)
	hash2 := testService.GenerateHash(email)
	if hash1 == hash2 {
		t.Error("Хеши не должны быть одинаковыми при разных вызовах")
	}
	if len(hash1) != 10 {
		t.Error("Длина хеша должна быть 10 символов")
	}
}

func TestCRUDVerification(t *testing.T) {
	email := "test@example.com"
	hash := "abc1234567"

	verification := &model.Verification{
		Email:     email,
		Hash:      hash,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
		Verified:  false,
	}

	err := testService.Repository.Create(verification)
	if err != nil {
		t.Fatal(err)
	}

	found, err := testService.Repository.FindByHash(hash)
	if err != nil {
		t.Fatal(err)
	}
	if found == nil {
		t.Fatal("Запись не найдена")
	}
	if found.Email != email {
		t.Error("Email не совпадает")
	}

	err = testService.Repository.MarkAsVerified(hash)
	if err != nil {
		t.Fatal(err)
	}

	found, _ = testService.Repository.FindByHash(hash)
	if !found.Verified {
		t.Error("Должно быть verified = true")
	}
}
