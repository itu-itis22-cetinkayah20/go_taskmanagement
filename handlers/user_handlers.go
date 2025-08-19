package handlers

import (
	"encoding/json"
	"go_taskmanagement/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("gizliAnahtar")

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Geçersiz veri"}`))
		return
	}
	if user.Username == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Kullanıcı adı ve şifre zorunlu"}`))
		return
	}
	// Kullanıcı adı benzersiz mi kontrolü
	for _, u := range models.Users {
		if u.Username == user.Username {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"Kullanıcı adı zaten mevcut"}`))
			return
		}
	}
	// Şifreyi hashle
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Şifre hashlenemedi"}`))
		return
	}
	user.Password = string(hash)
	user.ID = len(models.Users) + 1
	models.Users = append(models.Users, user)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Geçersiz veri"}`))
		return
	}
	var user models.User
	found := false
	for _, u := range models.Users {
		if u.Username == creds.Username {
			user = u
			found = true
			break
		}
	}
	if !found {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Kullanıcı adı veya şifre yanlış"}`))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Kullanıcı adı veya şifre yanlış"}`))
		return
	}
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Token oluşturulamadı"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Çıkış başarılı. Token client tarafından silinmeli."}`))
}
