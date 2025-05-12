document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('kycForm');
    const idProofInput = document.getElementById('idProof');
    const addressProofInput = document.getElementById('addressProof');
    const idProofPreview = document.getElementById('idProofPreview');
    const addressProofPreview = document.getElementById('addressProofPreview');

    // Handle file preview for ID Proof
    idProofInput.addEventListener('change', (e) => {
        handleFilePreview(e.target.files[0], idProofPreview);
    });

    // Handle file preview for Address Proof
    addressProofInput.addEventListener('change', (e) => {
        handleFilePreview(e.target.files[0], addressProofPreview);
    });

    // Function to handle file preview
    function handleFilePreview(file, previewElement) {
        if (!file) {
            previewElement.style.display = 'none';
            return;
        }

        previewElement.style.display = 'block';
        previewElement.innerHTML = '';

        if (file.type.startsWith('image/')) {
            const img = document.createElement('img');
            img.src = URL.createObjectURL(file);
            img.onload = () => URL.revokeObjectURL(img.src);
            previewElement.appendChild(img);
        } else {
            const fileInfo = document.createElement('div');
            fileInfo.innerHTML = `
                <p><strong>File Name:</strong> ${file.name}</p>
                <p><strong>File Type:</strong> ${file.type || 'Unknown'}</p>
                <p><strong>File Size:</strong> ${(file.size / (1024 * 1024)).toFixed(2)} MB</p>
            `;
            previewElement.appendChild(fileInfo);
        }
    }

    // Handle form submission
    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        // Basic form validation
        let isValid = true;
        let errorMessages = [];

        // Validate file sizes (max 5MB each)
        const maxSize = 5 * 1024 * 1024; // 5MB in bytes
        const files = [idProofInput.files[0], addressProofInput.files[0]];
        
        files.forEach((file, index) => {
            if (!file) {
                isValid = false;
                errorMessages.push(`${index === 0 ? 'ID Proof' : 'Address Proof'} is required`);
            } else if (file.size > maxSize) {
                isValid = false;
                errorMessages.push(`${index === 0 ? 'ID Proof' : 'Address Proof'} file must be less than 5MB`);
            }
        });

        if (!isValid) {
            alert('Please fix the following errors:\n' + errorMessages.join('\n'));
            return;
        }

        // Show loading state
        const submitBtn = form.querySelector('.submit-btn');
        const originalBtnText = submitBtn.textContent;
        submitBtn.textContent = 'Submitting...';
        submitBtn.disabled = true;

        try {
            const formData = new FormData();
            formData.append('files', idProofInput.files[0]);
            formData.append('files', addressProofInput.files[0]);

            const response = await fetch('http://127.0.0.1:3000/kyc-verify/swanhtetaungp@gmail.com', {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const result = await response.json();
            alert('KYC verification submitted successfully!');
            form.reset();
            idProofPreview.style.display = 'none';
            addressProofPreview.style.display = 'none';
        } catch (error) {
            console.error('Error:', error);
            alert('An error occurred while submitting the form. Please try again.');
        } finally {
            submitBtn.textContent = originalBtnText;
            submitBtn.disabled = false;
        }
    });
}); 