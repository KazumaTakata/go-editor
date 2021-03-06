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

	
func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {


    if len(os.Args) < 2 {
        panic("please specify your filename.")        
    }
    filename := os.Args[1]

	input_data, _ := ioutil.ReadFile(filename)
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

				case tcell.KeyEscape:
					close(quit)
					return
                case tcell.KeyEnter:
                    newline_data := text_data[cursor_pos.y][cursor_pos.x:] 
                    text_data[cursor_pos.y] = text_data[cursor_pos.y][:cursor_pos.x] 
                    text_data = append(text_data[: cursor_pos.y + 1 ], append([]string{newline_data}, text_data[cursor_pos.y + 1:]...)...)          
                    cursor_pos.x = 0
                    cursor_pos.y += 1
                	keydown <- struct{}{}

				case tcell.KeyCtrlL:
					screen.Sync()

				case tcell.KeyLeft, tcell.KeyRight, tcell.KeyUp, tcell.KeyDown:
					move_cursor(ev.Key(), &cursor_pos, text_data, keydown)
                case tcell.KeyRune:
                    text_data[cursor_pos.y] = text_data[cursor_pos.y][:cursor_pos.x] + string(ev.Rune()) +  text_data[cursor_pos.y][cursor_pos.x:]                    
                    cursor_pos.x += 1 
                	keydown <- struct{}{}
                case tcell.KeyBackspace2:
                    if cursor_pos.x > 0 {
                    text_data[cursor_pos.y] = text_data[cursor_pos.y][:cursor_pos.x - 1] +  text_data[cursor_pos.y][cursor_pos.x:]                    
                    cursor_pos.x -= 1 
                    keydown <- struct{}{}
                }
                case tcell.KeyCtrlS:
                     output_data := strings.Join(text_data, "\n") 
                     err := ioutil.WriteFile(filename, []byte(output_data), 0644)
                     check(err)


                   

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
            screen.Clear()
		}

	}

	screen.Fini()

}
