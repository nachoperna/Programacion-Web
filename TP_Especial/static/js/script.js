// let sort_order = 'desc';
//
// document.addEventListener("DOMContentLoaded", () => {
//       document.getElementById('tabla-pedidos-container').addEventListener("click", (event) => {
//             // Si se realiza un cick dentro de esta clase sabemos que la tabla-pedidos esta desplegada
//             const link = event.target.closest('a.descartar-pedido-recibido');
//             // Si se encontró el enlace, ejecuta la lógica
//             if (link) {
//                   event.preventDefault(); // Evita que el enlace recargue la página
//                   const alias_from = link.dataset.aliasFrom; // Obtiene el alias desde el data-attribute
//                   const alias_to = link.dataset.aliasTo; // Obtiene el alias desde el data-attribute
//                   fetch(`/deleteRequestsTo?from_alias=${alias_from}&to_alias=${alias_to}`)
//                         .then(Response => {
//                               if (!Response.ok) {
//                                     throw new Error('Hubo un error en la respuesta del servidor');
//                               }
//                               return Response.text();
//                         })
//                         .then(tablaPedidos => {
//                               document.getElementById('tabla-pedidos-container').innerHTML = tablaPedidos;
//                               document.getElementById("tabla-pedidos-container").classList.remove('oculto');
//                         })
//                         .catch(error => {
//                               console.error('Error obteniendo los datos de los pedidos; ', error)
//                               document.getElementById('tabla-pedidos-container').innerHTML = "<th><tr><td>Error obteniendo datos</td></tr></th>"
//                               document.getElementById("tabla-pedidos-container").classList.remove('oculto');
//                         })
//                   return;
//             }
//             // Obtenemos el elemento si el usuario hizo cick en una columna ordenable
//             const th_to = event.target.closest('th.ordenable-to');
//             const th_from = event.target.closest('th.ordenable-from');
//             if (th_to || th_from) {
//                   event.preventDefault();
//                   alias = new URLSearchParams(window.location.search).get('alias')
//                   let ordenables = null;
//                   if (th_to) {
//                         ordenables = document.querySelectorAll('th.ordenable-to');
//                   } else {
//                         ordenables = document.querySelectorAll('th.ordenable-from');
//                   }
//                   // Limpiamos la flecha de odenamiento de todas las columnas ordenables
//                   ordenables.forEach(ordenable => {
//                         ordenable.classList.remove('sort-asc');
//                         ordenable.classList.remove('sort-desc');
//                   })
//
//                   let sort_by = "";
//                   if (th_to) {
//                         sort_by = th_to.dataset.sortBy; // Obtenemos el tipo de ordenamiento del atributo del header
//                   } else {
//                         sort_by = th_from.dataset.sortBy; // Obtenemos el tipo de ordenamiento del atributo del header
//                   }
//                   let nueva_clase = '';
//                   if (sort_order === 'asc') {
//                         sort_order = 'desc';
//                         nueva_clase = 'sort-desc';
//                   } else {
//                         sort_order = 'asc';
//                         nueva_clase = 'sort-asc';
//                   }
//
//                   if (th_to) {
//                         fetch(`/listRequestsTo?to_alias=${alias}&sort_by=${sort_by}&sort_order=${sort_order}`)
//                               .then(Response => {
//                                     if (!Response.ok) {
//                                           throw new Error('Hubo un error en la respuesta del servidor');
//                                     }
//                                     return Response.text();
//                               })
//                               .then(tablaPedidos => {
//                                     document.getElementById('tabla-pedidos-container').innerHTML = tablaPedidos;
//                                     const nuevoheader = document.querySelector(`th[data-sort-by=${sort_by}]`); // Volvemos a obtener la columna que cickeo el usuario
//                                     document.querySelector(`th[data-sort-by=amount]`).classList.remove('sort-desc');
//                                     if (nuevoheader) {
//                                           nuevoheader.classList.add(nueva_clase); // Le asignamos la clase con el icono correspondiente
//                                     }
//                               })
//                               .catch(error => {
//                                     console.error('Error obteniendo los datos de los pedidos; ', error)
//                                     document.getElementById('tabla-pedidos-container').innerHTML = "<th><tr><td>Error obteniendo datos</td></tr></th>"
//                               })
//                   } else {
//                         fetch(`/listRequestsFrom?from_alias=${alias}&sort_by=${sort_by}&sort_order=${sort_order}`)
//                               .then(Response => {
//                                     if (!Response.ok) {
//                                           throw new Error('Hubo un error en la respuesta del servidor');
//                                     }
//                                     return Response.text();
//                               })
//                               .then(tablaPedidos => {
//                                     document.getElementById('tabla-pedidos-container').innerHTML = tablaPedidos;
//                                     const nuevoheader = document.querySelector(`th[data-sort-by=${sort_by}]`); // Volvemos a obtener la columna que cickeo el usuario
//                                     document.querySelector(`th[data-sort-by=amount]`).classList.remove('sort-desc');
//                                     if (nuevoheader) {
//                                           nuevoheader.classList.add(nueva_clase); // Le asignamos la clase con el icono correspondiente
//                                     }
//                               })
//                               .catch(error => {
//                                     console.error('Error obteniendo los datos de los pedidos; ', error)
//                                     document.getElementById('tabla-pedidos-container').innerHTML = "<th><tr><td>Error obteniendo datos</td></tr></th>"
//                               })
//                   }
//                   return;
//             }
//       })
// })
//
//
