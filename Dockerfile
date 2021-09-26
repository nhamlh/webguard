FROM golang:1.16 AS BUILDER

RUN mkdir /build
WORKDIR /build

COPY . /build

RUN go mod tidy && go build -o webguard ./cmd/

FROM golang:1.16

LABEL author="Hoai-Nham Le <lehoainham@gmail.com>"

ENV WEBGUARD_CONFIG_FILE="config.json"

USER root

RUN mkdir /webguard
WORKDIR /webguard

COPY --from=BUILDER /build/webguard .

RUN ./webguard genconf > ./config.json

ENTRYPOINT ["sh", "-c", "./webguard -c ${WEBGUARD_CONFIG_FILE} db migrate && ./webguard -c ${WEBGUARD_CONFIG_FILE} start"]
