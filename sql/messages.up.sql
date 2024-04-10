create or replace function type_exists(type_name text)
returns boolean
  return exists(
    select 1
    from pg_type
    where typname = type_name
  );

do $$ begin
  if not type_exists('address') then
    execute 'create domain address as bytea check(octet_length(value) = 20)';
  end if;
  if not type_exists('uint256') then
    execute 'create domain uint256 as bytea check(octet_length(value) <= 32)';
  end if;
  if not type_exists('signature') then
    execute 'create domain signature as bytea check(octet_length(value) = 65)';
  end if;
end $$;

create table if not exists messages (
  from_ address not null,
  nonce uint256 not null,
  order_ uint256 not null,
  start uint256 not null,
  to_ address not null,
  gas uint256 not null,
  data bytea not null,
  orig_sign signature not null,
  valid_sign signature not null unique,
  id serial primary key not null
);
