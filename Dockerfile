# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -o /out/cashtro ./cmd/cashtro \
 && CGO_ENABLED=0 go build -o /out/ultron ./cmd/ultron

FROM alpine:3.20
RUN adduser -D -H cashtro
USER cashtro
WORKDIR /home/cashtro
COPY --from=build /out/cashtro /usr/local/bin/cashtro
COPY --from=build /out/ultron /usr/local/bin/ultron
EXPOSE 8080 9090
VOLUME ["/home/cashtro/data"]
ENTRYPOINT ["cashtro"]
CMD ["-addr", ":8080", "-data", "/home/cashtro/data/cashtro.json"]
