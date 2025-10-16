function displayLogin(){
      document.getElementById("register-form").style.display = "none";
      document.getElementById("login-form").style.display = "block";
}
function displayRegister(){
      document.getElementById("login-form").style.display = "none";
      document.getElementById("register-form").style.display = "block";
}

function displayNewPass(){
      if (!document.getElementById("new-password").classList.contains('new-pass-unactive')){
            document.getElementById("new-password").classList.add('new-pass-unactive');
            document.getElementById("new-pass").disabled = true;
            document.getElementById("password").required = true;
      }else{
            document.getElementById("new-password").classList.remove('new-pass-unactive');
            document.getElementById("new-pass").disabled = false;
            document.getElementById("password").required = false;
      }
}

window.onload = function() {
      const urlParams = new URLSearchParams(window.location.search);
      const error = urlParams.get('error');
 
      switch (error) {
            case 'alias_not_found':
                  alert("El Alias ingresado no se encuentra registrado en la base"); 
                  break
            case 'password_incorrect':
                  alert("La contraseña ingresada es incorrecta")
                  break
            case 'deposit_ok':
                  alert("El depósito se realizó correctamente")
                  break
            case 'alias_usado':
                  alert("El alias ingresado ya se encuentra registrado en la base.")
                  break
            case 'not_enough_balance':
                  alert("No tienes fondos suficientes para realizar la operación.")
                  break
      }
};
// Función única para mostrar el formulario seleccionado y ocultar los demás
function mostrarFormulario(formId) {
const forms = document.querySelectorAll('.operation-form');
forms.forEach(form => {
      if (form.id === formId) {
            // Si es el formulario que queremos mostrar, le añadimos la clase
            form.classList.add('form-visible');
      } else {
            // A todos los demás, se la quitamos
            form.classList.remove('form-visible');
      }
});
}

function recargarPagina(){
      location.reload();
}
