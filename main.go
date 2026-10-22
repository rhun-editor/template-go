// A rhun extension, as small as one gets: commands, a toast and an edit, and greet, which knows
// where it was run from (a key, a menu item, the status bar) and gets that place's arguments.
// extension.ini declares its key, its menu item, its status item and language keys.
// Copy this repository, rename "hello" everywhere, and change what Invoke does.
// Protocol: https://rhun.app/extensions/v1/ and v1.1/ (the rhunext package does the framing).
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

// Invoke gets every run of a command, with where it came from (inv.From) and what that place
// passes (inv.Body): a key binding's arguments, a menu item's target and arguments, a command
// reference's arguments (the status item). It needs the command-args feature (main).
func (hello) Invoke(ext *rhunext.Ext, inv rhunext.Invocation) {
	switch inv.Name() {
	case "say":
		ext.Toast("Hello from an extension")
	case "upper":
		upper(ext, inv.Doc)
	case "greet":
		who := "you"
		if len(inv.Body) > 0 {
			who = string(inv.Body)
		}
		if target, args, ok := inv.Target(); ok {
			// a menu item: the target line ("editor VERSION L1 C1 L2 C2" here), then its args
			who = string(args)
			if target.Kind == "editor" && len(target.Fields) > 1 {
				who += " at line " + target.Fields[1]
			}
		}
		ext.Toast("Hello, " + who + " (from " + inv.From + ")")
	}
}

// Run is how a rhun without command-args runs commands: by their full name, NAME.key.
func (h hello) Run(ext *rhunext.Ext, command string, doc rhunext.Doc) {
	h.Invoke(ext, rhunext.Invocation{Command: command, Doc: doc})
}

func upper(ext *rhunext.Ext, doc rhunext.Doc) {
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

func (hello) Event(ext *rhunext.Ext, ev rhunext.Event) {
	// opened, changed, saved, cursor, selection, closed, active, project: see the spec
}

func main() {
	if err := rhunext.RunWith([]string{rhunext.FeatureCommandArgs}, hello{}); err != nil {
		os.Exit(1)
	}
}
