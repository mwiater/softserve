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
