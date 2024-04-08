# Prepaid Gas Server

How to interact with validator?

`curl -d '{"origSign":"0x6cec8506276e7b6edd424941174a2c8fa6a53ed26adb5b668a600de9540e3c106559f5c7a6d0dd1bfe11b4550b7a5a047f83d84a894de0bc0faa38159fe043571c","message":{"from":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8","nonce":"0x00","order":"0x00","start":"0xffeeffee","to":"0x0000000000000000000000000000000000000000","gas":"0x00","data":"0x"}}' -X POST http://localhost:8001/validate`

`curl 'http://localhost:8001/load?offset=0&reverse=false'`

How to run psql?

`source .env && docker run -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD -e POSTGRES_USER=$POSTGRES_USER -p 5432:5432 postgres:15.4`
