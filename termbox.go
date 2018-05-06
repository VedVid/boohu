// +build !tcell,!ansi,!js

package main

import (
	"errors"

	termbox "github.com/nsf/termbox-go"
)

type termui struct {
	cursor position
}

func (ui *termui) Init() error {
	return termbox.Init()
}

func (ui *termui) Close() {
	termbox.Close()
}

func (ui *termui) PostInit() {
	FixColor()
	termbox.SetOutputMode(termbox.Output256)
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.HideCursor()
	ui.HideCursor()
}

func (ui *termui) Clear() {
	termbox.Clear(termbox.Attribute(ColorFg), termbox.Attribute(ColorBg))
}

func (ui *termui) Flush() {
	termbox.Flush()
}

func (ui *termui) Interrupt() {
	termbox.Interrupt()
}

func (ui *termui) HideCursor() {
	ui.cursor = position{-1, -1}
}

func (ui *termui) SetCursor(pos position) {
	ui.cursor = pos
}

func (ui *termui) SetCell(x, y int, r rune, fg, bg uicolor) {
	termbox.SetCell(x, y, r, termbox.Attribute(fg), termbox.Attribute(bg))
}

func (ui *termui) WaitForContinue(g *game) {
loop:
	for {
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			if tev.Ch == 0 {
				switch tev.Key {
				case termbox.KeyEsc, termbox.KeySpace:
					break loop
				}
			}
			if tev.Ch == ' ' {
				break loop
			}
		case termbox.EventMouse:
			if tev.Ch == 0 {
				switch tev.Key {
				case termbox.MouseMiddle:
					break loop
				}
			}
		}
	}
}

func (ui *termui) PromptConfirmation(g *game) bool {
	for {
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			if tev.Ch == 'Y' || tev.Ch == 'y' {
				return true
			}
		}
		return false
	}
}

func (ui *termui) PressAnyKey() error {
	for {
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			return nil
		case termbox.EventInterrupt:
			return errors.New("interrupted")
		case termbox.EventMouse:
			if tev.Ch == 0 && tev.Key == termbox.MouseLeft ||
				tev.Key == termbox.MouseMiddle || tev.Key == termbox.MouseRight {
				return nil
			}
		}
	}
}

func (ui *termui) PlayerTurnEvent(g *game, ev event) (err error, again, quit bool) {
	again = true
	switch tev := termbox.PollEvent(); tev.Type {
	case termbox.EventKey:
		again = false
		if tev.Ch == 0 {
			switch tev.Key {
			case termbox.KeyArrowUp:
				tev.Ch = 'k'
			case termbox.KeyArrowRight:
				tev.Ch = 'l'
			case termbox.KeyArrowDown:
				tev.Ch = 'j'
			case termbox.KeyArrowLeft:
				tev.Ch = 'h'
			case termbox.KeyCtrlW:
				ui.EnterWizard(g)
				return nil, true, false
			case termbox.KeyCtrlQ:
				if ui.Quit(g) {
					return nil, false, true
				}
				return nil, true, false
			case termbox.KeyCtrlP:
				tev.Ch = 'm'
			}
		}
		err, again, quit = ui.HandleCharacter(g, ev, tev.Ch)
	case termbox.EventMouse:
		if tev.Ch == 0 {
			switch tev.Key {
			case termbox.MouseLeft:
				pos := position{X: tev.MouseX, Y: tev.MouseY}
				err, again = ui.GoToPos(g, ev, pos)
			case termbox.MouseRight:
				pos := position{X: tev.MouseX, Y: tev.MouseY}
				again = ui.ExaminePos(g, ev, pos)
			}
		}
	}
	if err != nil {
		again = true
	}
	return err, again, quit
}

func (ui *termui) Scroll(n int) (m int, quit bool) {
	switch tev := termbox.PollEvent(); tev.Type {
	case termbox.EventKey:
		if tev.Ch == 0 {
			switch tev.Key {
			case termbox.KeyEsc, termbox.KeySpace:
				quit = true
				return n, quit
			}
		}
		switch tev.Ch {
		case 'u':
			n -= 12
		case 'd':
			n += 12
		case 'j':
			n++
		case 'k':
			n--
		case ' ':
			quit = true
		}
	case termbox.EventMouse:
		if tev.Ch == 0 {
			switch tev.Key {
			case termbox.MouseMiddle:
				quit = true
			case termbox.MouseWheelUp:
				n -= 2
			case termbox.MouseWheelDown:
				n += 2
			}
		}
	}
	return n, quit
}

func (ui *termui) TargetModeEvent(g *game, targ Targeter, pos position, data *examineData) bool {
	switch tev := termbox.PollEvent(); tev.Type {
	case termbox.EventKey:
		if tev.Ch == 0 {
			switch tev.Key {
			case termbox.KeyArrowUp:
				tev.Ch = 'k'
			case termbox.KeyArrowRight:
				tev.Ch = 'l'
			case termbox.KeyArrowDown:
				tev.Ch = 'j'
			case termbox.KeyArrowLeft:
				tev.Ch = 'h'
			case termbox.KeyEsc, termbox.KeySpace:
				return true
			case termbox.KeyEnter:
				tev.Ch = '.'
			}
		}
		if ui.CursorCharAction(g, targ, tev.Ch, pos, data) {
			return true
		}
	case termbox.EventMouse:
		if tev.Ch == 0 {
			switch tev.Key {
			case termbox.MouseLeft:
				if ui.CursorMouseLeft(g, targ, pos) {
					return true
				}
			case termbox.MouseRight:
				data.npos = position{X: tev.MouseX, Y: tev.MouseY}
			case termbox.MouseMiddle:
				return true
			}
		}
	}
	return false
}

func (ui *termui) Select(g *game, ev event, l int) (index int, alternate bool, err error) {
	for {
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			if tev.Ch == 0 {
				switch tev.Key {
				case termbox.KeyEsc, termbox.KeySpace:
					return -1, false, errors.New(DoNothing)
				}
			}
			if 97 <= tev.Ch && int(tev.Ch) < 97+l {
				return int(tev.Ch - 97), false, nil
			}
			if tev.Ch == '?' {
				return -1, true, nil
			}
			if tev.Ch == ' ' {
				return -1, false, errors.New(DoNothing)
			}
		case termbox.EventMouse:
			if tev.Ch == 0 {
				switch tev.Key {
				case termbox.MouseLeft:
					y := tev.MouseY
					if y > 0 && y <= l {
						return y - 1, false, nil
					}
				case termbox.MouseRight:
					return -1, true, nil
				case termbox.MouseMiddle:
					return -1, false, errors.New(DoNothing)
				}
			}
		}
	}
}
