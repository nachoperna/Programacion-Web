# CandyPay

## Integrante
- Ignacio Agustín Perna

## Descripción
Este proyecto consta de una billetera virtual que tendrá todas las funcionalidades principales de la misma:
- Creación de cuenta con formularios de registro e inicios de sesión posteriores
- Depositar y retirar dinero
- Transferr a cualquier usuario con un alias registrado
- Pedir préstamos con fechas de pago
- Invertir dinero y generar rendimientos

## Modo de Uso
1. Abrir una terminal localizada dentro de la carpeta donde se encuentra este archivo README.
2. Ejecutar dentro de la terminal **go run main.go** o **go run .** para iniciar la página web que recibirá paquetes en el
   puerto 8080 de tu computadora. 
3. Dirígete a un navegador web y pon en el buscador **localhost:8080**
4. Ya podés acceder a las funciones actuales.

## Funcionalidades Actuales
### Servicio de página principal
Se sirve un [index.html](./static/index.html) que organiza una página de bienvenida con un **Título**, **Logos** y **Diseños personalizados** donde encontramos información de lo que podremos hacer dentro de la página y una sección para darnos de alta como usuario si no somos clientes ya, y una seccion para ingresar a nuestra cuenta si ya somos clientes.

### Inicio de sesión y Registro
Debemos primero registrar un mail junto con un alias y contraseña en la página para darnos de alta como usuario. **Tanto el alias como el mail deben ser únicos** dentro del registro de la página, caso contrario se lo notará con un aviso al usuario para que pueda cambiarlo las veces necesarias (proximámente en la implementación de una base de datos funcional)
Luego de ingresar nuestros datos correctos nos iremos redirigidos a una sección de
[bienvenida](./static/bienvenida.html) donde se mostrarán nuestros datos recién ingresados notando que el codigo en Go
los reconoció correctamente.

### Ruta inválida
En el caso de que el usuario ingrese una url no reconocida en el código de la página, se servirá un [ruta_invalida.html](./static/ruta_invalida.html) que le indique un error 404 significando que esa sección no se encuentra en la página, con la posibilidad de poder volver al inicio.
