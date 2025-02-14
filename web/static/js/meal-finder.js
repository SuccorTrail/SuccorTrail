document.addEventListener('DOMContentLoaded', function() {
    const receiverData = JSON.parse(localStorage.getItem('receiverData'));
    if (!receiverData) {
        window.location.href = '/receiver';
        return;
    }

    // Update profile information
    document.getElementById('receiverName').textContent = receiverData.name;
    document.getElementById('receiverLocation').textContent = receiverData.location;
    document.getElementById('receiverPhone').textContent = receiverData.phone;
    document.getElementById('familySize').textContent = receiverData.familySize;
    document.getElementById('dietaryRestrictions').textContent = 
        receiverData.dietaryRestrictions.length > 0 
            ? receiverData.dietaryRestrictions.join(', ') 
            : 'None';

    // Enhance location description based on location
    function getLocationDescription(location) {
        const locationDescriptions = {
            'urban': 'Urban areas often have diverse meal options and multiple distribution points.',
            'suburban': 'Suburban locations may have community centers and local food banks.',
            'rural': 'Rural areas might have fewer but more personalized meal distribution services.',
            'default': 'Your location helps us connect you with nearby meal resources.'
        };

        // Basic location type detection (very simple, can be expanded)
        const locationLower = location.toLowerCase();
        if (locationLower.includes('city') || locationLower.includes('urban')) return locationDescriptions['urban'];
        if (locationLower.includes('suburb')) return locationDescriptions['suburban'];
        if (locationLower.includes('rural') || locationLower.includes('country')) return locationDescriptions['rural'];
        return locationDescriptions['default'];
    }

    document.getElementById('locationDescription').textContent = 
        getLocationDescription(receiverData.location);

    // Fetch available meals
    async function fetchAvailableMeals(location) {
        const mealsList = document.getElementById('mealsList');
        
        try {
            const response = await fetch(`/api/meals?location=${encodeURIComponent(location)}`);
            
            if (!response.ok) {
                throw new Error(`Server error: ${response.status}`);
            }

            const responseText = await response.text();

            if (!responseText.trim()) {
                mealsList.innerHTML = `
                    <div class="info-message">
                        <i class="fas fa-info-circle"></i>
                        <p>No meals are currently available in ${location}.</p>
                        <p>Our team is working to bring more resources to your area.</p>
                    </div>`;
                return;
            }

            let meals;
            try {
                meals = JSON.parse(responseText);
            } catch (parseError) {
                throw new Error(`Failed to parse server response: ${parseError.message}`);
            }

            if (!Array.isArray(meals) || meals.length === 0) {
                mealsList.innerHTML = `
                    <div class="info-message">
                        <i class="fas fa-utensils"></i>
                        <p>No meals are currently available in ${location}.</p>
                        <p>Check back soon or contact local support.</p>
                    </div>`;
                return;
            }

            const mealsHtml = meals.map(meal => {
                const expiryDate = new Date(meal.expiryDate);
                const daysUntilExpiry = Math.ceil((expiryDate - new Date()) / (1000 * 60 * 60 * 24));

                return `
                    <div class="card meal-card">
                        <h3><i class="fas fa-drumstick-bite"></i> ${meal.type}</h3>
                        <p><strong>Available:</strong> ${meal.quantity} servings</p>
                        <p><strong>Location:</strong> ${meal.location}</p>
                        <p><strong>Expires:</strong> 
                            ${expiryDate.toLocaleString()} 
                            <span class="${daysUntilExpiry <= 1 ? 'text-danger' : ''}">
                                (${daysUntilExpiry} days left)
                            </span>
                        </p>
                        <button class="btn btn-primary request-meal" data-meal-id="${meal.id}">
                            <i class="fas fa-hand-holding-heart"></i> Request Meal
                        </button>
                    </div>`;
            }).join('');

            mealsList.innerHTML = mealsHtml;

            // Add event listeners to request buttons
            document.querySelectorAll('.request-meal').forEach(button => {
                button.addEventListener('click', function() {
                    const mealId = this.dataset.mealId;
                    requestMeal(mealId);
                });
            });

        } catch (error) {
            mealsList.innerHTML = `
                <div class="error-message">
                    <i class="fas fa-exclamation-triangle"></i>
                    <p>Unable to load meals at this time.</p>
                    <p>Error: ${error.message}</p>
                    <p>Please try again later or contact support.</p>
                </div>`;
        }
    }

    // Simulate meal history and stats (would be replaced with actual backend data)
    function updateMealStats() {
        const mealsReceivedCount = document.getElementById('mealsReceivedCount');
        const daysSupportedCount = document.getElementById('daysSupportedCount');
        const mealHistory = document.getElementById('mealHistory');

        // Mock data - replace with actual backend data
        const mockMealsReceived = Math.floor(Math.random() * 10);
        const mockDaysSupported = Math.floor(mockMealsReceived * 1.5);

        mealsReceivedCount.textContent = mockMealsReceived;
        daysSupportedCount.textContent = mockDaysSupported;

        // Mock meal history
        const mockMealHistory = [
            { date: '2025-02-10', type: 'Vegetarian Pasta', location: 'Community Center' },
            { date: '2025-02-05', type: 'Chicken Curry', location: 'Local Church' }
        ];

        if (mockMealHistory.length > 0) {
            mealHistory.innerHTML = mockMealHistory.map(meal => `
                <div class="meal-history-item">
                    <i class="fas fa-utensils"></i>
                    <span>${meal.date}: ${meal.type} from ${meal.location}</span>
                </div>
            `).join('');
        }
    }

    // Initial meals fetch and stats update
    fetchAvailableMeals(receiverData.location);
    updateMealStats();

    // Auto-refresh meals every 5 minutes
    setInterval(() => {
        fetchAvailableMeals(receiverData.location);
        updateMealStats();
    }, 5 * 60 * 1000);
});
