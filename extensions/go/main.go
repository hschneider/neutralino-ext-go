package main

import (
	"fmt"
	"neutralino-extension/neutralino-extension"
	"time"
)

const extDebug = true

var ext = new(neutralino_extension.WSClient)

func processAppEvent(data neutralino_extension.EventMessage) {

	// Check if the frontend requests a function-call via runGo-event.
	// If so, extract the embedded function-data.
	//
	if ext.IsEvent(data, "runGo") {
		if d, ok := data.Data.(map[string]interface{}); ok {

			if d["function"] == "ping" {

				// Sends back a map with a result field.
				// You can also use SendMessageString() to send plain strings or stringified JSON.
				//
				var out = make(map[string]interface{})
				out["result"] = fmt.Sprintf("Go says PONG in reply to '%s'", d["parameter"])
				ext.Send("pingResult", out)
			}

			if d["function"] == "longRun" {

				// This starts a long-running background-task, which reports
				// its progress to the frontend;
				//
				go longRun()
			}
		}
	}
}

func longRun() {
	for i := 1; i <= 10; i++ {
		s := fmt.Sprintf("Long running task progress %d / 10", i)
		// Send a plain string
		ext.SendMessageString("pingResult", s)
		time.Sleep(time.Second * 1)
	}
}

func main() {
	ext.Run(processAppEvent, extDebug)
}
