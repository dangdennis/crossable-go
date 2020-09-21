#!/bin/bash
set -e

if [[ ${CROSSING_ENV} == "development" ]]; then
    echo "NOTE: this script will DROP your database!"
    if [[ $1 != "-f" ]]; then
        read -p "press any key to continue or ctrl-C to exit"
    fi
else
    echo "ERROR: this script can only run in \"development\" environment"
    exit 1
fi

HOST=localhost
DATABASE=crossing_dev
USER=postgres

echo "### dropping $DATABASE database ..."
PGPASSWORD=postgres psql -h ${HOST} -d ${DATABASE} -U ${USER} -f sqls/drop_all.sql

echo "### init and migrate with Prisma"
go run github.com/prisma/prisma-client-go migrate save --experimental --create-db --name ${DATABASE}
go run github.com/prisma/prisma-client-go migrate up --experimental
go run github.com/prisma/prisma-client-go generate