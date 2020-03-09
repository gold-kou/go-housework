package handler

import (
	"encoding/json"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// get request parameter
	var createUser schemamodel.RequestCreateUser
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn(err)
		common.ResponseBadRequest(w, err.Error())
		return
	}
	defer r.Body.Close()
	if err := json.Unmarshal(b, &createUser); err != nil {
		log.Warn(err)
		common.ResponseBadRequest(w, err.Error())
		return
	}

	// validation
	if err := createUser.ValidateParam(); err != nil {
		log.Warn(err)
		common.ResponseBadRequest(w, err.Error())
		return
	}

	// service layer
	var u *db.User
	err = common.Transact(func(tx *gorm.DB) (err error) {
		userRepo := repository.NewUserRepository(tx)
		u, err = service.NewCreateUser(tx, &createUser, userRepo).Execute()
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
	if err := json.NewEncoder(w).Encode(&schemamodel.ResponseCreateUser{
		User:    schemamodel.User{UserId: int64(u.ID), UserName: u.Name},
		Message: "new user created"}); err != nil {
		log.Error(err)
		panic(err)
	}
}
