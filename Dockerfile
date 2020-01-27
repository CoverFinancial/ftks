FROM golang:1.13-alpine AS builder

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /src
COPY . .

RUN go build \
  -a \
  -ldflags "-s -w -extldflags 'static'" \
  -installsuffix cgo \
  -o /bin/app \
  .


FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -g 1001 appgroup && \
  adduser -H -D -s /bin/false -G appgroup -u 1001 appuser

USER 1001:1001
COPY --from=builder /bin/app /bin/app
CMD ["/bin/app"]
