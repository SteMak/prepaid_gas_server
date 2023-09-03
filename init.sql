CREATE TABLE messages (
  id bigint,
  signer varchar(40),
  nonce varchar(64),
  gas_order varchar(64),
  on_behalf varchar(40),
  deadline varchar(64),
  endpoint varchar(40),
  gas varchar(64),
  data text,
  signature text
);