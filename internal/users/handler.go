package users

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service userService) *UserHandler {
	return &UserHandler{service: &service}
}
func (h *UserHandler) RegisterRoutes(e *echo.Echo) {

	e.Post("/users", h.service.Register)
	e.Get("/users", h.service.GetUsers)
	e.Get("/users/:id", h.service.GetUser)
	e.Put("/users/:id", h.service.UpdateUser)
	e.Delete("/users/:id", h.service.DeleteUser)
}
func (h *UserHandler) Register(c echo.context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	newUser, err := h.service.Register(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, newUser)
}
func (h *UserHandler) GetUsers(c echo.Context) error {
	users, _ := h.service.GetUsers()
	return c.JSON(http.StatusOK, users)
}
func (h *UserHandler) GetUser(c echo.context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.service.GetUser(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var user User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})

	}
	return c.JSON(http.StatusOK, h.UpdateUser)
}
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.service.DeleteUser(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})

	}
	return c.Json(http.StatusNoContent)
}
