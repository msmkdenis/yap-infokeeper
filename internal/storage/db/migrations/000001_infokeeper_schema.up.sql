begin transaction;

create schema if not exists infokeeper;

create table if not exists infokeeper.user
(
    id                      text,
    login                   text unique not null,
    password                bytea not null,
    constraint pk_user primary key (id)
);

commit transaction;