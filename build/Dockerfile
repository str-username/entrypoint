ARG BUILD_IMAGE="golang:1.22.1-alpine3.18"

FROM $BUILD_IMAGE as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GARCH=amd64 go build -a -installsuffix cgo -o build/entrypoint cmd/main.go


FROM alpine:3.18

WORKDIR /app

COPY --from=build /app/etc /app/etc
COPY --from=build /app/build/ /usr/bin

USER nobody:nobody

ENTRYPOINT ["entrypoint"]