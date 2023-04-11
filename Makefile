.SILENT: start_db stop_db build run check_race_condition

start_db:
	- docker run --name $(db_name) -e MYSQL_ROOT_PASSWORD=$(db_pass) -e MYSQL_DATABASE=$(db_name) -p 3323:3306 -d mysql:latest --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

stop_db:
	- docker stop test_horus
	- docker rm test_horus

build:
	go build -o ./bin/test_horus

run: build
	- ./bin/test_horus

check_race_condition:
	- go run -race main.go