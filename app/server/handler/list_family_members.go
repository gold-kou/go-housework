package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gold-kou/go-housework/app/server/service"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// ListFamilyMembers handler top
func ListFamilyMembers(w http.ResponseWriter, r *http.Request) {
	common.Transact(func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		familyRepo := repository.NewFamilyRepository(tx)
		memberFamilyRepo := repository.NewMemberFamilyRepository(tx)
		h := ListFamilyMembersHandler{srv: service.NewListFamilyMembers(userRepo, familyRepo, memberFamilyRepo)}
		resp, status, err := h.ListFamilyMembers(w, r)
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

// ListFamilyMembersHandler struct
type ListFamilyMembersHandler struct {
	srv service.ListFamilyMembersServiceInterface
}

// ListFamilyMembers handler
func (h ListFamilyMembersHandler) ListFamilyMembers(w http.ResponseWriter, r *http.Request) (resp interface{}, status int, err error) {
	// verify header token
	authUser, err := middleware.VerifyHeaderToken(r)
	if err != nil {
		return common.NewAuthorizationError(err.Error()), http.StatusUnauthorized, err
	}

	// service layer
	f, us, err := h.srv.Execute(authUser)

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

	var members []schemamodel.Member
	for _, u := range us {
		m := schemamodel.Member{MemberId: int64(u.ID), MemberName: u.Name}
		members = append(members, m)
	}
	return &schemamodel.ResponseListFamilyMembers{Family: schemamodel.Family{FamilyId: int64(f.ID), FamilyName: f.Name}, Members: members}, http.StatusOK, nil
}
