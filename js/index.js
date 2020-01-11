import * as PIXI from 'pixi.js'

const app = new PIXI.Application({
    backgroundColor: "0xffffff",
});

let mech;

function initMechVars() {
    mech.scale.set(0.2, 0.2);
    mech.scale.y *= -1;
    mech.anchor.set(0.5);
    mech.x = app.screen.width / 2 - mech.width / 2;
    mech.y = app.screen.height / 2 - mech.height / 2;
    mech.vx = 0;
    mech.vy = 0;
    mech.vr = 0;
    mech.throttle = 0;
    mech.rotation = 0;
}

function parseResponse(result) {
    console.log(JSON.stringify(result, null, 2))
    let errorContainer = document.getElementById("errorsContainer");
    if (result.error) {
        let errorTextContainer = document.getElementById("errorsText");
        errorTextContainer.innerHTML = result.error.replace(/\n/g, '<br/>');

        errorContainer.style.display = 'block';
    } else {
        errorContainer.style.display = 'none';
        updateMechVars(result)
    }
}

function updateMechVars(result) {
    if (result.vr) {
        let vr = parseFloat(result.vr);
        if (vr === vr) {
            mech.vr = vr;
        }
    }
    if (result.rotation) {
        let rotation = parseFloat(result.rotation);
        if (rotation === rotation) {
            mech.rotation = rotation;
        }
    }
    if (result.throttle) {
        let throttle = parseFloat(result.throttle);
        if (throttle === throttle) {
            mech.throttle = throttle
        }
    }
}

function gameLoop(delta) {
    if (mech.throttle) {
        mech.vx = Math.cos(mech.rotation + Math.PI/2) * mech.throttle;
        mech.vy = Math.sin(mech.rotation + Math.PI/2) * mech.throttle;
    }
    mech.x += mech.vx;
    mech.y += mech.vy;
    mech.rotation += mech.vr;
}

function resetVelocity() {
    mech.vx = 0;
    mech.vy = 0;
    mech.vr = 0;
    mech.throttle = 0;
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
        .then(result => parseResponse(result))
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
