# Backend Coding Challenge

We are looking for engineers who can build web services which are simple, stable and secure and who have the user in mind.

To test these requirements, we have a challenge which should not take more than 2 hours of your time.

You should build something with which you are happy with and you think that it can be deployed to production.

## Challenge

### User Management API

The user service already comes with 3 working routes:

- create User
- get User
- delete User 

### Authentication API

Please extend the API the following functionality:

- login a User (return a Bearer token with a TTL of 24 hours and the logged in user)
- logout a User
- a User should only be able to delete itself

###  Technical Requirements
During your development, always remind youself of developing simple, stable and secure software.

1. all functionality should be covered by unit and integration tests
2. errors should be handled including 401s, 404s and 500s
3. don’t forget code comments and logging
4. code as simple / easy to read and functional / modular as possible
5. you can use Golang, Postgres and Redis
6. Use `make test-local` for a local test setup of your solution
7. Publish your code on a public Github repo and use a generic, hard to guess project name
8. Push to Github all the time. We are especially interested in the Git history and not only the final commit

## Development

### Dependencies

The development setup depends on the following tools:

- go
- docker
- docker-compose
- golint
    - please run `go get -u go get golang.org/x/lint/golint` to install golint

### Operations

Make provides a interface to common development operations

- `make test-local` spins up a testing environment and executes the tests