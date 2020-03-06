package common

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // load db driver
	log "github.com/sirupsen/logrus"
)

// LogInit ログの初期化処理
func LogInit() {
	// 出力先は標準出力とする
	log.SetOutput(os.Stdout)

	// JSONで出力する
	log.SetFormatter(&log.JSONFormatter{})

	// 開発時は全てのログ（TRACE以上）、それ以外（テスト、ステージング、本番環境）はINFO以上のログレベルを出力する
	log.SetLevel(log.InfoLevel)
	if os.Getenv("RUNSERVER") == "LOCAL" {
		log.SetLevel(log.TraceLevel)
	}
}

// callback that logs sql error
func logSQLError(scope *gorm.Scope) {
	err := scope.DB().Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		// RecordNotFoundErr以外のエラー発生時は、エラー内容、SQL、プレースホルダをログに出力する
		log.WithFields(log.Fields{"module": "gorm"}).Error(err, scope.SQL, scope.SQLVars)
	}
}
