/**
	Taken and repurposed from https://github.com/nsf/termbox-go/tree/master/_demos/editbox.go
	author: termbox-go authors & nan0 authors
 */
package nan0chat

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"unicode/utf8"
	"time"
	"math"
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func tbprintbounded(x, y, maxX int, fg, bg termbox.Attribute, msg string) (lastx, lasty int) {
	for i, c := range msg {
		if i == maxX {
			y++
			x = 0
		}
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
	return x, y
}

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func rune_advance_len(r rune, pos int) int {
	if r == '\t' {
		return tabstop_length - pos%tabstop_length
	}
	return runewidth.RuneWidth(r)
}

func voffset_coffset(text []byte, boffset int) (voffset, coffset int) {
	text = text[:boffset]
	for len(text) > 0 {
		r, size := utf8.DecodeRune(text)
		text = text[size:]
		coffset += 1
		voffset += rune_advance_len(r, voffset)
	}
	return
}

func byte_slice_grow(s []byte, desired_cap int) []byte {
	if cap(s) < desired_cap {
		ns := make([]byte, len(s), desired_cap)
		copy(ns, s)
		return ns
	}
	return s
}

func byte_slice_remove(text []byte, from, to int) []byte {
	size := to - from
	copy(text[from:], text[to:])
	text = text[:len(text)-size]
	return text
}

func byte_slice_insert(text []byte, offset int, what []byte) []byte {
	n := len(text) + len(what)
	text = byte_slice_grow(text, n)
	text = text[:n]
	copy(text[offset+len(what):], text[offset:])
	copy(text[offset:], what)
	return text
}

const preferred_horizontal_threshold = 5
const tabstop_length = 8

type ChatClientUI struct {
	editBox       EditBox
	outputBox     OutputBox
	editBoxWidth  int
	editBoxPrefix string
}

type EditBox struct {
	text           []byte
	line_voffset   int
	cursor_boffset int // cursor offset in bytes
	cursor_voffset int // visual cursor offset in termbox cells
	cursor_coffset int // cursor offset in unicode code points
}

type OutputBox struct {
	messages          []string
	width             int
	height            int
	windowTopIndex    int
	windowBottomIndex int
}

// Draws the EditBox in the given location, 'h' is not used at the moment
func (eb *EditBox) Draw(x, y, w, h int) {
	eb.AdjustVOffset(w)

	const coldef = termbox.ColorDefault
	fill(x, y, w, h, termbox.Cell{Ch: ' '})

	t := eb.text
	lx := 0
	tabstop := 0
	for {
		rx := lx - eb.line_voffset
		if len(t) == 0 {
			break
		}

		if lx == tabstop {
			tabstop += tabstop_length
		}

		if rx >= w {
			termbox.SetCell(x+w-1, y, '→',
				coldef, coldef)
			break
		}

		r, size := utf8.DecodeRune(t)
		if r == '\t' {
			for ; lx < tabstop; lx++ {
				rx = lx - eb.line_voffset
				if rx >= w {
					goto next
				}

				if rx >= 0 {
					termbox.SetCell(x+rx, y, ' ', coldef, coldef)
				}
			}
		} else {
			if rx >= 0 {
				termbox.SetCell(x+rx, y, r, coldef, coldef)
			}
			lx += runewidth.RuneWidth(r)
		}
	next:
		t = t[size:]
	}

	if eb.line_voffset != 0 {
		termbox.SetCell(x, y, '←', coldef, coldef)
	}
}

// Adjusts line visual offset to a proper value depending on width
func (eb *EditBox) AdjustVOffset(width int) {
	ht := preferred_horizontal_threshold
	max_h_threshold := (width - 1) / 2
	if ht > max_h_threshold {
		ht = max_h_threshold
	}

	threshold := width - 1
	if eb.line_voffset != 0 {
		threshold = width - ht
	}
	if eb.cursor_voffset-eb.line_voffset >= threshold {
		eb.line_voffset = eb.cursor_voffset + (ht - width + 1)
	}

	if eb.line_voffset != 0 && eb.cursor_voffset-eb.line_voffset < ht {
		eb.line_voffset = eb.cursor_voffset - ht
		if eb.line_voffset < 0 {
			eb.line_voffset = 0
		}
	}
}

func (eb *EditBox) MoveCursorTo(boffset int) {
	eb.cursor_boffset = boffset
	eb.cursor_voffset, eb.cursor_coffset = voffset_coffset(eb.text, boffset)
}

func (eb *EditBox) RuneUnderCursor() (rune, int) {
	return utf8.DecodeRune(eb.text[eb.cursor_boffset:])
}

func (eb *EditBox) RuneBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(eb.text[:eb.cursor_boffset])
}

func (eb *EditBox) MoveCursorOneRuneBackward() {
	if eb.cursor_boffset == 0 {
		return
	}
	_, size := eb.RuneBeforeCursor()
	eb.MoveCursorTo(eb.cursor_boffset - size)
}

func (eb *EditBox) MoveCursorOneRuneForward() {
	if eb.cursor_boffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.MoveCursorTo(eb.cursor_boffset + size)
}

func (eb *EditBox) MoveCursorToBeginningOfTheLine() {
	eb.MoveCursorTo(0)
}

func (eb *EditBox) MoveCursorToEndOfTheLine() {
	eb.MoveCursorTo(len(eb.text))
}

func (eb *EditBox) DeleteRuneBackward() {
	if eb.cursor_boffset == 0 {
		return
	}

	eb.MoveCursorOneRuneBackward()
	_, size := eb.RuneUnderCursor()
	eb.text = byte_slice_remove(eb.text, eb.cursor_boffset, eb.cursor_boffset+size)
}

func (eb *EditBox) DeleteRuneForward() {
	if eb.cursor_boffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.text = byte_slice_remove(eb.text, eb.cursor_boffset, eb.cursor_boffset+size)
}

func (eb *EditBox) DeleteTheRestOfTheLine() {
	eb.text = eb.text[:eb.cursor_boffset]
}

func (eb *EditBox) InsertRune(r rune) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	eb.text = byte_slice_insert(eb.text, eb.cursor_boffset, buf[:n])
	eb.MoveCursorOneRuneForward()
}

// Please, keep in mind that cursor depends on the value of line_voffset, which
// is being set on Draw() call, so.. call this method after Draw() one.
func (eb *EditBox) CursorX() int {
	return eb.cursor_voffset - eb.line_voffset
}

func (eb *EditBox) Clear() {
	eb.cursor_boffset = 0
	eb.cursor_coffset = 0
	eb.cursor_voffset = 0
	eb.line_voffset = 0
	eb.text = nil
}

func (chatUi *ChatClientUI) Start(prefix string, messageChannel chan<- string) {
	// configure output box
	chatUi.outputBox.width = 90
	chatUi.outputBox.height = 20
	chatUi.outputBox.windowTopIndex = 0
	chatUi.outputBox.windowBottomIndex = chatUi.outputBox.height - 1

	// configure edit box
	chatUi.editBoxWidth = chatUi.outputBox.width
	chatUi.editBoxPrefix = prefix

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	chatUi.redraw_all()

	// constantly redraw output box so that we can update output box when new
	// strings are written to the output box message list
	go func() {
		for {
			time.Sleep(30 * time.Millisecond)
			chatUi.redraw_all()
		}
	}()
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyArrowUp:
				chatUi.outputBox.windowUp()
			case termbox.KeyArrowDown:
				chatUi.outputBox.windowDown()
			case termbox.KeyCtrlC:
				chatUi.outputBox.clearMessages()
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				chatUi.editBox.MoveCursorOneRuneBackward()
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				chatUi.editBox.MoveCursorOneRuneForward()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				chatUi.editBox.DeleteRuneBackward()
			case termbox.KeyDelete, termbox.KeyCtrlD:
				chatUi.editBox.DeleteRuneForward()
			case termbox.KeyTab:
				chatUi.editBox.InsertRune('\t')
			case termbox.KeySpace:
				chatUi.editBox.InsertRune(' ')
			case termbox.KeyCtrlK:
				chatUi.editBox.DeleteTheRestOfTheLine()
			case termbox.KeyHome, termbox.KeyCtrlA:
				chatUi.editBox.MoveCursorToBeginningOfTheLine()
			case termbox.KeyEnd, termbox.KeyCtrlE:
				chatUi.editBox.MoveCursorToEndOfTheLine()
			case termbox.KeyEnter:
				if len(chatUi.editBox.text) > 0 {
					// add the prefix to the message before sending
					fullMsg := chatUi.editBoxPrefix + string(chatUi.editBox.text)
					messageChannel <- fullMsg
					chatUi.outputBox.addMessage(fullMsg)
					chatUi.editBox.Clear()
				}
			default:
				if ev.Ch != 0 {
					chatUi.editBox.InsertRune(ev.Ch)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func (chatUi *ChatClientUI) redraw_all() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)

	// setup initial locations for all ui elements
	// initial output box element location (top-left)
	outputx := 1
	outputy := 1
	// initial edit box element location (top-left)
	midy := chatUi.outputBox.height + 3
	midx := outputy

	// unicode box drawing chars around the edit box
	termbox.SetCell(midx-1, midy, '│', coldef, coldef)
	termbox.SetCell(midx+chatUi.editBoxWidth, midy, '│', coldef, coldef)
	termbox.SetCell(midx-1, midy-1, '┌', coldef, coldef)
	termbox.SetCell(midx-1, midy+1, '└', coldef, coldef)
	termbox.SetCell(midx+chatUi.editBoxWidth, midy-1, '┐', coldef, coldef)
	termbox.SetCell(midx+chatUi.editBoxWidth, midy+1, '┘', coldef, coldef)
	fill(midx, midy-1, chatUi.editBoxWidth, 1, termbox.Cell{Ch: '─'})
	fill(midx, midy+1, chatUi.editBoxWidth, 1, termbox.Cell{Ch: '─'})

	// draw unicode output box
	fill(outputx-1, outputy, 1, chatUi.outputBox.height, termbox.Cell{Ch: '|'})
	fill(outputx+chatUi.outputBox.width, outputy, 1, chatUi.outputBox.height, termbox.Cell{Ch: '|'})
	termbox.SetCell(outputx-1, outputy-1, '┌', coldef, coldef)
	termbox.SetCell(outputx-1, outputy+chatUi.outputBox.height, '└', coldef, coldef)
	termbox.SetCell(outputx+chatUi.outputBox.width, outputy-1, '┐', coldef, coldef)
	termbox.SetCell(outputx+chatUi.outputBox.width, outputy+chatUi.outputBox.height, '┘', coldef, coldef)
	fill(outputx, outputy-1, chatUi.outputBox.width, 1, termbox.Cell{Ch: '─'})
	fill(outputx, outputy+chatUi.outputBox.height, chatUi.outputBox.width, 1, termbox.Cell{Ch: '─'})

	// finishing touches on edit box
	chatUi.editBox.Draw(midx, midy, chatUi.editBoxWidth, 1)
	termbox.SetCursor(midx+chatUi.editBox.CursorX(), midy)

	// write instructions
	tbprint(midx+6, midy+3, coldef, coldef, "Press ESC to quit")

	// write all messages to the ouptut box
	lasty := outputy
	for i, msg := range chatUi.outputBox.messages {
		if i >= chatUi.outputBox.windowTopIndex && i < chatUi.outputBox.windowBottomIndex {
			_, lasty = tbprintbounded(outputx+1, lasty+1, chatUi.outputBox.width-1, coldef, coldef, msg)
		}
	}

	termbox.Flush()
}

// Adds a new message to the output box, shifting the draw window if there are too many messages to be safely added
func (outputBox *OutputBox) addMessage(message string) {
	messageLength := float64(len(message))
	modifiedWidth := float64(outputBox.width - 1)
	// if the message length is longer than the width of the outputbox, we must wrap by breaking up the message in
	// sizes equal to the width
	if messageLength > modifiedWidth {
		for i := 0.0; i < math.Ceil(messageLength/modifiedWidth); i++ {
			lowIdex := int(i * modifiedWidth)
			// the min in this expression is to bound the idex to the actual message size
			highIdex := int(math.Min((i*modifiedWidth)+modifiedWidth, messageLength))
			outputBox.messages = append(outputBox.messages, message[lowIdex:highIdex])
		}
	} else {
		outputBox.messages = append(outputBox.messages, message)
	}
	// adjust window height til message fits
	for len(outputBox.messages)-1 >= outputBox.windowBottomIndex {
		outputBox.windowDown()
	}
}

// Shifts the drawing window for messages up, so that earlier messages can be viewed
func (outputBox *OutputBox) windowUp() {
	if outputBox.windowTopIndex > 0 {
		outputBox.windowTopIndex--
		outputBox.windowBottomIndex--
	}
}

// Shifts the drawing window for messages down, so that the latest messages can be viewed
func (outputBox *OutputBox) windowDown() {
	if outputBox.windowBottomIndex < len(outputBox.messages) {
		outputBox.windowTopIndex++
		outputBox.windowBottomIndex++
	}
}

// Clears all messages from the output box and resets the window
func (outputBox *OutputBox) clearMessages() {
	outputBox.messages = make([]string, outputBox.height)
	outputBox.windowBottomIndex = outputBox.height
	outputBox.windowTopIndex = 0
}
