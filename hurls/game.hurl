POST http://localhost:8080/games
{
    "login": "login",
    "password": "pass"
}
HTTP 201
[Captures]
id: body

GET http://localhost:8080/games/{{id}}

GET http://localhost:8080/games
