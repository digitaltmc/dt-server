#### Git repository for article https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b

## How to connect to CI

heroku.com
owen263+digitaltmc@gmail.com

https://dt-be.herokuapp.com/

## Backend Interface
- Register
```
mutation {
  register(person:{name:"Wow",password:"world",email:"aaa@sap.com",mobile:"888888"})
}
```
- Loigin
```
{
  login(user:"Wow",password:"world")
}
```
- Book
<br />Please set the date value as the example format like "2019-03-11T00:00:00Z".
```
mutation {
  book(date: "2019-03-11T00:00:00Z", role: Speaker, title: "Hey buddy") {
    id
    date
    agenda {
      role
      title
      duration
      member{
        name
        mobile
        email
        id
      }
    }
  }
}
```