FROM alpine:latest

RUN mkdir /app

COPY binary_file/loggerServiceApp /app

CMD [ "/app/loggerServiceApp" ]
