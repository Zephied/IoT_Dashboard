let localStream = null;
let mockCameraActive = false; // mémorise l'état du flux vidéo mock
let lastDevicesSnapshot = null; // snapshot pour comparer l'état réseau

// Icônes simples selon type
const icons = {
    camera: '📷',
    door: '🚪',
    light: '💡',
    default: '🔌'
};

async function fetchDevices() {
    const res = await fetch('/api/devices');
    const devices = await res.json();

    // Si c'est le premier chargement, on affiche normalement
    if (lastDevicesSnapshot === null) {
        renderDevices(devices);
        lastDevicesSnapshot = JSON.stringify(devices);
        return;
    }

    // Sinon, on compare l'état courant avec le précédent
    const newSnapshot = JSON.stringify(devices);
    if (newSnapshot !== lastDevicesSnapshot) {
        showNetworkChangeNotification();
    }
}

function renderDevices(devices) {
    const container = document.getElementById('device-container');
    container.innerHTML = '';
    devices.forEach(d => {
        const icon = icons[d.type] || icons.default;

        // Création carte device
        const card = document.createElement('div');
        card.className = 'device-card';
        card.dataset.deviceId = d.id;

        // Header carte
        const header = document.createElement('div');
        header.className = 'device-header';
        header.innerHTML = `<div class="device-icon">${icon}</div><div class="device-name">${d.name}</div>`;
        card.appendChild(header);

        // Affichage de la description si présente
        if (d.desc && d.desc.trim() !== "") {
            const descDiv = document.createElement('div');
            descDiv.className = 'device-desc';
            descDiv.textContent = d.desc;
            descDiv.style.fontStyle = 'italic';
            descDiv.style.marginBottom = '4px';
            card.appendChild(descDiv);
        }

        // Infos IP et MAC
        const info = document.createElement('div');
        info.className = 'device-info';
        info.textContent = `IP: ${d.ip} | MAC: ${d.mac}`;
        card.appendChild(info);

        // Affichage de l'état online/offline
        const status = document.createElement('div');
        status.className = 'device-status';
        status.textContent = d.online ? '🟢 En ligne' : '🔴 Hors ligne';
        status.style.fontWeight = 'bold';
        status.style.marginBottom = '5px';
        card.appendChild(status);

        // Actions container
        const actionsDiv = document.createElement('div');
        actionsDiv.className = 'actions';

        // Bouton suppression
        const deleteBtn = document.createElement('button');
        deleteBtn.textContent = '🗑️ Supprimer';
        deleteBtn.style.marginRight = '8px';
        deleteBtn.onclick = (e) => {
            e.stopPropagation();
            if (confirm('Supprimer cet objet ?')) {
                deleteDevice(d.id);
            }
        };
        actionsDiv.appendChild(deleteBtn);

        // Bouton édition
        const editBtn = document.createElement('button');
        editBtn.textContent = '✏️ Modifier';
        editBtn.style.marginRight = '8px';
        editBtn.onclick = (e) => {
            e.stopPropagation();
            showEditForm(card, d);
        };
        actionsDiv.appendChild(editBtn);

        // Actions dynamiques existantes
        if (d.actions) {
            try {
                const actions = JSON.parse(d.actions);
                for (const actionName in actions) {
                    const btn = document.createElement('button');
                    btn.textContent = actionName;
                    btn.onclick = (e) => {
                        e.stopPropagation();
                        sendAction(d.id, actionName, d, card);
                    };
                    actionsDiv.appendChild(btn);
                }
            } catch (e) {
                const errMsg = document.createElement('div');
                errMsg.textContent = 'Actions invalides';
                errMsg.style.color = 'red';
                actionsDiv.appendChild(errMsg);
            }
        }

        // Ajout d'un conteneur vidéo dédié dans la carte
        const videoContainer = document.createElement('div');
        videoContainer.className = 'video-container';
        card.appendChild(videoContainer);

        card.appendChild(actionsDiv);

        // Toggle affichage actions au clic sur la carte
        card.onclick = () => {
            card.classList.toggle('expanded');
        };

        container.appendChild(card);

        // Si c'est le mock caméra et le flux était actif, le réafficher
        if (d.type === 'camera' && d.ip === '127.0.0.2' && mockCameraActive) {
            showLocalCamera(videoContainer);
        }
    });
}

function showNetworkChangeNotification() {
    let notif = document.getElementById('network-notif');
    if (!notif) {
        notif = document.createElement('div');
        notif.id = 'network-notif';
        notif.style.position = 'fixed';
        notif.style.bottom = '20px';
        notif.style.left = '50%';
        notif.style.transform = 'translateX(-50%)';
        notif.style.background = '#ffecb3';
        notif.style.color = '#333';
        notif.style.padding = '1em 2em';
        notif.style.border = '1px solid #e0b800';
        notif.style.borderRadius = '8px';
        notif.style.zIndex = 1000;
        notif.innerHTML = 'Des changements ont été détectés sur le réseau. <button onclick="location.reload()">Rafraîchir la page</button>';
        document.body.appendChild(notif);
    }
}

// Adapter sendAction pour gérer la réponse 'voir la caméra'
function sendAction(id, action, device = null, card = null) {
    if (action === 'voir la caméra' && device && device.type === 'camera' && device.ip === '127.0.0.2' && card) {
        toggleLocalCamera(card.querySelector('.video-container'));
        return;
    }
    fetch(`/api/devices/${id}/action`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ action })
    })
        .then(res => res.text())
        .then(alert)
        .catch(err => alert('Erreur : ' + err));
}

// Affiche ou coupe le flux webcam local dans le conteneur passé en paramètre
async function toggleLocalCamera(container) {
    if (!container) return;
    let video = container.querySelector('video#localCam');
    if (video && localStream) {
        // Si déjà affiché, on coupe le flux
        localStream.getTracks().forEach(track => track.stop());
        video.srcObject = null;
        localStream = null;
        container.innerHTML = '';
        mockCameraActive = false;
        return;
    }
    // Sinon, on affiche le flux
    container.innerHTML = '';
    video = document.createElement('video');
    video.id = 'localCam';
    video.controls = true;
    video.autoplay = true;
    video.style.display = 'block';
    video.style.width = '100%';
    container.appendChild(video);
    try {
        const stream = await navigator.mediaDevices.getUserMedia({ video: true });
        video.srcObject = stream;
        localStream = stream;
        mockCameraActive = true;
    } catch (err) {
        alert("Erreur d'accès à la caméra : " + err.message);
        mockCameraActive = false;
    }
}

async function toggleCamera() {
    const video = document.getElementById('localCam');
    if (!localStream) {
        try {
            const stream = await navigator.mediaDevices.getUserMedia({ video: true });
            video.srcObject = stream;
            localStream = stream;
            video.style.display = 'block';
        } catch (err) {
            alert('Erreur d\'accès à la caméra : ' + err.message);
        }
    } else {
        localStream.getTracks().forEach(track => track.stop());
        video.srcObject = null;
        localStream = null;
        video.style.display = 'none';
    }
}

async function scanNetwork() {
    await fetch('/api/scan');
    fetchDevices();
}

async function addMockDevice() {
    await fetch('/api/add-mock-device', { method: 'POST' });
    fetchDevices();
}

// Fonction suppression device
async function deleteDevice(id) {
    await fetch(`/api/devices/${id}`, { method: 'DELETE' });
    fetchDevices();
}

// Affiche un formulaire d'édition inline
function showEditForm(card, device) {
    // Empêche plusieurs formulaires
    if (card.querySelector('.edit-form')) return;
    const form = document.createElement('form');
    form.className = 'edit-form';
    form.style.marginTop = '10px';
    form.innerHTML = `
        <input type="text" name="name" value="${device.name}" placeholder="Nom" required style="margin-right:8px;">
        <input type="text" name="desc" value="${device.desc || ''}" placeholder="Description" style="margin-right:8px;">
        <button type="submit">💾 Enregistrer</button>
        <button type="button" class="cancel-btn">Annuler</button>
    `;
    form.onsubmit = async (e) => {
        e.preventDefault();
        e.stopPropagation();
        const name = form.elements['name'].value;
        const desc = form.elements['desc'].value;
        await updateDevice(device.id, name, desc);
        fetchDevices();
    };
    form.onclick = (e) => e.stopPropagation();
    form.querySelector('.cancel-btn').onclick = (e) => {
        e.stopPropagation();
        form.remove();
    };
    card.appendChild(form);
}

// Fonction édition device
async function updateDevice(id, name, desc) {
    await fetch(`/api/devices/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, desc })
    });
}

// Chargement initial
fetchDevices();
setInterval(fetchDevices, 10000);