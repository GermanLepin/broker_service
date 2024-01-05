FROM alpine:latest

RUN mkdir /app

COPY binary_file/listenerServiceApp /app

CMD [ "/app/listenerServiceApp" ]
