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

// UpdateTask handler top
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	common.Transact(func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		familyRepo := repository.NewFamilyRepository(tx)
		memberFamilyRepo := repository.NewMemberFamilyRepository(tx)
		taskRepo := repository.NewTaskRepository(tx)
		h := UpdateTaskHandler{tok: middleware.NewTokenStruct(), srv: service.NewUpdateTask(userRepo, familyRepo, memberFamilyRepo, taskRepo)}
		resp, status, err := h.UpdateTask(w, r)
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

// UpdateTaskHandler struct
type UpdateTaskHandler struct {
	tok middleware.TokenInterface
	srv service.UpdateTaskServiceInterface
}

// UpdateTask handler
func (h UpdateTaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) (resp interface{}, status int, err error) {
	// verify header token
	authUser, err := h.tok.VerifyHeaderToken(r)
	if err != nil {
		return common.NewAuthorizationError(err.Error()), http.StatusUnauthorized, err
	}

	// get request parameter
	var reqUpdateTask schemamodel.RequestUpdateTask
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}
	defer r.Body.Close()
	if err := json.Unmarshal(b, &reqUpdateTask); err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}

	// validation
	if err := reqUpdateTask.ValidateParam(); err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}

	// service layer
	t, f, u, err := h.srv.Execute(authUser, &reqUpdateTask)

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

	return &schemamodel.ResponseUpdateTask{Family: schemamodel.Family{FamilyId: int64(f.ID), FamilyName: f.Name},
		Task: schemamodel.Task{TaskId: int64(t.ID), TaskName: t.Name, MemberName: u.Name, Status: t.Status, Date: t.Date},
	}, http.StatusOK, nil
}
