INSERT INTO currencies (type, code, bank_code) VALUES ('classic', 'RUB', 'TinkoffNew');
INSERT INTO currencies (type, code, bank_code) VALUES ('crypto', 'USDT', '');
INSERT INTO currencies (type, code, bank_code) VALUES ('classic', 'KZT', 'KaspiBank');

INSERT
INTO
    converter_pairs
(level, first_currency_id, second_currency_id, third_currency_id)
VALUES
    (3, 1, 2, 3);

INSERT
INTO
    converter_pairs
(level, first_currency_id, second_currency_id)
VALUES
    (2, 1, 2);

INSERT
INTO
    converter_pairs
(level, first_currency_id, second_currency_id)
VALUES
    (2, 2, 3);