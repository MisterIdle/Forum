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

toggleForm('loginForm');