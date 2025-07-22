let dataTable = {};

$(document).ready(function () {
  // Create an array to hold all your promises
  const promises = [];

  // JQUERY $.getJSON REQUEST
  // Wrap $.getJSON in a Promise
  const getJsonPromise = new Promise((resolve, reject) => {
    $.getJSON("/api/users")
      .done(function (data) { // Use .done for success with $.getJSON
        dataTable.getJSON = {
          Method: 'GET',
          Endpoint: '/api/clients', // Note: Your endpoint was /api/users, but you set /api/clients in dataTable
          Type: 'JSON',
          Data: JSON.stringify(data)
        };
        resolve(); // Resolve the promise when data is processed
      })
      .fail(function (err) {
        console.error("jQuery getJSON error:", err);
        dataTable.getJSON = { // Still add an entry, but indicate error
          Method: 'GET',
          Endpoint: '/api/clients',
          Type: 'JSON',
          Data: 'ERROR: ' + JSON.stringify(err.responseJSON || err.statusText)
        };
        resolve(); // Resolve even on error to allow Promise.all to complete
      });
  });
  promises.push(getJsonPromise);

  // AXIOS REQUEST
  const axiosPromise = axios({
    method: 'get',
    url: '/api/users'
  })
  .then(function (response) {
    dataTable.axios = {
      Method: 'GET',
      Endpoint: '/api/clients', // Note: Your endpoint was /api/users, but you set /api/clients in dataTable
      Type: 'JSON',
      Data: JSON.stringify(response.data)
    };
  })
  .catch(function (error) {
    console.error('Axios error:', error);
    dataTable.axios = { // Add an entry for error
      Method: 'GET',
      Endpoint: '/api/clients',
      Type: 'JSON',
      Data: 'ERROR: ' + JSON.stringify(error.response ? error.response.data : error.message)
    };
  });
  promises.push(axiosPromise);

  // XMLHttpRequest REQUEST
  const xhrPromise = new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open('GET', '/api/users', true);
    xhr.responseType = 'json';

    xhr.onload = function () {
      if (xhr.status === 200) {
        const data = xhr.response;
        dataTable.XMLHttpRequest = {
          Method: 'GET',
          Endpoint: '/api/clients', // Note: Your endpoint was /api/users, but you set /api/clients in dataTable
          Type: 'JSON',
          Data: JSON.stringify(data)
        };
        resolve();
      } else {
        console.error('XMLHttpRequest error. Status:', xhr.status);
        dataTable.XMLHttpRequest = {
          Method: 'GET',
          Endpoint: '/api/clients',
          Type: 'JSON',
          Data: 'ERROR: Status ' + xhr.status
        };
        resolve(); // Resolve even on error
      }
    };

    xhr.onerror = function () {
      console.error('XMLHttpRequest network error occurred.');
      dataTable.XMLHttpRequest = {
        Method: 'GET',
        Endpoint: '/api/clients',
        Type: 'JSON',
        Data: 'ERROR: Network error'
      };
      resolve(); // Resolve even on error
    };
    xhr.send();
  });
  promises.push(xhrPromise);

  // FETCH REQUEST (first one)
  const fetchPromise1 = fetch("/api/users")
    .then((res) => {
        if (!res.ok) { // Check for HTTP errors
            throw new Error(`HTTP error! status: ${res.status}`);
        }
        return res.json();
    })
    .then((data) => {
      dataTable.fetch = { // Renamed from 'fetch' to 'fetch1' to avoid conflict if you add more
        Method: 'GET',
        Endpoint: '/api/clients', // Note: Your endpoint was /api/users, but you set /api/clients in dataTable
        Type: 'JSON',
        Data: JSON.stringify(data)
      };
    })
    .catch((err) => {
      console.error('Fetch error:', err);
      dataTable.fetch = { // Add an entry for error
        Method: 'GET',
        Endpoint: '/api/clients',
        Type: 'JSON',
        Data: 'ERROR: ' + err.message
      };
    });
  promises.push(fetchPromise1);

  // Use Promise.all to wait for all promises to resolve
  Promise.all(promises)
    .then(() => {
      console.log("All requests complete. Printing dataTable:");
      // Convert dataTable object to an array for console.table to show a clean index
      const tableData = Object.keys(dataTable).map(key => ({
          RequestType: key,
          ...dataTable[key]
      }));
      console.table(tableData);
    })
    .catch((error) => {
      console.error("One or more requests failed:", error);
      // Even if some failed, you might still want to print what you have
      const tableData = Object.keys(dataTable).map(key => ({
          RequestType: key,
          ...dataTable[key]
      }));
      console.table(tableData);
    });
});

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
