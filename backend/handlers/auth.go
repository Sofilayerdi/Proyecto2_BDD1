package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"proy2-bck/db"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Rol      string `json:"rol"`
	Username string `json:"username"`
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Rol      string `json:"rol"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(getSecret())

func getSecret() string {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		return "bloom_secret_key"
	}
	return s
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if input.Username == "" || input.Password == "" {
		http.Error(w, "username y password son requeridos", http.StatusBadRequest)
		return
	}

	var userID int
	var hashedPassword, rol string
	err := db.DB.QueryRow(
		`SELECT id_usuario, password, rol FROM usuario WHERE username = $1`,
		input.Username,
	).Scan(&userID, &hashedPassword, &rol)

	if err != nil {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password)); err != nil {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	// generar JWT
	claims := &Claims{
		UserID:   userID,
		Username: input.Username,
		Rol:      rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Error generando token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Token:    tokenString,
		Rol:      rol,
		Username: input.Username,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// con JWT el logout es del lado del cliente (eliminar el token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Sesión cerrada"})
}
