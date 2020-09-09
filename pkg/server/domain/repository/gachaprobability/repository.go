package gachaprobability

import model "20dojo-online/pkg/server/domain/model/gachaprobability"

// GachaProbRepo GachaProb におけるRepository のインターフェース
type GachaProbRepo interface {
	SelectAllGachaProb() ([]*model.GachaProb, error)
}
