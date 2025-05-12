package handler

import (
	"library/internal/facade"
	"library/responder"
	"net/http"
	"strconv"
)

type Rentaler interface {
	RentBook(w http.ResponseWriter, r *http.Request)
	ReturnBook(w http.ResponseWriter, r *http.Request)
}

type RentalHandler struct {
	rentUC    facade.Facader
	responder responder.Responder
}

func NewRentHandler(rentUC facade.Facader, responder responder.Responder) Rentaler {
	return &RentalHandler{
		rentUC:    rentUC,
		responder: responder,
	}
}

// @Summary			rental book
// @Description		rental book
// @Tags			rental
// @Accept			json
// @Produce			json
// @Param			bookId   path	string	true  "bookId"
// @Param			userId   path	string	true  "userID"
// @Success			200		{object}	Response
// @Router			/rental/{bookId}/{userId} [post]
func (h *RentalHandler) RentBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(r.PathValue("bookId"))
	if err != nil {
		h.responder.ErrorBadRequest(w, err)
		return
	}

	userID, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		h.responder.ErrorBadRequest(w, err)
		return
	}

	if err := h.rentUC.RentBook(r.Context(), bookID, userID); err != nil {
		h.responder.ErrorInternal(w, err)
		return
	}

	h.responder.OutputJSON(w, Response{
		Success: true,
		Data: Data{
			Message: "book has been leased",
		},
	})
}

// @Summary			return book
// @Description		return book
// @Tags			rental
// @Accept			json
// @Produce			json
// @Param			bookId   path	string	true  "bookId"
// @Success			200		{object}	Response
// @Router			/rental/{bookId} [delete]
func (h *RentalHandler) ReturnBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(r.PathValue("bookId"))
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	if err := h.rentUC.ReturnBook(r.Context(), bookID); err != nil {
		h.responder.ErrorInternal(w, err)
		return
	}

	h.responder.OutputJSON(w, Response{
		Success: true,
		Data: Data{
			Message: "book rental completed",
		},
	})
}
