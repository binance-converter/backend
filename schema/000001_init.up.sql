CREATE TABLE users
(
    id            serial primary key,
    chat_id       int,
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
    id              serial primary key,
    level           int                                              not null,
    first_currency  int references currencies (id) on delete cascade not null,
    second_currency int references currencies (id) on delete cascade not null,
    third_currency  int references currencies (id) on delete cascade
);

CREATE TABLE user_currencies
(
    id       serial primary key,
    user_id  int references users (id) on delete cascade      not null,
    currency int references currencies (id) on delete cascade not null
);

CREATE TABLE user_converter_pairs
(
    id             serial primary key,
    user_id        int references users (id) on delete cascade           not null,
    converter_pair int references converter_pairs (id) on delete cascade not null
);

CREATE TABLE user_converter_pair_thresholds
(
    id             serial primary key,
    user_id        int references users (id) on delete cascade           not null,
    converter_pair int references converter_pairs (id) on delete cascade not null,
    threshold      int                                                   not null
);