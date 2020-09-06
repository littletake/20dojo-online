package usecase

import (
	"fmt"

	"20dojo-online/pkg/server/domain/model"
	"20dojo-online/pkg/server/domain/repository"
	"20dojo-online/pkg/server/interface/myerror"
)

// RankingUseCase UserにおけるUseCaseのインターフェース
type RankingUseCase interface {
	GetUsersByHighScore(startNum int32) ([]*model.UserL, *myerror.MyErr)
}

type rankingUseCase struct {
	userRepository repository.UserRepository
}

// NewRankingUseCase Userデータに関するUseCaseを生成
func NewRankingUseCase(ur repository.UserRepository) RankingUseCase {
	return &rankingUseCase{
		userRepository: ur,
	}
}

// GetUsersByHighScore Userデータを条件抽出
func (ru rankingUseCase) GetUsersByHighScore(startNum int32) ([]*model.UserL, *myerror.MyErr) {
	// idと照合するユーザを取得
	userSlice, err := ru.userRepository.SelectUsersByHighScore(startNum)
	if err != nil {
		myErr := myerror.NewMyErr(err, 500)
		return nil, myErr
	}
	// TODO: 順位範囲外の処理
	if len(userSlice) == 0 {
		myErr := myerror.NewMyErr(
			fmt.Errorf("user not found. rank=%d", startNum),
			400,
		)
		return nil, myErr
	}
	return userSlice, nil
}
