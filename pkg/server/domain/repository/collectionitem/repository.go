//go:generate mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

package collectionitem

import model "20dojo-online/pkg/server/domain/model/collectionitem"

// CollectionItemRepo CollectionItem におけるRepository のインターフェース
type CollectionItemRepo interface {
	SelectAllCollectionItem() ([]*model.CollectionItem, error)
}
