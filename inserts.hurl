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

POST http://localhost:8080/users/operations?alias=alias1
Content-type: application/json
[
      {"type": "Deposito", "amount": "1000"},
      {"type": "Deposito", "amount": "2000"},
      {"type": "Deposito", "amount": "3000"},
      {"type": "Deposito", "amount": "4000"},
      {"type": "Deposito", "amount": "5000"},
      {"type": "Retiro", "amount": "500"},
      {"type": "Retiro", "amount": "750"},
      {"type": "Transferencia", "amount": "500", "to_alias": "alias2"},
      {"type": "Transferencia", "amount": "250", "to_alias": "alias3"}
]
