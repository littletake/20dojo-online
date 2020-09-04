package usecase

import (
	"fmt"

	"20dojo-online/pkg/server/domain/model"
	"20dojo-online/pkg/server/domain/repository"
	"20dojo-online/pkg/server/interface/myerror"
)

// CollectionUseCase UseCaseのインターフェース
type CollectionUseCase interface {
	GetCollectionSlice(userID string) ([]*model.CollectionItemResult, *myerror.MyErr)
}

type collectionUseCase struct {
	userRepository   repository.UserRepository
	cItemRepository  repository.CItemRepository
	ucItemRepository repository.UCItemRepository
}

// NewCollectionUseCase Userデータに関するUseCaseを生成
func NewCollectionUseCase(ur repository.UserRepository, cr repository.CItemRepository,
	ucr repository.UCItemRepository) CollectionUseCase {
	return &collectionUseCase{
		userRepository:   ur,
		cItemRepository:  cr,
		ucItemRepository: ucr,
	}
}

// collectionItem コレクションアイテム一覧
type collectionItem struct {
	CollectionID string `json:"collectionID"`
	ItemName     string `json:"name"`
	Rarity       int32  `json:"rarity"`
	HasItem      bool   `json:"hasItem"`
}

// GetUsersByHighScore Userデータを条件抽出
func (cu collectionUseCase) GetCollectionSlice(userID string) ([]*model.CollectionItemResult, *myerror.MyErr) {
	// ユーザデータの取得処理と存在チェックを実装
	user, err := cu.userRepository.SelectUserByUserID(userID)
	if err != nil {
		myErr := myerror.MyErr{err, 500}
		return nil, &myErr
	}
	if user == nil {
		myErr := myerror.MyErr{
			fmt.Errorf("user not found"),
			500,
		}
		return nil, &myErr
	}
	// 現ユーザが保持しているアイテムの情報をまとめる -> hasGotItemMap
	// table: user_collection_itemに対してuserIDのものを取得
	ucItemSlice, err := cu.ucItemRepository.SelectUCItemSliceByUserID(userID)
	if err != nil {
		myErr := myerror.MyErr{err, 500}
		return nil, &myErr
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
			myErr := myerror.MyErr{err, 500}
			return nil, &myErr
		}
		hasGotcItemSlice = true
	}
	// 二つのtableを合わせてresponseを作成
	cItemResult := make([]*model.CollectionItemResult, len(cItemSlice))
	for i, cItem := range cItemSlice {
		result := model.CollectionItemResult{
			CollectionID: cItem.ItemID,
			ItemName:     cItem.ItemName,
			Rarity:       cItem.Rarity,
			HasItem:      hasGotItemMap[cItem.ItemID],
		}
		cItemResult[i] = &result
	}
	return cItemResult, nil
}
