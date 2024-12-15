const form = document.getElementById('form');
const username = document.getElementById('username');
const email = document.getElementById('email');
const password = document.getElementById('password');
const password2 = document.getElementById('password2');
const yourUrl = "{Your url here}"

console.log("Javascript script started.")

form.addEventListener('submit', async e => {
    e.preventDefault();
    validateInputs();

    // Check if all inputs are valid before proceeding
    const isValid =
        username.parentElement.classList.contains('success') &&
        email.parentElement.classList.contains('success') &&
        password.parentElement.classList.contains('success') &&
        password2.parentElement.classList.contains('success');

    if (isValid) {
        // Prepare the form data
        const formData = {
            username: username.value.trim(),
            email: email.value.trim(),
            password: password.value.trim(),
        };

        try {
            // Send the POST request
            const response = await fetch(yourUrl + '/register/newuser', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData),
            });

            if (response.ok) {
                // Handle success
                console.log('Registration successful!');
                window.location.href = "/login.html"
            } else {
                // Handle errors from the server
                const errorData = await response.json();
                console.log('Server error:', errorData.message || 'An error occurred');
            }
        } catch (error) {
            // Handle network or other unexpected errors
            console.error('Request failed:', error);
        }
    }
});

const setError = (element, message) => {
    const inputControl = element.parentElement;
    const errorDisplay = inputControl.querySelector('.error');

    errorDisplay.innerText = message;
    inputControl.classList.add('error');
    inputControl.classList.remove('success')
}

const setSuccess = element => {
    const inputControl = element.parentElement;
    const errorDisplay = inputControl.querySelector('.error');

    errorDisplay.innerText = '';
    inputControl.classList.add('success');
    inputControl.classList.remove('error');
};

const isValidEmail = email => {
    const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(String(email).toLowerCase());
}

const validateInputs = () => {
    const usernameValue = username.value.trim();
    const emailValue = email.value.trim();
    const passwordValue = password.value.trim();
    const password2Value = password2.value.trim();

    if(usernameValue === '') {
        setError(username, 'Username is required');
    } else if (usernameValue.length > 16) {
        setError(username, 'Username must not be longer than 15 characters.')
    } else {
        setSuccess(username);
    }

    if(emailValue === '') {
        setError(email, 'Email is required');
    } else if (emailValue.length > 128) {
        setError(email, 'Please use an email shorter than 128 characters.')
    } else if (!isValidEmail(emailValue)) {
        setError(email, 'Provide a valid email address');
    } else {
        setSuccess(email);
    }

    if(passwordValue === '') {
        setError(password, 'Password is required');
    } else if (passwordValue.length < 8 ) {
        setError(password, 'Password must be at least 8 characters.')
    } else if (passwordValue.length > 128 ) {
        setError(password, 'Password must be less than 128 characters.')
    } else {
        setSuccess(password);
    }

    if(password2Value === '') {
        setError(password2, 'Please confirm your password');
    } else if (password2Value !== passwordValue) {
        setError(password2, "Passwords doesn't match");
    } else {
        setSuccess(password2);
    }

};
