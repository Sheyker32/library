package handler

import (
	"encoding/json"
	"library/internal/domain"
	"library/internal/usecase"
	"library/responder"
	"net/http"
	"strconv"
)

type Booker interface {
	AddBook(w http.ResponseWriter, r *http.Request)
	GetBook(w http.ResponseWriter, r *http.Request)
	DeleteBook(w http.ResponseWriter, r *http.Request)
}

type BookHandler struct {
	bookUC    usecase.Booker
	responder responder.Responder
}

func NewBookHandler(bookUC usecase.Booker, responder responder.Responder) Booker {
	return &BookHandler{
		bookUC:    bookUC,
		responder: responder,
	}
}

// @Summary			add book
// @Description		add book
// @Tags			book
// @Accept			json
// @Produce			json
// @Param			book   body	domain.Book	true  "book"
// @Success			200		{object}	Response
// @Router			/book [post]
func (h *BookHandler) AddBook(w http.ResponseWriter, r *http.Request) {
	var book domain.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		h.responder.ErrorBadRequest(w, err)
		return
	}

	if err := h.bookUC.AddBook(r.Context(), &book); err != nil {
		h.responder.ErrorInternal(w, err)
		return
	}

	h.responder.OutputJSON(w, Response{
		Success: true,
		Data:    book,
	})
}

// @Summary			get book
// @Description		get book
// @Tags			book
// @Accept			json
// @Produce			json
// @Param			bookId   path	string	true  "id book"
// @Success			200		{object}	Response
// @Router			/book/{bookId} [get]
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(r.PathValue("bookId"))
	if err != nil {
		h.responder.ErrorBadRequest(w, err)
		return
	}

	book, err := h.bookUC.GetBook(r.Context(), bookID)
	if err != nil {
		h.responder.ErrorInternal(w, err)
		return
	}

	h.responder.OutputJSON(w, Response{
		Success: true,
		Data:    book,
	})
}

// @Summary			delete book
// @Description		delete book
// @Tags			book
// @Accept			json
// @Produce			json
// @Param			bookId   path	string	true  "id book"
// @Success			200		{object}	Response
// @Router			/book/{bookId} [delete]
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(r.PathValue("bookId"))
	if err != nil {
		h.responder.ErrorBadRequest(w, err)
		return
	}

	err = h.bookUC.DeleteBook(r.Context(), bookID)
	if err != nil {
		h.responder.ErrorInternal(w, err)
		return
	}

	h.responder.OutputJSON(w, Response{
		Success: true,
		Data:    "book id delete",
	})
}
