// GoExtension
//
// Run GoExtension functions by sending dispatched event messages.
//
// (c)2024 Harald Schneider - marketmix.com

class GoExtension {
    constructor(debug=false) {
        this.version = '1.0.0';
        this.debug = debug;

        if(NL_MODE !== 'window') {
            window.addEventListener('beforeunload', function (e) {
                e.preventDefault();
                e.returnValue = '';
                GO.stop();
                return '';
            });
        }
    }
    async run(f, p=null) {
        //
        // Call a GoExtension function.

        let ext = 'extGo';
        let event = 'runGo';

        let data = {
            function: f,
            parameter: p
        }

        if(this.debug) {
            console.log(`EXT_GO: Calling ${ext}.${event} : ` + JSON.stringify(data));
        }

        await Neutralino.extensions.dispatch(ext, event, data);
    }

    async stop() {
        //
        // Stop and quit the Bun-extension and its parent app.
        // Use this if Neutralino runs in Cloud-Mode.

        let ext = 'extGo';
        let event = 'appClose';

        if(this.debug) {
            console.log(`EXT_GO: Calling ${ext}.${event}`);
        }
        await Neutralino.extensions.dispatch(ext, event, "");
        await Neutralino.app.exit();
    }
}