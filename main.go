// A rhun extension, as small as one gets: two commands, a toast and an edit.
// Copy this repository, rename "hello" everywhere, and change what Run does.
// Protocol: https://rhun.app/extensions/v1/ (the rhunext package does the framing).
package main

import (
	"bytes"
	"os"

	"github.com/rhun-editor/extension-protocol/go/rhunext"
)

type hello struct{}

func (hello) Ready(ext *rhunext.Ext, project string) {
	ext.Logf("ready, project %q", project)
}

func (hello) Run(ext *rhunext.Ext, command string, doc rhunext.Doc) {
	switch command {
	case "hello.say":
		ext.Toast("Hello from an extension")
	case "hello.upper":
		if doc == 0 {
			ext.Toast("No document")
			return
		}
		text, version, err := ext.Text(doc)
		if err != nil {
			ext.Logf("text: %v", err)
			return
		}
		upper := bytes.ToUpper(text)
		if bytes.Equal(upper, text) {
			return
		}
		if _, err := ext.Replace(doc, version, rhunext.WholeRange(text), upper); err != nil {
			ext.Logf("replace: %v", err) // stale: the user typed meanwhile; just give up
		}
	}
}

func (hello) Event(ext *rhunext.Ext, ev rhunext.Event) {
	// opened, changed, saved, cursor, selection, closed, active, project: see the spec
}

func main() {
	if err := rhunext.Run(hello{}); err != nil {
		os.Exit(1)
	}
}
