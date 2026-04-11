FROM golang:1.25-alpine AS builder
WORKDIR /app                                                                                                                                                                    
COPY go.mod go.sum ./                                                                                                                                                           
RUN go mod download
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest                                                                                         

COPY . .                                                                                                                                                                        
RUN go build -o warden ./cmd/server                                                                                                                                             
                                                                                                                                                                                
FROM alpine:latest
WORKDIR /app                                                                                                                                                                    
COPY --from=builder /app/warden .
COPY --from=builder /app/migrations/ ./migrations/
COPY --from=builder /go/bin/migrate .
EXPOSE 1323
CMD ["./warden"]