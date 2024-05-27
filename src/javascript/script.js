function toggleForm(formId) {
    var registerForm = document.getElementById('registerForm');
    var loginForm = document.getElementById('loginForm');
    if (formId === 'loginForm') {
        registerForm.style.display = 'none';
        loginForm.style.display = 'block';
    } else {
        loginForm.style.display = 'none';
        registerForm.style.display = 'block';
    }
}

// si il y a /login dans l'url, on affiche une popup de connexion
if (window.location.pathname === '/login') {
    toggleForm('loginForm');
} else if (window.location.pathname === '/register') {
    toggleForm('registerForm');
}