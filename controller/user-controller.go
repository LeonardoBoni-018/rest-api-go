package controller

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"go-api/usecase"
)

type authController struct {
	authUseCase *usecase.AuthUseCase
	secret      string
	ttlSeconds  int64
}

func NewAuthController(u *usecase.AuthUseCase) *authController {
	ttl := int64(3600) // 1 hora
	if v := os.Getenv("JWT_TTL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			ttl = int64(n)
		}
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret" // valor padrão para desenvolvimento
	}
	return &authController{
		authUseCase: u,
		secret:      secret,
		ttlSeconds:  ttl,
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ac *authController) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	user, err := ac.authUseCase.Authenticate(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// gerar token JWT
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ac.ttlSeconds) * time.Second)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(ac.secret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}

	ctx.SetCookie("session_token", signed, int(ac.ttlSeconds), "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"token": signed})
}

func (ac *authController) Logout(ctx *gin.Context) {
	ctx.SetCookie("session_token", "", -1, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
