package handler

import (
	"encoding/json"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	common.Transact(func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		c := Controller{srv: service.NewCreateUser(userRepo)}
		res, status, err := c.CreateUser(w, r)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(status)
		if encodeErr := json.NewEncoder(w).Encode(res); encodeErr != nil {
			log.Error(encodeErr)
			panic(encodeErr.Error())
		}
		return err
	})
}

type Controller struct {
	srv service.CreateUserServiceInterface
}

func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn(err)
		return schemamodel.ResponseBadRequest{Message: err.Error()}, http.StatusBadRequest, err
	}
	defer r.Body.Close()

	var user schemamodel.RequestCreateUser
	if err := json.Unmarshal(b, &user); err != nil {
		log.Warn(err)
		return schemamodel.ResponseBadRequest{Message: err.Error()}, http.StatusBadRequest, err
	}

	// validation
	if err := user.ValidateParam(); err != nil {
		log.Warn(err)
		return schemamodel.ResponseBadRequest{Message: err.Error()}, http.StatusBadRequest, err
	}

	// service layer
	u, err := c.srv.Execute(&user)

	// error handling
	switch err := err.(type) {
	case nil:
	case *common.BadRequestError:
		log.Warn(err)
		return schemamodel.ResponseBadRequest{Message: err.Error()}, http.StatusBadRequest, err
	case *common.AuthorizationError:
		log.Warn(err)
		return schemamodel.ResponseUnauthorized{Message: err.Error()}, http.StatusUnauthorized, err
	default:
		log.Error(err)
		return schemamodel.ResponseInternalServerError{Message: err.Error()}, http.StatusInternalServerError, err
	}

	return schemamodel.ResponseCreateUser{
		User:    schemamodel.User{UserId: int64(u.ID), UserName: u.Name},
		Message: "new user created",
	}, http.StatusOK, nil
}
