package transaction

import (
	"database/sql"
	"log"

	"20dojo-online/pkg/db"
	repository "20dojo-online/pkg/server/domain/repository/transaction"
)

type txPersistence struct{}

// NewTxPersistence Tx に関するPersistenceを生成
func NewTxPersistence() repository.TxRepository {
	return &txPersistence{}
}

// Transaction トランザクション処理
func (tp txPersistence) Transaction(function func(any interface{}, tx *sql.Tx) error) error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}

	// TODO: 書き方再検討
	defer func() {
		// panic
		if err := recover(); err != nil {
			log.Println("!! PANIC !!")
			log.Println(err)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println("failed to Rollback")
				log.Println(rollbackErr)
			}
		}
		// TODO: ここにpanic()を書くとAPIが二回実行される
		// panic(err)
	}()

	// 実行
	// TODO: エラーの書き方を考える
	if err = function(any, tx); err != nil {
		log.Fatal(err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Println("failed to Rollback")
			log.Println(rollbackErr)
			return rollbackErr
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// func (tp txPersistence) Transaction(newItemSlice []*model.UserCollectionItem, user *model.UserL, tx *sql.Tx) error {
// 	// 3-1. バルクインサート
// 	if len(newItemSlice) != 0 {
// 		if err := gu.ucItemRepository.BulkInsertUCItemSlice(newItemSlice, tx); err != nil {
// 			return err
// 		}
// 	}
// 	// 3-2. ユーザの保持コイン更新
// 	if err := gu.userRepository.UpdateUserByUserInTx(user, tx); err != nil {
// 		return err
// 	}
// 	return nil
// }
