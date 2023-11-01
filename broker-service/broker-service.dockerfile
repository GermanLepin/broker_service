FROM alpine:latest

RUN mkdir /app

COPY binary_file/brokerApp /app

CMD [ "/app/brokerApp" ]
