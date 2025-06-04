FROM alpine:latest

WORKDIR /app

COPY apiApp .
COPY .env.docker .

EXPOSE 3000

CMD ["./apiApp", ".env.docker"]
