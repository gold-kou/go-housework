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

// DeleteTask handler
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	// verify header token
	authUser, err := middleware.VerifyHeaderToken(r)
	if err != nil {
		common.ResponseUnauthorized(w, err.Error())
		return
	}

	// get path parameter
	vars := mux.Vars(r)
	taskIDStr, ok := vars["task_id"]
	if !ok {
		common.ResponseBadRequest(w, "Missing required parameter: task_id")
		return
	}
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		common.ResponseInternalServerError(w, err.Error())
	}

	// service layer
	err = common.Transact(func(tx *gorm.DB) (err error) {
		taskRepo := repository.NewTaskRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		err = service.NewDeleteTask(tx, uint64(taskID), taskRepo, userRepo).Execute(authUser)
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
	if err := json.NewEncoder(w).Encode(&schemamodel.ResponseDeleteTask{
		Message: "the task is deleted"}); err != nil {
		log.Error(err)
		panic(err)
	}
}
