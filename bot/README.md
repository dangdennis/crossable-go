https://github.com/prisma/prisma-client-go/blob/master/docs/quickstart.md

1. create the migration script(s)

go run github.com/prisma/prisma-client-go migrate save --experimental --create-db --name "crossing_dev"

2. apply the migration script(s)

go run github.com/prisma/prisma-client-go migrate up --experimental

3. generate the client

go run github.com/prisma/prisma-client-go generate
