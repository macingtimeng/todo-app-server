# TodoKu - API
Save your notes with TodoKu

# Demo
- [Swagger(Soon)]()
- [API(Soon)]()
- [Demo(Alpha Version)](http://todoku.netlify.app/)

# Schema
| Domain       | Method   | Endpoint                     | Middleware                     | Description          |
|--------------|----------|------------------------------|--------------------------------|----------------------|
| Users        | POST     | /users/register              | -                              | User register        |
| Users        | POST     | /users/login                 | -                              | User login           |
| Users        | PATCH    | /users/modify                | Authentication                 | User modify          |
| Users        | GET      | /users/profile               | Authentication                 | User profile         |
| Todos       | POST      | /todos/                      | Authentication                 | Add Todo             |
| Todos       | GET       | /todos/                      | Authentication                 | Get Todos            |
| Todos       | PATCH     | /todos/:todoId               | Authentication & Authorization | Update Todo          |
| Todos       | GET       | /todos/:todoId               | Authentication & Authorization | Detail Todo          |
| Todos       | DELETE    | /todos/:todoId               | Authentication & Authorization | Delete Todo          |

# Tech Stack
- [Go](https://go.dev/)
- [GORM](https://gorm.io/)
- [Fiber](https://gofiber.io/)
- [Govalidator](https://github.com/asaskevich/govalidator)
- [Jwt-go](https://github.com/golang-jwt/jwt)
- [Crypto](https://pkg.go.dev/crypto)
- [Swagger Documentation](https://github.com/swaggo)
- [GORM Postgres Driver](https://github.com/go-gorm/postgres)
- [PostgreSQL](https://www.postgresql.org/)
- [Godotenv](https://github.com/joho/godotenv)
