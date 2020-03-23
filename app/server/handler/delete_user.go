package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/jinzhu/gorm"

	log "github.com/sirupsen/logrus"
)

// DeleteUser - delete user API
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// get jwt from header
	authHeader := r.Header.Get("Authorization")
	// Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM3MjIwNTMsImlhdCI6IjIwMjAtMDMtMDhUMTE6NDc6MzMuMTc4NjU5MyswOTowMCIsIm5hbWUiOiJ0ZXN0In0.YIyT1RJGcYbdynx1V4-6MhiosmTlHmKiyiG_GjxQeuw
	bearerToken := strings.Split(authHeader, " ")[1]

	// verify jwt
	authUser, err := middleware.VerifyToken(bearerToken)
	if err != nil {
		common.ResponseUnauthorized(w, err.Error())
		return
	}

	if err == nil {
		// service layer
		err = common.Transact(func(tx *gorm.DB) (err error) {
			userRepo := repository.NewUserRepository(tx)
			err = service.NewDeleteUser(tx, userRepo).Execute(authUser)
			return
		})
	}

	// errors handling
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

	// HTTP response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&schemamodel.ResponseDeleteUser{Message: "the user deleted"}); err != nil {
		log.Error(err)
		panic(err)
	}
}
