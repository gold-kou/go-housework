package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/jinzhu/gorm"

	log "github.com/sirupsen/logrus"
)

// CreateTask handler
func CreateTask(w http.ResponseWriter, r *http.Request) {
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

	// get request parameter
	var reqCreateTask schemamodel.RequestCreateTask
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn(err)
		common.ResponseBadRequest(w, err.Error())
		return
	}
	defer r.Body.Close()
	if err := json.Unmarshal(b, &reqCreateTask); err != nil {
		log.Warn(err)
		common.ResponseBadRequest(w, err.Error())
		return
	}

	// validation
	if err := reqCreateTask.ValidateParam(); err != nil {
		log.Warn(err)
		common.ResponseBadRequest(w, err.Error())
		return
	}

	// service layer
	var t *db.Task
	var f *db.Family
	var u *db.User
	err = common.Transact(func(tx *gorm.DB) (err error) {
		taskRepo := repository.NewTaskRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		familyRepo := repository.NewFamilyRepository(tx)
		memberFamilyRepo := repository.NewMemberFamilyRepository(tx)
		t, f, u, err = service.NewCreateTask(tx, &reqCreateTask, taskRepo, userRepo, familyRepo, memberFamilyRepo).Execute(authUser)
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
	if err := json.NewEncoder(w).Encode(&schemamodel.ResponseCreateTask{
		Family: schemamodel.Family{FamilyId: int64(f.ID), FamilyName: f.Name},
		Task: schemamodel.Task{TaskId: int64(t.ID), TaskName: t.Name,
			MemberName: u.Name, Status: t.Status, Date: t.Date},
	}); err != nil {
		log.Error(err)
		panic(err)
	}
}
