package transaction

import "database/sql"

// TxRepository TxにおけるRepository のインターフェース
type TxRepository interface {
	Transaction(f func(tx *sql.Tx) error) error
}
