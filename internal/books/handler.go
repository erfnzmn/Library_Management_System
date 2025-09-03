package books

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

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/books")
	g.POST("", h.CreateBook)
	g.GET("", h.GetAllBooks)
	g.GET("/:id", h.GetBookByID)
	g.PUT("/:id", h.UpdateBook)
	g.DELETE("/:id", h.DeleteBook)
}

func (h *Handler) CreateBook(c echo.Context) error {
	var book Book
	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := h.service.CreateBook(&book); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, book)
}

func (h *Handler) GetAllBooks(c echo.Context) error {
	books, err := h.service.GetAllBooks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, books)
}

func (h *Handler) GetBookByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	book, err := h.service.GetBookByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "book not found"})
	}
	return c.JSON(http.StatusOK, book)
}

func (h *Handler) UpdateBook(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var book Book
	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	book.ID = uint(id)
	if err := h.service.UpdateBook(&book); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, book)
}

func (h *Handler) DeleteBook(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.DeleteBook(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "book deleted"})
}
