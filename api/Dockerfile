FROM golang:latest as builder

WORKDIR /usr/src/devbook/api
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags '-w -s' -o api .
FROM scratch
COPY --from=builder /usr/src/devbook/api .

CMD ["./api"]
