package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"velvet-archive-api/internal/db"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Structural shapes for incoming auth payloads
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token     string    `json:"token"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at"`
}

// POST /api/v1/auth/register - Register the initial admin user securely
func (bh *BaseHandler) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON payload structure")
		return
	}

	if req.Email == "" || len(req.Password) < 6 {
		respondWithError(w, http.StatusBadRequest, "Valid email and a password of at least 6 characters are required")
		return
	}

	// Hash the password securely using bcrypt (simulating standard production practices)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to process encryption security credentials")
		return
	}

	// Save user details to Postgres using generated SQLC code
	_, err = bh.DB.CreateAdminUser(r.Context(), db.CreateAdminUserParams{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		// Handle database unique violations cleanly
		respondWithError(w, http.StatusConflict, "An administrative profile with this email already exists")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Administrative profile registered successfully."})
}

// POST /api/v1/auth/login - Validate user credentials and issue signed JWT
func (bh *BaseHandler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON input payload")
		return
	}

	// Fetch user records matching requested email address
	admin, err := bh.DB.GetAdminByEmail(r.Context(), req.Email)
	if err != nil {
		// Intentionally generic error to avoid exposing user presence maps to attackers
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password credentials")
		return
	}

	// Compare stored hash against requested text password string
	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password credentials")
		return
	}

	// Create JWT token payload claims
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.MapClaims{
		"email": admin.Email,
		"exp":   expirationTime.Unix(),
		"iat":   time.Now().Unix(),
	}

	// Sign token using HMAC SHA256 and the JWT_SECRET from environmental config
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate authorization signature layer")
		return
	}

	respondWithJSON(w, http.StatusOK, AuthResponse{
		Token:     tokenString,
		Email:     admin.Email,
		ExpiresAt: expirationTime,
	})
}
