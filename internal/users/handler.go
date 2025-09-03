package users

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Handler struct {
	svc       *Service
	jwtSecret string
	jwtTTL    time.Duration
}

func RegisterUserRoutes(r *gin.Engine, db *gorm.DB, jwtSecret string, jwtTTL time.Duration) {
	_ = db.AutoMigrate(&User{}) 
	repo := NewRepository(db)
	svc := NewService(repo)
	h := &Handler{svc: svc, jwtSecret: jwtSecret, jwtTTL: jwtTTL}

	r.POST("/users/signup", h.Signup)
	r.POST("/users/login", h.Login)

}

func (h *Handler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "detail": err.Error()})
		return
	}

	u, err := h.svc.Signup(req.Name, req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if err == ErrEmailInUse {
			status = http.StatusConflict
		}
		if err == ErrWeakPassword {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	token, expSec, err := h.createJWT(u.ID, u.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   expSec,
		"user":         u,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "detail": err.Error()})
		return
	}

	u, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		status := http.StatusUnauthorized
		if err.Error() != ErrInvalidLogin.Error() {
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	token, expSec, err := h.createJWT(u.ID, u.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   expSec,
		"user":         u,
	})
}


func (h *Handler) createJWT(userID uint, role string) (string, int64, error) {
	now := time.Now()
	exp := now.Add(h.jwtTTL)

	claims := jwt.MapClaims{
		"sub":  strconv.Itoa(int(userID)),
		"role": role,
		"iat":  now.Unix(),
		"exp":  exp.Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return "", 0, err
	}
	return signed, exp.Unix() - now.Unix(), nil
}
