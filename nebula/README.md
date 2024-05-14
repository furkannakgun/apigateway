# Nebula

NaaP API Gateway.

## Running Nebula

Nebula API gateway has several dependencies to have the end to end functionality running:
- Microservice to proxy requests to
- Configuration DB (ETCD or Consul)
- Redis for distributed Nebula instances
- Nebula Management API

To ease the development job, there are two different ways to run Nebula in your PC.

### Standalone

Running Nebula standalone locally only requires having two files in the root directory of the repo: `config.yaml`, that already is in the repo (a template), and `openAPIspec.yaml`, that it is the specification of the service behind Nebula used to validate the requests with. The latter has to be created by the user and stored in the root directory.

To run Nebula in standalone, just run `make run`, as the Makefile will set up the environment to let Nebula know to use the two aforementioned local files instead of the configuration database to obtain the data.

Logs when running Nebula standalone are stored in `naap-service.log` file in the root directory.

To test the proxying and response management functionalities, a microservice has to be running and Nebula has to point to it in the configuration file. [go-test-service](https://github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/go-test-service) is used for testing Nebula. The [dev-env](https://github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/dev-env) can be used to deploy it with Docker.

### Real Environment Nebula

Although the standalone deployment can be used to test some of the functionalities provided by Nebula, we can run full-stack Nebula locally too. To deploy the dependencies the best way is to use the [dev-env](https://github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/dev-env). The following instructions assume the use of the dev environment dev-env to deploy the dependencies, using go-test-service as the service behind Nebula.

First, we need to configure the enviroment, i.e. environment variables:
```bash
export CONFIG_DB_URL=http://localhost:2379
export CONFIG_DB_KEY=go-test-service
```
By default, we use ETCD as configuration database, but if by any change you need to use consul, also set `CONFIG_DB_PROVIDER`:
```bash
export CONFIG_DB_PROVIDER=consul
```

If there is a need to test distributed rate limiters, the Redis configuration must be indicated in the configuration of Nebula. `config.yaml` can be use to store the configuration that we will later send to the configuration database. To check the configuration options, please check the [configuration section](#configuration).

After configuring, the next step is to deploy all the dependencies. Following [dev-env](https://github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/dev-env) `README.md` file, we start running in our docker everything needed:
```bash
docker-compose up etcd
docker-compose run --rm -p 8081:8080 go-test-service # Config should include http://localhost:8081 as target URL
docker-compose run --rm -p 8082:8080 nebula-management # If we want to use Nebula Management API
docker-compose up redis # Only if needed
```
**NOTE**: Nebula runs by default on port 8080, so the dependencies must run in different ports.

Before running Nebula, we need to store the configuration and the open API specification in our configuration database. If we are using [Nebula Management API](https://github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/NaaP_Nebula_Management) (follow its instructions if needed), we use it to store the configuration in our config database:
```bash
POST /nebula-management/config/{service_key} # Configuration
POST /nebula-management/spec/{service_key} # Open API specification
```

If we are not using Nebula management API, we need to install etcdctl and run the following commands:
```bash
etcdctl put go-test-service/config < config.yaml
etcdctl put go-test-service/spec < openAPIspec.yaml
```
**NOTE**: Both files shall be modified according to the needs of the user and the configuration stated above.


Finally, now that Nebula can read the config, we run it from outside the Makefile. From the root directory of the repo, run:
```bash
go run .
```

### Configuration

The description of Nebula's configuration can be found in the [configuration wiki page](https://github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/wiki/Configuration).

## Copyright

Copyright 2024 Vodafone.
