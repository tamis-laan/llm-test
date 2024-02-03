# .PHONY: dev load build postgres createdb dropdb psql

dev-api:
	fd -e go -e yaml -e html | entr -nrc go run  cmd/api/main.go --migrate

dev-etl:
	fd -e go -e yaml -e html | entr -nrc go run  cmd/etl/main.go --migrate

build:
	go mod tidy
	docker build -t api:latest --target deployment . 
	docker build -t ctl:latest --target ctl . 

load: build
	kind load docker-image api:latest
	kind load docker-image ctl:latest

psql:
	docker compose exec database psql -U myuser -d mydb

api-shell:
	docker compose exec -it api /bin/sh

etl-shell:
	docker compose exec -it etl /bin/sh
