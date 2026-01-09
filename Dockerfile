FROM golang:1.25.5-alpine3.23 AS builder

WORKDIR /gymbuddy

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY base_page.html .
COPY exercise_page.html .
COPY handlers.go .
COPY main.go .
COPY schemas.go .
COPY styles.css .
COPY workout_page.html .

RUN go build -o gymbuddy .

FROM scratch

WORKDIR /gymbuddy

COPY --from=builder /gymbuddy/gymbuddy .

EXPOSE 8080

ENV GYMBUDDY_CONFIG_PATH=/gymbuddy/config.yml

CMD ["./gymbuddy"]
