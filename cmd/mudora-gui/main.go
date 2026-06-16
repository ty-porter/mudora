package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ty-porter/mudora/internal"
	"github.com/ty-porter/mudora/internal/alttp"
	"github.com/ty-porter/mudora/internal/icons"
	"github.com/ty-porter/mudora/internal/rom"
	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

const (
	sectionBg = "#3f3f3f"
	sectionFg = "#e6e6e6"
	surround  = "#2b2b2b"
)

var windowTitle string = fmt.Sprintf("Mudora %s - ALttPR ROM Inspector", internal.Version)

func main() {
	App.WmTitle(windowTitle)

	list := buildUI()

	ActivateTheme("azure dark")

	if len(os.Args) > 1 {
		list.load(os.Args[1])
	}
	if png, ok := icons.PNG("Book of Mudora"); ok {
		App.IconPhoto(NewPhoto(Data(png)))
	}
	WmGeometry(App, "1000x800")
	App.Wait()
}

type regionList struct {
	canvas  *CanvasWidget
	inner   *FrameWidget
	search  *TEntryWidget
	groups  []alttp.Group
	row     int
	created []*Window
}

type regionSection struct {
	name     string
	expanded bool
	toggle   *LabelWidget
	rows     [][3]*Window
	rowIndex []int
}

func buildUI() *regionList {
	fr := TFrame()
	cv := fr.Canvas(Background("#2b2b2b"), Highlightthickness(0))
	sb := fr.TScrollbar()
	inner := cv.Frame(Background(surround))
	cv.CreateWindow(0, 0, ItemWindow(inner.Window), Anchor("nw"))
	GridColumnConfigure(inner.Window, 0, Weight(1))

	l := &regionList{canvas: cv, inner: inner}

	bar := fr.TFrame()
	Pack(bar, Side("top"), Fill("x"), Padx("2m"), Pady("2m"))
	open := bar.TButton(Txt("Open ROM..."), Command(func() { chooseAndLoad(l) }))
	Pack(open, Side("left"))
	l.search = bar.TEntry(Width(30), Textvariable(""))
	Pack(l.search, Side("right"))
	Pack(bar.TLabel(Txt("Search item: ")), Side("right"))
	Bind(l.search, "<KeyRelease>", Command(func() { l.render(l.search.Textvariable()) }))

	Pack(sb, Side("right"), Fill("y"))
	Pack(cv, Side("left"), Expand(true), Fill("both"))
	cv.Configure(Yscrollcommand(func(e *Event) { e.ScrollSet(sb) }))
	sb.Configure(Command(func(e *Event) { e.Yview(cv.Window) }))

	Bind(inner, "<Configure>", Command(func() {
		cv.Configure(Scrollregion(strings.Join(cv.Bbox("all"), " ")))
	}))

	Pack(fr, Expand(true), Fill("both"), Padx("2m"), Pady("2m"))
	return l
}

func chooseAndLoad(l *regionList) {
	paths := GetOpenFile(
		Title("Open ALttPR ROM"),
		Filetypes([]FileType{
			{TypeName: "SNES ROM", Extensions: []string{".sfc", ".smc"}},
			{TypeName: "All files", Extensions: []string{"*"}},
		}),
	)
	if len(paths) == 0 {
		return
	}
	l.load(paths[0])
}

func (l *regionList) load(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		MessageBox(Icon("error"), Title("Mudora"),
			Msg(fmt.Sprintf("Could not read ROM:\n%s", err)))
		return
	}

	l.groups = alttp.Grouped(rom.Inspect(data))
	l.render(l.search.Textvariable())
	App.WmTitle(windowTitle + " (inspecting " + path + ")")
}

func (l *regionList) render(query string) {
	l.clear()
	collapsed := strings.TrimSpace(query) == ""
	for _, g := range alttp.Filter(l.groups, query) {
		l.addRegion(g, collapsed)
	}
}

func (l *regionList) clear() {
	for _, w := range l.created {
		Destroy(w)
	}
	l.created = nil
	l.row = 0
}

func (l *regionList) addRegion(g alttp.Group, collapsed bool) {
	sec := &regionSection{name: g.Region, expanded: true}

	box := l.inner.Frame(Background(sectionBg), Relief("solid"), Borderwidth(1), Padx(6), Pady(4))
	l.created = append(l.created, box.Window)
	Grid(box, Row(l.row), Column(0), Sticky("we"), Padx("2m"), Pady("1m"))
	l.row++

	hdr := box.Frame(Background(sectionBg))
	Grid(hdr, Row(0), Column(0), Columnspan(3), Sticky("we"))

	sec.toggle = hdr.Label(Txt("▾  "+g.Region), Background(sectionBg), Foreground(sectionFg), Width(regionNameWidth), Anchor("w"))
	Pack(sec.toggle, Side("left"))
	for _, p := range g.Locations {
		if alttp.IsProgression(p.Item) {
			if img := iconFor(p.Item); img != nil {
				Pack(hdr.Label(Image(img), Background(sectionBg)), Side("left"), Padx(1))
			}
		}
	}
	Bind(hdr, "<Button-1>", Command(func() { l.toggle(sec) }))
	Bind(sec.toggle, "<Button-1>", Command(func() { l.toggle(sec) }))

	for i, p := range g.Locations {
		r := i + 1

		locLbl := box.Label(Txt(p.Location), Background(sectionBg), Foreground(sectionFg), Anchor("w"))
		iconLbl := box.Label(Background(sectionBg))
		if img := iconFor(p.Item); img != nil {
			iconLbl.Configure(Image(img))
		}
		itemLbl := box.Label(Txt(p.Item), Background(sectionBg), Foreground(sectionFg), Anchor("w"))

		Grid(locLbl, Row(r), Column(0), Sticky("w"), Padx("48 0"))
		Grid(iconLbl, Row(r), Column(1), Padx(6))
		Grid(itemLbl, Row(r), Column(2), Sticky("w"))

		sec.rows = append(sec.rows, [3]*Window{locLbl.Window, iconLbl.Window, itemLbl.Window})
		sec.rowIndex = append(sec.rowIndex, r)
	}

	if collapsed {
		l.toggle(sec)
	}
}

func (l *regionList) toggle(sec *regionSection) {
	sec.expanded = !sec.expanded
	if !sec.expanded {
		sec.toggle.Configure(Txt("▸  " + sec.name))
		for _, row := range sec.rows {
			GridForget(row[0], row[1], row[2])
		}
		return
	}
	sec.toggle.Configure(Txt("▾  " + sec.name))
	for i, row := range sec.rows {
		r := sec.rowIndex[i]
		Grid(row[0], Row(r), Column(0), Sticky("w"), Padx("48 0"))
		Grid(row[1], Row(r), Column(1), Padx(6))
		Grid(row[2], Row(r), Column(2), Sticky("w"))
	}
}

var regionNameWidth = func() int {
	max := 0
	for _, r := range alttp.RegionOrder {
		if len(r) > max {
			max = len(r)
		}
	}
	return max + 5
}()

var photoCache = map[string]*Img{}

func iconFor(item string) *Img {
	if img, seen := photoCache[item]; seen {
		return img
	}
	var img *Img
	if data, ok := icons.PNG(item); ok {
		img = NewPhoto(Data(data))
	}
	photoCache[item] = img
	return img
}
