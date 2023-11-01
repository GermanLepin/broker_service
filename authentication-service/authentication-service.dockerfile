FROM alpine:latest

RUN mkdir /app

COPY binary_file/authenticationApp /app

CMD [ "/app/authenticationApp" ]
