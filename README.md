<p align="center">
<img src="https://marketmix.com/git-assets/neutralino-ext-go/neutralino-go-header-2.svg">
</p>

# neutralino-ext-go

**A Go Extension for Neutralino >= 5.0.0**

This extension adds a Go backend to Neutralino with the following features:
- Requires only a few lines of code on both ends.
- Read all events from the Neutralino app in your Go code.
- Run Go functions from Neutralino.
- Run Neutralino functions from Go.
- All communication between Neutralino and Go runs asynchronously.
- All events are queued, so none will be missed during processing.
- Track the data flow between Neutralino and Go in realtime.
- Works in Window- and headless Cloud-Mode.
- Terminates the Go Runtime when the Neutralino app quits.

![Neutralino Go Extension](https://marketmix.com/git-assets/neutralino-ext-go/go-neutralino.gif)

## Run the demo

The demo opens a Neutralino app. Clicking on the blue link sends a Ping to Go, which replies with Pong.
This illustrates the data-flow in both directions. 

Before running the demo, the Go extension needs to be compiled with Go. Make this folder the project root for your Go-compiler:
```bash
/extensions/go
```
Then build with:
```bash
# macOS and Linux
./build.sh
# Windows
build.cmd
```
The demo is configured to launch the Go-extension binary directly from the source-folder.

After this, run these commands in the **ext-go folder:
```commandline
neu update
neu run
```

## Integrate into your own project
Follow these steps:
- Adapt the Go code in **extensions/go/main.go** to your needs.
- Build the Go-binary.
- Create an empty **/extensions/go** folder, used by your installer.
- Copy the Go-binary to **/extensions/go**
- Copy this **/extensions** folder to your project.
- Copy **resources/js/go-extension.js** to **resources/js**.
- Add `<script src="js/go-extension.js"></script>` to your **index.html**
- Add `const GO = new GoExtension(true)` to your **main.js**
- Add **GO.run(function_name, data) to main.js** to run Go-functions from Neutralino.
- Add **event listeners to main.js**, to fetch result data from Go.
- Modify **neutralino.config.json** (see below).

Make sure that **neutralino.config.json** contains this, adapted to your environment:
```json
  "extensions": [
    {
      "id": "extGo",
      "commandDarwin": "${NL_PATH}/extensions/go/go ${NL_PATH}",
      "commandLinux": "${NL_PATH}/extensions/go/go ${NL_PATH}",
      "commandWindows": "${NL_PATH}/extensions/go/go.exe ${NL_PATH}"
    }
  ],
```

## ./extensions/go/main.go explained

```js
package main

import (
	"fmt"
	"neutralino-extension/neutralino-extension"
	"time"
)

var ext = new(neutralino_extension.WSClient)

func processAppEvent(data neutralino_extension.EventMessage) {

	if ext.IsEvent(data, "runGo") {
		if d, ok := data.Data.(map[string]interface{}); ok {

			if d["function"] == "ping" {
				var out = make(map[string]interface{})
				out["result"] = fmt.Sprintf("Go says PONG in reply to '%s'", d["parameter"])
				ext.Send("pingResult", out)
			}

			// More functions here:
			...
		}
	}
}

func main() {
	ext.Run(processAppEvent, true)
}


```

The extension is activated in main(). 
**processAppEvent** is a callback function, which is triggered with each event coming from the Neutralino app.

In the callback function, you can process the incoming events by their name. In this case we react to the **"runGo"** event.
**d["function"]** holds the requested Rust-function and **d["parameter"]** its data payload as map derived from JSON.

If the requested function is named **ping**, we send back a message to the Neutralino frontend. 

**Send()** requires the following parameters:

- An event name, here "pingResult"
- The data package as a hash-map to send.

YOu can also use **SendString()** if you want to reply with a plain string.

## ./resources/js/main.js explained

```JS

async function onPingResult(e) {
...
}

// Init Neutralino
//
Neutralino.init();
...
Neutralino.events.on("pingResult", onPingResult);
...
// Init Bun Extension
const GO = new GoExtension(true)
```

The last line initializes the JavaScript part of the Go-extension. It's important to place this after Neutralino.init() and after all event handlers have been installed. Put it in the last line of your code and you are good to go. The const **GO** is accessible globally and **must not be renamed** to something else.

The **GoExtension class** takes only 1 argument which instructs it to run in debug mode (here true). In this mode, all data coming from the extension is printed to the dev-console:

![Debug Meutralino](https://marketmix.com/git-assets/neutralino-ext-go/go-console.jpg)

The **pingResult event handler** listens to messages with the same name, sent by send_message() on Go's side. 

In **index.html**, you can see how to send data from Neutralino to Go, which is dead simple:
```html
<a href="#" onclick="GO.run('ping', 'Neutralino says PING!');">Send PING to Go</a><br>
```

**GO.run()** takes 2 arguments:

- The Go function to call, here "ping"
- The data package to submit, either as string or JSON.

### Long-running tasks and their progress

For details how to start a long-running background task in Go and how to poll its progress,
see the comments in `extensions/go/main.go`and `resources/js/main.js`.

## Modules & Classes Overview

### neutralino-extension.go

| Method                  | Description                                                  |
| ----------------------- | ------------------------------------------------------------ |
| run(callback, debug)    | Starts the extensions main processing loop. Each incoming message triggers the `callback` function. If `debug` is true, all events are printed to the IDE's console. |
| callback(e)             | The callback function referenced by `ext.run(callback, debug)`.<br>e: The event message of type 'EventMessage'. |
| IsEvent(d, e)           | Checks the incoming event data-package for a particular event-name.<br>d: Data-package of type `EventMessage`<br />e: Event-name as `string` |
| SendMessage(e, d)       | Send an event-message to Neutralino. <br>e: Event-name as `string`<br>d: Data package as `map[string]interface{}` |
| SendMessageString(e, d) | Send an event-message to Neutralino. Uses a plain message-string instead of mapped data.<br/>e: Event-name as `string`<br/>d: Data package as `string`. |

### go-extension.js

GoExtension Class:

| Method               | Description                                                  |
| -------------------- | ------------------------------------------------------------ |
| async run(f, p=null) | Call a Go-function.<br>f: Function-name.<br>p: Parameter data package as string or JSON. |
| async stop()         | Stop and quit the Go-extension and its parent app. Use this if Neutralino runs in Cloud-Mode. This is called automatically, when the browser tab is closed. |

| Property | Description                                        |
| -------- | -------------------------------------------------- |
| debug    | If true,  data flow is printed to the dev-console. |

Events, sent from the frontend to the extension:

| Event    | Description                                                  |
| -------- | ------------------------------------------------------------ |
| appClose | Notifies the extension, that the app will close. This quits the extension. |

## More about Neutralino

- [NeutralinoJS Home](https://neutralino.js.org) 

- [Neutralino Build Automation for macOS, Windows, Linux](https://github.com/hschneider/neutralino-build-scripts)

- [Neutralino related blog posts at marketmix.com](https://marketmix.com/de/tag/neutralinojs/)



<img src="https://marketmix.com/git-assets/star-me-2.svg">

