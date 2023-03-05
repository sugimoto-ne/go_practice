package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sugimoto-ne/go_practice.git/domain"
	"github.com/sugimoto-ne/go_practice.git/interface/controllers"
	usecases "github.com/sugimoto-ne/go_practice.git/usecases/user"
)

type ShowUser struct {
	Service usecases.GetUser
}

func (su *ShowUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		controllers.RespondJSON(w, &controllers.ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)

		return
	}
	userID := domain.UserID(id)
	user, err := su.Service.GetUser(userID)
	if err != nil {
		controllers.RespondJSON(w, &controllers.ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)

		return
	}

	controllers.RespondJSON(w, user, http.StatusOK)
}
