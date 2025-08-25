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
		t.Fatal("Не удалось создать запись:", err)
	}

	found, err := testService.Repository.FindByHash(hash)
	if err != nil {
		t.Fatal("Ошибка при поиске:", err)
	}
	if found == nil {
		t.Fatal("Запись не найдена после создания")
	}
	if found.Email != email {
		t.Errorf("Email не совпадает: ожидается %s, получено %s", email, found.Email)
	}

	err = testService.Repository.DeleteByHash(hash)
	if err != nil {
		t.Fatal("Не удалось удалить запись:", err)
	}

	foundAfter, err := testService.Repository.FindByHash(hash)
	if err != nil {
		t.Fatal("Ошибка при поиске после удаления:", err)
	}
	if foundAfter != nil {
		t.Fatal("Запись найдена после удаления — ожидается nil")
	}
}
