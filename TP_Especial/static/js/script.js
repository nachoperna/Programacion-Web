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
      }else{
            document.getElementById("new-password").classList.remove('new-pass-unactive');
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
      }
};
