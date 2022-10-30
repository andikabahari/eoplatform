DB_USER="root"
DB_PASS="root"
DB_HOST="localhost"
DB_PORT="3306"
DB_NAME="eoplatform"

DSN="${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?charset=utf8&parseTime=True&loc=Local"

migrateup:
	goose -dir migration mysql ${DSN} up

migratedown:
	goose -dir migration mysql ${DSN} down
