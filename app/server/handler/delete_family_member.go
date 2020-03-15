package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// DeleteFamilyMember handler
func DeleteFamilyMember(w http.ResponseWriter, r *http.Request) {
	// get jwt from header
	authHeader := r.Header.Get("Authorization")
	// Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM3MjIwNTMsImlhdCI6IjIwMjAtMDMtMDhUMTE6NDc6MzMuMTc4NjU5MyswOTowMCIsIm5hbWUiOiJ0ZXN0In0.YIyT1RJGcYbdynx1V4-6MhiosmTlHmKiyiG_GjxQeuw
	bearerToken := strings.Split(authHeader, " ")[1]

	// verify jwt
	authUser, err := middleware.VerifyToken(bearerToken)

	// PathParameter
	vars := mux.Vars(r)
	memberIDStr, ok := vars["member_id"]
	if !ok {
		common.ResponseBadRequest(w, "Missing required parameter: member_id")
		return
	}
	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		common.ResponseInternalServerError(w, err.Error())
	}

	// service layer
	err = common.Transact(func(tx *gorm.DB) (err error) {
		familyRepo := repository.NewFamilyRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		memberFamilyRepo := repository.NewMemberFamilyRepository(tx)
		err = service.NewDeleteFamilyMember(tx, uint64(memberID), familyRepo, userRepo, *memberFamilyRepo).Execute(authUser)
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
	if err := json.NewEncoder(w).Encode(&schemamodel.ResponseDeleteFamilyMember{
		Message: "deleted the member from the family"}); err != nil {
		log.Error(err)
		panic(err)
	}
}
