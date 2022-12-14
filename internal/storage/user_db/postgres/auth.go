package userDbPostgres

import (
	"github.com/binance-converter/backend/core"
	"github.com/jackc/pgconn"
	"github.com/sirupsen/logrus"
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
				logrus.WithFields(logrus.Fields{
					"massage": pgErr.Message,
					"where":   pgErr.Where,
					"detail":  pgErr.Detail,
					"code":    pgErr.Code,
					"query":   query,
					"user":    user,
				}).Error("user already has")
				return 0, core.ErrorAuthServiceAuthUserAlreadyExists
			default:
				logrus.WithFields(logrus.Fields{
					"massage": pgErr.Message,
					"where":   pgErr.Where,
					"detail":  pgErr.Detail,
					"code":    pgErr.Code,
					"query":   query,
					"user":    user,
				}).Error("error add user to postgres")
				return 0, err
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"query": query,
				"error": err,
				"user":  user,
			}).Error("error add user to postgres")
			return 0, err
		}
	}

	return userId, nil
}

func (u *UserDb) ValidateUser(ctx context.Context, chatId int) (int, error) {
	db := u.dbDriver
	tx, ok := u.transactionDB.ExtractTx(ctx)
	if ok {
		db = tx
	}

	query := `	SELECT 
					id
                FROM 
                    users
                WHERE 
                    chat_id = $1`

	row := db.QueryRow(ctx, query, chatId)

	var userId int
	if err := row.Scan(&userId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "no rows in result set":
				return 0, core.ErrorAuthServiceUserNotFound
			default:
				return 0, err
			}
		} else {
			return 0, err
		}
	}
	return userId, nil
}
