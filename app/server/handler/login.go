package handler

import (
	"encoding/json"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/jinzhu/gorm"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Login - ログインAPI
func Login(w http.ResponseWriter, r *http.Request) {
	// クエリパラメータの取得
	userName := r.URL.Query().Get("user_name")
	password := r.URL.Query().Get("password")

	// バリデーションチェック
	if userName == "" {
		common.ResponseBadRequest(w, "user_name parameter is required")
		return
	}
	if password == "" {
		common.ResponseBadRequest(w, "password parameter is required")
		return
	}

	// service層へ処理を移譲
	var tokenString string
	err := common.Transact(func(tx *gorm.DB) (err error) {
		userRepo := repository.NewUserRepository(tx)
		tokenString, err = service.NewLogin(tx, userName, password, userRepo).Execute()
		return
	})

	// エラーハンドリング
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
