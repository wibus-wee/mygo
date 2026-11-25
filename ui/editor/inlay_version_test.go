package editor

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// Inlays of an older version move with the edits since, as a language
// server's that arrive late.
func TestInlaysMoveWithEditsSinceTheirVersion(t *testing.T) {
	st := &State{Buffer: NewBuffer("abcde\n")}
	v0 := st.Buffer.Version()
	tt := ui.NewTester(func(c *ui.Context) { View(c, st).Fill() }, 200, 100)
	tt.Frame()

	// Two bytes inserted before the inlay after it was requested.
	st.Buffer.Edit(Edit{Start: 0, End: 0, Text: "XY"})
	tt.Frame()

	st.SetInlays(v0, []Inlay{{Offset: 3, Text: ": T"}})
	tt.Frame()
	if got := st.ed.view.inlays[0].at; got != 5 {
		t.Fatalf("inlay at %d, want 5 after a two byte insertion before it", got)
	}
}

// Inlays of a version the buffer's log no longer keeps leave the inlays
// as they are rather than move to where the text has moved on.
func TestInlaysOfAnUnknownVersionAreDropped(t *testing.T) {
	st := &State{Buffer: NewBuffer("abcde\n")}
	tt := ui.NewTester(func(c *ui.Context) { View(c, st).Fill() }, 200, 100)
	tt.Frame()

	// An inlay of the text as it is shows and moves with the edits.
	v0 := st.Buffer.Version()
	st.SetInlays(v0, []Inlay{{Offset: 3, Text: ": T"}})
	tt.Frame()
	if len(st.ed.view.inlays) != 1 {
		t.Fatalf("no inlay shown: %d", len(st.ed.view.inlays))
	}

	// Outrun the log: each edit is a version, and the log keeps 256.
	for range 300 {
		st.Buffer.Edit(Edit{Start: 0, End: 0, Text: "x"})
	}
	tt.Frame()

	// A reply so late its version fell out of the log neither replaces
	// nor clears what the editor shows.
	st.SetInlays(0, []Inlay{{Offset: 100, Text: ": stale"}})
	tt.Frame()
	if len(st.ed.view.inlays) != 1 || st.ed.view.inlays[0].text != ": T" {
		t.Fatalf("a stale version changed the inlays: %+v", st.ed.view.inlays)
	}
}
