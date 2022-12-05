package userDbPostgres

import (
	"github.com/binance-converter/backend/core"
	"github.com/jackc/pgconn"
	"golang.org/x/net/context"
)

func (u *UserDb) AddConverterPair(ctx context.Context, converterPair core.ConverterPair) (int,
	error) {

	if len(converterPair.Currencies) != 2 && len(converterPair.Currencies) != 3 {
		return 0, core.ErrorConverterInvalidConverterPair
	}

	db := u.dbDriver
	tx, ok := u.transactionDB.ExtractTx(ctx)
	if ok {
		db = tx
	}

	query := `	INSERT INTO
        			converter_pairs 
	    			(level, first_currency_id, second_currency_id, third_currency_id)
        		VALUES 
            		($1, $2, $3, $4)
				RETURNING 
					id`

	var additionalArgs []interface{}

	for _, currency := range converterPair.Currencies {
		currencyId, err := u.AddCurrencyIfHasNot(ctx, currency)
		if err != nil {
			return 0, err
		}
		additionalArgs = append(additionalArgs, currencyId)
	}

	if len(additionalArgs) == 2 {
		additionalArgs = append(additionalArgs, nil)
	}

	row := db.QueryRow(ctx, query, len(additionalArgs), additionalArgs)

	var converterPairId int
	if err := row.Scan(&converterPairId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505":
				return 0, core.ErrorConverterConverterPairAlreadyExists
			default:
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	return converterPairId, nil
}

func (u *UserDb) CheckConverterPair(ctx context.Context, converterPair core.ConverterPair) (int,
	error) {

	if len(converterPair.Currencies) != 2 && len(converterPair.Currencies) != 3 {
		return 0, core.ErrorConverterInvalidConverterPair
	}

	db := u.dbDriver
	tx, ok := u.transactionDB.ExtractTx(ctx)
	if ok {
		db = tx
	}
	query := `	SELECT
                    id
                FROM
                    converter_pairs
                WHERE
                    level = $1 AND
                    first_currency_id = $2 AND
                    second_currency_id = $3 AND
                    third_currency_id = $4`

	var additionalArgs []interface{}
	for _, currency := range converterPair.Currencies {
		currencyId, err := u.CheckCurrency(ctx, currency)
		if err != nil {
			return 0, core.ErrorConverterInvalidConverterPair
		}
		additionalArgs = append(additionalArgs, currencyId)
	}
	if len(converterPair.Currencies) == 2 {
		additionalArgs = append(additionalArgs, nil)
	}
	row := db.QueryRow(ctx, query, additionalArgs)
	var converterPairId int
	if err := row.Scan(&converterPairId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "no rows in result set":
				return 0, core.ErrorConverterConverterPairNotFound
			default:
				return 0, err
			}
		} else {
			return 0, err
		}
	}
	return converterPairId, nil
}

func (u *UserDb) AddConverterPairIfHasNot(ctx context.Context,
	converterPair core.ConverterPair) (int, error) {
	id, err := u.CheckConverterPair(ctx, converterPair)
	if err != nil {
		switch err {
		case core.ErrorConverterConverterPairNotFound:
			id, err = u.AddConverterPair(ctx, converterPair)
			if err != nil {
				return 0, err
			}
			break
		case core.ErrorCurrencyNotFound:
			id, err = u.AddConverterPair(ctx, converterPair)
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

func (u *UserDb) SetUserConverterPair(ctx context.Context, userId int,
	converterPair core.ConverterPair) (int, error) {
	//TODO: move to service
	converterPairId, err := u.AddConverterPairIfHasNot(ctx, converterPair)
	if err != nil {
		return 0, err
	}

	db := u.dbDriver
	tx, ok := u.transactionDB.ExtractTx(ctx)
	if ok {
		db = tx
	}

	query := `	INSERT INTO 
	    			user_converter_pairs
					(user_id, converter_pair_id)
        		VALUES
        		    ($1, $2)
				RETURNING 
					id`

	row := db.QueryRow(ctx, query, userId, converterPairId)

	var userConverterPairId int
	if err := row.Scan(&userConverterPairId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505":
				return 0, core.ErrorConverterConverterPairAlreadyExists
			default:
				return 0, err
			}
		} else {
			return 0, err
		}
	}
	return userConverterPairId, nil
}

func (u *UserDb) GetUserConverterPairs(ctx context.Context, userId int) ([]core.ConverterPair,
	error) {

	db := u.dbDriver
	tx, ok := u.transactionDB.ExtractTx(ctx)
	if ok {
		db = tx
	}

	query := `	SELECT
	    			(level, first_currency_id, second_currency_id, third_currency_id)
				FROM
				    converter_pairs
                WHERE
                	id IN
                		(SELECT 
                		     converter_pair_id
                		 FROM
                		     user_converter_pairs
                		 WHERE 
                		     user_id = $1)`

	rows, err := db.Query(ctx, query, userId)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			default:
				return nil, core.ErrorConverterConverterPairNotFound
			}
		}
	}

	var converterPairs []core.ConverterPair

	for rows.Next() {
		var level, firstId, secondId, thirdId int
		if err := rows.Scan(&level, &firstId, &secondId, &thirdId); err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				switch pgErr.Code {
				default:
					return nil, err
				}
			}
		}
		var converterPair core.ConverterPair
		firstCurrency, err := u.GetCurrency(ctx, firstId)
		if err != nil {
			return nil, err
		}
		converterPair.Currencies = append(converterPair.Currencies, *firstCurrency)

		secondCurrency, err := u.GetCurrency(ctx, secondId)
		if err != nil {
			return nil, err
		}
		converterPair.Currencies = append(converterPair.Currencies, *secondCurrency)

		if level == 3 {
			thirdCurrency, err := u.GetCurrency(ctx, thirdId)
			if err != nil {
				return nil, err
			}
			converterPair.Currencies = append(converterPair.Currencies, *thirdCurrency)
		}
	}
	return converterPairs, nil
}

func (u *UserDb) SetThresholdConvertPair(ctx context.Context, userId int,
	threshold core.ThresholdConvertPair) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserDb) GetThresholdConvertPair(ctx context.Context,
	userId int) ([]core.ThresholdConvertPair, error) {
	//TODO implement me
	panic("implement me")
}
