FROM golang:1.20-alpine as build

WORKDIR /build

COPY . .
# RUN go mod download

RUN go build -o ./app ./main.go

FROM alpine as app

WORKDIR /app

RUN apk update && apk upgrade

COPY --from=build /build/app .

ENTRYPOINT [ "app" ]
