package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/ty-porter/mudora/internal/icons"
	"github.com/ty-porter/mudora/internal/rom"
	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

// iconColWidth is the pixel width of the #0 icon column, and the canvas width
// each sprite is centered onto. ttk's tree column left-aligns its image, so
// centering the sprite within a column-wide transparent canvas is what makes it
// appear centered in the column.
const iconColWidth = 40

func main() {
	App.WmTitle("Mudora — ROM Inspector")

	tv := buildUI()

	if len(os.Args) > 1 {
		load(tv, os.Args[1])
	}

	ActivateTheme("azure dark")

	// Each row is a flat, non-expandable list item, so drop the tree-column
	// disclosure indicator (the arrow). The #0 column holds only the icon, so
	// give the image the whole cell with no -side and empty -sticky, which
	// centers it instead of left-packing. Must run after ActivateTheme.
	StyleLayout("Treeview.Item",
		"Treeitem.padding", Sticky("nswe"), Children(
			"Treeitem.image", Sticky(""),
		),
	)

	App.Wait()
}

func buildUI() (tv *TTreeviewWidget) {
	fr := TFrame()
	sb := fr.TScrollbar()

	// Only the #0 tree column can hold an image, and it is always leftmost, so
	// the icon sits at the far left: icon | Location | Item.
	tv = fr.TTreeview(
		Selectmode("browse"),
		Columns("loc item"),
		Height(25),
		Yscrollcommand(func(e *Event) { e.ScrollSet(sb) }),
	)

	open := TButton(Txt("Open ROM..."), Command(func() { chooseAndLoad(tv) }))
	Pack(open, Side("top"), Anchor("w"), Padx("2m"), Pady("2m"))

	Pack(sb, Side("right"), Fill("y"))
	Pack(tv, Expand(true), Fill("both"))
	sb.Configure(Command(func(e *Event) { e.Yview(tv) }))

	tv.Column("#0", Anchor("center"), Width(iconColWidth), Stretch(false))
	tv.Column("loc", Anchor("w"), Width(320))
	tv.Column("item", Anchor("w"), Width(220))
	tv.Heading("#0", Txt(""))
	tv.Heading("loc", Txt("Location"))
	tv.Heading("item", Txt("Item"))

	Pack(fr, Expand(true), Fill("both"), Padx("2m"), Pady("2m"))
	return tv
}

func chooseAndLoad(tv *TTreeviewWidget) {
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
	load(tv, paths[0])
}

func load(tv *TTreeviewWidget, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		MessageBox(Icon("error"), Title("Mudora"),
			Msg(fmt.Sprintf("Could not read ROM:\n%s", err)))
		return
	}

	tv.Delete(tv.Children(""))
	for _, e := range rom.Inspect(data) {
		opts := []Opt{Values([]string{
			e.Location,
			e.Item,
		})}
		if img := iconFor(e.Item); img != nil {
			opts = append(opts, Image(img))
		}
		tv.Insert("", "end", opts...)
	}
	App.WmTitle("Mudora (inspecting " + path + ")")
}

var photoCache = map[string]*Img{}

func iconFor(item string) *Img {
	if img, seen := photoCache[item]; seen {
		return img
	}
	var img *Img
	if data, ok := icons.PNG(item); ok {
		img = NewPhoto(Data(center(data, iconColWidth)))
	}
	photoCache[item] = img
	return img
}

// center returns a copy of the PNG with the sprite horizontally centered on a
// transparent canvas of the given width (height unchanged). If decoding fails
// or the sprite is already at least that wide, the original bytes are returned.
func center(data []byte, width int) []byte {
	src, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return data
	}
	b := src.Bounds()
	if b.Dx() >= width {
		return data
	}
	canvas := image.NewNRGBA(image.Rect(0, 0, width, b.Dy()))
	off := (width - b.Dx()) / 2
	draw.Draw(canvas, image.Rect(off, 0, off+b.Dx(), b.Dy()), src, b.Min, draw.Src)

	var buf bytes.Buffer
	if err := png.Encode(&buf, canvas); err != nil {
		return data
	}
	return buf.Bytes()
}
