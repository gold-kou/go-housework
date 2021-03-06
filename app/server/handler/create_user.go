package handler

import (
	"encoding/json"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"io/ioutil"
	"net/http"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// CreateUser handler top
func CreateUser(w http.ResponseWriter, r *http.Request) {
	common.Transact(func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		h := CreateUserHandler{tok: middleware.NewTokenStruct(), srv: service.NewCreateUser(userRepo)}
		resp, status, err := h.CreateUser(w, r)
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

// CreateUserHandler struct
type CreateUserHandler struct {
	tok middleware.TokenInterface
	srv service.CreateUserServiceInterface
}

// CreateUser handler
func (h *CreateUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) (resp interface{}, status int, err error) {
	// get request parameter
	var createUser schemamodel.RequestCreateUser
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}
	defer r.Body.Close()
	if err := json.Unmarshal(b, &createUser); err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}

	// validation
	if err := createUser.ValidateParam(); err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}

	// service layer
	u, err := h.srv.Execute(&createUser)
	switch err := err.(type) {
	case nil:
	case *common.BadRequestError:
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	case *common.AuthorizationError:
		return common.NewAuthorizationError(err.Error()), http.StatusUnauthorized, err
	default:
		return common.NewInternalServerError(err.Error()), http.StatusInternalServerError, err
	}

	return &schemamodel.ResponseCreateUser{User: schemamodel.User{UserId: int64(u.ID), UserName: u.Name}}, http.StatusOK, nil
}
