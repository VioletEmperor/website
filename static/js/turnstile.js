let turnstileToken = null;

function onTurnstileCallback(token) {
    turnstileToken = token;
    // Enable submit button when Turnstile is completed
    document.querySelector('button[type="submit"]').disabled = false;
}

// Disable submit button initially
document.addEventListener('DOMContentLoaded', function() {
    document.querySelector('button[type="submit"]').disabled = true;
});

// Add Turnstile token to form submission
document.addEventListener('htmx:configRequest', function(evt) {
    if (evt.detail.path === '/contact' && turnstileToken) {
        evt.detail.parameters['cf-turnstile-response'] = turnstileToken;
    }
});