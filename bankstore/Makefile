pg_start:
	docker-compose start

pg_stop:
	docker-compose stop

createdb:
	docker-compose exec db createdb --username=postgres --owner=postgres bankstoredb

dropdb:
	docker-compose exec db dropdb -U postgres bankstoredb

migrateup:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5435/bankstoredb?sslmode=disable" up

migratedown:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5435/bankstoredb?sslmode=disable" down

sqlc:
	sqlc generate

test:
	go test -timeout 180s -v -cover ./... -count=1

server:
	go run main.go

.PHONY: pg_start pg_stop createdb dropdb migrateup migratedown sqlc test server