# Ozon Tech Golang School Homework
This repo contains homework from Route256 course from Ozon Tech

## Startup
Put API_KEY and port settings for GRPC and Gateway servers into config.yaml file

"make run-server" - buildand run GRPC server
"make run-bot" - build and run bot

## Supported bot commands

Create person
/person create {LastName} {FirstName}

Update person
/person update {personID} {LastName} {FirstName}

Delete person
/person delete {personID}

Person list
/person list


For GRPC
Create person
POST http://localhost:xxxx/v1/persons
{
    "lastName": "A",
    "firstName": "B"
}

Update person
PUT http://localhost:xxxx/v1/persons
{
    "id": 1,
    "lastName": "A",
    "firstName": "B"
}

Delete person
DELETE http://localhost:xxxx/v1/persons/{id}

Person list
GET http://localhost:xxxx/v1/persons

Swagger
http://localhost:xxxx/swagger/index.html

## Project structure

- entity
- service
- handlers
- storage