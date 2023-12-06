FROM alpine:latest

RUN mkdir /app

COPY binary_file/mailServiceApp /app
COPY templates /templates

CMD [ "/app/mailServiceApp" ]
