document.addEventListener('DOMContentLoaded', function() {
    const donorForm = document.getElementById('donorForm');
    const qrContainer = document.getElementById('qrContainer');
    const printQRButton = document.getElementById('printQR');

    // Set minimum date for expiry date input to current date
    const expiryDateInput = document.getElementById('expiryDate');
    const today = new Date().toISOString().split('T')[0];
    expiryDateInput.min = today;

    donorForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const formData = new FormData(donorForm);
        const data = {
            name: formData.get('name'),
            email: formData.get('email'),
            phone: formData.get('phone'),
            mealType: formData.get('mealType'),
            quantity: parseInt(formData.get('quantity'), 10),
            expiryDate: new Date(formData.get('expiryDate')).toISOString(),
            location: formData.get('location'),
            notes: formData.get('notes') || ''
        };

        try {
            const response = await fetch('/api/donations', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText);
            }

            const result = await response.json();
            console.log('API Response:', result);
            
            // Clear previous QR code if exists
            const qrCodeDiv = document.getElementById('qrCode');
            if (!qrCodeDiv) {
                console.error('QR code div not found');
                return;
            }
            qrCodeDiv.innerHTML = '';
            
            // Generate QR Code
            qrContainer.style.display = 'block';
            console.log('Generating QR code for donation ID:', result.donationId);
            new QRCode(qrCodeDiv, {
                text: result.donationId || 'test',
                width: 200,
                height: 200
            });
            console.log('QR code generated');

            // Clear form
            donorForm.reset();

            // Show success message
            showNotification('Donation registered successfully!', 'success');

        } catch (error) {
            console.error('Error:', error);
            showNotification(error.message || 'Failed to register donation. Please try again.', 'error');
        }
    });

    printQRButton.addEventListener('click', function() {
        const qrCodeImg = document.getElementById('qrCode').querySelector('img');
        if (!qrCodeImg) {
            showNotification('No QR code to print', 'error');
            return;
        }
        
        const printWindow = window.open('', '', 'width=600,height=600');
        printWindow.document.open();
        printWindow.document.write('<html><head><title>Print QR Code</title></head><body>');
        printWindow.document.write('<img src="' + qrCodeImg.src + '" style="width: 100%; max-width: 400px;">');
        printWindow.document.write('</body></html>');
        printWindow.document.close();
        printWindow.print();
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
