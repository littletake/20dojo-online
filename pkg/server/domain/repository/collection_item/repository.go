package collection_item

import "20dojo-online/pkg/server/domain/model"

// CItemRepository CollectionItem におけるRepository のインターフェース
type CItemRepository interface {
	SelectAllCollectionItem() ([]*model.CollectionItem, error)
}
