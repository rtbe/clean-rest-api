# Clean Rest API

## Description

This is example of sales application, build around ideas of clean architecture (hexagonal architecture + layered technique).

This project is build for the sake of practice of different approaches to building golang applications.

[Основные идеи чистой архитектуры](https://gist.github.com/rtbe/7b4f0c5be5369c91942656197c91a7fc)

[Паттерн репозиторий](https://gist.github.com/rtbe/3705b5b3b9dcd0fb34a276d09a5cd93c)

**Some neat tricks&tips applied in this project:**

- Automatization of boring stuff with Makefile (To run a project just type ```make run```).
- Centralized error handling of errors. As well as allowing decide whenever we need to show exact error message to user or show generic one. That lets us hide errors with implementation details from users and log them internally.
- Built in OpenApi v2 (Swagger) documentation.
- More effective kind of pagination [do not use offset for pagination](https://use-the-index-luke.com/no-offset).
- JWT token based authentication.
- Persistent storage tests without mocks using docker containers (To run repository tests you should stop postgresql service: ```sudo systemctl stop postgresql```)
- Two staged docker build to make the size of final docker image small.
- Database migrations.
- Validations for incoming data.
- Built it debugging with debug/pprof.

**Required envirionment variables for the application are located inside .env file**

### Dependencies

When i built this application i stumble upon the set of problems that is to difficult to me to handle alone, so i decide not to reinvent the wheel, but to use existing solutions.

#### Documentation

- [go-swagger/go-swagger](https://github.com/go-swagger/go-swagger) - I used this tool to generate and serve swagger (OpenAPI 2.0) documentation. For generating and updating documentation run ```make swagger```.

#### Storage of data

I used postgreSQL as main database store for an application because it's very popular open-source solution for storing standardized data, with a lot of constraints and connections between data items. These particular needs predeclared my use of relational database for main database store.

#### Storage of authentication data

I used mongoDB as database to store authentication data because ... . Honestly i can used postgreSQL for this need but i wanted to practice to use different (non-relational) approach to storing data.

#### Routing

- [go-chi/chi](https://github.com/go-chi/chi) - I think this is pretty balanced router (Not so hardcore/barebone as [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) while not so bloated as [gorrila/mux](https://github.com/gorilla/mux)).

### Tests

You can run tests with ```make tests``` command. Tests are done with table tests technique, for database related tests i also have been used ory/dockertest to run them in a separate docker container what is very handy (You don`t have to use database mocks so you testing real interaction with real database).
To run database tests inside a separated docker container i used [ory/dockertest](https://github.com/ory/dockertest).

#### Validating incoming data

- [go-playground/validator](https://github.com/go-playground/validator) - Validator for incoming data, uses struct fields for validation.

#### Logging

- [uber-go/zap](https://github.com/uber-go/zap) - I think this is pretty bloated logger, with a lot of stuff is going inside.  Also you don`t have to use custom logger at all [Dave Cheney post about logging](https://dave.cheney.net/2015/11/05/lets-talk-about-logging)). But for large projects or for projects that's planning to be large i think it can be reasonable to use zap logger as it can please a lot of possible use cases of a large projects and be allocation efficient. (You have to keep in mind that needs of Uber company are probably quite different then needs of your personal/company project).

#### Database migrations

- [golang-migrate/migrate](https://github.com/golang-migrate/migrate) - Migration library that runs on wide variety of databases, can use different sources as inputs for migrations, and also has CLI.

## Links

- [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

- [A Little Architecture](https://blog.cleancoder.com/uncle-bob/2016/01/04/ALittleArchitecture.html)
