package common

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // load db driver
	log "github.com/sirupsen/logrus"
)

// LogInit log init set up
func LogInit() {
	// stdout
	log.SetOutput(os.Stdout)

	// JSON format
	log.SetFormatter(&log.JSONFormatter{})

	// local: over TRACE, stg/prd: over INFO
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
