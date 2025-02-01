FROM golang:1.23-alpine

RUN apk add --no-cache bash gcc musl-dev postgresql-dev

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY ./migrations /app/migrations

COPY migrations/entrypoint_up_all.sh /app/entrypoint_up_all.sh
COPY migrations/entrypoint_up_one.sh /app/entrypoint_up_one.sh
COPY migrations/entrypoint_down_one.sh /app/entrypoint_down_one.sh

RUN chmod +x /app/entrypoint_up_all.sh
RUN chmod +x /app/entrypoint_up_one.sh
RUN chmod +x /app/entrypoint_down_one.sh

ENTRYPOINT ["/app/entrypoint_up_all.sh"]
