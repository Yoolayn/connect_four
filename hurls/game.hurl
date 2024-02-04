POST http://localhost:8080/games
{
    "login": "login",
    "password": "pass"
}
HTTP 201
[Captures]
id: body

POST http://localhost:8080/games/{{id}}
{
    "credentials": {
        "login": "login",
        "password": "pass"
    },
    "color": "red"
}

DELETE http://localhost:8080/games/{{id}}/leave
{
    "login": "login",
    "password": "pass"
}

GET http://localhost:8080/games

PUT http://localhost:8080/games/{{id}}
{
    "credentials": {
        "login": "login",
        "password": "pass"
    },
    "title": "updated title"
}

GET http://localhost:8080/games/{{id}}

DELETE http://localhost:8080/games/{{id}}
{
    "login": "login",
    "password": "pass"
}

GET http://localhost:8080/games
