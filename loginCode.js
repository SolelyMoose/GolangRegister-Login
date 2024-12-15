const form = document.getElementById('form');
const email = document.getElementById('email');
const password = document.getElementById('password');
const savelogin = document.getElementById("savelogin")
const yourUrl = "{Your url here}"

console.log("Javascript script started.")
form.addEventListener('submit', async e => {
    e.preventDefault();
    validateInputs();

    // Check if all inputs are valid before proceeding
    const isValid =
        email.parentElement.classList.contains('success') &&
        password.parentElement.classList.contains('success');

    if (isValid) {
        // Prepare the form data
        const formData = {
            email: email.value.trim(),
            password: password.value.trim(),
            savelogin: savelogin.checked,
        };

        console.log(formData);
        console.log(JSON.stringify(formData));

        try {
            // Send the POST request
            const response = await fetch(yourUrl + '/login/newrequest', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData),
            });
        
            // Check if the response is okay
            if (response.ok) {
                const data = await response.json(); // Parse the JSON from the response
                console.log(data.token); // Access the token here
        
                // Handle success
                console.log('Login successful!');
                window.location.href = "../index.html"
                setCorrect(email, "");
                setCorrect(password, "");
                console.log(savelogin.checked);
                
                // Only save token to cookies if the "savelogin" checkbox is checked
                if (savelogin.checked) {
                    document.cookie = `token=${data.token}; path=/; secure; samesite=strict`;
                    console.log("Saved token to cookies!");
                }
            } else {
                // Handle errors from the server
                const errorData = await response.json();
                console.log('Server error:', errorData.message || 'An error occurred');
                setError(email, "Email or password is incorrect.");
                setError(password, "Email or password is incorrect.");
            }
        } catch (error) {
            // Handle network or other unexpected errors
            console.error('Request failed:', error);
        }        
    }
});

const setIncorrect = (element, message) => {
    const inputControl = element.parentElement;
    const errorDisplay = inputControl.querySelector('.error');

    errorDisplay.innerText = message;
    inputControl.classList.add('incorrect');
    inputControl.classList.remove('correct')
}

const setCorrect = element => {
    const inputControl = element.parentElement;
    const errorDisplay = inputControl.querySelector('.error');

    errorDisplay.innerText = '';
    inputControl.classList.add('correct');
    inputControl.classList.remove('incorrect');
};

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
    const emailValue = email.value.trim();
    const passwordValue = password.value.trim();

    if(emailValue === '') {
        setError(email, 'Email is required');
    } else if (emailValue.length > 128) {
        setError(email, 'Incorrect email.')
    } else if (!isValidEmail(emailValue)) {
        setError(email, 'Provide a valid email address.');
    } else {
        setSuccess(email);
    }

    if(passwordValue === '') {
        setError(password, 'Password is required.');
    } else if (passwordValue.length < 8 ) {
        setError(password, 'Incorrect password.')
    } else if (passwordValue.length > 128 ) {
        setError(password, 'Incorrect password.')
    } else {
        setSuccess(password);
    }

};
