package handler

import (
	"encoding/json"
	"library/internal/domain"
	"library/internal/usecase"
	"library/responder"
	"net/http"
	"strconv"
)

type Authorer interface {
	CreateAuthor(w http.ResponseWriter, r *http.Request)
	GetAuthor(w http.ResponseWriter, r *http.Request)
	GetTopAuthors(w http.ResponseWriter, r *http.Request)
	GetAllAuthors(w http.ResponseWriter, r *http.Request)
	DeleteAuthor(w http.ResponseWriter, r *http.Request)
	GetByBooksAuthor(w http.ResponseWriter, r *http.Request)
}

type AuthorHandler struct {
	authorUC  usecase.Authorer
	responder responder.Responder
}

func NewAuthorHandler(authorUC usecase.Authorer, responder responder.Responder) Authorer {
	return &AuthorHandler{
		authorUC:  authorUC,
		responder: responder,
	}
}

// @Summary			create author
// @Description		create author
// @Tags			author
// @Accept			json
// @Produce			json
// @Param			author   body	domain.Author	true  "author"
// @Success			200		{object}	Response
// @Router			/author [post]
func (h *AuthorHandler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var author domain.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.authorUC.CreateAuthor(r.Context(), &author); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.responder.OutputJSON(w, Response{
		Success: true,
		Data:    author,
	})
}

// @Summary			get author
// @Description		get author
// @Tags			author
// @Accept			json
// @Produce			json
// @Param			authorId   path	string	true  "id author"
// @Success			200		{object}	Response
// @Router			/author/{authorId} [get]
func (h *AuthorHandler) GetAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("authorId"))
	if err != nil {
		h.responder.ErrorBadRequest(w, err)
		return
	}

	author, err := h.authorUC.GetAuthor(r.Context(), id)
	if err != nil {
		h.responder.ErrorInternal(w, err)
		return
	}

	h.responder.OutputJSON(w, Response{
		Success: true,
		Data:    author,
	})
}

// @Summary			get top authors
// @Description		get top
// @Tags			author
// @Accept			json
// @Produce			json
// @Param			limit   query	string	true  "limit"
// @Success			200		{object}	Response
// @Router			/author/top [get]
func (h *AuthorHandler) GetTopAuthors(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if l, err := strconv.Atoi(l); err == nil && l > 0 {
			limit = l
		}
	}

	authors, err := h.authorUC.GetTopAuthors(r.Context(), limit)
	if err != nil {
		h.responder.ErrorInternal(w, err)
		return
	}

	h.responder.OutputJSON(w, Response{
		Success: true,
		Data:    authors,
	})
}

// @Summary			get all authors
// @Description		get all
// @Tags			author
// @Accept			json
// @Produce			json
// @Success			200		{object}	Response
// @Router			/author/all [get]
func (h *AuthorHandler) GetAllAuthors(w http.ResponseWriter, r *http.Request) {
	authors, err := h.authorUC.ListAuthors(r.Context())
	if err != nil {
		h.responder.ErrorBadRequest(w, err)
		return
	}

	h.responder.OutputJSON(w, Response{
		Success: true,
		Data:    authors,
	})
}

// @Summary			delete author
// @Description		delete author
// @Tags			author
// @Accept			json
// @Produce			json
// @Param			authorId   path	string	true  "id author"
// @Success			200		{object}	Response
// @Router			/author/{authorId} [delete]
func (h *AuthorHandler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("authorId"))
	if err != nil {
		h.responder.ErrorBadRequest(w, err)
		return
	}

	err = h.authorUC.DeleteAuthor(r.Context(), id)
	if err != nil {
		h.responder.ErrorInternal(w, err)
		return
	}

	h.responder.OutputJSON(w, Response{
		Success: true,
		Data:    "author is delete",
	})
}

// @Summary			get by books author
// @Description		get books
// @Tags			book
// @Accept			json
// @Produce			json
// @Param			authorId   path	string	true  "authorId"
// @Success			200		{object}	Response
// @Router			/author/books/{authorId} [get]
func (h *AuthorHandler) GetByBooksAuthor(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(r.PathValue("authorId"))
	if err != nil {
		h.responder.ErrorBadRequest(w, err)
		return
	}

	books, err := h.authorUC.GetByBooksAuthor(r.Context(), bookID)
	if err != nil {
		h.responder.ErrorInternal(w, err)
		return
	}
	h.responder.OutputJSON(w, Response{
		Success: true,
		Data:    books,
	})
}
