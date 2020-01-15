import * as PIXI from 'pixi.js'
import { Viewport } from 'pixi-viewport'

const app = new PIXI.Application({
    backgroundColor: "0xffffff",
    width: 880,
});

let viewport;

let mech, mechBase, mechWeaponCannon, terra;
let xShift = 1000;
let yShift = 1000;

function initMechVars() {
    mech.scale.y *= -1;
    mech.pivot.set(0.5);
    mech.x = xShift;
    mech.y = yShift;
    mech.vx = 0;
    mech.vy = 0;
    mech.vr = 0;
    mech.throttle = 0;
    mech.rotation = 0;

    // смещаем башню немного, потому что она не по центру меха
    mechWeaponCannon.y = 6;
    mechWeaponCannon.x = 45;
    mechWeaponCannon.vr = 0;
    mechWeaponCannon.rotation = 0;
}

let changelogToRun = [];
let timeShiftForPrediction = 2000;
function parseChangelog(changelog) {
    // console.log(changelog);
    changelog.forEach(function (changeByTime) {
        let changesByObject = changeByTime.changesByObject;
        changesByObject.forEach(function (changeByObj) {
            if (changeByObj.objId !== userId) {
                return;
            }
            let change = {timeId: changeByTime.timeId};
            if (changeByObj.pos) {
                change.x = changeByObj.pos.x + xShift;
                change.y = changeByObj.pos.y + yShift;
            }
            if (changeByObj.angle) {
                change.rotation = changeByObj.angle;
            }
            changelogToRun.push(change);
        });
        if (!currTimeId) {
            // use time shift for more smooth prediction: we need changelogToRun always be not empty on run
            currTimeId = changeByTime.timeId - timeShiftForPrediction;
        }
    });
}

let timer = new Date();
let currTimeId;
function gameLoop(delta) {
    mech.x += mech.vx;
    mech.y += mech.vy;
    mech.rotation += mech.vr;

    let now = new Date();
    let timeDelta = now.getTime() - timer.getTime();
    timer = now;
    if (currTimeId) {
        currTimeId += timeDelta;
        if (changelogToRun.length) {
            if (changelogToRun[0].timeId < currTimeId) {
                let timeId = changelogToRun[0].timeId;
                let change = changelogToRun.shift();
                if (change.x) {
                    mech.x = change.x;
                }
                if (change.y) {
                    mech.y = change.y;
                }
                if (change.rotation) {
                    mech.rotation = change.rotation;
                }

                // prediction for smooth moving
                if (changelogToRun.length) {
                    let nextChange = changelogToRun[0];
                    let nextTimeIdDelta = nextChange.timeId - timeId;
                    let futureGameTicks = nextTimeIdDelta / timeDelta;
                    mech.vx = !nextChange.x ? 0 : (nextChange.x - mech.x) / futureGameTicks;
                    mech.vy = !nextChange.y ? 0 : (nextChange.y - mech.y) / futureGameTicks;
                    mech.vr = !nextChange.rotation ? 0 : (nextChange.rotation - mech.rotation) / futureGameTicks;
                }
            }
        } else {
            // stop prediction
            mech.vx = 0;
            mech.vy = 0;
            mech.vr = 0;
        }
    }
}

function viewportSetup() {
    viewport = new Viewport({
        screenWidth: 880,
        screenHeight: 600,
        worldWidth: 3000,
        worldHeight: 2000,

        // the interaction module is important for wheel to work properly when renderer.view is placed or scaled
        interaction: app.renderer.plugins.interaction
    });
    viewport.clampZoom({
        minWidth: 300,
        maxWidth: 3000,
    }).bounce({
        time: 400
    })
        .moveCenter(xShift, yShift)
        .drag()
        .pinch()
        .wheel()
        .decelerate();
}

function mechSetup(resources) {
    let sheet = resources["/images/spritesheet.json"];
    mechBase = new PIXI.Sprite(sheet.textures['mech_base_128.png']);
    mechWeaponCannon = new PIXI.Sprite(sheet.textures['cannon_128.png']);
    terra = new PIXI.TilingSprite(sheet.textures['terra_256.png'], 2800, 2000);
    terra.anchor.set(0);
    viewport.addChild(terra);
    mechBase.anchor.set(0.5);
    mechWeaponCannon.anchor.set(0.5, 0.6);
    mech = new PIXI.Container();
    mech.addChild(mechBase);
    mech.addChild(mechWeaponCannon);
    initMechVars();
}

window.onload = function() {
    document.getElementById('pixiDiv').appendChild(app.view);
    viewportSetup();

    app.loader
        .add('/images/spritesheet.json')
        .load((loader, resources) => {
            mechSetup(resources);
            app.stage.addChild(viewport);
            viewport.addChild(mech);

            app.ticker.add(delta => gameLoop(delta));
        });

    let resetVarsButton = document.getElementById('resetVars');
    resetVarsButton.onclick = initMechVars;

    let sourceCodeEl = document.getElementById('sourceCode');
    let sourceCodeFromLocalStorage = localStorage.getItem('sourceCode');
    if (sourceCodeFromLocalStorage && sourceCodeFromLocalStorage.length > 0) {
        sourceCodeEl.value = sourceCodeFromLocalStorage;
    }

    let saveCodeButton = document.getElementById('saveCode');
    saveCodeButton.onclick = saveCode;

    let runProgramButton = document.getElementById('runProgram');
    runProgramButton.onclick = runProgram;

    let stopProgramButton = document.getElementById('stopProgram');
    stopProgramButton.onclick = stopProgram;

    let autoSaveCheckbox = document.getElementById('autoSaveCheckbox');
    autoSaveCheckbox.onchange = function () {
        if (this.checked) {
            localStorage.setItem('autoSave', "true")
            document.getElementById('autoStartSpan').classList.remove('disabled');
            autoStartCheckbox.disabled = '';
        } else {
            document.getElementById('autoStartSpan').classList.add('disabled');
            autoStartCheckbox.disabled = 'disabled';
            localStorage.removeItem('autoSave')
        }
    };

    let autoStartCheckbox = document.getElementById('autoStartCheckbox');
    autoStartCheckbox.onchange = function () {
        if (this.checked) {
            localStorage.setItem('autoStart', "true")
        } else {
            localStorage.removeItem('autoStart')
        }
    };

    if (localStorage.getItem('autoSave')) {
        autoSaveCheckbox.checked = true;
        setTimeout(saveCode, 500);
        if (localStorage.getItem('autoStart')) {
            autoStartCheckbox.checked = true;
            setTimeout(runProgram, 1500);
        }
    } else {
        autoStartCheckbox.disabled = 'disabled';
        document.getElementById('autoStartSpan').classList.add('disabled');
    }


};

function saveCode() {
    document.getElementById("errorsContainer").style.display = 'none';

    let sourceCode = document.getElementById('sourceCode').value;
    localStorage.setItem('sourceCode', sourceCode);
    fetch("save_source_code", {
        method: "POST",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            userId: userId,
            sourceCode: sourceCode
        })
    }).then(function (response) {

    })
}

function parseError(payload) {
    let errorContainer = document.getElementById("errorsContainer");
    let errorTextContainer = document.getElementById("errorsText");
    errorTextContainer.innerHTML = payload.message.replace(/\n/g, '<br/>');
    errorContainer.style.display = 'block';
}

function runProgram() {
    programFlow(1)
}

function stopProgram() {
    programFlow(0)
}

function programFlow(flowCmd) {
    fetch("program_flow", {
        method: "POST",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            userId: userId,
            flowCmd: flowCmd
        })
    }).then(function (response) {

    })
}

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
            if (data.type === 'worldChanges') {
                parseChangelog(payload)
            } else if (data.type === 'error') {
                parseError(payload)
            }
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