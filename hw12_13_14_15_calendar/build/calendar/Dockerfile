# Собираем в гошке
FROM golang:1.19 as build

ENV BIN_CALENDAR_FILE /opt/calendar/calendar-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_CALENDAR_FILE} cmd/calendar/*

# На выходе тонкий образ
FROM alpine:latest

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="dmitrijj.kuzmin.v@vseinstrumenti.ru"

ENV BIN_CALENDAR_FILE /opt/calendar/calendar-app
COPY --from=build ${BIN_CALENDAR_FILE} ${BIN_CALENDAR_FILE}

ENV CONFIG_CALENDAR /etc/calendar/config_calendar.toml
COPY ./configs/config_calendar.toml ${CONFIG_CALENDAR}

CMD ${BIN_CALENDAR_FILE} -config ${CONFIG_CALENDAR}
