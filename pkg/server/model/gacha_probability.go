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

// GachaProbSlice GachaProbの複数形
type GachaProbSlice []*GachaProb

// SelectAllGachaProb table:Gacha_probabilityの全件取得
func SelectAllGachaProb() (GachaProbSlice, error) {
	rows, err := db.Conn.Query("SELECT * FROM gacha_probability")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return convertToGachaProbSlice(rows)
}

// convertToGachaProbSlice rowsデータをGachaProbSliceデータへ変換する
func convertToGachaProbSlice(rows *sql.Rows) (GachaProbSlice, error) {
	gachaProbSlice := GachaProbSlice{}
	for rows.Next() {
		gachaProb := GachaProb{}
		err := rows.Scan(&gachaProb.CollectionItemID, &gachaProb.Ratio)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		gachaProbSlice = append(gachaProbSlice, &gachaProb)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return gachaProbSlice, nil
}
