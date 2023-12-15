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

function websocket_onopen(evt) {

}

function websocket_onmessage(evt) {

}

function websocket_onclose(evt) {

}

let websocket = new WebSocket("ws://" + location.host + "/events")
websocket.onopen = function(event) { 
    appendLogTable([{"m":"Websocket connected to " + websocket.url}])
};
websocket.onmessage = function(event) { 
    // logMessages = JSON.parse(event.data)
    // appendLogTable(logMessages); 
    appendLogTable([{"m":"Websocket connected to " + event.data}])
};
websocket.onclose = function(event) { 
    appendLogTable([{"m":"Websocket connection closed to " + websocket.url}]); 
};

function isScrolledToBottom() {
    var currScrollPos = window.scrollY || document.documentElement.scrollTop;
    var totalHeight = document.documentElement.scrollHeight;
    var windowHeight = window.innerHeight || document.documentElement.clientHeight;
    return currScrollPos + windowHeight >= totalHeight;
}

function scrollToBottom() {
    window.scrollTo({
        top: document.documentElement.scrollHeight,
        behavior: 'smooth'
    });
}



function populateLogTable(logMessages) {

    const logTable = document.getElementById("logTable");

    while(logTable.firstChild) {
        logTable.removeChild(logTable.firstChild)
    }

    if (logMessages.length === 0) {
        logMessages.push({"m":"No log messages"})
    }
    
    for (const logMsg of logMessages) {
        var row = logTable.insertRow();
        var cell = row.insertCell();
        cell.innerHTML = logMsg.m;
    }
};

function appendLogTable(logMessages) {
    isScrollingNeeded = isScrolledToBottom();

    const logTable = document.getElementById("logTable");
    if (0 < logMessages.length) {
    
        for (const logMsg of logMessages) {
            var row = logTable.insertRow();
            var cell = row.insertCell();
            cell.innerHTML = logMsg.m;
        }    
    }

    if (isScrollingNeeded) {
        scrollToBottom();
    }
}

let searchForm = document.getElementById("searchForm");
searchForm.addEventListener("submit", (e) => {
    e.preventDefault();

    let searchInput = document.getElementById("searchInput");
    submitQuery(searchInput.value);
})

const body = document.querySelector('body');
body.onload =  submitQuery();