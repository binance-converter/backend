package userDbPostgres

import (
	"github.com/binance-converter/backend/core"
	"github.com/jackc/pgconn"
	"golang.org/x/net/context"
)

func (u *UserDb) AddCurrency(ctx context.Context,
	currency core.FullCurrency) (int, error) {
	db := u.dbDriver
	tx, ok := u.transactionDB.ExtractTx(ctx)
	if ok {
		db = tx
	}

	query := `
				INSERT INTO currencies
				    (type, code, bank_code) 
				VALUES 
				    ($1, $2, $3) 
				RETURNING 
				    id
				`

	currencyType, err := u.convertCoreCurrencyTypeToPostgres(currency.CurrencyType)
	if err != nil {
		return 0, err
	}

	row := db.QueryRow(ctx, query, currencyType, currency.CurrencyCode, currency.BankCode)

	var currencyId int
	if err := row.Scan(&currencyId); err != nil {
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
				return 0, core.ErrorCurrencyAlreadyHas
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

	return currencyId, nil
}

func (u *UserDb) CheckCurrency(ctx context.Context,
	currency core.FullCurrency) (int, error) {
	db := u.dbDriver
	tx, ok := u.transactionDB.ExtractTx(ctx)
	if ok {
		db = tx
	}

	query := `	SELECT 
	    			id
	    		FROM 
	           		currencies
				WHERE
				    type = $1 AND
				    code = $2 AND
                    bank_code = $3`

	currencyType, err := u.convertCoreCurrencyTypeToPostgres(currency.CurrencyType)
	if err != nil {
		return 0, err
	}

	row := db.QueryRow(ctx, query, currencyType, currency.CurrencyCode, currency.BankCode)

	var rows int
	if err := row.Scan(&rows); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "no rows in result set":
				//logrus.WithFields(logrus.Fields{
				//	"base":   logBase,
				//	"userId": userId,
				//	"query":  query,
				//	"error":  err,
				//}).Error("user not found")
				return 0, core.ErrorCurrencyNotFound
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

	return rows, nil
}

func (u *UserDb) GetCurrency(ctx context.Context, currencyId int) (*core.FullCurrency, error) {
	db := u.dbDriver
	tx, ok := u.transactionDB.ExtractTx(ctx)
	if ok {
		db = tx
	}
	query := `	SELECT 
                	(type, code, bank_code)
                FROM 
                    currencies
                WHERE
                    id = $1`

	row := db.QueryRow(ctx, query, currencyId)

	var currency core.FullCurrency
	var currencyType string

	if err := row.Scan(&currencyType, &currency.CurrencyCode, &currency.CurrencyType); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "no rows in result set":
				return nil, core.ErrorCurrencyNotFound
			default:
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var err error
	currency.CurrencyType, err = u.convertPostgresCurrencyTypeToCore(currencyType)
	if err != nil {
		return nil, err
	}
	return &currency, nil
}

func (u *UserDb) AddCurrencyIfHasNot(ctx context.Context, currency core.FullCurrency) (int, error) {
	id, err := u.CheckCurrency(ctx, currency)
	if err != nil {
		switch err {
		case core.ErrorCurrencyNotFound:
			id, err = u.AddCurrency(ctx, currency)
			if err != nil {
				return 0, err
			}
			break
		default:
			return 0, err
		}
	}
	return id, nil
}

func (u *UserDb) AddUserCurrency(ctx context.Context, userId int,
	currency core.FullCurrency) (int, error) {
	//TODO: move to service
	currencyId, err := u.AddCurrencyIfHasNot(ctx, currency)
	if err != nil {
		return 0, err
	}

	db := u.dbDriver
	tx, ok := u.transactionDB.ExtractTx(ctx)
	if ok {
		db = tx
	}

	query := `
				INSERT INTO user_currencies
				    (user_id, currency_id) 
				VALUES 
				    ($1, $2)
				RETURNING
                    id
				`

	row := db.QueryRow(ctx, query, userId, currencyId)

	var userCurrencyId int
	if err := row.Scan(&userCurrencyId); err != nil {
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
				return 0, core.ErrorCurrencyAlreadyHas
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

	return userCurrencyId, nil
}

func (u *UserDb) GetUserCurrencies(ctx context.Context, userId int,
	currencyType *core.CurrencyType) ([]core.FullCurrency, error) {

	db := u.dbDriver
	tx, ok := u.transactionDB.ExtractTx(ctx)
	if ok {
		db = tx
	}

	query := `	SELECT
   					(type, code, bank_code)
				FROM
    				currencies
				WHERE
    				id IN (
    				SELECT
      		   			id
    			 	FROM
         				user_currencies
     				WHERE
         				user_id = $1)`

	var additionalArgs []interface{}

	if currencyType != nil {
		query += " AND type = $2"
		additionalArgs = append(additionalArgs, *currencyType)
	}

	rows, err := db.Query(ctx, query, userId, additionalArgs)
	if err != nil {
		return nil, err
	}

	var currencies []core.FullCurrency

	for rows.Next() {
		var currency core.FullCurrency
		var currencyCode string
		err = rows.Scan(&currencyCode, &currency.CurrencyCode, &currency.BankCode)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				switch pgErr.Code {
				default:
					//logrus.WithFields(logrus.Fields{
					//	"base":    logBase,
					//	"massage": pgErr.Message,
					//	"where":   pgErr.Where,
					//	"detail":  pgErr.Detail,
					//	"code":    pgErr.Code,
					//	"query":   logQuery(query),
					//}).Error("error scan rows")
					return nil, err
				}
			} else {
				//logrus.WithFields(logrus.Fields{
				//	"base":  logBase,
				//	"query": logQuery(query),
				//	"error": err,
				//}).Error("error scan rows")
				return nil, err
			}
		}
		currency.CurrencyType, err = u.convertPostgresCurrencyTypeToCore(currencyCode)
		if err == nil {
			currencies = append(currencies, currency)
		}
	}
	return currencies, nil
}

func (u *UserDb) DeleteUserCurrency(ctx context.Context, userId int,
	currency core.CurrencyCode) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserDb) convertCoreCurrencyTypeToPostgres(currencyType core.CurrencyType) (string, error) {
	switch currencyType {
	case core.CurrencyTypeCrypto:
		return "crypto", nil
	case core.CurrencyTypeClassic:
		return "classic", nil
	default:
		return "", core.ErrorCurrencyInvalidCurrencyType
	}
}

func (u *UserDb) convertPostgresCurrencyTypeToCore(currencyType string) (core.CurrencyType, error) {
	switch currencyType {
	case "crypto":
		return core.CurrencyTypeCrypto, nil
	case "classic":
		return core.CurrencyTypeClassic, nil
	default:
		return 0, core.ErrorCurrencyInvalidCurrencyType
	}
}
