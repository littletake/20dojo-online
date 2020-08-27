package model

import (
	"database/sql"
	"log"

	"20dojo-online/pkg/db"
)

// GachaProb GachaProbabilityのデータ
type GachaProb struct {
	CollectionItemID string
	Ratio            int32
}

// GachaProbList GachaProbの複数形
type GachaProbList []*GachaProb

// SelectAllGachaProb table:Gacha_probabilityの全件取得
// TODO: 未確認
func SelectAllGachaProb() (GachaProbList, error) {
	rows, err := db.Conn.Query("SELECT * FROM gacha_probability")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}
	return convertToGachaProbList(rows)
}

// convertToGachaProbList rowsデータをGachaProbListデータへ変換する
func convertToGachaProbList(rows *sql.Rows) (GachaProbList, error) {
	gachaProbList := GachaProbList{}
	for rows.Next() {
		gachaProb := GachaProb{}
		err := rows.Scan(&gachaProb.CollectionItemID, &gachaProb.Ratio)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		gachaProbList = append(gachaProbList, &gachaProb)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return gachaProbList, nil
}
