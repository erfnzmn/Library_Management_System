package loans

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/erfnzmn/Library_Management_System/pkg/middleware"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

type Handler struct {
	service *Service
	rabbitChannel *amqp.Channel
	jwtSecret []byte
}

func NewHandler(service *Service, rabbitChannel *amqp.Channel, jwtSecret string) *Handler {
	return &Handler{
		service:       service,
		rabbitChannel: rabbitChannel,
		jwtSecret:     []byte(jwtSecret),
	}
}
// RegisterRoutes 
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	log.Printf("JWT secret in handler: %s", string(h.jwtSecret))

	g := e.Group("/api/loans")

	g.Use(echojwt.JWT([]byte(h.jwtSecret)))

	g.POST("/reserve", h.ReserveBook)
	g.POST("/:id/confirm", h.ConfirmBorrow)
	g.POST("/:id/return", h.ReturnBook)
	g.POST("/:id/cancel", h.CancelReservation)
	g.GET("/user/:userID", h.GetUserLoans)
}

// ReserveBook 
func (h *Handler) ReserveBook(c echo.Context) error {
	var req struct {
		BookID uint `json:"book_id" binding:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid or missing token"})
	}
	body, err := json.Marshal(map[string]uint{
		"user_id": userID,
		"book_id": req.BookID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "encode error"})
	}

	err = h.rabbitChannel.Publish(
		"", "reserve_requests", false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to queue reservation"})
	}

	return c.JSON(http.StatusAccepted, echo.Map{
		"message": "reservation request queued",
	})
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
