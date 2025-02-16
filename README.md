# avito_shop
## Запуск
```bash
git clone git@github.com:StarNuik/avito_shop.git
cd avito_shop

sudo docker compose -f prod.compose.yaml build
sudo docker compose -f prod.compose.yaml up -d

http POST localhost:8080/api/auth username=user#1 password=password#1
http GET localhost:8080/api/buy/book "Authorization:Bearer {insert auth token here}"
http GET localhost:8080/api/info  "Authorization:Bearer {insert auth token here}"
```
## Тест
```bash
git clone git@github.com:StarNuik/avito_shop.git
cd avito_shop

sudo docker compose -f test.compose.yaml up -d

DATABASE_PASSWORD=password DATABASE_NAME=shop go test -p 1 -coverpkg=./... -coverprofile=./coverage/all.cov ./...
go tool cover -html=./coverage/all.cov -o=./coverage/all.html
go tool cover -func=./coverage/all.cov | grep total | awk '{print $3}'
```