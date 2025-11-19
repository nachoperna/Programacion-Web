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
