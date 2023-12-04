postgres: 
	docker start faztcode
createdb:
	docker exec -it faztcode createdb tasks
dropdb:
	docker exec -it faztcode dropdb tasks
migrateup:
	migrate -path db/migration -database "postgresql://root:83postgres19@localhost:5432/tasks?sslmode=disable" -verbose up 
migratedown:
	migrate -path db/migration -database "postgresql://root:83postgres19@localhost:5432/tasks?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go clean -testcache
	go test -v --cover ./...
server:
	air
.PHONY: postgres createdb dropdb migrateup migratedown server