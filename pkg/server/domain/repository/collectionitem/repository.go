package collectionitem

import model "20dojo-online/pkg/server/domain/model/collectionitem"

// CollectionItemRepo CollectionItem におけるRepository のインターフェース
type CollectionItemRepo interface {
	SelectAllCollectionItem() ([]*model.CollectionItem, error)
}
