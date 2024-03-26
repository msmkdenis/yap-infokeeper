begin transaction;

create schema if not exists infokeeper;

create table if not exists infokeeper.user
(
    id                      text,
    login                   text unique not null,
    password                bytea not null,
    created_at              timestamp not null default now(),
    constraint pk_user primary key (id)
);

create table if not exists infokeeper.credit_card
(
    id                      text,
    number                  text unique not null,
    owner_id                text not null,
    expires_at              timestamp not null,
    cvv_code                text not null,
    pin_code                text not null,
    created_at              timestamp not null default now(),
    constraint pk_credit_card primary key (id),
    constraint fk_owner_id foreign key (owner_id) references infokeeper.user (id)
);

commit transaction;