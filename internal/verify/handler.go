package verify

import (
	"email-verification-api/internal/model"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (s *VerificationService) SendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email обязателен", http.StatusBadRequest)
		return
	}

	hash := s.GenerateHash(req.Email)

	verification := &model.Verification{
		Email:     req.Email,
		Hash:      hash,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Verified:  false,
	}

	if err := s.Repository.Create(verification); err != nil {
		log.Printf("Ошибка сохранения в БД: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	if err := s.SendVerificationEmail(req.Email, hash); err != nil {
		log.Printf("Ошибка отправки email: %v", err)
		http.Error(w, "Не удалось отправить письмо", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Письмо отправлено",
		"hash":    hash,
	})
}

func (s *VerificationService) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	hash := r.URL.Path[len("/verify/"):]

	verification, err := s.Repository.FindByHash(hash)
	if err != nil || verification == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"valid": false})
		return
	}

	if err := s.Repository.DeleteByHash(hash); err != nil {
		log.Printf("Ошибка удаления записи: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"valid": true})
}
