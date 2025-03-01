document.addEventListener('DOMContentLoaded', function() {
    const signupForm = document.getElementById('signupForm');
    
    if (signupForm) {
        signupForm.addEventListener('submit', function(event) {
            event.preventDefault();
            
            // Get form data
            const email = document.getElementById('email').value;
            const name = document.getElementById('name').value;
            const userType = document.getElementById('userType').value;
            const password = document.getElementById('password')?.value || 'temporary-password-12345';
            
            // Validate form
            if (!email || !name || !userType) {
                showMessage('Please fill in all required fields', 'error');
                return;
            }
            
            // Simple email validation
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(email)) {
                showMessage('Please enter a valid email address', 'error');
                return;
            }
            
            // Submit as JSON instead of FormData
            const userData = {
                name: name,
                email: email,
                userType: userType,
                password: password
            };
            
            console.log("Submitting form data:", {
                name, email, userType, password: "********" // Log for debugging
            });
            
            // Submit form data as JSON
            fetch('/api/auth/signup', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(userData)
            })
            .then(response => {
                console.log("Response status:", response.status);
                return response.json();
            })
            .then(data => {
                console.log("Response data:", data);
                if (data.success) {
                    showMessage(data.message || 'Account created successfully!', 'success');
                } else {
                    showMessage(data.message || 'Error creating account', 'error');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                showMessage('An error occurred. Please try again later.', 'error');
            });
        });
    }
    
    function showMessage(message, type) {
        const messageContainer = document.getElementById('messageContainer');
        if (!messageContainer) return;
        
        messageContainer.innerHTML = `<div class="message ${type}">${message}</div>`;
        
        // Auto-remove message after 5 seconds for success messages
        if (type === 'success') {
            setTimeout(() => {
                messageContainer.innerHTML = '';
            }, 5000);
        }
    }
}); 