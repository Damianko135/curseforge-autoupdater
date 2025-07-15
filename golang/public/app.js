function updateTime() {
    const timeElement = document.getElementById('time-display');
    const now = new Date();

    // Format date to YYYY-MM-DD HH:MM:SS
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, '0');
    const day = String(now.getDate()).padStart(2, '0');
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');
    const seconds = String(now.getSeconds()).padStart(2, '0');

    const formatted = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    timeElement.textContent = formatted;
}

// Update every second
setInterval(updateTime, 1000);

// Initial call
updateTime();
