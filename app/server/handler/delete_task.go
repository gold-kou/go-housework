package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	log "github.com/sirupsen/logrus"
)

// DeleteTask handler top
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	common.Transact(func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		taskRepo := repository.NewTaskRepository(tx)
		h := DeleteTaskHandler{srv: service.NewDeleteTask(userRepo, taskRepo)}
		resp, status, err := h.DeleteTask(w, r)
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

// DeleteTaskHandler struct
type DeleteTaskHandler struct {
	srv service.DeleteTaskServiceInterface
}

// DeleteTask handler
func (h DeleteTaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) (resp interface{}, status int, err error) {
	// verify header token
	authUser, err := middleware.VerifyHeaderToken(r)
	if err != nil {
		return common.NewAuthorizationError(err.Error()), http.StatusUnauthorized, err
	}

	// get path parameter
	vars := mux.Vars(r)
	taskIDStr, ok := vars["task_id"]
	if !ok {
		return common.NewBadRequestError(err.Error()), http.StatusBadRequest, err
	}
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return common.NewInternalServerError(err.Error()), http.StatusInternalServerError, err
	}

	// service layer
	err = h.srv.Execute(authUser, uint64(taskID))

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

	return &schemamodel.ResponseDeleteTask{Message: "the task is deleted"}, http.StatusOK, nil

}
