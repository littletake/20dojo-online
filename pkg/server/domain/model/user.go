package model

// UserL ユーザ情報
type UserL struct {
	ID        string
	AuthToken string
	Name      string
	HighScore int32
	Coin      int32
}
