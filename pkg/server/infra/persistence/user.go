// Package persistence 技術的関心事を扱う
package persistence

import (
	"database/sql"
	"log"

	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/db"
	"20dojo-online/pkg/server/domain/model"
	"20dojo-online/pkg/server/domain/repository"
)

type userPersistence struct{}

// NewUserPersistence User データに関するPersistence を生成
func NewUserPersistence() repository.UserRepository {
	return &userPersistence{}
}

func (up userPersistence) SelectUserByUserID(id string) (user *model.UserL, err error) {
	// auth_tokenを条件にSELECTを行うSQLを第1引数に入力する
	row := db.Conn.QueryRow("SELECT * FROM user WHERE id = ?", id)
	return convertToUser(row)
}

// SelectUsersByHighScore high_scoreを基準に降順にしたレコードを取得
func (up userPersistence) SelectUsersByHighScore(start int32) (user []*model.UserL, err error) {
	// 任意の順位(start)からRankingListNumber分取得
	rows, err := db.Conn.Query("SELECT * FROM user ORDER BY high_score DESC LIMIT ? OFFSET ?", constant.RankingListNumber, start-1)
	// TODO: 意味をしっかり理解してエラー処理を書く
	if err != nil {
		return nil, err
	}
	return convertToUsers(rows)
}

// InsertUser データベースにレコードを登録する
func (up userPersistence) InsertUser(record *model.UserL) (err error) {
	// userテーブルへのレコードの登録を行うSQLを入力する
	stmt, err := db.Conn.Prepare("INSERT INTO user (id, auth_token, name, high_score, coin) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.ID, record.AuthToken, record.Name, record.HighScore, record.Coin)
	return err
}

// UpdateUserByUser
func (up userPersistence) UpdateUserByUser(record *model.UserL) (err error) {
	// idを条件に指定した値で以下の値を更新するSQLを入力する
	// 更新カラム: name, coin, high_score
	stmt, err := db.Conn.Prepare("UPDATE user SET name = ?, coin = ?, high_score = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.Name, record.Coin, record.HighScore, record.ID)
	return err
}

// convertToUser rowデータをUserデータへ変換する
func convertToUser(row *sql.Row) (*model.UserL, error) {
	user := model.UserL{}
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

// convertToUsers rowsデータをUser型のslice（ポインタ）へ変換する
func convertToUsers(rows *sql.Rows) ([]*model.UserL, error) {
	users := []*model.UserL{}
	for rows.Next() {
		user := model.UserL{}
		err := rows.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore, &user.Coin)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			log.Println(err)
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
}
