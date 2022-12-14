CREATE TABLE users
(
    id            serial primary key,
    chat_id       bigint unique,
    user_name     varchar(255),
    first_name    varchar(255) not null,
    last_name     varchar(255) not null,
    language_code varchar(3)
);

CREATE INDEX user_chat_id ON users (chat_id);

CREATE TYPE currency_types as enum ('classic', 'crypto');

CREATE TABLE currencies
(
    id        serial primary key,
    type      currency_types not null,
    code      varchar(255)   not null,
    bank_code varchar(255),
    UNIQUE (code, type, bank_code)
);


CREATE TABLE converter_pairs
(
    id                 serial primary key,
    level              int                                              not null,
    first_currency_id  int references currencies (id) on delete cascade not null,
    second_currency_id int references currencies (id) on delete cascade not null,
    third_currency_id  int references currencies (id) on delete cascade,
    UNIQUE (first_currency_id, second_currency_id, third_currency_id)
);

CREATE TABLE user_currencies
(
    id          serial primary key,
    user_id     int references users (id) on delete cascade      not null,
    currency_id int references currencies (id) on delete cascade not null,
    UNIQUE (user_id, currency_id)
);

CREATE TABLE user_converter_pairs
(
    id                serial primary key,
    user_id           int references users (id) on delete cascade           not null,
    converter_pair_id int references converter_pairs (id) on delete cascade not null,
    UNIQUE (user_id, converter_pair_id)
);

CREATE TABLE user_converter_pair_thresholds
(
    id                     serial primary key,
    user_id                int references users (id) on delete cascade                not null,
    user_converter_pair_id int references user_converter_pairs (id) on delete cascade not null,
    threshold              int                                                        not null,
    UNIQUE (user_id, user_converter_pair_id, threshold)
);