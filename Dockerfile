FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
ADD static static
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /simple_sqli

EXPOSE 8080

CMD ["/simple_sqli"]
