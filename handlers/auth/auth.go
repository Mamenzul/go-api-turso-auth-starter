package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/mamenzul/go-api/models"
	"github.com/nrednav/cuid2"
)

func Router(db *sql.DB) *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("POST /signup", func(w http.ResponseWriter, r *http.Request) {
		var credentials models.User
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hashedPassword, err := hashPassword(credentials.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := cuid2.Generate()
		_, err = db.Exec("INSERT INTO users (id, username, password) VALUES (?, ?, ?)", id, credentials.Username, hashedPassword)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sessionToken := cuid2.Generate()
		expiresAt := time.Now().Add(120 * time.Second)
		_, err = db.Exec("INSERT INTO sessions (id, user_id, expiresAt) VALUES (?, ?, ?)", sessionToken, id, expiresAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expiresAt,
		})

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User created successfully",
		})
	})
	r.HandleFunc("POST /signin", func(w http.ResponseWriter, r *http.Request) {
		var credentials models.User
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var user models.User
		var user_id models.User_id
		err = db.QueryRow("SELECT id, password FROM users WHERE username = ?", credentials.Username).Scan(&user_id.Id, &user.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sessionToken := cuid2.Generate()
		expiresAt := time.Now().Add(120 * time.Second)
		_, err = db.Exec("INSERT INTO sessions (id, user_id, expiry) VALUES (?, ?, ?)", sessionToken, user_id.Id, expiresAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expiresAt,
		})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User logged in successfully",
		})
	})
	r.HandleFunc("POST /signout", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now(),
		})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User logged out successfully",
		})
	})
	r.HandleFunc("GET /refresh", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var session models.Session

		err = db.QueryRow("SELECT user_id, expiry FROM sessions WHERE id = ?", cookie.Value).Scan(&session.User_id, &session.Expiry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if session.IsExpired() {
			_, err = db.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   "",
				Expires: time.Now(),
			})
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		sessionToken := cuid2.Generate()
		expiresAt := time.Now().Add(120 * time.Second)
		_, err = db.Exec("UPDATE sessions SET expiry = ? WHERE id = ?", expiresAt, cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expiresAt,
		})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Session refreshed successfully",
		})
	})

	r.HandleFunc("GET /me", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var session models.Session
		err = db.QueryRow("SELECT user_id, expiry FROM sessions WHERE id = ?", cookie.Value).Scan(&session.User_id, &session.Expiry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if session.IsExpired() {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		var user models.User
		err = db.QueryRow("SELECT username FROM users WHERE id = ?", session.User_id).Scan(&user.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	})

	return r

}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
