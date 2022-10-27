# Собираем в гошке
FROM golang:1.18.3 as build

ENV BIN_FILE /bin/app
ENV CODE_DIR /go/src/

# Кэшируем слои с модулями
WORKDIR ${CODE_DIR}
COPY go.mod .
COPY go.sum .
RUN go mod download

# Копипуем код
COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o ${BIN_FILE} ./cmd/sender

## На выходе тонкий образ
FROM alpine:3.16

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="sender"
LABEL MAINTAINERS="shabandrew@mail.ru"

ENV BIN_FILE /bin/app
ENV CONFIG_FILE config.toml

COPY --from=build ${BIN_FILE} ${BIN_FILE}
COPY configs/sender/config.toml ${CONFIG_FILE}

CMD sleep 15 && ${BIN_FILE} -config ${CONFIG_FILE}
