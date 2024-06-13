function toggleForm(formId) {
    var registerForm = document.getElementById('registerForm');
    var loginForm = document.getElementById('loginForm');

    switch (formId) {
        case 'registerForm':
            registerForm.style.display = 'block';
            loginForm.style.display = 'none';
            break;
        case 'loginForm':
            registerForm.style.display = 'none';
            loginForm.style.display = 'block';
            break;
    }
    
}

if (window.location.pathname === '/register') {
    toggleForm('registerForm');
}

if (window.location.pathname === '/login') {
    toggleForm('loginForm');
}