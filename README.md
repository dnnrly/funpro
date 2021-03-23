# funpro
Access the infinite compute of functions using HTTP and web security

## Aims

* Allow you to access your function via HTTP calls
* Allow you to on-board your functions using the minimum of configuration
* Allow you to own the security model for accessing your lambdas

## Non Aims

* Be a general API Gateway

## Developing funpro

Well, it's written in Go so you should be alright. But to get you started you should probably know a few things. We're using `make` to help us run and coordinate all of the commands. It's a doddle to run the tasks that you need.

Building `funpro`:

`make build`

Running unit tests:

`make test`

Running unit tests with coverage data:

`make ci-test`

Running acceptance tests (fully built app, tested locally):

`make acceptance-test`

Running integration tests (fully built app, tested against mock AWS environment):

`make integration-test`

Run linting:

`make lint`