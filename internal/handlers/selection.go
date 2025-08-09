package handlers

import (
	"fmt"
	"net/http"
)

const deviceSelection = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>Cryonics Device Selection</title>
<style>
  /* Base & Reset */
  * {
    box-sizing: border-box;
  }
  body {
    margin: 0; padding: 2rem;
    min-height: 100vh;
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    background: linear-gradient(135deg, #1e3c72, #2a5298);
    color: #fff;
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
  h1 {
    font-weight: 700;
    font-size: 2.8rem;
    margin-bottom: 2rem;
    text-shadow: 0 2px 6px rgba(0,0,0,0.4);
  }
  #devicesList {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 1.25rem;
    width: 100%;
    max-width: 700px;
  }
  .device-item {
    background: rgba(255 255 255 / 0.15);
    border-radius: 12px;
    padding: 1.25rem 1.5rem;
    cursor: pointer;
    box-shadow: 0 8px 24px rgb(0 0 0 / 0.25);
    transition: background 0.3s ease, transform 0.2s ease;
    user-select: none;
    outline-offset: 4px;
  }
  .device-item:hover,
  .device-item:focus-visible {
    background: rgba(255 255 255 / 0.3);
    transform: translateY(-3px);
    outline: none;
  }
  .device-name {
    font-weight: 700;
    font-size: 1.4rem;
    margin-bottom: 0.25rem;
  }
  .device-type {
    font-size: 1rem;
    opacity: 0.85;
  }
  #loading,
  #error {
    margin-top: 3rem;
    font-style: italic;
    font-weight: 600;
  }
  #error {
    color: #ff6b6b;
  }
</style>
</head>
<body>
  <h1>Select Your Device</h1>
  <div id="devicesList" role="list"></div>
  <div id="loading" aria-live="polite">Loading devices...</div>
  <div id="error" role="alert"></div>

<script type="module">
  import { initializeApp } from "https://www.gstatic.com/firebasejs/10.12.0/firebase-app.js";
  import { getAuth, onAuthStateChanged } from "https://www.gstatic.com/firebasejs/10.12.0/firebase-auth.js";

  const firebaseConfig = {
    apiKey: "AIzaSyAy-ElSswFhDrejhXcs_MsUaWgG-2cnbGM",
    authDomain: "cryonics-em.firebaseapp.com",
    projectId: "cryonics-em",
    storageBucket: "cryonics-em.firebasestorage.app",
    messagingSenderId: "907895762193",
    appId: "1:907895762193:web:5122aecfcbefbf309d70b0"
  };

  const app = initializeApp(firebaseConfig);
  const auth = getAuth(app);

  const devicesList = document.getElementById('devicesList');
  const loadingEl = document.getElementById('loading');
  const errorEl = document.getElementById('error');

  function renderDevices(devices) {
    devicesList.innerHTML = '';
    if (!devices || Object.keys(devices).length === 0) {
      devicesList.innerHTML = '<p>No devices found for your account.</p>';
      return;
    }
    for (const deviceId in devices) {
      if (!devices.hasOwnProperty(deviceId)) continue;
      const device = devices[deviceId];

      const div = document.createElement('div');
      div.className = 'device-item';
      div.tabIndex = 0;
      div.setAttribute('role', 'listitem');
      div.setAttribute('aria-pressed', 'false');

      div.innerHTML =
        '<div class="device-name">' + (device.name || deviceId) + '</div>' +
        '<div class="device-type">' + (device.deviceType || 'Unknown Type') + '</div>';

      div.onclick = () => selectDevice(deviceId, device);
      div.onkeypress = e => { if(e.key === 'Enter' || e.key === ' ') selectDevice(deviceId, device); };

      devicesList.appendChild(div);
    }
  }

  async function selectDevice(deviceId, device) {
	const user = auth.currentUser;
    const formData = new FormData();
    formData.append("deviceId", deviceId);
    formData.append("name", device.name);
    formData.append("deviceType", device.deviceType);
	formData.append("email", user.email);
  	formData.append("uid", user.uid);

    fetch("/receive-device", {
      method: "POST",
      body: formData
    }).then(res => {
      if (res.redirected) {
        window.location.href = res.url; // follow the redirect
      }
  	});
  }

  async function fetchDevices(idToken, uid) {
    loadingEl.style.display = 'block';
    errorEl.textContent = '';
    devicesList.innerHTML = '';

    const url = 'https://cryonics-em-default-rtdb.asia-southeast1.firebasedatabase.app/users/' + uid + '.json?auth=' + idToken;
    try {
      const res = await fetch(url);
      if (!res.ok) throw new Error('Failed to fetch devices: ' + res.statusText);
      const data = await res.json();
      renderDevices(data);
    } catch (e) {
      errorEl.textContent = e.message;
    } finally {
      loadingEl.style.display = 'none';
    }
  }

  onAuthStateChanged(auth, async user => {
    if (user) {
      const idToken = await user.getIdToken();
      const uid = user.uid;
      fetchDevices(idToken, uid);
    } else {
      errorEl.textContent = 'User not authenticated. Please log in first.';
      loadingEl.style.display = 'none';
    }
  });
</script>
</body>
</html>
`

func ServeSelection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, deviceSelection)
}
