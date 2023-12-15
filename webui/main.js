const Utils = {
    isScrolledToBottom: function() {
        var currScrollPos = window.scrollY || document.documentElement.scrollTop;
        var totalHeight = document.documentElement.scrollHeight;
        var windowHeight = window.innerHeight || document.documentElement.clientHeight;
        return currScrollPos + windowHeight >= totalHeight;
    },

    scrollToBottom: function() {
        window.scrollTo({
            top: document.documentElement.scrollHeight,
            behavior: 'smooth'
        });
    },
}

const LogTableManager = {

    _logTable: document.getElementById("logTable"),
    _lastLogMessageId: null,

    init: function(logMessages) {

        while(this._logTable.firstChild) {
            this._logTable.removeChild(this._logTable.firstChild)
        }
    
        if (logMessages.length === 0) {
            logMessages.push({"m":"No log messages"});
            this.lastLogMessageId = null;
        }
        
        for (const logMsg of logMessages) {
            var row = this._logTable.insertRow();
            var cell = row.insertCell();
            cell.innerHTML = logMsg.m;

            if (this._lastLogMessageId < logMsg.i) {
                this._lastLogMessageId = logMsg.i;
            }            
        }
    },

    lastLogMessageId: function() {
        return this._lastLogMessageId;
    },

    appendLogs: function(logMessages) {

        var scrollToBottom = Utils.isScrolledToBottom();

        if (0 < logMessages.length) {
        
            for (const logMsg of logMessages) {
                var row = this._logTable.insertRow();
                var cell = row.insertCell();
                cell.innerHTML = logMsg.m;

                if (this._lastLogMessageId < logMsg.i) {
                    this._lastLogMessageId = logMsg.i;
                }
            }    
        }

        if (scrollToBottom) {
            Utils.scrollToBottom();
        }        
    }, 
};

const WebSocketManager = {

    socket: null,

    init: function(callback) {

        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }

        this.socket = new WebSocket("ws://" + location.host + "/events");

        this.socket.onopen = function(event) { 
            console.log('WebSocket connection opened:', event);

            if (callback) {
                console.log(this.socket.readyState);
                callback();
            }
        }.bind(this);

        this.socket.onmessage = function(event) { 

            console.log('Message from server:', event.data);

            const receivedData = JSON.parse(event.data);

            LogTableManager.appendLogs(receivedData);
        };

        this.socket.onclose = function(event) { 
            console.log('WebSocket connection closed:', event);
            LogTableManager.appendLogs([{"m":"WebSocket connection closed: " + event}]); 
        };

        this.socket.onerror = function(event) {
            console.error('WebSocket error:', event);
            LogTableManager.appendLogs([{"m":"Websocket error: " + event}]);
        };

    },

    sendData: function(data) {
        if (this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(data));
        } else {
            console.error("Failed to send data: Websocket connection not open.")
        }
    },
};

//WebSocketManager.init();

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
    .then((json) => {
        LogTableManager.init(json);
        
        dataToSend = { 
            lastLogMessageId: LogTableManager.lastLogMessageId(), 
            queryString: queryMsg,
            timestamp: Date.now(),
        };

        WebSocketManager.init(() => {
            WebSocketManager.sendData(dataToSend);
        });        

        WebSocketManager.sendData(dataToSend);
    })
    .catch((error) => {
        footer.textContent = `Could not fetch log messages: ${error}`;
    })
}

let searchForm = document.getElementById("searchForm");
searchForm.addEventListener("submit", (e) => {
    e.preventDefault();

    let searchInput = document.getElementById("searchInput");
    submitQuery(searchInput.value);
})

const body = document.querySelector('body');
const footer = document.querySelector('footer');

body.onload =  function() {

    submitQuery();
};