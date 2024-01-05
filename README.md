My broker
=============

In this repository, I have designed six services. The front-end part is already created, and the authentication service, broker service, mail service, listener service, and logger service are also written.

More information about every service can be found below:
- [Front-end-service](/front-end-service/README.md)
- [Broker-service](/broker-service/README.md)
- [Authentication-service](/authentication-service/README.md)
- [Logger-service](/logger-service/README.md)
- [Mail-service](/mail-service/README.md)
- [Listener-service](/listener-service/README.md)

In the `project` can be found all files with basic commands and the main `docker-compose.yml`

To start all back-end services, you need to clone this repository and go to `project`:
```
git clone git@github.com:GermanLepin/my_broker.git
cd project/
```

Let's start all back-end services with the command:
```
make up_build
```

To start front-end service, use the command:
```
make front_end_start
```

address of the main page: http://localhost:80
address of the mailHog: http://localhost:8025

To stop front-end service, use the command:
```
make front_end_stop
```

To stop back-end services, use the command:
```
make down
```


