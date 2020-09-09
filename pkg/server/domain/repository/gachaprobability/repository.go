package gachaprobability

import model "20dojo-online/pkg/server/domain/model/gachaprobability"

// GachaProbRepository GachaProb におけるRepository のインターフェース
type GachaProbRepository interface {
	SelectAllGachaProb() ([]*model.GachaProb, error)
}
