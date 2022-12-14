package userDbPostgres

import (
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

type dbDriverUserDB interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type transactionDBUserDB interface {
	ExtractTx(ctx context.Context) (pgx.Tx, bool)
}

type UserDb struct {
	dbDriver      dbDriverUserDB
	transactionDB transactionDBUserDB
}

func NewUserDB(dbDriver dbDriverUserDB, transactionDB transactionDBUserDB) *UserDb {
	return &UserDb{
		dbDriver:      dbDriver,
		transactionDB: transactionDB,
	}
}
