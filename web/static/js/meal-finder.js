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

    // Fetch available meals
    async function fetchAvailableMeals(location) {
        const mealsList = document.getElementById('mealsList');
        
        try {
            console.log('Fetching meals for location:', location);
            const response = await fetch(`/api/meals?location=${encodeURIComponent(location)}`);
            
            console.log('Response status:', response.status);
            console.log('Response headers:', Object.fromEntries([...response.headers]));
            
            const responseText = await response.text();
            console.log('Raw response text:', responseText);

            if (!response.ok) {
                throw new Error(`Server error: ${response.status} - ${responseText}`);
            }

            if (!responseText.trim()) {
                console.log('Empty response, treating as empty array');
                mealsList.innerHTML = `
                    <div class="info-message">
                        <p>No meals are currently available in ${location}.</p>
                        <p>Please check back later or try a different location.</p>
                    </div>`;
                return;
            }

            let meals;
            try {
                meals = JSON.parse(responseText);
                console.log('Parsed meals data:', meals);
            } catch (parseError) {
                console.error('JSON parse error:', parseError);
                throw new Error(`Failed to parse server response: ${parseError.message}`);
            }

            if (!Array.isArray(meals)) {
                console.error('Non-array response:', meals);
                meals = []; // Convert null or undefined to empty array
            }

            mealsList.innerHTML = '';

            if (meals.length === 0) {
                console.log('No meals found for location:', location);
                mealsList.innerHTML = `
                    <div class="info-message">
                        <p>No meals are currently available in ${location}.</p>
                        <p>Please check back later or try a different location.</p>
                    </div>`;
                return;
            }

            console.log(`Found ${meals.length} meals for location:`, location);

            const mealsHtml = meals.map((meal, index) => {
                console.log(`Processing meal ${index}:`, meal);
                
                // Validate meal object
                if (!meal || typeof meal !== 'object') {
                    console.error(`Invalid meal at index ${index}:`, meal);
                    return '';
                }

                // Check required fields
                const requiredFields = ['id', 'type', 'quantity', 'location', 'expiryDate'];
                const missingFields = requiredFields.filter(field => !meal[field]);
                if (missingFields.length > 0) {
                    console.error(`Meal at index ${index} is missing fields:`, missingFields);
                    return '';
                }

                try {
                    const expiryDate = new Date(meal.expiryDate);
                    if (isNaN(expiryDate.getTime())) {
                        console.error(`Invalid expiry date for meal ${index}:`, meal.expiryDate);
                        return '';
                    }

                    return `
                        <div class="card meal-card">
                            <h3>${meal.type}</h3>
                            <p><strong>Available:</strong> ${meal.quantity} servings</p>
                            <p><strong>Location:</strong> ${meal.location}</p>
                            <p><strong>Expires:</strong> ${expiryDate.toLocaleString()}</p>
                            <button class="btn btn-primary request-meal" data-meal-id="${meal.id}">
                                Request Meal
                            </button>
                        </div>`;
                } catch (error) {
                    console.error(`Error processing meal ${index}:`, error);
                    return '';
                }
            }).filter(html => html !== '').join('');

            if (!mealsHtml) {
                mealsList.innerHTML = `
                    <div class="info-message">
                        <p>No valid meals found in ${location}.</p>
                        <p>Please check back later or try a different location.</p>
                    </div>`;
                return;
            }

            mealsList.innerHTML = mealsHtml;

            // Add event listeners to request buttons
            document.querySelectorAll('.request-meal').forEach(button => {
                button.addEventListener('click', function() {
                    const mealId = this.dataset.mealId;
                    requestMeal(mealId);
                });
            });

        } catch (error) {
            console.error('Error in fetchAvailableMeals:', error);
            mealsList.innerHTML = `
                <div class="error-message">
                    <p>Unable to load meals at this time.</p>
                    <p>Error: ${error.message}</p>
                    <p>Please check the console for more details.</p>
                </div>`;
        }
    }

    // Change location handler
    document.getElementById('changeLocation').addEventListener('click', function() {
        const newLocation = prompt('Enter new location:');
        if (newLocation) {
            receiverData.location = newLocation;
            localStorage.setItem('receiverData', JSON.stringify(receiverData));
            document.getElementById('receiverLocation').textContent = newLocation;
            fetchAvailableMeals(newLocation);
        }
    });

    // Initial meals fetch
    fetchAvailableMeals(receiverData.location);

    // Auto-refresh meals every 5 minutes
    setInterval(() => {
        fetchAvailableMeals(receiverData.location);
    }, 5 * 60 * 1000);
});
