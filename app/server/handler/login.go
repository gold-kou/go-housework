package handler

import (
	"encoding/json"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"net/http"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/jinzhu/gorm"

	validation "github.com/go-ozzo/ozzo-validation"
	log "github.com/sirupsen/logrus"
)

// Login handler top
func Login(w http.ResponseWriter, r *http.Request) {
	common.Transact(func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		h := LoginHandler{tok: middleware.NewTokenStruct(), srv: service.NewLogin(userRepo)}
		resp, status, err := h.Login(w, r)
		if err != nil {
			log.Error(err)
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(status)
		if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
			log.Error(encodeErr)
			panic(encodeErr.Error())
		}
		return err
	})
}

// LoginHandler struct
type LoginHandler struct {
	tok middleware.TokenInterface
	srv service.LoginServiceInterface
}

// Login - ログインAPI
func (h LoginHandler) Login(w http.ResponseWriter, r *http.Request) (resp interface{}, status int, err error) {
	// get query parameters
	userName := r.URL.Query().Get("user_name")
	password := r.URL.Query().Get("password")

	// validation
	if err := validation.Validate(userName, validation.Required); err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}
	if err := validation.Validate(password, validation.Required); err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}

	// service layer
	tokenString, err := h.srv.Execute(userName, password)

	// error handling
	switch err := err.(type) {
	case nil:
	case *common.BadRequestError:
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	case *common.AuthorizationError:
		return common.NewAuthorizationError(err.Error()), http.StatusNonAuthoritativeInfo, err
	default:
		return common.NewInternalServerError(err.Error()), http.StatusInternalServerError, err
	}

	return &schemamodel.ResponseLogin{Token: tokenString}, http.StatusOK, nil
}
