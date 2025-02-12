document.addEventListener('DOMContentLoaded', function() {
    const donorForm = document.getElementById('donorForm');
    const qrContainer = document.getElementById('qrContainer');
    const printQRButton = document.getElementById('printQR');

    // Set minimum date for expiry date input to current date/time
    const expiryDateInput = document.getElementById('expiryDate');
    const now = new Date();
    const tzoffset = now.getTimezoneOffset() * 60000;
    const localISOTime = (new Date(Date.now() - tzoffset)).toISOString().slice(0, 16);
    expiryDateInput.min = localISOTime;

    donorForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const formData = new FormData(donorForm);
        const data = Object.fromEntries(formData.entries());

        try {
            const response = await fetch('/api/donations', {
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
            
            // Generate QR Code
            qrContainer.style.display = 'block';
            const qrCode = new QRCode(document.getElementById('qrCode'), {
                text: result.donationId,
                width: 200,
                height: 200
            });

            // Clear form
            donorForm.reset();

            // Show success message
            showNotification('Donation registered successfully!', 'success');

        } catch (error) {
            console.error('Error:', error);
            showNotification('Failed to register donation. Please try again.', 'error');
        }
    });

    printQRButton.addEventListener('click', function() {
        const qrCodeImg = document.getElementById('qrCode').querySelector('img');
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
