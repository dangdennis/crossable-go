DATABASE=crossing_dev

go run github.com/prisma/prisma-client-go migrate save --experimental --create-db --name ${DATABASE}
go run github.com/prisma/prisma-client-go migrate up --experimental
go run github.com/prisma/prisma-client-go generate