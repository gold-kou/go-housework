package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/jinzhu/gorm"

	validation "github.com/go-ozzo/ozzo-validation"
	log "github.com/sirupsen/logrus"
)

// Login - ログインAPI
func Login(w http.ResponseWriter, r *http.Request) {
	// get query parameters
	userName := r.URL.Query().Get("user_name")
	password := r.URL.Query().Get("password")

	// validation
	if err := validation.Validate(userName, validation.Required); err != nil {
		common.ResponseBadRequest(w, err.Error())
		return
	}
	if err := validation.Validate(password, validation.Required); err != nil {
		common.ResponseBadRequest(w, err.Error())
		return
	}

	// service layer
	var tokenString string
	err := common.Transact(func(tx *gorm.DB) (err error) {
		userRepo := repository.NewUserRepository(tx)
		tokenString, err = service.NewLogin(tx, userName, password, userRepo).Execute()
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

	// HTTPレスポンス作成
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&schemamodel.ResponseLogin{Token: tokenString}); err != nil {
		log.Error(err)
		panic(err)
	}
}
