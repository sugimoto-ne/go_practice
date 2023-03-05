package controllers

import (
	"net/http"

	"github.com/sugimoto-ne/go_practice.git/interface/controllers"
	usecases "github.com/sugimoto-ne/go_practice.git/usecases/user"
)

type ListUser struct {
	Service usecases.ListUser
}

func (su *ListUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	user, err := su.Service.ListUser()
	if err != nil {
		controllers.RespondJSON(w, &controllers.ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)

		return
	}

	controllers.RespondJSON(w, user, http.StatusOK)
}
