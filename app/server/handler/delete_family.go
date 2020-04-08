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

// DeleteFamily handler top
func DeleteFamily(w http.ResponseWriter, r *http.Request) {
	common.Transact(func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		familyRepo := repository.NewFamilyRepository(tx)
		memberFamilyRepo := repository.NewMemberFamilyRepository(tx)
		h := DeleteFamilyHandler{tok: middleware.NewTokenStruct(), srv: service.NewDeleteFamily(userRepo, familyRepo, memberFamilyRepo)}
		resp, status, err := h.DeleteFamily(w, r)
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

// DeleteFamilyHandler struct
type DeleteFamilyHandler struct {
	tok middleware.TokenInterface
	srv service.DeleteFamilyServiceInterface
}

// DeleteFamily handler
func (h DeleteFamilyHandler) DeleteFamily(w http.ResponseWriter, r *http.Request) (resp interface{}, status int, err error) {
	// verify header token
	authUser, err := h.tok.VerifyHeaderToken(r)
	if err != nil {
		return common.NewAuthorizationError(err.Error()), http.StatusUnauthorized, err
	}

	// service layer
	err = h.srv.Execute(authUser)

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

	return &schemamodel.ResponseDeleteFamily{Message: "delete complete"}, http.StatusOK, nil
}
