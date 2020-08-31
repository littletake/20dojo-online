// Package persistence 技術的関心事を扱う
package persistence

import (
	"database/sql"
	"log"

	"20dojo-online/pkg/db"
	"20dojo-online/pkg/server/domain/model"
	"20dojo-online/pkg/server/domain/repository"
)

type userPersistence struct{}

// NewUserPersistence User データに関するPersistence を生成
func NewUserPersistence() repository.UserRepository {
	return &userPersistence{}
}

func (up userPersistence) SelectUserByAuth(auth string) (*model.UserD, error) {
	// auth_tokenを条件にSELECTを行うSQLを第1引数に入力する
	row := db.Conn.QueryRow("SELECT * FROM user WHERE auth_token = ?", auth)
	return convertToUser(row)
}

// convertToUser rowデータをUserデータへ変換する
func convertToUser(row *sql.Row) (*model.UserD, error) {
	user := model.UserD{}
	err := row.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore, &user.Coin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}
	return &user, nil
}
