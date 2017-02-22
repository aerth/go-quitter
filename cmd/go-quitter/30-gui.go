package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	quitter "github.com/aerth/go-quitter"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	runewidth "github.com/mattn/go-runewidth"
)

var row = 1
var style = tcell.StyleDefault

func initgui() {
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	// This is just so if we have an error, we can exit cleanly and not completely
	// mess up the terminal being worked in
	// In other words we need to shut down tcell before the program crashes
	defer func() {
		if err := recover(); err != nil {
			s.Fini()
			fmt.Println("go-quitter encountered an error:", err)
			// Print the stack trace too
			fmt.Println(err)
			os.Exit(1)
		}
	}()
	encoding.Register()

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	plain := tcell.StyleDefault
	bold := style.Bold(true)

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()
	quit := make(chan struct{})
	style = bold
	putln(s, "Press ESC to Exit")
	//putln(s, "Character set: "+s.CharacterSet())
	style = plain
	drawUserBox(s)

	// for i := 1; i < len(q.Username); i++ {
	// 	putln(s, string([]rune{
	// 		tcell.RuneLTee,
	// 		tcell.RuneHLine,
	// 		tcell.RunePlus,
	// 		tcell.RuneHLine,
	// 		tcell.RuneRTee,
	// 	}))
	//
	// }
	// putln(s, string([]rune{
	// 	tcell.RuneVLine,
	// 	tcell.RuneDiamond,
	// 	tcell.RuneVLine,
	// 	tcell.RuneUArrow,
	// 	tcell.RuneVLine,
	// })+"  (diamond, up arrow)")
	// putln(s, string([]rune{
	// 	tcell.RuneLLCorner,
	// 	tcell.RuneHLine,
	// 	tcell.RuneBTee,
	// 	tcell.RuneHLine,
	// 	tcell.RuneLRCorner,
	// }))

	s.Show()
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlD:
					drawFakeTweets(s)
					s.Sync()
				case tcell.KeyCtrlL:
					s.Sync()
				case tcell.KeyUp:

					if bufYindex < len(buf)-1 {
						bufYindex++
					}
					if col < len(buf)-1 {
						col++
					}
					redrawBuf(s)
					s.Sync()
				case tcell.KeyDown:
					if bufYindex > 0 && bufYindex != 1 {
						bufYindex--
					}
					if col > 0 {
						col--
					}
					redrawBuf(s)
					s.Sync()

				case tcell.KeyCtrlT:
					s.Clear()
					quips := []quitter.Quip{dummyQuip(), dummyQuip(), dummyQuip(), dummyQuip(), dummyQuip(), dummyQuip(), dummyQuip(), dummyQuip()}
					//var bz []string
					for _, quip := range quips {
						b1 := "@" + quip.User.Screenname
						buf = append(buf, b1)
						b2 := quip.Text
						maxwidth, _ := s.Size()
						lines := cutline(maxwidth-2, b2)
						if len(lines) == 0 {
							putln(s, "woah")
							break
						}

						for _, line := range lines {
							buf = append(buf, line)
						}

						for _, line := range buf {
							putln(s, line)
						}

						switch ev.Key() {
						case tcell.KeyUp:

							if bufYindex < len(buf)-1 {
								bufYindex++
							}
							if col < len(buf)-1 {
								col++
							}
							redrawBuf(s)
							s.Sync()
						case tcell.KeyDown:
							if bufYindex > 0 && bufYindex != 1 {
								bufYindex--
							}
							if col > 0 {
								col--
							}
							redrawBuf(s)
							s.Sync()
						}

						row++

					}
					s.Sync()

				case tcell.KeyCtrlC:

					s.Clear()
					row = 1
					style = bold
					putln(s, "Press ESC to Exit")
					//putln(s, "Character set: "+s.CharacterSet())
					style = plain
					drawUserBox(s)
					quips := []quitter.Quip{dummyQuip(), dummyQuip(), dummyQuip(), dummyQuip(), dummyQuip(), dummyQuip(), dummyQuip(), dummyQuip()}
					drawTweetBox(s, quips)
					s.Sync()
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	<-quit

	s.Fini()
}
func putln(s tcell.Screen, str string, style ...tcell.Style) {
	_, y := s.Size()
	if row > y-1 {
		row = 1
	}
	if style == nil {
		style = []tcell.Style{tcell.StyleDefault}
	}
	puts(s, style[0], 1, row, str)
	row++
}

var col int

func putstuff(s tcell.Screen, str string, style ...tcell.Style) {
	if style == nil {
		style = []tcell.Style{tcell.StyleDefault}
	}
	puts(s, style[0], col, row, str)
	col++
}
func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	i := 0
	var deferred []rune
	dwidth := 0
	//_, ry := s.Size()

	for _, r := range str {
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}
func greatest(ints ...int) int {
	var v = 0
	for _, i := range ints {
		if i > v {
			v = i
		}
	}
	return v
}
func dummyQuip() quitter.Quip {
	var q quitter.Quip
	q.User.Name = "Joe"
	q.User.Screenname = "JoeBlowtorch"
	q.Text = "Wow this works!"
	return q
}
func dummyQuipLong() quitter.Quip {
	var q quitter.Quip
	q.User.Name = "Jodfgldskjglkflkjge"
	q.User.Screenname = "JoeBlowtorchfsdf"
	q.Text = "Wow this works! fkjsldkfjlkfj slkfjd lkfjd flkdjf lkdfj dlkfjd flkjdf lkdjf ldkfjdlkf jdflkjdf lkdjf  one two three"
	return q
}
func drawFakeTweets(s tcell.Screen) {
	putln(s, "@lol lolksdokfokdokfdogkg")

}

func cutline(size int, s string) []string {
	times := len(s) / size
	if times == 0 {
		return []string{s}
	}
	list := []string{s[:size]}
	s = strings.TrimPrefix(s, list[0])
	for i := 1; i < times; i++ {
		list = append(list, s[:size])
		s = strings.TrimPrefix(s, list[i])
	}
	if s != "" {
		list = append(list, s)
	}
	return list
}

// Bust Tweet into Lines
func bustTweet(s tcell.Screen, q quitter.Quip) []string {
	maxwidth, _ := s.Size()
	lines := cutline(maxwidth-10, q.Text)
	return lines
}
func drawTweetBox(s tcell.Screen, quips []quitter.Quip) {

	maxwidth, maxheight := s.Size()
	width, height := maxwidth-4, maxheight-10
	if width == 0 || height == 0 {
		return
	}
	putln(s,
		string([]rune{tcell.RuneULCorner})+
			strings.Repeat(string([]rune{tcell.RuneHLine}), width-2)+
			string([]rune{tcell.RuneURCorner}),
	)

	if len(quips) == 0 {
		putln(s, "no quips")
	}
	//putln(s, strconv.Itoa(len(quips)))
	// for _, quip := range quips {
	//
	// 	putln(s, string([]rune{tcell.RuneHLine})+quip.Text)
	//
	// }
	for _, quip := range quips {
		bust := bustTweet(s, quip)
		//	putln(s, string([]rune{tcell.RuneHLine})+strconv.Itoa(len(bust)))
		putln(s, "@"+quip.User.Screenname)
		for _, line := range bust {
			putln(s, string([]rune{tcell.RuneVLine})+line)
		}
		putln(s, strings.Repeat(".", width-4))
	}

	putln(s,
		string([]rune{tcell.RuneLLCorner})+
			strings.Repeat(string([]rune{tcell.RuneHLine}), width-2)+
			string([]rune{tcell.RuneLRCorner}),
	)

}
func drawUserBox(s tcell.Screen) {

	width := greatest(len("username"), len(q.Username), len(q.Node))
	putln(s,
		string([]rune{tcell.RuneULCorner})+
			strings.Repeat(string([]rune{tcell.RuneHLine}), len("username"))+
			string([]rune{tcell.RuneHLine, tcell.RuneTTee})+
			strings.Repeat(string([]rune{tcell.RuneHLine}), len(q.Username))+
			strings.Repeat(string([]rune{tcell.RuneHLine}), width-len(q.Username))+
			string([]rune{tcell.RuneURCorner}),
	)
	putln(s, string([]rune{
		tcell.RuneVLine,
	})+"Username "+string([]rune{
		tcell.RuneVLine})+
		q.Username+strings.Repeat(" ", len(q.Node)-len(q.Username))+
		string([]rune{
			tcell.RuneVLine,
		}))
	putln(s, string([]rune{
		tcell.RuneVLine,
	})+"Node     "+string([]rune{
		tcell.RuneVLine})+
		q.Node+
		string([]rune{
			tcell.RuneVLine,
		}))
	putln(s,
		string([]rune{tcell.RuneLLCorner})+
			strings.Repeat(string([]rune{tcell.RuneHLine}), len("username"))+
			string([]rune{tcell.RuneHLine, tcell.RuneBTee})+
			strings.Repeat(string([]rune{tcell.RuneHLine}), len(q.Username))+
			strings.Repeat(string([]rune{tcell.RuneHLine}), width-len(q.Username))+
			string([]rune{tcell.RuneLRCorner}),
	)
}

type View struct {
	Topline       int
	leftCol       int
	widthPercent  int
	heightPercent int
	width         int
	height        int
	x, y          int
	lineNumOffset int
	Buf           *Buffer
}
type Buffer struct {
	*LineArray
	NumLines int
}
type LineArray struct {
	lines [][]byte
}

// ScrollUp scrolls the view up n lines (if possible)
func (v *View) ScrollUp(n int) {
	// Try to scroll by n but if it would overflow, scroll by 1
	if v.Topline-n >= 0 {
		v.Topline -= n
	} else if v.Topline > 0 {
		v.Topline--
	}
}

// ScrollDown scrolls the view down n lines (if possible)
func (v *View) ScrollDown(n int) {
	// Try to scroll by n but if it would overflow, scroll by 1
	if v.Topline+n <= v.Buf.NumLines-v.height {
		v.Topline += n
	} else if v.Topline < v.Buf.NumLines-v.height {
		v.Topline++
	}
}

var buf []string
var bufYindex int
var bufXindex int
var maxheight int
var maxwidth int

func redrawBuf(s tcell.Screen) {
	if bufYindex > len(buf) {
		bufYindex = len(buf) - 1
	}
	for _, line := range buf[bufYindex:] {
		putln(s, "Col: "+strconv.Itoa(col)+" Row:"+strconv.Itoa(row)+line)
	}
}
