document.addEventListener('DOMContentLoaded', function() {
    const receiverForm = document.getElementById('receiverForm');
    const availableMeals = document.getElementById('availableMeals');
    const mealsList = document.getElementById('mealsList');
    const scanQRButton = document.getElementById('scanQR');
    const feedbackForm = document.getElementById('feedbackForm');
    let scanner = null;

    receiverForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        console.log('Form submission started');

        const formData = new FormData(receiverForm);
        const dietarySelect = document.getElementById('dietaryRestrictions');
        const selectedDietary = Array.from(dietarySelect.selectedOptions).map(option => option.value);

        const data = {
            name: formData.get('name'),
            phone: formData.get('phone'),
            location: formData.get('location'),
            familySize: parseInt(formData.get('familySize'), 10),
            dietaryRestrictions: selectedDietary
        };

        // Validate the data
        if (!data.name || !data.phone || !data.location || !data.familySize || data.familySize <= 0) {
            showNotification('Please fill in all required fields and ensure family size is positive', 'error');
            return;
        }

        console.log('Sending registration data:', data);

        try {
            const response = await fetch('/api/receivers', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });

            const responseText = await response.text();
            console.log('Raw server response:', responseText);

            let result;
            try {
                result = JSON.parse(responseText);
            } catch (parseError) {
                console.error('Error parsing server response:', parseError);
                throw new Error('Invalid server response');
            }

            if (!response.ok) {
                throw new Error(result.message || 'Registration failed');
            }

            console.log('Registration successful:', result);
            
            // Store receiver data in localStorage
            data.id = result.receiverId;
            localStorage.setItem('receiverData', JSON.stringify(data));
            
            // Show success message and redirect
            showNotification('Registration successful! Redirecting to meal finder...', 'success');
            
            // Redirect to meal finder after a short delay
            setTimeout(() => {
                window.location.href = '/meal-finder';
            }, 1500);

        } catch (error) {
            console.error('Registration error:', error);
            showNotification(error.message || 'Failed to register. Please try again.', 'error');
        }
    });

    async function fetchAvailableMeals(location) {
        try {
            console.log('Fetching meals for location:', location);
            const response = await fetch(`/api/meals?location=${encodeURIComponent(location)}`);
            
            if (!response.ok) {
                const errorText = await response.text();
                console.error('Server error:', errorText);
                throw new Error('Failed to fetch meals');
            }

            const meals = await response.json();
            console.log('Received meals:', meals);

            mealsList.innerHTML = '';
            availableMeals.style.display = 'block';

            if (!Array.isArray(meals)) {
                throw new Error('Server returned unexpected data format');
            }

            if (meals.length === 0) {
                mealsList.innerHTML = `
                    <div class="info-message">
                        <p>No meals are currently available in ${location}.</p>
                        <p>Please check back later or try a different location.</p>
                    </div>`;
                return;
            }

            const mealsHtml = meals.map(meal => {
                const expiryDate = new Date(meal.expiryDate);
                return `
                    <div class="card meal-card">
                        <h3>${meal.type}</h3>
                        <p><strong>Available:</strong> ${meal.quantity} servings</p>
                        <p><strong>Location:</strong> ${meal.location}</p>
                        <p><strong>Expires:</strong> ${expiryDate.toLocaleString()}</p>
                    </div>`;
            }).join('');

            mealsList.innerHTML = mealsHtml;

        } catch (error) {
            console.error('Error fetching meals:', error);
            mealsList.innerHTML = `
                <div class="error-message">
                    <p>Unable to load meals at this time.</p>
                    <p>Error: ${error.message}</p>
                    <p>Please try again later.</p>
                </div>`;
            availableMeals.style.display = 'block';
        }
    }

    scanQRButton.addEventListener('click', function() {
        if (scanner) {
            scanner.stop();
        }

        const scannerContainer = document.createElement('div');
        scannerContainer.innerHTML = '<video id="preview"></video>';
        document.body.appendChild(scannerContainer);

        scanner = new Instascan.Scanner({ video: document.getElementById('preview') });
        
        scanner.addListener('scan', async function(donationId) {
            try {
                const response = await fetch('/api/verify-meal', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ donationId })
                });

                if (!response.ok) {
                    throw new Error('Failed to verify meal');
                }

                showNotification('Meal verified successfully!', 'success');
                scanner.stop();
                scannerContainer.remove();
                feedbackForm.style.display = 'block';

            } catch (error) {
                console.error('Error:', error);
                showNotification('Failed to verify meal. Please try again.', 'error');
            }
        });

        Instascan.Camera.getCameras().then(function(cameras) {
            if (cameras.length > 0) {
                scanner.start(cameras[0]);
            } else {
                showNotification('No cameras found on your device.', 'error');
            }
        });
    });

    document.getElementById('mealFeedback').addEventListener('submit', async function(e) {
        e.preventDefault();

        const formData = new FormData(this);
        const data = Object.fromEntries(formData.entries());

        try {
            const response = await fetch('/api/feedback', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });

            if (!response.ok) {
                throw new Error('Failed to submit feedback');
            }

            showNotification('Thank you for your feedback!', 'success');
            this.reset();
            feedbackForm.style.display = 'none';

        } catch (error) {
            console.error('Error:', error);
            showNotification('Failed to submit feedback. Please try again.', 'error');
        }
    });

    function showNotification(message, type) {
        const notification = document.createElement('div');
        notification.className = `notification notification-${type}`;
        notification.textContent = message;
        
        // Remove any existing notifications
        const existingNotifications = document.querySelectorAll('.notification');
        existingNotifications.forEach(n => n.remove());
        
        document.body.appendChild(notification);
        
        setTimeout(() => {
            notification.remove();
        }, 5000);
    }
});
