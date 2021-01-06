FROM golang:1.15-alpine AS builder
WORKDIR /src
COPY . .
RUN GO111MODULE=on go build

FROM alpine:latest AS server
WORKDIR /app
COPY --from=builder /src/maria ./maria
COPY --from=builder /src/index.html ./index.html

EXPOSE 3201

CMD [ "./maria" ]