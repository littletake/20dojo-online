package transaction

import "database/sql"

// TxRepo TxにおけるRepository のインターフェース
type TxRepo interface {
	Transaction(f func(any interface{}, tx *sql.Tx) error) error
}
