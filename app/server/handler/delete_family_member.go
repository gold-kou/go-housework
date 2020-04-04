package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// DeleteFamilyMember handler top
func DeleteFamilyMember(w http.ResponseWriter, r *http.Request) {
	common.Transact(func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		familyRepo := repository.NewFamilyRepository(tx)
		memberFamilyRepo := repository.NewMemberFamilyRepository(tx)
		h := DeleteFamilyMemberHandler{srv: service.NewDeleteFamilyMember(userRepo, familyRepo, memberFamilyRepo)}
		resp, status, err := h.DeleteFamilyMember(w, r)
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

// DeleteFamilyMemberHandler struct
type DeleteFamilyMemberHandler struct {
	srv service.DeleteFamilyMemberServiceInterface
}

// DeleteFamilyMember handler
func (h DeleteFamilyMemberHandler) DeleteFamilyMember(w http.ResponseWriter, r *http.Request) (resp interface{}, status int, err error) {
	// verify header token
	authUser, err := middleware.VerifyHeaderToken(r)
	if err != nil {
		return common.NewAuthorizationError(err.Error()), http.StatusUnauthorized, err
	}

	// get PathParameter
	vars := mux.Vars(r)
	memberIDStr, ok := vars["member_id"]
	if !ok {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}
	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}

	// service layer
	err = h.srv.Execute(authUser, uint64(memberID))

	// error handling
	switch err := err.(type) {
	case nil:
	case *common.BadRequestError:
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	case *common.AuthorizationError:
		return common.NewAuthorizationError(err.Error()), http.StatusBadRequest, err
	default:
		return common.NewInternalServerError(err.Error()), http.StatusBadRequest, err
	}

	return &schemamodel.ResponseDeleteFamilyMember{Message: "deleted the member from the family"}, http.StatusOK, nil
}
