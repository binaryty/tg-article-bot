run:
	go run ./...
build:
	go build -v ./...

storageup:
	docker-compose up -d

storagedown:
	docker-compose down

createdb:
	docker exec -it article-bot-db-1 createdb --username=postgres --owner=postgres articles

dropdb:
	docker exec -it article-bot-db-1 dropdb articles

mgrup:
	migrate -path migration -database "postgresql://postgres:postgres@localhost:5432/articles?sslmode=disable" -verbose up

mgrdown:
	migrate -path migration -database "postgresql://postgres:postgres@localhost:5432/articles?sslmode=disable" -verbose down

.PHONY: run build storageup storagedown createdb dropdb mgrup mgrdown