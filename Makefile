db_login:
	psql ${DATABASE_URL}
db_create_migrate:
	migrate create -ext sql -dir migrations -seq init_mg
db_migrate:
	migrate -path migrations -database ${DATABASE_URL} -verbose up
db_migrate_up:
	migrate -database ${DATABASE_URL_TEST} -path ${DATABASE_URL_TEST} up

