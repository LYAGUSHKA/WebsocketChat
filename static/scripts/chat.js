document.querySelector('.menu__container').addEventListener('click', () => {
  document.querySelector('.drop-menu').classList.remove('disabled');  
  document.querySelector('.chats').classList.add('disabled'); 
});

window.onload = function () {
    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");
    var chat_items = document.querySelectorAll('.chats__item'),
    index, chat_item

    for (index = 0; index < chat_items.length; index++) {
        chat_item = chat_items[index];
        chat_item.addEventListener('click', clickHandler);
    }
    
    function clickHandler(event) {
        document.getElementById('log').innerHTML = '';
        var chatID = this.getAttribute("chat-id")
        console.log('ChatID: ', chatID) 
        connectToChat(chatID);
        event.preventDefault();
    }
    
    function connectToChat(chatID) {
        
        conn = new WebSocket("ws://" + document.location.host + "/ws/" + chatID);
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            log.appendChild(item);
        };
        conn.onmessage = function (evt) {
            console.log(evt.data);
            var message = JSON.parse(evt.data);
            var item = document.createElement("div");
            item.innerText = message.data;
            log.appendChild(item);
            
        };
    }

    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        conn.send(JSON.stringify({"data":msg.value}));
        msg.value = "";
        return false;
    };

    if (!window["WebSocket"]) {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        log.appendChild(item);
    }
}


