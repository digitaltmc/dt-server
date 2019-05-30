#### Git repository for article https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b

## How to connect to CI

heroku.com
owen263+digitaltmc@gmail.com

https://dt-be.herokuapp.com/

## Backend Interface

### Register

```
mutation {
  register(person:{name:"Wow",password:"world",email:"aaa@sap.com",mobile:"888888"})
}
```

### Login

```
{
  login(user:"Wow",password:"world")
}
```
You will receive a jwt token here.

### Book

Please set the date value as the example format like "2019-03-11T00:00:00Z".

```
mutation {
  book(token: "THE_JWT_TOKEN_YOU_RECEIVED", date: "2011-01-01", roleName: "ie", title: "gege")
}
```
There will be a unique index constraint on data+roleName, e.g. user can't rebook the same role in the same day.
<pre>
digitaltmc.meetingItems
index: date_1_roleName_1
</pre>
The unique index will be manually created in mlab portal.

