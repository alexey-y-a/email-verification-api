package main

import (
	"database/sql"
	"log"
	"net/http"

	"email-verification-api/config"
	"email-verification-api/db"
	"email-verification-api/internal/verify"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	connStr := "host=" + cfg.DBHost +
		" port=" + cfg.DBPort +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" sslmode=" + cfg.DBSSLMode

	dbConn, err := sql.Open("postgres", connStr)
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

	log.Println("Подключение к PostgreSQL установлено")

	service := verify.NewVerificationService(cfg, dbConn)

	http.HandleFunc("/send", service.SendHandler)
	http.HandleFunc("/verify/", service.VerifyHandler)

	log.Printf("Сервис запущен на порту %s", cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, nil))
}
