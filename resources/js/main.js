
// Exit gracefully
//
function onWindowClose() {
    Neutralino.app.exit();
}

// Handle incoming PONGs
//
async function onPingResult(e) {
    console.log("DBG RECEIVED: " + e.detail.result );
    let msg = document.getElementById("msg");
    msg.innerHTML += e.detail.result + '<br>';
}

document.getElementById('link-long-run')
    .addEventListener('click', () => {
        GO.run('longRun');
    });


// Init Neutralino
//
Neutralino.init();
Neutralino.events.on("windowClose", onWindowClose);
Neutralino.events.on("pingResult", onPingResult);

// Set title
//
(async () => {
    await Neutralino.window.setTitle(`Neutralino RustExtension ${NL_APPVERSION}`);
    await Neutralino.window.show();
})();

// Init Rust Extension
const GO = new GoExtension(true);