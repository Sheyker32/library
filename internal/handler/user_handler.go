package handler

import (
	"library/internal/domain"
	"library/internal/usecase"
	"library/responder"
	"net/http"
	"strconv"
	"time"
)

type Userer interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	userUC    usecase.Userer
	responder responder.Responder
}

func NewUserHandler(userUC usecase.Userer, responder responder.Responder) Userer {
	return &UserHandler{
		userUC:    userUC,
		responder: responder,
	}
}

// @Summary			add user
// @Description		add user
// @Tags			user
// @Accept			x-www-form-urlencoded
// @Produce			json
// @Param name   	formData	string	true  "name"
// @Param email   	formData	string	true  "email"
// @Success			200		{object}	Response
// @Router			/user [post]
func (u *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	user := domain.User{
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}
	if err := u.userUC.CreateUser(r.Context(), &user); err != nil {
		u.responder.ErrorInternal(w, err)
		return
	}

	u.responder.OutputJSON(w, Response{
		Success: true,
		Data:    user,
	})
}

// @Summary			get user
// @Description		get user
// @Tags			user
// @Accept			json
// @Produce			json
// @Param			userId   path	string	true  "get user"
// @Success			200		{object}	Response
// @Router			/user/{userId} [get]
func (u *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		u.responder.ErrorBadRequest(w, err)
		return
	}

	user, err := u.userUC.GetByIDUser(r.Context(), userID)
	if err != nil {
		u.responder.ErrorInternal(w, err)
		return
	}

	u.responder.OutputJSON(w, Response{
		Success: true,
		Data:    user,
	})
}

// @Summary			delete user
// @Description		delete user
// @Tags			user
// @Accept			json
// @Produce			json
// @Param			userId   path	string	true  "id user"
// @Success			200		{object}	Response
// @Router			/user/{userId} [delete]
func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		u.responder.ErrorBadRequest(w, err)
		return
	}

	err = u.userUC.DeleteUser(r.Context(), userID)
	if err != nil {
		u.responder.ErrorInternal(w, err)
		return
	}

	u.responder.OutputJSON(w, Response{
		Success: true,
		Data:    "user id delete",
	})
}

// @Summary			get all user
// @Description		get all user
// @Tags			user
// @Accept			json
// @Produce			json
// @Success			200		{object}	Response
// @Router			/user/all [get]
func (u *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := u.userUC.GetAllUsers(r.Context())
	if err != nil {
		u.responder.ErrorBadRequest(w, err)
		return
	}

	u.responder.OutputJSON(w, Response{
		Success: true,
		Data:    users,
	})
}
