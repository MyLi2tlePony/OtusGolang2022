FROM migrate/migrate

ENV MIGRATION_PATH /migrations

COPY ${MIGRATION_PATH} ${MIGRATION_PATH}

CMD migrate -path ${MIGRATION_PATH} -database postgresql://postgres:1234512345@postgres:/postgres?sslmode=disable up