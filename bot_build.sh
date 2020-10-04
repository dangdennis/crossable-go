./db_migrate.sh
cd seeder && go build -o seeder
cd ../
cd bot && go build -o bot