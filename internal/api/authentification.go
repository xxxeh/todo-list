package api

//Файл содержит хендлер обрабатывающий запросы на аутентификацию пользователя
//и функцию auth проверки аутентификации.

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// authHandler обрабатывает запросы на аутентификацию, проверяет правильность введённого пароля,
// и, в случае успешной проверки, генерирует JWT-токен с использованием секретного ключа и возвращает его в ответе.
// Ключ для подписания токена и пароль должны хранится в переменных окружения TODO_SECRET_KEY и TODO_PASSWORD соответственно.
// Если не указана переменная окружения TODO_PASSWORD, то процесс авторизации пропускается.
func authHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	//Пропускаем авторизацию если не определена переменная TODO_PASSWORD
	todo_pass := os.Getenv("TODO_PASSWORD")
	if len(todo_pass) == 0 {
		return
	}

	secret := os.Getenv("TODO_SECRET_KEY")
	if len(secret) == 0 {
		writeJson(w, map[string]string{"error": "Не определена переменная окружения TODO_SECRET_KEY"}, http.StatusInternalServerError)
		return
	}

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	log.Printf("Requset body - %s", buf.String())

	var data map[string]string
	err = json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	pass := data["password"]
	passHash := sha256.Sum256([]byte(pass))
	if todo_pass != pass {
		writeJson(w, map[string]string{"error": "Неверный пароль"}, http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"hash": hex.EncodeToString(passHash[:]),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{"token": signedToken}, http.StatusOK)
}

// auth проверяет перед началом обработки запроса валидность JWT-токена в куках.
// Если пользователь авторизован, то управление передается следующему обработчику.
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Пропускаем авторизацию если не определена переменная TODO_PASSWORD
		todo_pass := os.Getenv("TODO_PASSWORD")
		if len(todo_pass) == 0 {
			next(w, r)
			return
		}

		var signedToken string
		cookie, err := r.Cookie("token")
		if err != nil {
			writeJson(w, map[string]string{"error": "Authentification required"}, http.StatusUnauthorized)
			return
		}

		signedToken = cookie.Value

		jwtToken, err := jwt.Parse(signedToken, func(t *jwt.Token) (interface{}, error) {
			secret := os.Getenv("TODO_SECRET_KEY")
			if len(secret) == 0 {
				return nil, fmt.Errorf("Не определена переменная окружения TODO_SECRET_KEY")
			}
			return []byte(secret), nil
		})

		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}

		if !jwtToken.Valid {
			writeJson(w, map[string]string{"error": "Authentification required"}, http.StatusUnauthorized)
			return
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			writeJson(w, map[string]string{"error": "Authentification required"}, http.StatusUnauthorized)
			return
		}

		passHash := claims["hash"]
		todoPassHash := sha256.Sum256([]byte(todo_pass))
		if passHash != hex.EncodeToString(todoPassHash[:]) {
			writeJson(w, map[string]string{"error": "Authentification required: Pass_hash != todo_pass_hash"}, http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}
