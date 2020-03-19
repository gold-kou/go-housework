package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/jinzhu/gorm"

	validation "github.com/go-ozzo/ozzo-validation"
	log "github.com/sirupsen/logrus"
)

// ListTasks handler
func ListTasks(w http.ResponseWriter, r *http.Request) {
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

	// get query parameter
	targetDate := r.URL.Query().Get("date")

	// validation
	if err := validation.Validate(targetDate, validation.Required, validation.Date("2006-01-02")); err != nil {
		common.ResponseBadRequest(w, err.Error())
		return
	}

	// service layer
	var dbTasks []*db.Task
	var dbFamily *db.Family
	var dbUsers []*db.User
	err = common.Transact(func(tx *gorm.DB) (err error) {
		taskRepo := repository.NewTaskRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		familyRepo := repository.NewFamilyRepository(tx)
		memberFamilyRepo := repository.NewMemberFamilyRepository(tx)
		dbTasks, dbFamily, dbUsers, err = service.NewListTasks(tx, targetDate, taskRepo, userRepo, familyRepo, memberFamilyRepo).Execute(authUser)
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

	// TODO RDBテーブル変更した影響で、dbTaskにfamily情報も入っているから、create/updateでfamilyをserviceで返す必要がなくなった。その対応。
	// http response
	var tasks []schemamodel.Task
	for i, t := range dbTasks {
		tasks = append(tasks, schemamodel.Task{TaskId: int64(t.ID), TaskName: t.Name,
			MemberName: dbUsers[i].Name, Status: t.Status, Date: t.Date})
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&schemamodel.ResponseListTasks{
		Family: schemamodel.Family{FamilyId: int64(dbFamily.ID), FamilyName: dbFamily.Name},
		Tasks:  tasks,
	}); err != nil {
		log.Error(err)
		panic(err)
	}
}
