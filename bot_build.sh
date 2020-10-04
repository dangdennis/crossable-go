./db_migrate.sh
cd seeder && go run main.go
cd ../
cd bot && go build -o bot