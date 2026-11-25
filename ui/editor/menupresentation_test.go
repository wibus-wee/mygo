package editor

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
)

// The editor's context menu draws in the window while the window draws its
// menus (ui.Context.SetMenuPresentation).
func TestCodeEditorMenuDrawnInTheWindow(t *testing.T) {
	st := &State{Buffer: NewBuffer("hello world")}
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetMenuPresentation(ui.PopoverMenus)
		View(c, st).Fill().AutoFocus()
	}, 600, 400)
	tt.Frame()
	tt.RightClickAt(300, 200)
	if menu := tt.Menu(); !slices.Contains(menu, "Select All") || !tt.HasText("Select All") {
		t.Fatalf("Menu() %q, drawn %v", menu, tt.HasText("Select All"))
	}
	if err := tt.ChooseMenuItem("Select All"); err != nil {
		t.Fatal(err)
	}
	ct := &editorTest{t: t, st: st, tt: tt}
	ct.want("[hello world]")
	if tt.HasText("Select All") {
		t.Fatal("the menu stayed open after a choice")
	}
}

func TestEditorContextMenuReplacement(t *testing.T) {
	for _, test := range []struct {
		name         string
		replace      bool
		keyboard     bool
		defaultItems bool
	}{
		{name: "appends by default", defaultItems: true},
		{name: "replaces on pointer", replace: true},
		{name: "replaces from keyboard", replace: true, keyboard: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			st := &State{
				Buffer:             NewBuffer("hello world"),
				ReplaceContextMenu: test.replace,
				ContextMenu:        func(m *ui.Menu) { m.Item("App action") },
			}
			tt := ui.NewTester(func(c *ui.Context) {
				View(c, st).Fill().AutoFocus()
			}, 600, 400)
			tt.Frame()
			if test.keyboard {
				tt.Key(0, ui.KeyContextMenu)
			} else {
				tt.RightClickAt(300, 200)
			}
			menu := tt.Menu()
			if !slices.Contains(menu, "App action") {
				t.Fatalf("menu %q does not contain the app item", menu)
			}
			if test.replace && !slices.Equal(menu, []string{"App action"}) {
				t.Fatalf("replacement menu = %q, want only the app item", menu)
			}
			if !test.replace && slices.Index(menu, "Cut") > slices.Index(menu, "App action") {
				t.Fatalf("default items did not precede the app item: %q", menu)
			}
			if got := slices.Contains(menu, "Cut"); got != test.defaultItems {
				t.Fatalf("default Cut item present = %v, want %v; menu %q", got, test.defaultItems, menu)
			}
		})
	}
}
