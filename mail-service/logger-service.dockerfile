FROM alpine:latest

RUN mkdir /app

COPY binary_file/mailServiceApp /app

CMD [ "/app/mailServiceApp" ]
