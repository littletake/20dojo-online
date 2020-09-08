package transaction

import "database/sql"

// TxRepository TxにおけるRepository のインターフェース
type TxRepository interface {
	Transaction(f func(any interface{}, tx *sql.Tx) error) error
}
