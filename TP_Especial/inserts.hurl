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

POST http://localhost:8080/requests
Content-type: application/json
[
    {
        "from_alias": "alias2",
        "to_alias": "alias1",
        "amount": "5000",
        "message": "prueba de primer pedido"
    },
    {
        "from_alias": "alias3",
        "to_alias": "alias1",
        "amount": "7500",
        "message": "prueba de segundo pedido"
    },
    {
        "from_alias": "alias2",
        "to_alias": "alias1",
        "amount": "15000",
        "message": "prueba de tercer pedido"
    }
]
