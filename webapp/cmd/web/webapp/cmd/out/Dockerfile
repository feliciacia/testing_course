FROM golang:1.17-alpine
WORKDIR /app
COPY go.mod .
COPY go.sum .
COPY . .
RUN go mod download
RUN go build -o ./webapp/cmd/web .
CMD ./webapp/cmd/web
