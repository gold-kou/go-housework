package common

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// MockDB DBモックのヘルパー関数。
// 例：
// MockDB(func(db *gorm.DB, mock sqlmock.Sqlmock) {
// 	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "order_item_measurement" ("key","value","created_at","updated_at","order_item_id") VALUES (?,?,?,?,?)`)).
// 		WithArgs("neckline", 21.0, &testNow, &testNow, 1).
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// 	err := createOrderItemMeasurement(db, "neckline", 21, 1)
// 	assert.NoError(err)
// 	assert.NoError(mock.ExpectationsWereMet())
// })
func MockDB(mockFunc func(db *gorm.DB, mock sqlmock.Sqlmock)) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	// テスト実行時、「`sqlmock` is not officially supported, running under compatibility mode.」
	// という警告が出るがsqlmockを使用しているために発生しているものなので無視している
	gdb, _ := gorm.Open("sqlmock", db)
	defer gdb.Close()
	mockFunc(gdb, mock)
}

// Transact DBトランザクションのヘルパー関数。
// 例:
// err := common.Transact(func(tx *gorm.DB) (err error) {
//   // トランザクション内のハンドリング
//   return
// })
// if err != nil {
//   // エラーハンドリング
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
