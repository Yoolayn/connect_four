GET http://localhost:8080/users
HTTP 200

POST http://localhost:8080/users
{
    "login": "user2",
    "password": "password"
}
HTTP 201

PUT http://localhost:8080/users/user2/password
{
    "credentials": {
        "login": "user2",
        "password": "password"
    },
    "newpassword": "password2"
}

PUT http://localhost:8080/users/user2/name
{
    "credentials": {
        "login": "user2",
        "password": "password"
    },
    "newname": "madie"
}

GET http://localhost:8080/users/user2
HTTP 200

POST http://localhost:8080/admins/user2
{
    "login": "user2",
    "password": "password"
}
HTTP 403

POST http://localhost:8080/admins/user2
{
    "login": "login",
    "password": "pass"
}
HTTP 200

DELETE http://localhost:8080/admins/user2
{
    "login": "login",
    "password": "pass"
}
HTTP 200

DELETE http://localhost:8080/users/user2
{
    "login": "user2",
    "password": "password"
}
HTTP 200
