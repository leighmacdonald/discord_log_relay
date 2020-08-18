FROM golang:1.14-alpine as build
LABEL maintainer="Leigh MacDonald <leigh.macdonald@gmail.com>"
WORKDIR /build
RUN apk add make git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make linux

FROM alpine:latest
LABEL maintainer="Leigh MacDonald <leigh.macdonald@gmail.com>"
EXPOSE 7777
ARG CHANNEL_ID
ARG TOKEN
RUN apk add bash
WORKDIR /app
COPY --from=build /build/build/linux64/discord_log_relay .
COPY docker.sh .
ENTRYPOINT ["./docker.sh"]
