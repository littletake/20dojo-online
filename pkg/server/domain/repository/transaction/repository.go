//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

package transaction

import "database/sql"

// TxRepo TxにおけるRepository のインターフェース
type TxRepo interface {
	Transaction(f func(tx *sql.Tx) error) error
}
