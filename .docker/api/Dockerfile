#build stage
FROM golang:1.20 AS builder
RUN apt-get update -qq
RUN apt-get install -y -qq \
    git \
    openssl \
    libtesseract-dev \
    libleptonica-dev \
    tesseract-ocr-eng \
    tesseract-ocr-rus
WORKDIR /app
COPY . .
RUN go get -d -v ./...
RUN GOOS=linux go build -o /bin/app -v ./cmd/muerta/
RUN openssl genrsa -out ./cert/id_rsa 4096
RUN openssl rsa -in ./cert/id_rsa -pubout -out ./cert/id_rsa.pub
COPY ./cert /cert
COPY ./config/conf.dev.yaml /cfg/conf.yaml
ENV CONFIG_PATH=/cfg/conf.yaml
ENV CERT_PATH=/cert
# EXPOSE 8000
ENTRYPOINT [ "/bin/app" ]
