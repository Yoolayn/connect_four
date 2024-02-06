GET http://localhost:8080/search?users=login

GET http://localhost:8080/users

POST http://localhost:8080/games
{
    "login": "login",
    "password": "pass"
}

GET http://localhost:8080/games
GET http://localhost:8080/search?games=new
GET http://localhost:8080/search?games=new%20game
