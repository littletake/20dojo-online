package model

// // User userテーブルデータ
// type User struct {
// 	ID        string
// 	AuthToken string
// 	Name      string
// 	HighScore int32
// 	Coin      int32
// }

// Users Userの複数形
// type Users []*User

// // GachaProb GachaProbabilityのデータ
// type GachaProb struct {
// 	ItemID int32
// 	Ratio  int32
// }

// // GachaProbList GachaProbの複数形
// type GachaProbList []*GachaProb

// // GachaItem Gachaアイテムの情報
// type GachaItem struct {
// 	ItemID   int32
// 	ItemName string
// 	Rarity   int32
// }

// TODO: DBアクセス後にエラーハンドリングを追加する

// InsertUser データベースにレコードを登録する
// func InsertUser(record *User) error {
// 	// userテーブルへのレコードの登録を行うSQLを入力する
// 	stmt, err := db.Conn.Prepare("INSERT INTO user (id, auth_token, name, high_score, coin) VALUES (?, ?, ?, ?, ?)")
// 	if err != nil {
// 		return err
// 	}
// 	_, err = stmt.Exec(record.ID, record.AuthToken, record.Name, record.HighScore, record.Coin)
// 	return err
// }

// SelectItemByItemID ガチャを一回回す
// TODO: itemIDはstring?int?
// func SelectItemByItemID(itemID string) (*Users, error) {
// // auth_tokenを条件にSELECTを行うSQLを第1引数に入力する
// row := db.Conn.QueryRow("SELECT * FROM gacha_probability WHERE auth_token = ?", authToken)
// return convertToUser(row)
// }

// // SelectAllGachaProb table:Gacha_probabilityの全件取得
// // TODO: 未確認
// func SelectAllGachaProb() (GachaProbList, error) {
// 	rows, err := db.Conn.Query("SELECT * FROM gacha_probability")
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil
// 		}
// 		log.Println(err)
// 		return nil, err
// 	}
// 	return convertToList(rows)
// }

// // SelectGachaItemByItemID IDを条件にレコードを取得する
// func SelectGachaItemByItemID(itemID int32) (*GachaItem, error) {
// 	// idを条件にSELECTを行うSQLを第1引数に入力する
// 	row := db.Conn.QueryRow("SELECT * FROM collection_item WHERE id = ?", itemID)
// 	return convertToItem(row)
// }

// TODO: 関数を分けるか統一するかを考える

// // SelectMyItemByItemID itemIDとuserIDを条件にアイテムが既出か否かを判定する
// func SelectMyItemByItemID(itemID int32, userID int32) {
// 	row := db.Conn.QueryRow("SELECT * FROM user_collection_item WHERE user_id = ? AND collection_item_id = ?", userID, itemID)
// 	if row == nil {
// 		// 登録

// 	} else {

// 	}
// }

// // InsertUser データベースにレコードを登録する
// func InsertItem(record *User) error {
// 	// userテーブルへのレコードの登録を行うSQLを入力する
// 	stmt, err := db.Conn.Prepare("INSERT INTO user (id, auth_token, name, high_score, coin) VALUES (?, ?, ?, ?, ?)")
// 	if err != nil {
// 		return err
// 	}
// // 	_, err = stmt.Exec(record.ID, record.AuthToken, record.Name, record.HighScore, record.Coin)
// // 	return err
// // }

// // SelectUsersByHighScore high_scoreを基準に降順にしたレコードを取得
// func SelectUsersByHighScore(start int) (Users, error) {
// 	// 任意の順位(start)からRankingListNumber分取得
// rows, err := db.Conn.Query("SELECT * FROM user ORDER BY high_score DESC LIMIT ? OFFSET ?", constant.RankingListNumber, start-1)
// 	// TODO: 意味をしっかり理解してエラー処理を書く
// if err != nil {
// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}
// 	log.Println(err)
// 	return nil, err
// }
// 	return convertToUsers(rows)
// }

// // UpdateUserByPrimaryKey 主キーを条件にレコードを更新する
// func UpdateUserByPrimaryKey(record *User) error {
// 	// idを条件に指定した値で以下の値を更新するSQLを入力する
// 	// 更新カラム: name, coin, high_score
// 	stmt, err := db.Conn.Prepare("UPDATE user SET name = ?, coin = ?, high_score = ? WHERE id = ?")
// 	if err != nil {
// 		return err
// 	}
// 	_, err = stmt.Exec(record.Name, record.Coin, record.HighScore, record.ID)
// 	return err
// }

// // convertToItem rowデータをGachaItemデータへ変換する
// func convertToItem(row *sql.Row) (*GachaItem, error) {
// 	gachaItem := GachaItem{}
// 	err := row.Scan(&gachaItem.ItemID, &gachaItem.ItemName, &gachaItem.Rarity)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil
// 		}
// 		log.Println(err)
// 		return nil, err
// 	}
// 	return &gachaItem, nil
// }

// // convertToList rowsデータをGachaProbListデータへ変換する
// func convertToList(rows *sql.Rows) (GachaProbList, error) {
// 	gachaProbList := GachaProbList{}
// 	for rows.Next() {
// 		gachaProb := GachaProb{}
// 		err := rows.Scan(&gachaProb.ItemID, &gachaProb.Ratio)
// 		if err != nil {
// 			log.Println(err)
// 			return nil, err
// 		}
// 		gachaProbList = append(gachaProbList, &gachaProb)
// 	}
// 	if err := rows.Err(); err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}
// 	return gachaProbList, nil
// }
