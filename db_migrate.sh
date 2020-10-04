cd common || exit
go run github.com/prisma/prisma-client-go migrate save --experimental --create-db --name ${CROSSING_DATABASE_NAME}
go run github.com/prisma/prisma-client-go migrate up --experimental
go run github.com/prisma/prisma-client-go generate