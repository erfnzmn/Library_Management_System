package middleware

import(
	"errors"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CurrentUserID(c echo.Context) (uint, error){
	u := c.Get("user")
	if u == nil {
		return 0, errors.New("no jwt user in context (missing JWT middleware)")
	}

	token, ok := u.(*jwt.Token)
	if !ok {
		return 0, errors.New("invalid token type in context")

	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid JWT claims type")
	}

	rawSub , ok := claims["sub"]
	if !ok {
		return 0, errors.New("sub claim not found")
	}

	switch v := rawSub.(type){
	case string:
		n, err := strconv.ParseUint(v, 10 ,64)
		if err != nil {
			return 0 ,errors.New("invalid sub claim (not a uint)")
		}
		return uint(n), nil
	case float64:
		return uint(v), nil
	default:
		return 0, errors.New("unsupported sub claim type")
	}
}