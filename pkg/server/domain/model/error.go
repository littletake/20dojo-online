package model

// MyErr エラーコードを追加した独自の型
type MyErr struct {
	// TODO: pkg/errorsを使いたい
	ErrMsg  error
	ErrCode int32
}
