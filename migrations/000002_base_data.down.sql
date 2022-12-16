DELETE FROM
           converter_pairs
       WHERE
           level = 3 AND
           first_currency_id IN (
               SELECT
                   id
               FROM
                   currencies
               WHERE
                       type = 'classic' AND
                       code = 'RUB' AND
                       bank_code = 'TinkoffNew') AND
           second_currency_id IN (
               SELECT
                   id
               FROM
                   currencies
               WHERE
                       type = 'crypto' AND
                       code = 'USDT') AND
           third_currency_id IN (
               SELECT
                   id
               FROM
                   currencies
               WHERE
                       type = 'classic' AND
                       code = 'KZT' AND
                       bank_code = 'KaspiBank');

DELETE FROM
    converter_pairs
WHERE
        level = 2 AND
        first_currency_id IN (
        SELECT
            id
        FROM
            currencies
        WHERE
                type = 'classic' AND
                code = 'RUB' AND
                bank_code = 'TinkoffNew') AND
        second_currency_id IN (
        SELECT
            id
        FROM
            currencies
        WHERE
                type = 'crypto' AND
                code = 'USDT');

DELETE FROM
    converter_pairs
WHERE
        level = 2 AND
        first_currency_id IN (
        SELECT
            id
        FROM
            currencies
        WHERE
                type = 'crypto' AND
                code = 'USDT') AND
        second_currency_id IN (
        SELECT
            id
        FROM
            currencies
        WHERE
                type = 'classic' AND
                code = 'KZT' AND
                bank_code = 'KaspiBank');

DELETE FROM
           currencies
       WHERE
           type = 'classic' AND
           code = 'RUB' AND
           bank_code = 'TinkoffNew';

DELETE FROM
           currencies
       WHERE
           type = 'crypto' AND
           code = 'USDT';

DELETE FROM
           currencies
       WHERE
           type = 'classic' AND
           code = 'KZT' AND
           bank_code = 'KaspiBank';