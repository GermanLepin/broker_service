FROM alpine:latest

RUN mkdir /app

COPY binary_file/brokerServiceApp /app

CMD [ "/app/brokerServiceApp" ]
