package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"io/ioutil"
	"os"
	"strings"
)

type CursorPos struct {
	x int
	y int
}

func cursor_movable_place(text_data []string, cursor_pos *CursorPos, dx, dy int) bool {
	if dx == 0 {

		height := len(text_data)
        next_y := cursor_pos.y + dy
            
		if next_y < height && 0 <= next_y {
            cursor_pos.y = next_y

            width := len(text_data[next_y])
            if width < cursor_pos.x {
                cursor_pos.x = width
            }

			return true 
		} else {
            return false
        }

	} else {


		width := len(text_data[cursor_pos.y])
        next_x := cursor_pos.x + dx

		if 0 <= next_x && next_x <= width{
            cursor_pos.x = next_x
			return true 
		} else {
            return false
        }

	}

	return true

}

func move_cursor(key tcell.Key, cursor_pos *CursorPos, text_data []string, keydown chan struct{}) {

	dx := 0
	dy := 0
	switch key {
	case tcell.KeyLeft:
		dx = -1
	case tcell.KeyRight:
		dx = 1
	case tcell.KeyUp:
		dy = -1
	case tcell.KeyDown:
		dy = 1
	}

	if cursor_movable_place(text_data, cursor_pos, dx, dy) {
		keydown <- struct{}{}
	}

}

func main() {

	input_data, _ := ioutil.ReadFile("sample.txt")
	text_data := strings.Split(string(input_data), "\n")
   
      


	screen, err := tcell.NewScreen()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)

	}
	if err = screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	screen.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))

	screen.Clear()

	cursor_pos := CursorPos{x: 2, y: 3}

	quit := make(chan struct{})

	keydown := make(chan struct{})

	go func() {
		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					screen.Sync()

				case tcell.KeyLeft, tcell.KeyRight, tcell.KeyUp, tcell.KeyDown:
					move_cursor(ev.Key(), &cursor_pos, text_data, keydown)

				}
			case *tcell.EventResize:
				screen.Sync()

			}
		}
	}()
loop:
	for {

		st := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)

		for row, data := range text_data {
			for col, ch := range data {
				screen.SetContent(col, row, ch, nil, st)

			}
		}

		screen.ShowCursor(cursor_pos.x, cursor_pos.y)
		screen.Sync()

		select {

		case <-quit:
			break loop
		case <-keydown:
		}

	}

	screen.Fini()

}
