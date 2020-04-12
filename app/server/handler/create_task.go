package handler

import (
	"encoding/json"
	"github.com/gold-kou/go-housework/app/model"
	"io/ioutil"
	"net/http"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/jinzhu/gorm"

	log "github.com/sirupsen/logrus"
)

// CreateTask handler top
func CreateTask(w http.ResponseWriter, r *http.Request) {
	common.Transact(func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		familyRepo := repository.NewFamilyRepository(tx)
		memberFamilyRepo := repository.NewMemberFamilyRepository(tx)
		taskRepo := repository.NewTaskRepository(tx)
		h := CreateTaskHandler{tok: middleware.NewTokenStruct(), srv: service.NewCreateTask(userRepo, familyRepo, memberFamilyRepo, taskRepo)}
		resp, status, err := h.CreateTask(w, r)
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

// CreateTaskHandler struct
type CreateTaskHandler struct {
	tok      middleware.TokenInterface
	authUser *model.Auth
	srv      service.CreateTaskServiceInterface
}

// CreateTask handler
func (h CreateTaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) (resp interface{}, status int, err error) {
	// verify header token
	authUser, err := h.tok.VerifyHeaderToken(r)
	if err != nil {
		return common.NewAuthorizationError(err.Error()), http.StatusUnauthorized, err
	}

	// get request parameter
	var reqCreateTask schemamodel.RequestCreateTask
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}
	defer r.Body.Close()
	if err := json.Unmarshal(b, &reqCreateTask); err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}

	// validation
	if err := reqCreateTask.ValidateParam(); err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}

	// service layer
	var u *db.User
	var f *db.Family
	var t *db.Task
	u, f, t, err = h.srv.Execute(authUser, &reqCreateTask)

	// error handling
	switch err := err.(type) {
	case nil:
	case *common.BadRequestError:
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	case *common.AuthorizationError:
		return common.NewAuthorizationError(err.Error()), http.StatusUnauthorized, err
	default:
		return common.NewInternalServerError(err.Error()), http.StatusInternalServerError, err
	}

	return &schemamodel.ResponseCreateTask{Family: schemamodel.Family{FamilyId: int64(f.ID), FamilyName: f.Name}, Task: schemamodel.Task{TaskId: int64(t.ID), TaskName: t.Name, MemberName: u.Name, Status: t.Status, Date: t.Date}},
		http.StatusOK, nil
}
