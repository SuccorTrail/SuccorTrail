document.addEventListener('DOMContentLoaded', function() {
    // Get URL parameters
    const urlParams = new URLSearchParams(window.location.search);
    const nameParam = urlParams.get('name');
    const emailParam = urlParams.get('email');
    const userTypeParam = urlParams.get('userType');
    
    // Prefill form if parameters exist
    if (nameParam) document.getElementById('name').value = nameParam;
    if (emailParam) document.getElementById('email').value = emailParam;
    if (userTypeParam) document.getElementById('userType').value = userTypeParam;

    const signupForm = document.getElementById('signupForm');s
    
    if (signupForm) {
        signupForm.addEventListener('submit', function(event) {
            event.preventDefault();
            
            // Validate form
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            const name = document.getElementById('name').value;
            const userType = document.getElementById('userType').value;
            
            if (!email || !password || !name || !userType) {
                showMessage('Please fill in all required fields', 'error');
                return;
            }
        });
    }
}); 