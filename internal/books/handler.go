package books

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/books")
	g.GET("", h.ListBooks)
	g.GET("/:id", h.GetBook)
	g.GET("/search", h.SearchBooks)
	g.POST("/:id/favorite", h.AddFavorite)
	g.GET("/favorites/:userID", h.GetFavorites)
}

func (h *Handler) ListBooks(c echo.Context) error {
	books, err := h.svc.ListBooks(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, books)
}

func (h *Handler) GetBook(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	book, err := h.svc.GetBook(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, book)
}

func (h *Handler) SearchBooks(c echo.Context) error {
	q := c.QueryParam("q")
	books, err := h.svc.SearchBooks(c.Request().Context(), q)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, books)
}

func (h *Handler) AddFavorite(c echo.Context) error {
	userID := 1 // 👈 فعلاً هاردکد شده (بعداً از JWT یا سشن می‌گیریم)
	bookID, _ := strconv.Atoi(c.Param("id"))

	if err := h.svc.AddFavorite(c.Request().Context(), userID, bookID); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, "added to favorites")
}

func (h *Handler) GetFavorites(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("userID"))
	books, err := h.svc.GetFavorites(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, books)
}
