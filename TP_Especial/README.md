# CandyPay
## Integrantes
- Ignacio Agustin Perna

## Descripción
Este proyecto consta de una billetera virtual que tendrá todas las funcionalidades principales de la misma:
- Creación de cuenta con formularios de registro e inicios de sesión posteriores
- Depositar y retirar dinero
- Transferr a cualquier usuario con un alias registrado
- Pedir préstamos con fechas de pago
- Invertir dinero y generar rendimientos

## Modo de Uso
1. Abrir una terminal localizada dentro de la carpeta donde se encuentra este archivo README.
2. Ejecutar el comando **make reset** que hará la baja del contenedor de Docker, borrando todos los datos de la base,
   para luego levantarlo y migrar todos los cambios hechos hasta el momento para contener la última versión de la base
con la herramienta golang-migrate, y generar el código sqlc para poder conectar nuestra base de datos con funciones de
Go.
3. Ejecutar el comando **go run main.go** 
4. Dirígete a un navegador web y pon en el buscador **localhost:8080**
5. Ya podés acceder a las funciones actuales.

## Funcionalidades Actuales
### Servicio de página principal
Se sirve un [index.html](./static/index.html) que organiza una página de bienvenida con un **Título**, **Logos** y **Diseños personalizados** donde encontramos información de lo que podremos hacer dentro de la página y una sección para darnos de alta como usuario si no somos clientes ya, y una seccion para ingresar a nuestra cuenta si ya somos clientes.

### Inicio de sesión y Registro
Debemos primero registrar un mail junto con un alias y contraseña en la página para darnos de alta como usuario en la
base de datos. **Tanto el alias como el mail deben ser únicos** dentro del registro de la página, caso contrario se lo
notará con un aviso al usuario. Luego de ingresar nuestros datos correctos nos iremos redirigidos a una sección de
[bienvenida](./static/bienvenida.html) donde se mostrarán nuestros datos recién ingresados junto con nuestro balance de
cuenta actual y el ultimo movimiento que hicimos.

### Depósitos
En esta función el usuario tiene ingresado por defecto su propio alias y solo debe colocar el monto a depositar en su cuenta, donde una función en Go parsea
esos datos y verifica con la base que el alias ingresado este registrado y el monto sea un valor correcto, continuando
con la actualización del balance del usuario. Vamos a poder observar el cambio en
el balance del usuario y el ultimo tipo de movimiento hecho, gracias a una redirección que hacemos desde cada método
funcional en Go devuelta a la página de bienvenida y pasando parámetros a través de la URL donde el Handler que se
encarga de servir este html parsea los datos..

### Retiros
En esta función el usuario tiene ingresado por defecto su propio alias y solo debe colocar el monto a retirar en su cuenta, donde una función en Go parsea
esos datos y verifica con la base que el alias ingresado este registrado y el monto sea un valor correcto, continuando
con el descuento en el balance del usuario, que podemos observar automaticamente sin hacer nada.

### Transferencias
En esta función el usuario tiene ingresado por defecto su propio alias y debe colocar el alias al que quiera transferir y el monto, donde una función en Go parsea
esos datos y verifica con la base que ambos alias ingresados esten registrados y el monto sea un valor correcto, continuando
con el descuento en el balance de la cuenta origen y el deposito en el balance de la cuenta destino. El correcto
funcionamiento lo vemos comprobando que el balance del usuario disminuyó si todo salió
bien y haciendo un login en la cuenta a la que transferimos para verificar su nuevo balance.

### Pedidos de dinero 
En esta funcion el usuario puede solicitar un monto de dinero cualquiera (hasta $99.999.999,99) a cualquier alias
registrado en la plataforma, junto con un mensaje. Luego en el mismo Home, el usuario que recibe el pedido podra
visualizar con el boton "Pedidos de dinero" una tabla con toda la información necesaria de los pedidos que tenga, como
el usuario que lo solicita, el monto, el mensaje y la fecha de la solicitud. 

### Nueva contraseña
El usuario podrá desde la pagina de Login cambiar su contraseña, **por ahora sin verificacion de identidad**.

### Ruta inválida
En el caso de que el usuario ingrese una url no reconocida en el código de la página, se servirá un [ruta_invalida.html](./static/ruta_invalida.html) que le indique un error 404 significando que esa sección no se encuentra en la página, con la posibilidad de poder volver al inicio.

### Testing de Base de Datos
Es posible ejecutar un codigo de Testing implementando el comando **go test -v** en la terminal ubicada en la carpeta
actual y podemos revisar en el codigo [db_test.go](./db_test.go) como hacemos una conexión a la base de datos para
ejecutar operaciones básicas donde registramos usuarios en nuestra plataforma, actualizamos sus datos, los listamos, y
borramos con éxito, se hacen depósitos, retiros y transferencias con éxito.

## Estructura de la Base de Datos
En primer momento la base se organiza en dos tablas principales:

- Users: donde guardamos información básica del usuario que servirá de identificación en todos los momentos de ingreso a
  la plataforma, constatando los datos ingresados con los guardados en la Base para su autenticación.

- Accounts: donde se guarda toda la información monetaria del usuario, vinculado con sus datos a través del alias como
clave primaria.

- Triggers: se tienen por el momento dos triggers que se activan, cuando el usuario se registra en la base se crea
automaticamente una cuenta con los valores por defecto y el alias registrado, y cuando se borra un usuario del sistema,
se borra automaticamente la cuenta asociada que tenia ese usuario.

