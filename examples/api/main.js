function getUsers() {
  fetch("/api/users")
    .then((res) => res.json())
    .then((data) => {
      document.getElementById("output").textContent = JSON.stringify(data, null, 2);
    })
    .catch((err) => {
      document.getElementById("output").textContent = "Error: " + err;
    });
}

function getClients() {
  fetch("/api/clients")
    .then((res) => res.json())
    .then((data) => {
      document.getElementById("output").textContent = JSON.stringify(data, null, 2);
    })
    .catch((err) => {
      document.getElementById("output").textContent = "Error: 404. Ensure you have the '/api/clients' route moked in your api.yaml";
    });
}

function postLogin() {
  fetch("/api/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ username: "test", password: "wrong" })
  })
    .then((res) => res.json())
    .then((data) => {
      document.getElementById("output").textContent = JSON.stringify(data, null, 2);
    })
    .catch((err) => {
      document.getElementById("output").textContent = "Error: " + err;
    });
}
