FROM golang:alpine AS Builder
LABEL maintainer="Abdulsamet İleri <abdulsamet.ileri@ceng.deu.edu.tr>"
WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

FROM alpine:3

RUN apk update \
    && apk upgrade
RUN apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates 2>/dev/null || true

COPY --from=builder /app/.env .
COPY --from=builder /app/main .

EXPOSE 3000
CMD ["./main"]