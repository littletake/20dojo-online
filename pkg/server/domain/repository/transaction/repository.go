//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

package transaction

import "database/sql"

// TxRepository TxにおけるRepository のインターフェース
type TxRepository interface {
	Transaction(f func(tx *sql.Tx) error) error
}
