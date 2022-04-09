let username = ""
let ConnectionBtn = document.getElementById("connect")
let userInput = document.getElementById("username")
let form = document.getElementById("chat-form")
let socket

form?.addEventListener("submit", function(event) {
    event.preventDefault()
    const data = {
        sender: username ,
        target: form.target.value,
        body: form.message.value,
    }

    console.log("-> :::",data)

    socket.send(JSON.stringify(data))

})

ConnectionBtn?.addEventListener('click',() => {

    username = userInput.value;
    
    socket = new WebSocket(`ws://localhost:8080/chat?username=${username}`)
    
    socket.addEventListener('message', function(event) {
        console.log('message from server', JSON.parse(event.data))
    })
})

    //                 json:"id,omitempty"`
	// Sender string `json:"sender"`
	// Target string `json:"target"`
	// Body   string `json:"body"`