# CandyPay

## Integrantes
- Ignacio Agustin Perna

## Dependencias
Este proyecto necesita de las siguientes dependencias y sus versiones especificadas para funcionar correctamente:
- **Docker 28.5.1 o superior** (make install-docker)
- **Go 1.24.9 o superior** (make install-go)
- **templ v0.3.960 o superior** (make install-templ)
- **sqlc v1.30.0 o superior** (make install-sqlc)
- **hurl 4.2.0 o superior** (make install-hurl)

## Descripción
Este proyecto consta de una billetera virtual que tendrá todas las funcionalidades principales de la misma:
- Creación de cuenta con formularios de registro e inicios de sesión posteriores
- Depositar y retirar dinero
- Transferr a cualquier usuario con un alias registrado
- Pedir dinero a otro usuario
- Pedir préstamos con fechas de pago
- Invertir dinero y generar rendimientos

## Instrucciones de Uso
1. Abrir una terminal localizada dentro de la carpeta donde se encuentra este archivo README.

2. Ejecutar el comando **make run** que hará las siguientes acciones:
      * Baja del contenedor Docker correspondiente borrando todos los datos de la base.
      * Levantar el contenedor con la base.
      * Migrar todos los cambios hechos hasta el momento para contener la ultima version de la base usando golang-migrate
      * Generar código sqlc para poder conectar nuestra base de datos con funciones Go.
      * Generar código Go correspondiente a los templates html usados.
      * Ejecutar todos los archivos Go y escuchar en el puerto 8080 ** *EN SEGUNDO PLANO*  ** para poder ejecutar los
      test en el mismo comando.
      * Insertar datos de prueba con el archivo inserts.hurl donde se insertan 3 usuarios junto con 3 pedidos de dinero
        hechos hacia el usuario alias1.

3. Dirigirte hacia el navegador web y abrir localhost:8080.

4. Para finalizar la ejecución de la aplicacion y dar de baja el contenedor docker junto con la base de datos debe
   ejecutar el comando **make down** 

5. **Aclaración:** Si usted ejecuta esta aplicación en un sistema operativo Windows puede que el comando para cerrar el
   puerto 8080 no le funcione porque se realiza con comandos Bash. En ese caso, debe cerrar el puerto por su cuenta.

## Estructura de la Base de Datos
En primer momento la base se organiza en dos tablas principales:

- Users: donde guardamos información básica del usuario que servirá de identificación en todos los momentos de ingreso a la plataforma, constatando los datos ingresados con los guardados en la Base para su autenticación.

- Accounts: donde se guarda toda la información monetaria del usuario, vinculado con sus datos a través del alias como clave primaria.

- Money_requests: donde se almacenan todos los pedidos de dinero que hacen los usuarios con un alias origen y otro destino, junto con el monto pedido y un mensaje.

- Triggers: se tienen por el momento dos triggers que se activan, cuando el usuario se registra en la base se crea automaticamente una cuenta con los valores por defecto y el alias registrado, y cuando se borra un usuario del sistema, se borra automaticamente la cuenta asociada que tenia ese usuario.

## Funcionalidades Actuales
### Servicio de página principal
Se sirve un [index.html](./static/index.html) que organiza una página de bienvenida con un **Título**, **Logos** y **Diseños personalizados** donde encontramos información de lo que podremos hacer dentro de la página y una sección para darnos de alta como usuario si no somos clientes ya, y una seccion para ingresar a nuestra cuenta si ya somos clientes.

### Inicio de sesión y Registro
Debemos primero registrar un mail junto con un alias y contraseña en la página para darnos de alta como usuario en la base de datos. **Tanto el alias como el mail deben ser únicos** dentro del registro de la página, caso contrario se lo
notará con un aviso al usuario. Luego de ingresar nuestros datos correctos nos iremos redirigidos a una sección de [bienvenida](./static/bienvenida.html) donde se mostrarán nuestros datos recién ingresados junto con nuestro balance de cuenta actual y el ultimo movimiento que hicimos.

### Depósitos
En esta función el usuario tiene ingresado por defecto su propio alias y solo debe colocar el monto a depositar en su cuenta, donde una función en Go parsea
esos datos y verifica con la base que el alias ingresado este registrado y el monto sea un valor correcto, continuando con la actualización del balance del usuario. Vamos a poder observar el cambio en el balance del usuario y el ultimo tipo de movimiento hecho, gracias a una redirección que hacemos desde cada método funcional en Go devuelta a la página de bienvenida y pasando parámetros a través de la URL donde el Handler que se encarga de servir este html parsea los datos..

### Retiros
En esta función el usuario tiene ingresado por defecto su propio alias y solo debe colocar el monto a retirar en su cuenta, donde una función en Go parsea esos datos y verifica con la base que el alias ingresado este registrado y el monto sea un valor correcto, continuando con el descuento en el balance del usuario, que podemos observar automaticamente sin hacer nada.

### Transferencias
En esta función el usuario tiene ingresado por defecto su propio alias y debe colocar el alias al que quiera transferir y el monto, donde una función en Go parsea
esos datos y verifica con la base que ambos alias ingresados esten registrados y el monto sea un valor correcto, continuando con el descuento en el balance de la cuenta origen y el deposito en el balance de la cuenta destino. El correcto
funcionamiento lo vemos comprobando que el balance del usuario disminuyó si todo salió bien y haciendo un login en la cuenta a la que transferimos para verificar su nuevo balance.

### Pedidos de dinero 
En esta funcion el usuario puede solicitar un monto de dinero cualquiera (hasta $99.999.999,99) a cualquier alias registrado en la plataforma, junto con un mensaje. Luego en el mismo Home, el usuario que recibe el pedido podra visualizar con el boton "Pedidos de dinero" una tabla con toda la información necesaria de los pedidos que tenga, comoel usuario que lo solicita, el monto, el mensaje y la fecha de la solicitud. 

### Nueva contraseña
El usuario podrá desde la pagina de Login cambiar su contraseña, **por ahora sin verificacion de identidad**.

### Ruta inválida
En el caso de que el usuario ingrese una url no reconocida en el código de la página, se servirá un [ruta_invalida.html](./static/ruta_invalida.html) que le indique un error 404 significando que esa sección no se encuentra en la página, con la posibilidad de poder volver al inicio.

### Visualización de pedidos dinero
Tenemos dos botones en el home donde podemos ver una tabla donde figuran los pedidos de dinero que le llegaron al usuario y los que realizo el mismo, junto con toda la informacion necesaria y un boton de descarte que elimina esa peticion en particular.

### Ordenamiento de pedidos de dinero
Al visualizar la tabla de pedidos de dinero, por defecto se encuentra ordenada segun el monto de los pedidos en forma descendente, con su flecha que lo indica pero si clickeamos las columnas de **Pedidor** o **Hacia**, **Monto** y **Dia** vamos a poder re ordenar la tabla segun estas categorias en formato ascendente descendente por monto u orden alfabético, de manera dinámica donde por el servidor viaja el pedido de ordenamiento a través de un fetch a un endpoint y la base de datos devuelve la tabla ordenada, la cual el servidor inserta en el html sin alterar su entorno.

## TP5
### Componentes templ
#### Tabla de pedidos de dinero
En este proyecto los componentes templ existentes están creados para las cosas que el usuario no tiene que ver obligatoriamente como la tabla de pedidos de dinero, sino que se carga dinámicamente en el momento que se presiona el botón **Pedidos recibidos** y **Pedidos realizados**, donde al hacerlo, a través de **HTMX** es envía un request *GET* al servidor que se encarga de renderizar la tabla haciendo previamente una consulta a la base de datos para obtener de la tabla MoneyRequests todos los pedidos del usuario con el alias correspondiente pasado como parámetro en la URL. Luego se indica con htmx que la respuesta del servidor se coloque internamente en un container oculto por defecto donde luego se le cambia la visibilidad.

#### Ordenamiento de tabla
Dentro del templ de la tabla de pedidos dinero que se renderiza en el /home, las columnas **Pedidor Hacia Monto y Dia** son clickeables y realizan un pedido HX-GET al servidor con la información del orden actual de la tabla (ascendente o descendente) y la columna que lo solicitó para que el servidor renderice nuevamente la tabla pero ordenada como quiso el usuario y calculando el próximo orden. Para esto el pedido de ordenamiento viaja hacia la base de datos que se encarga de devolver la tabla ordenada. Además, cada columna tiene una flecha que le indica visualmente al usuario el orden actual de la tabla. Por defecto se encuentra ordenada por monto descendente.

### Obtención de datos desde base de datos y renderizado en layout / Manejo de formulario
Para este apartado, lo que ya tenia hecho desde las primeras entregas, es que el usuario al loguearse o registrarse en la página desde la ruta principal **/** con los datos correctos, se lo redirecciona a al **/home** donde se ejecuta el handler Showhome que toma los parametros de la URL como el alias y consulta a la base de datos por su balance actual y ultimo movimiento, para luego parsear estos datos e introducirlos en el html de bienvenida.html que se sirve para esta ruta.
