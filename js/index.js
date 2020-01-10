import * as PIXI from 'pixi.js'

const app = new PIXI.Application();

window.onload = function() {
    document.body.appendChild(app.view);

    app.loader.add('mech', '/images/mech_base.png').load((loader, resources) => {
        // This creates a texture from a 'mech.png' image
        const mech = new PIXI.Sprite(resources.mech.texture);

        mech.scale.set(0.3, 0.3);
        mech.anchor.set(0.5);
        mech.x = 100;
        mech.y = 100;
        mech.rotation = Math.PI;
        app.stage.addChild(mech);
    });

    let saveCodeButton = document.getElementById('saveCode');
    saveCodeButton.onclick = () => {
        fetch("save_source_code", {
            method: "POST",
            body: document.getElementById('sourceCode')
        }).then((response) => {
            if (response.status !== 200) {
                return Promise.reject();
            }
        }).catch(() => console.log("Sending failed"))
    };
};

let socket = new WebSocket("ws://localhost/ws");
console.log("Connection to websocket");

socket.onopen = () => {
    console.log("Connection success");
    socket.send("Hi from the client!");
};
socket.onmessage = (msg) => {
    console.log(msg);
};
socket.onclose = (event) => {
    console.log("Socket connection closed: ", event);
};
socket.onerror = (error) => {
    console.log("Socket error: ", error);
};
