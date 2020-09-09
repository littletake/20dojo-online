//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

package collectionitem

import model "20dojo-online/pkg/server/domain/model/collectionitem"

// CItemRepository CollectionItem におけるRepository のインターフェース
type CItemRepository interface {
	SelectAllCollectionItem() ([]*model.CollectionItem, error)
}
