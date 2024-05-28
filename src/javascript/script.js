function toggleForm(formId) {
    var registerForm = document.getElementById('registerForm');
    var loginForm = document.getElementById('loginForm');
    var forgotForm = document.getElementById('forgotForm');
    var resetForm = document.getElementById('resetForm');

    switch (formId) {
        case 'registerForm':
            registerForm.style.display = 'block';
            loginForm.style.display = 'none';
            forgotForm.style.display = 'none';
            resetForm.style.display = 'none';
            break;
        case 'loginForm':
            registerForm.style.display = 'none';
            loginForm.style.display = 'block';
            forgotForm.style.display = 'none';
            resetForm.style.display = 'none';
            break;
        case 'forgotForm':
            registerForm.style.display = 'none';
            loginForm.style.display = 'none';
            forgotForm.style.display = 'block';
            resetForm.style.display = 'none';
            break;
        case 'resetForm':
            registerForm.style.display = 'none';
            loginForm.style.display = 'none';
            forgotForm.style.display = 'none';
            resetForm.style.display = 'block';
            break;
    }
    
}

if (window.location.pathname === '/register') {
    toggleForm('registerForm');
}

if (window.location.pathname === '/login') {
    toggleForm('loginForm');
}

if (window.location.pathname === '/forgot') {
    toggleForm('forgotForm');
}

if (window.location.pathname === '/reset') {
    toggleForm('resetForm');
}