function submitQuery(queryMsg="") {
    let url = "/query";

    if (queryMsg !== "") {
        url = url + "?q=" + queryMsg;
    }
    

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
}

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

let searchForm = document.getElementById("searchForm");
searchForm.addEventListener("submit", (e) => {
    e.preventDefault();

    let searchInput = document.getElementById("searchInput");
    submitQuery(searchInput.value);
})

const body = document.querySelector('body');
body.onload =  submitQuery();