package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/jinzhu/gorm"

	validation "github.com/go-ozzo/ozzo-validation"
	log "github.com/sirupsen/logrus"
)

// ListTasks handler top
func ListTasks(w http.ResponseWriter, r *http.Request) {
	common.Transact(func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		familyRepo := repository.NewFamilyRepository(tx)
		memberFamilyRepo := repository.NewMemberFamilyRepository(tx)
		taskRepo := repository.NewTaskRepository(tx)
		h := ListTasksHandler{tok: middleware.NewTokenStruct(), srv: service.NewListTasks(userRepo, familyRepo, memberFamilyRepo, taskRepo)}
		resp, status, err := h.ListTasks(w, r)
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

// ListTasksHandler struct
type ListTasksHandler struct {
	tok middleware.TokenInterface
	srv service.ListTasksServiceInterface
}

// ListTasks handler
func (h ListTasksHandler) ListTasks(w http.ResponseWriter, r *http.Request) (resp interface{}, status int, err error) {
	// verify header token
	authUser, err := h.tok.VerifyHeaderToken(r)
	if err != nil {
		return common.NewAuthorizationError(err.Error()), http.StatusUnauthorized, err
	}

	// get query parameter
	targetDate := r.URL.Query().Get("date")

	// validation
	if err := validation.Validate(targetDate, validation.Required, validation.Date("2006-01-02")); err != nil {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}

	// service layer
	dbTasks, dbFamily, dbUsers, err := h.srv.Execute(authUser, targetDate)

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

	var tasks []schemamodel.Task
	for i, t := range dbTasks {
		tasks = append(tasks, schemamodel.Task{TaskId: int64(t.ID), TaskName: t.Name,
			MemberName: dbUsers[i].Name, Status: t.Status, Date: t.Date})
	}
	return &schemamodel.ResponseListTasks{Family: schemamodel.Family{FamilyId: int64(dbFamily.ID), FamilyName: dbFamily.Name}, Tasks: tasks}, http.StatusOK, nil
}
