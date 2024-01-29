POST http://localhost:8080/games
{
    "login": "login",
    "password": "pass"
}
HTTP 201
[Captures]
id: body

PUT http://localhost:8080/games/{{id}}
{
    "credentials": {
        "login": "login",
        "password": "pass"
    },
    "title": "updated title"
}

GET http://localhost:8080/games/{{id}}

GET http://localhost:8080/games
