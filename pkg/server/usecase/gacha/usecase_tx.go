package gacha

// import (
// 	"context"
// 	"database/sql"
// )

// type txGachaUseCase struct {
// 	db *sql.DB
// }

// // データベースのオブジェクトをもらって InputPort実装を返却する。
// func NewGachaUseCaseTx(db *sql.DB) GachaUseCase {
// 	return &txGachaUseCase{
// 		db: db,
// 	}
// }

// func (u *txGachaUseCase) BulkInsertAndUpdate(newItemSlice []*model.UserCollectionItem, user *model.UserL, tx *sql.Tx) error {
// 	// トランザクションを開始して本当のユースケースの実装作って呼び出すだけのうすーーーい処理。
// 	v, err := database.DoInTx(u.db, func(tx *sqlx.Tx) (interface{}, error) {
// 		ar := database.NewAccount(tx)
// 		return NewAccountInteractor(ar, dr).Store(ctx, in)
// 	})
// 	err := u.
// 	return v.(*entity.Account), err
// }
