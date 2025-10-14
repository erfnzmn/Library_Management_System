package loans

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes رذ
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/api/loans")
	g.POST("/reserve", h.ReserveBook)
	g.POST("/:id/confirm", h.ConfirmBorrow)
	g.POST("/:id/return", h.ReturnBook)
	g.POST("/:id/cancel", h.CancelReservation)
	g.GET("/user/:userID", h.GetUserLoans)
}

// ReserveBook 
func (h *Handler) ReserveBook(c echo.Context) error {
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
		BookID uint `json:"book_id" binding:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	if err := h.service.ReserveBook(c.Request().Context(), req.UserID, req.BookID); err != nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"message": "book reserved successfully"})
}

// ConfirmBorrow 
func (h *Handler) ConfirmBorrow(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid loan id"})
	}
	if err := h.service.ConfirmBorrow(c.Request().Context(), uint(id)); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "loan confirmed"})
}

// ReturnBook 
func (h *Handler) ReturnBook(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid loan id"})
	}
	if err := h.service.ReturnBook(c.Request().Context(), uint(id)); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "book returned successfully"})
}

// CancelReservation 
func (h *Handler) CancelReservation(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid loan id"})
	}
	if err := h.service.CancelReservation(c.Request().Context(), uint(id)); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "reservation cancelled"})
}

// GetUserLoans 
func (h *Handler) GetUserLoans(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid user id"})
	}
	loans, err := h.service.repo.GetLoansByUser(c.Request().Context(), uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, loans)
}
