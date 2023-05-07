# Собираем в гошке
FROM golang:1.19 as build

ENV BIN_SCHEDULER_FILE /opt/calendar/scheduler-app
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
        -o ${BIN_SCHEDULER_FILE} cmd/calendar_scheduler/*

# На выходе тонкий образ
FROM alpine:latest

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar-scheduler"
LABEL MAINTAINERS="dmitrijj.kuzmin.v@vseinstrumenti.ru"

ENV BIN_SCHEDULER_FILE /opt/calendar/scheduler-app
COPY --from=build ${BIN_SCHEDULER_FILE} ${BIN_SCHEDULER_FILE}

ENV CONFIG_SCHEDULER /etc/calendar/config_scheduler.toml
COPY ./configs/config_scheduler.toml ${CONFIG_SCHEDULER}

CMD ${BIN_SCHEDULER_FILE} -config ${CONFIG_SCHEDULER}
