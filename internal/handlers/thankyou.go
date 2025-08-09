package handlers

import (
	"fmt"
	"net/http"
)

var thankYouHTML = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Thank You</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background: #f2f2f2;
      display: flex;
      justify-content: center;
      align-items: center;
      height: 100vh;
      margin: 0;
    }
    .container {
      background: white;
      padding: 30px;
      border-radius: 12px;
      box-shadow: 0 4px 10px rgba(0,0,0,0.1);
      width: 400px;
      text-align: center;
    }
    h1 {
      color: #333;
    }
    .info {
      margin-top: 20px;
      text-align: left;
    }
    .info p {
      background: #f9f9f9;
      padding: 10px;
      border-radius: 6px;
      margin: 8px 0;
      font-size: 14px;
    }
    .button {
      margin-top: 20px;
      padding: 10px 20px;
      background: #4CAF50;
      color: white;
      border: none;
      border-radius: 6px;
      cursor: pointer;
      font-size: 14px;
    }
    .button:hover {
      background: #45a049;
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>Device Linked Successfully!</h1>
    <div class="info">
      <h3>Device Info:</h3>
      <!--<<p><b>ID:</b> <span id="deviceId"></span></p>>-->
      <p><b>Name:</b> <span id="deviceName"></span></p>
      <p><b>Type:</b> <span id="deviceType"></span></p>
    </div>
    <div class="info">
      <h3>Firebase Account:</h3>
      <p><b>Email:</b> <span id="firebaseEmail"></span></p>
      <!--<<p><b>UID:</b> <span id="firebaseUid"></span></p>>-->
    </div>
    <!--<<button class="button" onclick="goHome()">Go to Dashboard</button>>-->
  </div>
  <script>
    function goHome() {
      window.location.href = "/";
    }

    // Populate from query params
    const params = new URLSearchParams(window.location.search);
    // document.getElementById("deviceId").textContent = params.get("deviceId") || "";
    document.getElementById("deviceName").textContent = params.get("deviceName") || "";
    document.getElementById("deviceType").textContent = params.get("deviceType") || "";
    document.getElementById("firebaseEmail").textContent = params.get("email") || "";
    // document.getElementById("firebaseUid").textContent = params.get("uid") || "";
  </script>
</body>
</html>
`
var End = make(chan bool, 1)

func ServeThankYou(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	End <- true
	fmt.Fprint(w, thankYouHTML)
}
