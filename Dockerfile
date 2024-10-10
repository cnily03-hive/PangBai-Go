FROM golang:1.18-alpine AS builder

COPY . /app
RUN mkdir /out
WORKDIR /app
RUN go build -o /out/main .
RUN cp -r /app/assets /app/views /app/sign.txt -t /out

FROM alpine:3.20

COPY --from=builder /out /app

ENV FLAG="flag{test_flag}"

RUN apk add --update \
    curl \
    && rm -rf /var/cache/apk/*

RUN adduser -D ctf
USER ctf

WORKDIR /app
CMD ["/app/main"]

EXPOSE 8000