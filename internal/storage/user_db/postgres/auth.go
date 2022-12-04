package userDbPostgres

import (
	"github.com/binance-converter/backend/core"
	"github.com/jackc/pgconn"
	"golang.org/x/net/context"
)

func (u *UserDb) AddUser(ctx context.Context, user core.AddUser) (int, error) {
	db := u.dbDriver
	tx, ok := u.transactionDB.ExtractTx(ctx)
	if ok {
		db = tx
	}

	query := `	INSERT INTO users
				    (chat_id, user_name, first_name, last_name, language_code) 
				VALUES 
				    ($1, $2, $3, $4, $5) 
				RETURNING 
					id`

	row := db.QueryRow(ctx, query, user.ChatId, user.UserName, user.FirstName, user.LastName,
		user.LanguageCode)

	var userId int
	if err := row.Scan(&userId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505":
				//logrus.WithFields(logrus.Fields{
				//	"base":    logBase,
				//	"email":   user.Email,
				//	"massage": pgErr.Message,
				//	"where":   pgErr.Where,
				//	"detail":  pgErr.Detail,
				//	"code":    pgErr.Code,
				//	"query":   query,
				//}).Error("user already has")
				return 0, core.ErrorAuthServiceAuthUserAlreadyExists
			default:
				//logrus.WithFields(logrus.Fields{
				//	"base":    logBase,
				//	"email":   user.Email,
				//	"massage": pgErr.Message,
				//	"where":   pgErr.Where,
				//	"detail":  pgErr.Detail,
				//	"code":    pgErr.Code,
				//	"query":   query,
				//}).Error("error add user to postgres")
				return 0, err
			}
		} else {
			//logrus.WithFields(logrus.Fields{
			//	"base":  logBase,
			//	"query": query,
			//	"error": err,
			//}).Error("error add user to postgres")
			return 0, err
		}
	}

	return userId, nil
}
