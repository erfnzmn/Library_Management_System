package users

import (
	"errors"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailInUse   = errors.New("email already in use")
	ErrWeakPassword = errors.New("password does not meet policy requirements")
	ErrInvalidLogin = errors.New("invalid email or password")
	ErrInvalidRole  = errors.New("invalid role")


)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Signup: قوانین ثبت‌نام
func (s *Service) Signup(name, email, password, role string) (*User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	role = strings.ToLower(strings.TrimSpace(role))

	// 1) ایمیل تکراری
	exists, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if exists != nil {
		return nil, ErrEmailInUse
	}

	// 2) اعتبارسنجی نقش
	if !IsValidRole(role) {
		return nil, ErrInvalidRole
	}

	// 3) سیاست رمز
	if !passwordStrong(password) {
		return nil, ErrWeakPassword
	}

	// 4) هش رمز
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &User{
		Name:         strings.TrimSpace(name),
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}


// Loginمتد 
func (s *Service) Login(email, password string) (*User, error) {
    email = strings.ToLower(strings.TrimSpace(email))

    u, err := s.repo.FindByEmail(email)
    if err != nil {
        return nil, err
    }
    if u == nil {
        return nil, ErrInvalidLogin
    }
    if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
        return nil, ErrInvalidLogin
    }
    return u, nil
}

func passwordStrong(p string) bool {
	if len(p) < 8 {
		return false
	}
	var hasLower, hasUpper, hasDigit bool
	for _, r := range p {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	return hasLower && hasUpper && hasDigit
}
