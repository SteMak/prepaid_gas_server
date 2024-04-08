drop table if exists messages;
drop domain if exists address;
drop domain if exists uint256;
drop domain if exists signature;
create domain address as bytea check(octet_length(value) = 20);
create domain uint256 as bytea check(octet_length(value) <= 32);
create domain signature as bytea check(octet_length(value) = 65);
create table messages (
  from address not null,
  nonce uint256 not null,
  order uint256 not null,
  start uint256 not null,
  to address not null,
  gas uint256 not null,
  data bytea not null,
  orig_sign signature not null,
  valid_sign signature not null unique,
  id serial primary key not null
);