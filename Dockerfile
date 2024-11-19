ARG GO_VERSION=1.22.7

FROM golang:${GO_VERSION}-alpine AS builder

# no proxy for install go dependencies
# go direct to the dependencies in go.mod
RUN go env -w GOPROXY=direct

# for subdependencies git
RUN apk add --no-cache git

# get and update security certificates
RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /src

# download modules
COPY ./go.mod ./go.sum ./
RUN go mod download

# copy als the files
COPY ./ ./

# build the app
# CGO_ENABLED=0 not use the c++ compiler
# installsuffix 
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /go-rest-ws .

# ------------------------------------------------------------------------
# scratch image from run the app
FROM scratch AS runner

# copy the certificates from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /ect/ssl/certs/

# copy the env vars
COPY .env ./
# copy the binary from builder stage
COPY --from=builder /go-rest-ws /go-rest-ws

EXPOSE 5050

# execute
ENTRYPOINT ["/go-rest-ws"]
