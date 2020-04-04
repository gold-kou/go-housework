package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/jinzhu/gorm"

	log "github.com/sirupsen/logrus"
)

// DeleteUser handler top
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	common.Transact(func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		h := DeleteUserHandler{srv: service.NewDeleteUser(userRepo)}
		resp, status, err := h.DeleteUser(w, r)
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

// DeleteUserHandler struct
type DeleteUserHandler struct {
	srv service.DeleteUserServiceInterface
}

// DeleteUser - delete user API
func (h DeleteUserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) (resp interface{}, status int, err error) {
	// verify header token
	authUser, err := middleware.VerifyHeaderToken(r)
	if err != nil {
		return common.NewAuthorizationError(err.Error()), http.StatusUnauthorized, err
	}

	// service layer
	err = h.srv.Execute(authUser)

	// errors handling
	switch err := err.(type) {
	case nil:
	case *common.BadRequestError:
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	case *common.AuthorizationError:
		return common.NewAuthorizationError(err.Error()), http.StatusNonAuthoritativeInfo, err
	default:
		return common.NewInternalServerError(err.Error()), http.StatusInternalServerError, err
	}

	return &schemamodel.ResponseDeleteUser{Message: "the user deleted"}, http.StatusOK, nil
}
