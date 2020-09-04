package repository

import "20dojo-online/pkg/server/domain/model"

// GachaProbRepository GachaProb におけるRepository のインターフェース
type GachaProbRepository interface {
	SelectAllGachaProb() ([]*model.GachaProb, error)
}
