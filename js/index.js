import * as PIXI from 'pixi.js'

const app = new PIXI.Application({
    backgroundColor: "0xffffff"
});

let mech;

function initMechVars() {
    mech.scale.set(0.3, 0.3);
    mech.anchor.set(0.5);
    mech.x = 100;
    mech.y = 100;
    mech.vx = 0;
    mech.vy = 0;
    mech.vr = 0;
    mech.rotation = Math.PI;
}

function updateMechVars(result) {
    console.log(JSON.stringify(result, null, 2))
    if (result.vr) {
        let vr = parseFloat(result.vr);
        if (vr === vr) {
            mech.vr = vr;
        }
    }
    if (result.vx) {
        let vx = parseFloat(result.vx);
        if (vx === vx) {
            mech.vx = vx;
        }
    }
    if (result.vy) {
        let vy = parseFloat(result.vy);
        if (vy === vy) {
            mech.vy = vy;
        }
    }
}

function gameLoop(delta) {
    mech.x += mech.vx;
    mech.y += mech.vy;
    mech.rotation += mech.vr;
}

function resetVelocity() {
    mech.vx = 0;
    mech.vy = 0;
    mech.vr = 0;
}

window.onload = function() {
    document.getElementById('pixiDiv').appendChild(app.view);

    app.loader
        .add('mech', '/images/mech_base.png')
        .load((loader, resources) => {
        mech = new PIXI.Sprite(resources.mech.texture);
        initMechVars();
        app.stage.addChild(mech);
        app.ticker.add(delta => gameLoop(delta));
    });

    let stopMechButton = document.getElementById('stopMech');
    stopMechButton.onclick = resetVelocity;

    let resetVarsButton = document.getElementById('resetVars');
    resetVarsButton.onclick = initMechVars;

    let saveCodeButton = document.getElementById('saveCode');
    saveCodeButton.onclick = () => {
        fetch("http://localhost/save_source_code", {
            method: "POST",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(document.getElementById('sourceCode').value)
        }).then(response => response.json())
        .then(result => updateMechVars(result))
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
