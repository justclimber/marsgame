import * as PIXI from 'pixi.js'

const app = new PIXI.Application({
    backgroundColor: "0xffffff",
});

let mech, mechBase, mechWeaponCannon;

function initMechVars() {
    mech.scale.set(0.2, 0.2);
    mech.scale.y *= -1;
    mech.pivot.set(0.5);
    mech.x = app.screen.width / 2 - mech.width / 2;
    mech.y = app.screen.height / 2 - mech.height / 2;
    mech.vx = 0;
    mech.vy = 0;
    mech.vr = 0;
    mech.throttle = 0;
    mech.rotation = 0;

    // смещаем башню немного, потому что она не по центру меха
    mechWeaponCannon.y = 30;
    mechWeaponCannon.vr = 0;
    mechWeaponCannon.rotation = 0;
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
        mech.rotation = fetchFloatOr0(result.rotation);
    }
    if (result.throttle) {
        mech.throttle = fetchFloatOr0(result.throttle);
    }
    if (result.cannonVr) {
        mechWeaponCannon.vr = fetchFloatOr0(result.cannonVr);
    }
}

function fetchFloatOr0(value) {
    let floatVal = parseFloat(value);
    if (floatVal === floatVal) {
        return floatVal;
    } else {
        return 0;
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
    mechWeaponCannon.rotation += mechWeaponCannon.vr
}

function resetVelocity() {
    mech.vx = 0;
    mech.vy = 0;
    mech.vr = 0;
    mech.throttle = 0;
    mechWeaponCannon.vr = 0;
}

window.onload = function() {
    document.getElementById('pixiDiv').appendChild(app.view);

    app.loader
        .add('mechBase', '/images/mech_base.png')
        .add('mechWeaponCannon', '/images/mech_weapon_cannon.png')
        .load((loader, resources) => {
            mechBase = new PIXI.Sprite(resources.mechBase.texture);
            mechWeaponCannon = new PIXI.Sprite(resources.mechWeaponCannon.texture);
            mechBase.anchor.set(0.5);
            mechWeaponCannon.anchor.set(0.5, 0.6);
            mech = new PIXI.Container();
            mech.addChild(mechBase)
            mech.addChild(mechWeaponCannon)
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
        fetch("save_source_code", {
            method: "POST",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                userId: userId,
                sourceCode: document.getElementById('sourceCode').value
            })
        }).then(response => response.json())
        .then(result => parseResponse(result))
    };
};
let userId = getUserId();
let url = "ws://localhost/ws?id=" + userId;
let socket = new WebSocket(url);
console.log("Connection to websocket", url);

socket.onopen = () => {
    console.log("Connection success");
    let command = {
        "type": "greetings",
        "payload": "Hi from the client!",
    };
    socket.send(JSON.stringify(command));
};
socket.onmessage = (msg) => {
    if (msg.data) {
        let data = JSON.parse(msg.data);
        if (data.type && data.payload) {
            let payload = JSON.parse(data.payload);
            console.log(data.type, data.payload);
        } else {
            console.log(data);
        }
    } else {
        console.log(msg);
    }
};
socket.onclose = (event) => {
    console.log("Socket connection closed: ", event);
};
socket.onerror = (error) => {
    console.log("Socket error: ", error);
};

function getUserId() {
    return Math.random().toString(36).replace(/[^a-z]+/g, '').substr(0, 5);
}