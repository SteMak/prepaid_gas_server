# Prepaid Gas Server

How to interact with validator?

`curl -d '{"sign":"784e44edf1884edc8f73d355512642d25e0569ab7f67f5f65a8df433c1053ce140855a7f966b82e2cfdef9f8ca41b6538254dc844eb4aa6d4eaed9b5a9503b4501","message":{"signer":"0x3428B2b8384d024881445cc7fd6423065849CEA8","nonce":"0x123","gasOrder":"0x12","onBehalf":"0x0000656EC7ab88b098defB751B7401B5f6d8976F","deadline":"0xffeeffee","endpoint":"0x71C7656EC7ab88b098defB751B7401B5f6d8976F","gas":"0x13","data":"0f0000000000000000000000000000000000000000000000000000000000000000001113150d"}}' -X POST http://localhost:8001`

How to run psql?

`source .env && docker run -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD -e POSTGRES_USER=$POSTGRES_USER -p 5432:5432 postgres:15.4`