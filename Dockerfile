FROM golang:1.17-alpine AS GO_BUILD
ARG CGO_ENABLED=0
WORKDIR /app
ADD . ./
RUN go build -o server ./cmd/server

FROM alpine
WORKDIR /app
COPY --from=GO_BUILD /app/server ./
ENTRYPOINT ["./server"]