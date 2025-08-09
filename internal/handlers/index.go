package handlers

import (
	"fmt"
	"net/http"
)

const indexHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>Cryonics Firebase Login</title>
<style>
  /* Reset & base */
  * {
    box-sizing: border-box;
  }
  body {
    margin: 0;
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    background: linear-gradient(135deg, #1e3c72, #2a5298);
    color: #fff;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    padding: 2rem;
    text-align: center;
  }
  h1 {
    font-weight: 700;
    font-size: 2.8rem;
    margin-bottom: 1.5rem;
    text-shadow: 0 2px 6px rgba(0,0,0,0.4);
  }
  #loginBtn {
    background: #4caf50;
    border: none;
    border-radius: 8px;
    padding: 0.75rem 2.5rem;
    font-size: 1.25rem;
    font-weight: 600;
    color: white;
    cursor: pointer;
    box-shadow: 0 4px 12px rgb(76 175 80 / 0.5);
    transition: background 0.3s ease, box-shadow 0.3s ease;
  }
  #loginBtn:hover,
  #loginBtn:focus {
    background: #43a047;
    box-shadow: 0 6px 18px rgb(67 160 71 / 0.7);
    outline: none;
  }
  #loginBtn:active {
    background: #388e3c;
    box-shadow: none;
  }
  footer {
    margin-top: 3rem;
    font-size: 0.9rem;
    opacity: 0.7;
  }
</style>
<script type="module">
  import { initializeApp } from "https://www.gstatic.com/firebasejs/10.12.0/firebase-app.js";
  import { getAuth, GoogleAuthProvider, signInWithPopup } from "https://www.gstatic.com/firebasejs/10.12.0/firebase-auth.js";

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

  async function login() {
    const provider = new GoogleAuthProvider();
    const result = await signInWithPopup(auth, provider);
    const token = await result.user.getIdToken();
    console.log("ID Token:", token);

    const formData = new FormData();
    formData.append("idToken", token);

    const res = await fetch("/receive-token", {
      method: "POST",
      body: formData
    });
    const text = await res.text();
	window.location.href="/select-device"
  }

  window.onload = () => {
    document.getElementById("loginBtn").onclick = login;
  };
</script>
</head>
<body>
  <h1>Cryonics Firebase Google Login</h1>
  <button id="loginBtn" aria-label="Login with Google">Login with Google</button>
  <footer>
    &copy; 2025 Cryonics Project
  </footer>
</body>
</html>
`

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, indexHTML)
}
