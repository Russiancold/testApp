create database accounts;

create schema test;

create table test.accounts(
	id serial primary key,
	username varchar (30) not null,
	email varchar (30) not null,
	api_token varchar (30),
	is_del boolean default FALSE
);

create user accounts_app with encrypted password 'secret';
grant connect on database accounts to accounts_app;
grant usage on schema test to accounts_app;
grant all privileges an all sequences in schema test to accounts_app;
grant select, update, insert on test.accounts to accounts_app;