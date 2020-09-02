package model

// UserL ユーザ情報
type UserL struct {
	ID        string
	AuthToken string
	Name      string
	HighScore int32
	Coin      int32
}

// MyErr エラーコードを追加した独自の型
type MyErr struct {
	// TODO: pkg/errorsを使いたい
	ErrMsg  error
	ErrCode int32
}
