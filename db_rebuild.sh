./db_drop.sh

./db_migrate.sh

echo "### seed the database from our seeder"
cd seeder && go run main.go
