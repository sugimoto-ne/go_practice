package infrastracture

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sugimoto-ne/go_practice.git/config"
	"github.com/sugimoto-ne/go_practice.git/domain"
	"github.com/sugimoto-ne/go_practice.git/infrastracture/logger"
	userController "github.com/sugimoto-ne/go_practice.git/interface/controllers/user"
	"github.com/sugimoto-ne/go_practice.git/interface/store"
	usecases "github.com/sugimoto-ne/go_practice.git/usecases/user"
)

func NewMux(cfg *config.Config) (http.Handler, error) {
	// 各種ログの初期化
	reqLogger, err := logger.NewLogger("../logs", cfg.Env, "req", "json")
	if err != nil {
		return nil, err
	}
	appLogger, err := logger.NewLogger("../logs", cfg.Env, "application", "json")
	if err != nil {
		return nil, err
	}
	resLogger, err := logger.NewLogger("../logs", cfg.Env, "res", "json")
	if err != nil {
		return nil, err
	}

	mux := chi.NewRouter()

	mux.Get("/sample", func(w http.ResponseWriter, r *http.Request) {
		// リクエストログ
		reqLogger.Logger.Info("from hellohandler", "from", r.RemoteAddr, "method", r.Method)

		time.Sleep(10 * time.Second)

		appLogger.Logger.Info("success time sleep")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		rsp := struct {
			Message string `json:"message"`
		}{
			Message: "Hello, world!",
		}
		bodyBytes, _ := json.Marshal(rsp)
		if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
			appLogger.Logger.Error("write response error", err)
		}

		// レスポンスログ
		resLogger.Logger.Info("from hellohandler", "body", rsp, "to", r.RemoteAddr)
	})

	v := validator.New()

	UserRepo := &store.UserRepository{
		LastID: domain.UserID(len(store.Users)),
	}

	AddUserController := userController.AddUser{
		Service: usecases.AddUser{
			UserRepository: UserRepo,
		},
		Validator: v,
	}

	mux.Post("/users", AddUserController.ServeHTTP)

	GetUserController := userController.ShowUser{
		Service: usecases.GetUser{
			UserRepository: UserRepo,
		},
	}
	mux.Get("/users/{id}", GetUserController.ServeHTTP)

	return mux, nil
}
