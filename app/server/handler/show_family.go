package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gold-kou/go-housework/app/server/service"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// ShowFamily handler
func ShowFamily(w http.ResponseWriter, r *http.Request) {
	// verify header token
	authUser, err := middleware.VerifyHeaderToken(r)
	if err != nil {
		common.ResponseUnauthorized(w, err.Error())
		return
	}

	// service layer
	var f *db.Family
	err = common.Transact(func(tx *gorm.DB) (err error) {
		familyRepo := repository.NewFamilyRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		memberFamilyRepo := repository.NewMemberFamilyRepository(tx)
		f, err = service.NewShowFamily(tx, familyRepo, userRepo, *memberFamilyRepo).Execute(authUser)
		return
	})

	// error handling
	switch err := err.(type) {
	case nil:
	case *common.BadRequestError:
		log.Warn(err)
		common.ResponseBadRequest(w, err.Message)
		return
	case *common.AuthorizationError:
		log.Warn(err)
		common.ResponseUnauthorized(w, err.Message)
		return
	default:
		log.Error(err)
		common.ResponseInternalServerError(w, err.Error())
		return
	}

	// http response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&schemamodel.ResponseShowFamily{
		Family: schemamodel.Family{FamilyId: int64(f.ID), FamilyName: f.Name}}); err != nil {
		log.Error(err)
		panic(err)
	}
}
