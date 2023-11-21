const footer = document.querySelector("footer p");
const url = "http://localhost:8080/query";

fetch(url)
.then((response) => {
    if (!response.ok) {
        throw new Error(`HTTP error: ${response.status}`);
    }

    return response.json();
})
.then((json) => populateLogTable(json))
.catch((error) => {
    footer.textContent = `Could not fetch log messages: ${error}`;
})

function populateLogTable(logMessages) {

    const logTable = document.getElementById("logTable");

    while(logTable.firstChild) {
        logTable.removeChild(logTable.firstChild)
    }

    for (const logMsg of logMessages) {
        var row = logTable.insertRow();
        var cell = row.insertCell();
        cell.innerHTML = logMsg.m;
    }

};