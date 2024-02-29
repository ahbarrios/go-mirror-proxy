# syntax=docker/dockerfile:1

FROM golang:1.21-alpine

WORKDIR /app

# Download Go modules
# COPY go.mod go.sum ./
# RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./ 

# Trust local certificates
# COPY internal/email/testdata/server.pem /usr/local/share/ca-certificates/server.pem
# RUN cat /usr/local/share/ca-certificates/server.pem >> /etc/ssl/certs/ca-certificates.crt

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /proxy

CMD [ "/proxy" ]