//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

package ranking

import (
	"fmt"
	"net/http"

	"20dojo-online/pkg/constant"
	model "20dojo-online/pkg/server/domain/model/user"
	ur "20dojo-online/pkg/server/domain/repository/user"

	"20dojo-online/pkg/server/interface/myerror"
)

// RankingUseCase UserにおけるUseCaseのインターフェース
type RankingUseCase interface {
	GetUsersByHighScore(startNum int32) ([]*model.UserL, *myerror.MyErr)
}

type rankingUseCase struct {
	userRepository ur.UserRepo
}

// NewRankingUseCase Userデータに関するUseCaseを生成
func NewRankingUseCase(ur ur.UserRepo) RankingUseCase {
	return &rankingUseCase{
		userRepository: ur,
	}
}

// GetUsersByHighScore Userデータを条件抽出
func (ru *rankingUseCase) GetUsersByHighScore(startNum int32) ([]*model.UserL, *myerror.MyErr) {
	// idと照合するユーザを取得
	userSlice, err := ru.userRepository.SelectUsersByHighScore(constant.RankingListNumber, startNum)
	if err != nil {
		myErr := myerror.NewMyErr(
			err,
			http.StatusInternalServerError,
		)
		return nil, myErr
	}
	// アイテム数以上の順位が開始位置に指定された場合の処理
	if len(userSlice) == 0 {
		myErr := myerror.NewMyErr(
			fmt.Errorf("user not found. rank=%d", startNum),
			http.StatusBadRequest,
		)
		return nil, myErr
	}
	return userSlice, nil
}
