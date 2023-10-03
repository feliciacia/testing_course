
FROM golang:1.19.5 as build

WORKDIR /app

COPY . .

RUN go build -v -o app ./webapp/cmd/api

FROM gcr.io/distroless/base
COPY --from=build /app/app /app

CMD ["/app"]
