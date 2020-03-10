package common

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

// MockDB Helper for DB mock
// "`sqlmock` is not officially supported, running under compatibility mode." should be ignored
func MockDB(t *testing.T, mockFunc func(db *gorm.DB, mock sqlmock.Sqlmock)) {
	// connect to not real db but sqlmock
	db, mock, _ := sqlmock.New()
	defer db.Close()
	gdb, _ := gorm.Open("sqlmock", db)
	defer gdb.Close()
	mockFunc(gdb, mock)

	t.Helper()
	assert := assert.New(t)
	assert.NoError(mock.ExpectationsWereMet())
}

// Transact Helper for DB transaction
// ex:
// err := common.Transact(func(tx *gorm.DB) (err error) {
//   // something
//   return
// })
// if err != nil {
//   // error handling
// }
func Transact(txFunc func(*gorm.DB) error) (err error) {
	if db == nil {
		log.Error("must set DB Connection")
		panic("must set DB Connection")
	}
	tx := db.Begin()
	if err = tx.Error; err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			log.Error(p)
			if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
				// rollbackでエラーが発生した場合、後で追えるようログを出力する。（本来のpanicをrethrowするため、panicしない）
				log.Error(rollbackErr)
			}
			// panicは握りつぶさずにrethrowする
			panic(p)
		} else if err != nil {
			log.Warn(err)
			if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
				log.Error(rollbackErr)
				// rollbackでエラーが発生した場合、処理を継続できないためpanicする。
				panic(rollbackErr)
			}
		} else {
			err = tx.Commit().Error
		}
	}()
	err = txFunc(tx)
	return err
}
