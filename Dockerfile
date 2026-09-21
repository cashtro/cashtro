# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -o /out/cashtro ./cmd/cashtro

FROM alpine:3.20
RUN adduser -D -H cashtro
USER cashtro
WORKDIR /home/cashtro
COPY --from=build /out/cashtro /usr/local/bin/cashtro
ENV PORT=8080
ENV CASHTRO_DATA=/home/cashtro/data/cashtro.json
ENV CASHTRO_CLOSED=1
ENV CASHTRO_PULSE_EVERY=2m
EXPOSE 8080
ENTRYPOINT ["cashtro"]
