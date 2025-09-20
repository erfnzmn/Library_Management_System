package users

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	svc       *Service
	jwtSecret string
	jwtTTL    time.Duration
}

func RegisterUserRoutes(e *echo.Echo, db *gorm.DB, jwtSecret string, jwtTTL time.Duration) {
	_ = db.AutoMigrate(&User{})

	repo := NewRepository(db)
	svc := NewService(repo)
	h := &Handler{svc: svc, jwtSecret: jwtSecret, jwtTTL: jwtTTL}

	e.POST("/users/signup", h.Signup)
	e.POST("/users/login", h.Login)
}

//  Signup 
func (h *Handler) Signup(c echo.Context) error {
	var req SignupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request", "detail": err.Error()})
	}
	if req.Name == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request", "detail": "name/email/password/role required"})
	}
	if !IsValidRole(req.Role) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid role", "detail": "role must be 'member' or 'student'"})
	}

	u, err := h.svc.Signup(req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case ErrEmailInUse:
			status = http.StatusConflict
		case ErrWeakPassword, ErrInvalidRole:
			status = http.StatusBadRequest
		}
		return c.JSON(status, echo.Map{"error": err.Error()})
	}

	token, expSec, err := h.createJWT(u.ID, u.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "token generation failed"})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   expSec,
		"user":         u,
	})
}


// Login
func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request", "detail": err.Error()})
	}
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request", "detail": "email/password required"})
	}

	u, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		status := http.StatusUnauthorized
		if err.Error() != ErrInvalidLogin.Error() {
			status = http.StatusInternalServerError
		}
		return c.JSON(status, echo.Map{"error": err.Error()})
	}

	token, expSec, err := h.createJWT(u.ID, u.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "token generation failed"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   expSec,
		"user":         u,
	})
}

//JWT helper
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
