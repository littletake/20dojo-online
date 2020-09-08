package gachaprobability

import (
	"database/sql"

	"20dojo-online/pkg/db"
	model "20dojo-online/pkg/server/domain/model/gachaprobability"
	repository "20dojo-online/pkg/server/domain/repository/gachaprobability"
)

type gachaProbPersistence struct{}

// NewGachaProbPersistence Gachaprob データに関するPersistence を生成
func NewGachaProbPersistence() repository.GachaProbRepository {
	return &gachaProbPersistence{}
}

func (gp gachaProbPersistence) SelectAllGachaProb() ([]*model.GachaProb, error) {
	rows, err := db.Conn.Query("SELECT * FROM gacha_probability")
	if err != nil {
		return nil, err
	}
	return convertToGachaProb(rows)

}

// convertToGachaProb rowsデータをGachaProbデータへ変換する
func convertToGachaProb(rows *sql.Rows) ([]*model.GachaProb, error) {
	gachaProbSlice := []*model.GachaProb{}
	for rows.Next() {
		gachaProb := model.GachaProb{}
		err := rows.Scan(&gachaProb.CollectionItemID, &gachaProb.Ratio)
		if err != nil {
			return nil, err
		}
		gachaProbSlice = append(gachaProbSlice, &gachaProb)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return gachaProbSlice, nil
}
