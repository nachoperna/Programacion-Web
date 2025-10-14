DELETE http://localhost:8080/users
HTTP 200

POST http://localhost:8080/users
Content-type: application/json
[
      {
            "email": "persona1@mail",
            "alias": "alias1",
            "name": "persona1",
            "password": "persona1pass"
      },
      {
            "email": "persona2@mail",
            "alias": "alias2",
            "name": "persona2",
            "password": "persona2pass"
      },
      {
            "email": "persona3@mail",
            "alias": "alias3",
            "name": "persona3",
            "password": "persona3pass"
      }
]
HTTP 201
[Asserts]
jsonpath "$.usuarios_creados[0].alias" == "alias1"

GET http://localhost:8080/users
HTTP 200

GET http://localhost:8080/users/?alias=alias3
HTTP 200

PUT http://localhost:8080/users/?alias=alias4
Content-type: application/json
{
      "name": "persona1nueva"
}
HTTP 200

DELETE http://localhost:8080/users/?alias=alias1 
HTTP 204
