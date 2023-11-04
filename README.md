Broker
=============

In this repository, I have designed three services. The front-end part is already created, and the authentication service and broker service are written and function very well. 

More information about every service can be found below:
- [Front-end-service](/front-end-service/README.md) 
- [Broker-service](/broker-service/README.md) 
- [Authentication-service](/authentication-service/README.md) 

In the `project` can be found all files with basic commands and the main `docker-compose.yml`

To start broker service, authentication service, and front-end service, you need to clone this repository and go to project:
```
git clone git@github.com:GermanLepin/my_broker.git
cd project/
```

Lets start broker service and authentication service with the command:
```
make up_build
```

To start front-end service, use the command:
```
make start
```

To stop front-end service, use the command:
```
make stop
```

To stop broker service and authentication service, use the command:
```
make down
```

To be continue...

