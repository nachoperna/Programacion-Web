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

### Ruta inválida
En el caso de que el usuario ingrese una url no reconocida en el código de la página, se servirá un [ruta_invalida.html](./static/ruta_invalida.html) que le indique un error 404 significando que esa sección no se encuentra en la página, con la posibilidad de poder volver al inicio.

### Testing de Base de Datos
Es posible ejecutar un codigo de Testing implementando el comando **go test -v** en la terminal ubicada en la carpeta
actual y podemos revisar en el codigo [db_test.go](./db_test.go) como hacemos una conexión a la base de datos para
ejecutar operaciones básicas donde registramos usuarios en nuestra plataforma, actualizamos sus datos, los listamos, y
borramos con éxito, se hacen depósitos, retiros y transferencias con éxito.

### Lógica API REST
Para probar la lógica de negocio de la app podemos ejecutar un script de tests con la herramienta HURL donde tenemos una
lista de peticiones a realizar hacia los endpoints **/users** y **/users/**  en el archivo **requests.hurl** y ejecutando el comando **hurl --test requests.hurl**  donde podemos
aplicar reglas generales a todos los usuarios o en /users o especificar algun usuario en particular para la accion del
metodo que especificado.

####  POST http://localhost:8080/users
- La funcion que maneja esta peticion decodifica todos los objetos json en una estructura correspondiente a este tipo de
  datos.
- Se debe especificar obligatoriamente un **Content-type: application/json** seguido de un arreglo de objetos JSON donde
  se especifiquen los usuarios que se quieran registrar en el sistema. Por ejemplo:
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

- En esta peticion son obligatorios los 4 campos de informacion de cada usuario para el correcto registro en el sistema.
- La funcion en Go que maneja la peticion devuelve un arreglo de objetos JSON donde para no eliminar la operacion entera
  de la peticiion, agrega en un arreglo de "creados" los usuarios que fueron insertados correctamente, y en un arreglo
de "fallidos" aquellos que tuvieron algun problema y no pudieron insertarse en el sistema, indicando al lado de la clave
primaria "alias" el tipo de error ocurrido.
- Además dentro de la funcion que maneja la peticion, solo devolvemos el estado 201 si todos los usuarios ingresados
pertenecen a la lista de creados y no hay ningun fallido, caso contrario obtenemos un StatusBadRequest **400** 
- Dentro de la peticion una vez que se obtiene la respuesta, nos fijamos que el codigo de estado obtenido se corresponda
  a 201 (Created) y que el alias del usuario que recibimos se corresponda con el que ingresamos .

#### DELETE http://localhost:8080/users
- Se le indica a la app que se quiere borrar todos los usuarios registrados en la página.
- La funcion en Go devuelve un estado **404** si hubo algun error eliminando usuarios o un estado **200** si fue
exitoso.
- En la peticion exigimos que para saber que todo salio bien el codigo http devuelto sea 200.

            DELETE http://localhost:8080/users
            HTTP 200

#### GET http://localhost:8080/users 
- Se indica que se quiere listar todos los usuarios registrados en el sistema.
- Una funcion en nuestro codigo Go hace una invocacion al metodo sqlc que se encarga de realizar la query que retorna
todos los usuarios, y codifica esa salida en un formato json que envia como respuesta al usuario que realizo la
peticion.

            GET http://localhost:8080/users
            HTTP 200

#### GET http://localhost:8080/users/?alias=alias3
- En este tipo de peticiones donde debemos indicar un usuario o parametro en particular, decidi implementarlo a través
de parametros en la URL.
- Las peticiones /users/ usan otro handler diferente al anterior para manejar estas consultas singulares, donde cada
funcion correspondiente a cada metodo extrae los parametros necesarios de la URL que introdujo el usuario.
- En este caso es obligatorio introducir un alias a buscar porque es la clave primaria que tiene cada usuario registrado
- Si el usuario a buscar no se encuentra, se retorna un estado **404 NotFound**
- Si el usuario se encuentra en el sistema, se establece un estado **200 Ok** y se codifica la informacion del usuario a
  un objeto json.

#### PUT http://localhost:8080/users/?alias=alias1
- Se le indica al sistema que queremos actualizar los datos del usuario correspondiente al alias.
- Se debe especificar un objeto json en el cuerpo de la peticion donde se marca como clave el atributo a actualizar y
como valor el nuevo valor actualizado.
- La funcion en el codigo que maneja la peticion invoca al metodo sqlc encargado del update en la base de datos de la
tabla Users, donde la query controla *en todos los atributos que se puedan actualizar* si el campo del atributo json de
la peticion es nulo, entonces deja el valor actual del usuario, o si no es nulo, actualiza el atributo del usuario con
este nuevo valor.
- Se retorna solamente un codigo **200 Ok** si todo salio bien o un **400 BadRequest** si hubo cualquier error
actualizando al usuario.

            PUT http://localhost:8080/users/?alias=alias4
            Content-type: application/json
            {
                  "name": "persona1nueva"
            }
            HTTP 200


#### DELETE http://localhost:8080/users/?alias=alias1 
- Se indica que se quiere eliminar el usuario con clave primaria indicada, y la funcion correspondiente en Go invoca al
  metodo de sqlc correspondiente a la eliminacion de la tabla User, que devuelve el usuario eliminado y verificamos si
ese usuario existia previamente, retornando un estado **204 NoContent** si todo salio bien o un **404 BadRequest** si
hubo error al eliminar.

            DELETE http://localhost:8080/users/?alias=alias1 
            HTTP 204

## Estructura de la Base de Datos
En primer momento la base se organiza en dos tablas principales:

- Users: donde guardamos información básica del usuario que servirá de identificación en todos los momentos de ingreso a
  la plataforma, constatando los datos ingresados con los guardados en la Base para su autenticación.

- Accounts: donde se guarda toda la información monetaria del usuario, vinculado con sus datos a través del alias como
clave primaria.

- Triggers: se tienen por el momento dos triggers que se activan, cuando el usuario se registra en la base se crea
automaticamente una cuenta con los valores por defecto y el alias registrado, y cuando se borra un usuario del sistema,
se borra automaticamente la cuenta asociada que tenia ese usuario.

