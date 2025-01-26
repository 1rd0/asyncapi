db_login:
	psql ${DATABASE_URL}
db_create_migrate:
	migrate create -ext sql -dir migrations -seq init_mg
db_migrate:
	migrate -path migrations -database ${DATABASE_URL} -verbose up
db_migrate_test:
 	migrate -path migrations -database ${DATABASE_URL} -verbose up
