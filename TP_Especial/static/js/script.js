function displayLogin(){
      document.getElementById("register-form").style.display = "none";
      document.getElementById("login-form").style.display = "block";
}
function displayRegister(){
      document.getElementById("login-form").style.display = "none";
      document.getElementById("register-form").style.display = "block";
}
window.onload = function() {
      const urlParams = new URLSearchParams(window.location.search);
      const error = urlParams.get('error');
 
      if (error) {
            if (error === 'alias_not_found') {
                  alert("El Alias ingresado no se encuentra registrado en la base"); 
            } else if (error === 'password_incorrect') {
                  alert("La contraseña ingresada es incorrecta")
            } else if (error === 'deposit_ok'){
                  alert("El depósito se realizó correctamente")
            }
 
      }
};
