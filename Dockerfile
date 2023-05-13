#build stage
FROM golang:1.20 AS builder
RUN apt-get update -qq
RUN apt-get install -y -qq \
    git \
    openssl \
    curl \
    libtesseract-dev \
    libleptonica-dev \
    tesseract-ocr-eng \
    tesseract-ocr-rus
WORKDIR /app
COPY . .
RUN go get -d -v ./...
RUN GOOS=linux go build -o /bin/app -v ./cmd/muerta/
RUN openssl genrsa -out ./cert/access.pem 4096
RUN openssl rsa -in ./cert/access.pem -pubout -out ./cert/access.pub
RUN openssl genrsa -out ./cert/refresh.pem 4096
RUN openssl rsa -in ./cert/refresh.pem -pubout -out ./cert/refresh.pub
ENV CERT_PATH=/app/cert
EXPOSE ${PORT}
ENTRYPOINT [ "/bin/app" ]
