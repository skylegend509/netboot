document.addEventListener('DOMContentLoaded', () => {
    loadISOs();

    const fileInput = document.getElementById('file-uploader');
    fileInput.addEventListener('change', handleUpload);
});

async function loadISOs() {
    try {
        const response = await fetch('/api/isos');
        const isos = await response.json();
        renderTable(isos);
    } catch (error) {
        console.error('Error loading ISOs:', error);
        showNotification('error', 'Failed to load ISO list');
    }
}

function renderTable(isos) {
    const tbody = document.getElementById('iso-table-body');
    tbody.innerHTML = '';

    if (isos.length === 0) {
        tbody.innerHTML = '<tr><td colspan="4" style="text-align:center;">No ISOs found. Upload one to get started.</td></tr>';
        return;
    }

    isos.forEach(iso => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>${iso.name}</td>
            <td>${formatSize(iso.size)}</td>
            <td>${iso.path}</td>
            <td>
                <button class="bx--btn bx--btn--danger bx--btn--sm bx--btn--icon-only" 
                        onclick="deleteISO('${iso.name}')" 
                        aria-label="Delete">
                    <svg focusable="false" preserveAspectRatio="xMidYMid meet" xmlns="http://www.w3.org/2000/svg" fill="currentColor" width="16" height="16" viewBox="0 0 32 32" aria-hidden="true"><path d="M12 12H14V24H12zM18 12H20V24H18z"></path><path d="M4 6V8H6V28a2 2 0 002 2H24a2 2 0 002-2V8h2V6zM8 28V8H24V28zM12 2H20V4H12z"></path></svg>
                </button>
            </td>
        `;
        tbody.appendChild(tr);
    });
}

function formatSize(bytes) {
    const units = ['B', 'KB', 'MB', 'GB', 'TB'];
    let i = 0;
    while (bytes >= 1024 && i < units.length - 1) {
        bytes /= 1024;
        i++;
    }
    return `${bytes.toFixed(2)} ${units[i]}`;
}

async function handleUpload(event) {
    const file = event.target.files[0];
    if (!file) return;

    const statusArea = document.getElementById('upload-status');
    statusArea.innerHTML = `<div class="bx--inline-loading" role="alert" aria-live="assertive">
        <div class="bx--inline-loading__animation">
            <div data-loading class="bx--loading bx--loading--small">
                <svg class="bx--loading__svg" viewBox="-75 -75 150 150"><circle class="bx--loading__stroke" cx="0" cy="0" r="37.5" /></svg>
            </div>
        </div>
        <p class="bx--inline-loading__text">Uploading ${file.name}...</p>
    </div>`;

    const formData = new FormData();
    formData.append('file', file);

    try {
        const response = await fetch('/api/upload', {
            method: 'POST',
            body: formData
        });

        if (response.ok) {
            statusArea.innerHTML = `<div class="bx--inline-notification bx--inline-notification--success" role="alert">
                <div class="bx--inline-notification__details">
                    <p class="bx--inline-notification__title">Upload Complete</p>
                    <p class="bx--inline-notification__subtitle">${file.name} has been uploaded successfully.</p>
                </div>
            </div>`;
            loadISOs(); // Refresh list
        } else {
            const errText = await response.text();
            throw new Error(errText);
        }
    } catch (error) {
        statusArea.innerHTML = `<div class="bx--inline-notification bx--inline-notification--error" role="alert">
            <div class="bx--inline-notification__details">
                <p class="bx--inline-notification__title">Upload Failed</p>
                <p class="bx--inline-notification__subtitle">${error.message}</p>
            </div>
        </div>`;
    }

    // Reset input
    event.target.value = '';
}

async function deleteISO(name) {
    if (!confirm(`Are you sure you want to delete ${name}?`)) return;

    try {
        const response = await fetch(`/api/delete?name=${encodeURIComponent(name)}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            loadISOs();
        } else {
            alert('Failed to delete ISO');
        }
    } catch (error) {
        console.error('Error deleting ISO:', error);
    }
}

function showNotification(type, message) {
    // Simple alert replacement for now
    alert(message);
}