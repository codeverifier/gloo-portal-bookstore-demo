# ---------------------------
# 1st stage - do a build
# ---------------------------
FROM golang:1.19-alpine AS builder

RUN apk --no-cache add make
WORKDIR /app
COPY . .
RUN make build

# ---------------------------
# 2nd stage - copy binary
# ---------------------------
FROM alpine

RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/dist/* ./
CMD ["./main"]