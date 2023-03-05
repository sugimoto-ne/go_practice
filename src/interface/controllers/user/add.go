package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sugimoto-ne/go_practice.git/interface/controllers"
	usecases "github.com/sugimoto-ne/go_practice.git/usecases/user"
)

type AddUser struct {
	Service   usecases.AddUser
	Validator *validator.Validate
}

func (au *AddUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		controllers.RespondJSON(w, &controllers.ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)

		return
	}

	// バリデーションの検証
	err := au.Validator.Struct(body)
	if err != nil {
		controllers.RespondJSON(w, &controllers.ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)

		return
	}

	id, err := au.Service.AddUser(body.Name, body.Password)
	if err != nil {
		controllers.RespondJSON(w, &controllers.ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)

		return
	}

	rsp := struct {
		ID int64
	}{
		ID: int64(id),
	}

	controllers.RespondJSON(w, rsp, http.StatusOK)
}
