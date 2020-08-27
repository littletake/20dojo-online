package model

import (
	"database/sql"
	"log"

	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/db"
)

// User userテーブルデータ
type User struct {
	ID        string
	AuthToken string
	Name      string
	HighScore int32
	Coin      int32
}

// Users Userの複数形
type Users []*User

// TODO: DBアクセス後にエラーハンドリングを追加する

// InsertUser データベースにレコードを登録する
func InsertUser(record *User) error {
	// userテーブルへのレコードの登録を行うSQLを入力する
	stmt, err := db.Conn.Prepare("INSERT INTO user (id, auth_token, name, high_score, coin) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.ID, record.AuthToken, record.Name, record.HighScore, record.Coin)
	return err
}

// SelectUserByAuthToken auth_tokenを条件にレコードを取得する
func SelectUserByAuthToken(authToken string) (*User, error) {
	// auth_tokenを条件にSELECTを行うSQLを第1引数に入力する
	row := db.Conn.QueryRow("SELECT * FROM user WHERE auth_token = ?", authToken)
	return convertToUser(row)
}

// SelectUserByPrimaryKey 主キーを条件にレコードを取得する
func SelectUserByPrimaryKey(userID string) (*User, error) {
	// idを条件にSELECTを行うSQLを第1引数に入力する
	row := db.Conn.QueryRow("SELECT * FROM user WHERE id = ?", userID)
	return convertToUser(row)
}

// SelectUsersByHighScore high_scoreを基準に降順にしたレコードを取得
func SelectUsersByHighScore(start int) (Users, error) {
	// 任意の順位(start)からRankingListNumber分取得
	rows, err := db.Conn.Query("SELECT * FROM user ORDER BY high_score DESC LIMIT ? OFFSET ?", constant.RankingListNumber, start-1)
	// TODO: 意味をしっかり理解してエラー処理を書く
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}
	return convertToUsers(rows)
}

// UpdateUserByPrimaryKey 主キーを条件にレコードを更新する
func UpdateUserByPrimaryKey(record *User) error {
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
func convertToUser(row *sql.Row) (*User, error) {
	user := User{}
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

// convertToUsers rowsデータをUsersデータへ変換する
func convertToUsers(rows *sql.Rows) (Users, error) {
	users := Users{}
	for rows.Next() {
		user := User{}
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
	return users, nil
}
