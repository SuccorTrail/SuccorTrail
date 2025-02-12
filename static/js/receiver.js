document.addEventListener('DOMContentLoaded', function() {
    const receiverForm = document.getElementById('receiverForm');
    const availableMeals = document.getElementById('availableMeals');
    const mealsList = document.getElementById('mealsList');
    const scanQRButton = document.getElementById('scanQR');
    const feedbackForm = document.getElementById('feedbackForm');
    let scanner = null;

    receiverForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const formData = new FormData(receiverForm);
        const data = Object.fromEntries(formData.entries());

        try {
            const response = await fetch('/api/receivers', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });

            if (!response.ok) {
                throw new Error('Network response was not ok');
            }

            const result = await response.json();
            
            // Show available meals
            await fetchAvailableMeals(data.location);
            
            // Show success message
            showNotification('Registration successful! You will be notified when meals are available.', 'success');

        } catch (error) {
            console.error('Error:', error);
            showNotification('Failed to register. Please try again.', 'error');
        }
    });

    async function fetchAvailableMeals(location) {
        try {
            const response = await fetch(`/api/meals?location=${encodeURIComponent(location)}`);
            const meals = await response.json();

            if (meals.length > 0) {
                mealsList.innerHTML = '';
                meals.forEach(meal => {
                    const mealCard = document.createElement('div');
                    mealCard.className = 'card';
                    mealCard.innerHTML = `
                        <h3>${meal.type}</h3>
                        <p>Available: ${meal.quantity}</p>
                        <p>Pickup Location: ${meal.location}</p>
                        <p>Expires: ${new Date(meal.expiryDate).toLocaleString()}</p>
                    `;
                    mealsList.appendChild(mealCard);
                });
                availableMeals.style.display = 'block';
            }
        } catch (error) {
            console.error('Error fetching meals:', error);
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
        
        document.body.appendChild(notification);
        
        setTimeout(() => {
            notification.remove();
        }, 3000);
    }
});
