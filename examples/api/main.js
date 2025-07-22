function getUsers() {
  const fetchUrl = "/api/users";
  fetch(fetchUrl)
    .then((res) => res.json())
    .then((data) => {
      document.getElementById("output").textContent = JSON.stringify(data, null, 2);
    })
    .catch((err) => {
      document.getElementById("output").textContent = "Error: " + err;
    });
}

function getClients() {
  const fetchUrl = "/api/clients";
  fetch(fetchUrl)
    .then((res) => res.json())
    .then((data) => {
      document.getElementById("output").textContent = JSON.stringify(data, null, 2);
    })
    .catch((err) => {
      document.getElementById("output").textContent = "Error: " + err;
    });
}

function postLogin() {
  const fetchUrl = "/api/login";
  fetch(fetchUrl, {
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
