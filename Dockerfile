FROM golang:1.25-alpine AS builder
WORKDIR /app                                                                                                                                                                    
COPY go.mod go.sum ./                                                                                                                                                           
RUN go mod download
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest                                                                                         

COPY . .                                                                                                                                                                        
RUN go build -o warden ./cmd/server                                                                                                                                             
                                                                                                                                                                                
FROM alpine:3.23
RUN addgroup -S warden && adduser -S warden -G warden
WORKDIR /app
COPY --from=builder /app/warden .
COPY --from=builder /app/migrations/ ./migrations/
COPY --from=builder /go/bin/migrate .
RUN chown -R warden:warden /app
USER warden
EXPOSE 1323
CMD ["./warden"]