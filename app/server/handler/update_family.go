package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// UpdateFamily handler top
func UpdateFamily(w http.ResponseWriter, r *http.Request) {
	common.Transact(func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		familyRepo := repository.NewFamilyRepository(tx)
		memberFamilyRepo := repository.NewMemberFamilyRepository(tx)
		h := UpdateFamilyHandler{srv: service.NewUpdateFamily(userRepo, familyRepo, memberFamilyRepo)}
		resp, status, err := h.UpdateFamily(w, r)
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

// UpdateFamilyHandler struct
type UpdateFamilyHandler struct {
	srv service.UpdateFamilyServiceInterface
}

// UpdateFamily handler
func (h UpdateFamilyHandler) UpdateFamily(w http.ResponseWriter, r *http.Request) (resp interface{}, status int, err error) {
	// verify header token
	authUser, err := middleware.VerifyHeaderToken(r)
	if err != nil {
		return common.NewAuthorizationError(err.Error()), http.StatusUnauthorized, err
	}

	// get request parameter
	var updateFamily schemamodel.RequestUpdateFamily
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}
	defer r.Body.Close()
	if err := json.Unmarshal(b, &updateFamily); err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}

	// validation
	if err := updateFamily.ValidateParam(); err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}

	// service layer
	f, err := h.srv.Execute(authUser, &updateFamily)

	// error handling
	switch err := err.(type) {
	case nil:
	case *common.BadRequestError:
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	case *common.AuthorizationError:
		return common.NewAuthorizationError(err.Error()), http.StatusNonAuthoritativeInfo, err
	default:
		return common.NewAuthorizationError(err.Error()), http.StatusNonAuthoritativeInfo, err
	}

	return &schemamodel.ResponseUpdateFamily{Family: schemamodel.Family{FamilyId: int64(f.ID), FamilyName: f.Name}}, http.StatusOK, nil
}
