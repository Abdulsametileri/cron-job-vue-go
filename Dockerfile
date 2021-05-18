FROM node:lts-alpine as VueBuilder
COPY client/ client/
RUN cd client && yarn build

FROM golang:alpine AS GoBuilder
WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download
COPY --from=VueBuilder client/dist client/dist
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

FROM alpine:3
RUN apk update \
    && apk upgrade
RUN apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates 2>/dev/null || true
COPY --from=GoBuilder /app/.env .
COPY --from=GoBuilder /app/main .
EXPOSE 3000
CMD ["./main"]