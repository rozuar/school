package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/school-monitoring/backend/internal/models"
)

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "school-monitoring-secret-key-change-in-production"
	}
	jwtSecret = []byte(secret)
}

// Claims estructura de claims del JWT
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Rol    string    `json:"rol"`
	jwt.RegisteredClaims
}

// GenerateToken genera un token JWT para un usuario
func GenerateToken(usuario *models.Usuario) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: usuario.ID,
		Email:  usuario.Email,
		Rol:    usuario.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "school-monitoring",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken valida un token JWT y retorna los claims
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshToken genera un nuevo token a partir de uno existente
func RefreshToken(tokenString string) (string, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Verificar que el token no este muy cerca de expirar
	if time.Until(claims.ExpiresAt.Time) > 30*time.Minute {
		return "", errors.New("token still valid, refresh not needed")
	}

	// Crear nuevo token con claims actualizados
	expirationTime := time.Now().Add(24 * time.Hour)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
