//go:generate mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

package collection

import (
	"fmt"
	"net/http"

	model "20dojo-online/pkg/server/domain/model/collectionitem"
	cr "20dojo-online/pkg/server/domain/repository/collectionitem"
	ur "20dojo-online/pkg/server/domain/repository/user"
	ucr "20dojo-online/pkg/server/domain/repository/usercollectionitem"

	"20dojo-online/pkg/server/interface/myerror"
)

// CollectionUseCase UseCaseのインターフェース
type CollectionUseCase interface {
	GetCollectionSlice(userID string) ([]*CollectionItemResult, *myerror.MyErr)
}

type collectionUseCase struct {
	userRepository   ur.UserRepo
	cItemRepository  cr.CollectionItemRepo
	ucItemRepository ucr.UserCollectionItemRepo
}

// NewCollectionUseCase UseCaseを生成
func NewCollectionUseCase(ur ur.UserRepo, cr cr.CollectionItemRepo,
	ucr ucr.UserCollectionItemRepo) CollectionUseCase {
	return &collectionUseCase{
		userRepository:   ur,
		cItemRepository:  cr,
		ucItemRepository: ucr,
	}
}

// CollectionItemResult レスポンス用の構造体
type CollectionItemResult struct {
	CollectionID string `json:"collectionID"`
	ItemName     string `json:"name"`
	Rarity       int32  `json:"rarity"`
	HasItem      bool   `json:"hasItem"`
}

// cItemSlice collectionItemのスライス
var cItemSlice []*model.CollectionItem

// hasGotcItemSlice table:collection_itemの取得状況
var hasGotcItemSlice bool

// GetUsersByHighScore Userデータを条件抽出
func (cu *collectionUseCase) GetCollectionSlice(userID string) ([]*CollectionItemResult, *myerror.MyErr) {
	// ユーザデータの取得処理と存在チェックを実装
	user, err := cu.userRepository.SelectUserByUserID(userID)
	if err != nil {
		myErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)
		return nil, myErr
	}
	if user == nil {
		myErr := myerror.NewMyErr(
			fmt.Errorf("user not found. userID=%s", userID),
			http.StatusBadRequest,
		)
		return nil, myErr
	}
	// 現ユーザが保持しているアイテムの情報をまとめる -> hasGotItemMap
	// table: user_collection_itemに対してuserIDのものを取得
	ucItemSlice, err := cu.ucItemRepository.SelectSliceByUserID(userID)
	if err != nil {
		myErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)
		return nil, myErr
	}
	// hasGotItemMap 既出アイテム一覧map
	// [注意] ガチャ実行時も追加するので可変長指定
	hasGotItemMap := make(map[string]bool)
	for _, ucItem := range ucItemSlice {
		itemID := ucItem.CollectionItemID
		hasGotItemMap[itemID] = true
	}
	// 全アイテムの情報をまとめる -> cItemSlice
	// hasGotItemSlice 一度だけ実行するように制御
	if !hasGotcItemSlice {
		cItemSlice, err = cu.cItemRepository.SelectAllCollectionItem()
		if err != nil {
			myErr := myerror.NewMyErr(
				err,
				http.StatusInternalServerError,
			)
			return nil, myErr
		}
		hasGotcItemSlice = true
	}
	// 二つのtableを合わせてresponseを作成
	cItemResult := make([]*CollectionItemResult, len(cItemSlice))
	for i, cItem := range cItemSlice {
		result := CollectionItemResult{
			CollectionID: cItem.ItemID,
			ItemName:     cItem.ItemName,
			Rarity:       cItem.Rarity,
			HasItem:      hasGotItemMap[cItem.ItemID],
		}
		cItemResult[i] = &result
	}
	return cItemResult, nil
}
