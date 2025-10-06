package books

import (
	"context"
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
	e.POST("/books", h.CreateBook)
	e.GET("/books", h.ListBooks)
	e.GET("/books/:id", h.GetBookByID)
	e.PUT("/books/:id", h.UpdateBook)
	e.DELETE("/books/:id", h.DeleteBook)
	e.GET("/books/search", h.SearchBooks)

	e.POST("/books/:id/favorite/:user_id", h.AddToFavorites)
	e.GET("/books/favorites/:user_id", h.GetFavoritesByUser)
}

func (h *Handler) CreateBook(c echo.Context) error {
	var book Book
	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := h.service.CreateBook(context.Background(), &book); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, book)
}

func (h *Handler) UpdateBook(c echo.Context) error {
	var book Book
	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	book.ID = uint(id)
	if err := h.service.UpdateBook(context.Background(), &book); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, book)
}

func (h *Handler) DeleteBook(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.DeleteBook(context.Background(), uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) ListBooks(c echo.Context) error {
	books, err := h.service.ListBooks(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, books)
}

func (h *Handler) GetBookByID(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	book, err := h.service.GetBookByID(context.Background(), uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, book)
}

func (h *Handler) SearchBooks(c echo.Context) error {
	q := c.QueryParam("q")
	books, err := h.service.SearchBooks(context.Background(), q)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, books)
}

func (h *Handler) AddToFavorites(c echo.Context) error {
	bookID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	err := h.service.AddToFavorites(context.Background(), uint(userID), uint(bookID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "book added to favorites")
}

func (h *Handler) GetFavoritesByUser(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("user_id"))
	books, err := h.service.GetFavoritesByUser(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, books)
}
