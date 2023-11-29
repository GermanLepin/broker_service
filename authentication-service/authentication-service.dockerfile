FROM alpine:latest

RUN mkdir /app

COPY binary_file/authenticationServiceApp /app

CMD [ "/app/authenticationServiceApp" ]
