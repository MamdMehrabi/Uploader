document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('uploadForm');
    const fileInput = document.getElementById('fileInput');
    const fileName = document.getElementById('fileName');
    const uploadBtn = document.getElementById('uploadBtn');
    const loader = document.getElementById('loader');
    const messageDiv = document.getElementById('message');
    
    // Fetch max file size from server on page load
    let maxFileSizeBytes = 50 * 1024 * 1024; // Default 50MB
    fetch('/api/max-file-size')
        .then(response => response.json())
        .then(data => {
            maxFileSizeBytes = data.maxFileSizeBytes;
        })
        .catch(e => {
            console.warn('Could not fetch max file size, using default 50MB');
        });

    // Show selected file name
    fileInput.addEventListener('change', function(e) {
        const file = e.target.files[0];
        if (file) {
            fileName.textContent = `Selected: ${file.name} | Size: (${formatFileSize(file.size)})`;
            fileName.classList.add('show');
        } else {
            fileName.classList.remove('show');
        }
    });

    // Handle form submission
    form.addEventListener('submit', async function(e) {
        e.preventDefault();

        const chatId = document.getElementById('chatId').value.trim();
        const caption = document.getElementById('caption').value.trim();
        const file = fileInput.files[0];

        if (!file) {
            showMessage('Please select a file', 'error');
            return;
        }

        if (!chatId) {
            showMessage('Please enter a Chat ID', 'error');
            return;
        }

        // Check file size (using maxFileSizeBytes fetched on page load)
        if (file.size > maxFileSizeBytes) {
            const maxMB = (maxFileSizeBytes / (1024 * 1024)).toFixed(0);
            const fileMB = (file.size / (1024 * 1024)).toFixed(2);
            showMessage(`File size (${fileMB} MB) exceeds the limit of ${maxMB} MB`, 'error');
            return;
        }

        // Show loading state
        uploadBtn.disabled = true;
        loader.classList.add('show');
        uploadBtn.querySelector('.btn-text').textContent = 'Uploading...';
        hideMessage();

        // Create form data
        const formData = new FormData();
        formData.append('file', file);
        formData.append('chatId', chatId);
        if (caption) {
            formData.append('caption', caption);
        }

        try {
            const response = await fetch('/api/upload', {
                method: 'POST',
                body: formData
            });

            const data = await response.json();

            if (data.success) {
                showMessage(
                    `✅ ${data.message}\nFile ID: ${data.fileId}`,
                    'success'
                );
                form.reset();
                fileName.classList.remove('show');
            } else {
                let errorMsg = `❌ Error: ${data.error}`;
                // Provide helpful hints for common errors
                if (data.error && data.error.includes('Bot Token')) {
                    errorMsg += '\n\n💡 Tip: Create a .env file with TELEGRAM_BOT_TOKEN=your_token_here';
                } else if (data.error && data.error.includes('Unauthorized')) {
                    errorMsg += '\n\n💡 Tip: Your bot token is invalid. Get a new one from @BotFather on Telegram';
                } else if (data.error && data.error.includes('chat not found')) {
                    errorMsg += '\n\n💡 Tip: Make sure you have started a conversation with your bot first';
                } else if (data.error && data.error.includes('Cannot connect') || data.error.includes('blocked')) {
                    errorMsg += '\n\n💡 Tip: Telegram might be blocked in your country. Add HTTP_PROXY to your .env file (e.g., HTTP_PROXY=socks5://127.0.0.1:1080)';
                }
                showMessage(errorMsg, 'error');
            }
        } catch (error) {
            showMessage(`❌ Network error: ${error.message}`, 'error');
        } finally {
            // Reset loading state
            uploadBtn.disabled = false;
            loader.classList.remove('show');
            uploadBtn.querySelector('.btn-text').textContent = 'Upload to Telegram';
        }
    });

    function showMessage(text, type) {
        messageDiv.textContent = text;
        messageDiv.className = `message show ${type}`;
        
        // Auto-hide after 5 seconds
        setTimeout(() => {
            hideMessage();
        }, 5000);
    }

    function hideMessage() {
        messageDiv.classList.remove('show');
    }

    function formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
    }
});

