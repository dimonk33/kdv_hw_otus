# Собираем в гошке
FROM golang:1.19 as build

ENV BIN_SENDER_FILE /opt/calendar/sender-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY .. ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_SENDER_FILE} cmd/calendar_sender/*

# На выходе тонкий образ
FROM alpine:latest

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar-sender"
LABEL MAINTAINERS="dmitrijj.kuzmin.v@vseinstrumenti.ru"

ENV BIN_SENDER_FILE /opt/calendar/sender-app
COPY --from=build ${BIN_SENDER_FILE} ${BIN_SENDER_FILE}

ENV CONFIG_SENDER /etc/calendar/config_sender.toml
COPY ./configs/config_sender.toml ${CONFIG_SENDER}

CMD ${BIN_SENDER_FILE} -config ${CONFIG_SENDER}
