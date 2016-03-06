// The MIT License (MIT)
//
// Copyright (c) 2016 aerth
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

package main

import (
	"fmt"
	"log"
	"os"
  "time"
	"os/signal"
	"syscall"
  "github.com/jroimartin/gocui"
  qw "github.com/aerth/go-quitter"
)

var (
	sigs = make(chan os.Signal, 1)
)

type Window struct {
	Layout     string
	Controller string
}

func NewWindow() *Window {
	return &Window{Layout: "as", Controller: "es"}
}
/*
func (w *Window) Destroy() {

}
*/
type Terminal struct {
	Gui     *gocui.Gui
	Windows []*Window
}

func NewTerminal() *Terminal {
	return &Terminal{Gui: gocui.NewGui()}
}

func (term *Terminal) Init() {
	if err := term.Gui.Init(); err != nil {
		log.Panicln(err)
	}

	term.Gui.SetLayout(layout)
  term.Gui.SelBgColor = gocui.ColorGreen
	term.Gui.SelFgColor = gocui.ColorWhite
	term.Gui.Cursor = true

	if err := keybindings(term.Gui); err != nil {
		log.Panicln(err)
	}

	term.Gui.SelBgColor = gocui.ColorGreen
term.Gui.SelFgColor = gocui.ColorWhite
	term.Gui.Cursor = true
}

func (term *Terminal) Loop() {
	var err error = term.Gui.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (term *Terminal) AddWindow(w *Window) {
	// Add window to list of windows.
	term.Windows[0] = w
	// term.Gui.SetLayout(layout)
}

func (term *Terminal) Destroy() {
	term.Gui.Close()
}


func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyPgup, gocui.ModNone, pageUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyPgdn, gocui.ModNone, pageDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, refresh); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, reload); err != nil {
		return err
	}




	return nil
}

func main() {
	// Capture and manage signals.
	signal.Notify(sigs, os.Interrupt)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigs
		//cleanup()
		fmt.Println("Signal interruption.")
		fmt.Println(sig)
		os.Exit(1)
	}()
	term := NewTerminal()
  if err := keybindings(term.Gui); err != nil {
    log.Panicln(err)
  }
	defer term.Destroy()
	term.Init()
	term.Loop()
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v1, err := g.SetView("top", -1, -1, maxX, 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v1, "go-quitter v0.0.8")
	}

	if v3, err := g.SetView("shell", -1, maxY-3, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v3.Editable = false


	}
	if sb, err := g.SetView("statusbar", -1, maxY-3, maxX, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		sb.Editable = false
    fmt.Fprintln(sb, "Reading gs.sdf.org")
	}


  	if v2, err := g.SetView("main", 1, 1, maxX-2, maxY-3); err != nil {
  		if err != gocui.ErrUnknownView {
  			return err
  		}

  		v2.Frame = false
  		v2.Wrap = true

      //HomePage(g, v2)

  	}
	return nil
}


func cursorDown(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        cx, cy := v.Cursor()
        if err := v.SetCursor(cx, cy+1); err != nil {
            ox, oy := v.Origin()
            if err := v.SetOrigin(ox, oy+1); err != nil {
                return err
            }
        }
    }
    return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        ox, oy := v.Origin()
        cx, cy := v.Cursor()
        if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
            if err := v.SetOrigin(ox, oy-1); err != nil {
                    return err
            }
        }
    }
    return nil
}

func pageDown(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        cx, cy := v.Cursor()
        if err := v.SetCursor(cx, cy+20); err != nil {
            ox, oy := v.Origin()
            if err := v.SetOrigin(ox, oy+20); err != nil {
                return err
            }
        }
    }
    return nil
}

func pageUp(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        ox, oy := v.Origin()
        cx, cy := v.Cursor()
        if err := v.SetCursor(cx, cy-20); err != nil && oy > 0 {
            if err := v.SetOrigin(ox, oy-20); err != nil {
                    return err
            }
        }
    }
    return nil
}

func refreshPage(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        ox, oy := v.Origin()
        cx, cy := v.Cursor()
        if err := v.SetCursor(cx, cy-cy); err != nil && oy > 0 {
            if err := v.SetOrigin(ox, oy-oy); err != nil {
                    return err
            }
        }
    }
    return nil
}

func HomePage(g *gocui.Gui, v *gocui.View) error {
//  v.Clear()
  fmt.Fprintln(v, "Refreshing...")
  //v.Clear()
  if v != nil {

    g.Execute(func(g *gocui.Gui) error {

	//v.Clear()
  q2 := qw.NewAuth()
  q2.Node = "gs.sdf.org"
  quips, err := q2.GetPublic(true)
  if err != nil {
    fmt.Fprintln(v, err)
  }
  for i := range quips {
    if quips[i].User.Screenname == quips[i].User.Name {
      fmt.Fprintln(v, "[@"+quips[i].User.Screenname+"] "+quips[i].Text+"\n")
    } else {
      fmt.Fprintln(v, "@"+quips[i].User.Screenname+" ["+quips[i].User.Name+"] "+quips[i].Text+"\n")
    }
  }
	return nil
})


  }
  return nil
}


func refreshHome(g *gocui.Gui, v *gocui.View) error {
  v.Clear()
  fmt.Fprintln(v, "Refreshing...")
  v.Clear()
  if v != nil {

    g.Execute(func(g *gocui.Gui) error {

	//v.Clear()
  defer func() { q2 := qw.NewAuth()
  q2.Node = "gs.sdf.org"
  quips, err := q2.GetPublic(false)
  if err != nil {
    fmt.Fprintln(v, err)
  }
  for i := range quips {
    if quips[i].User.Screenname == quips[i].User.Name {
      fmt.Fprintln(v, "[@"+quips[i].User.Screenname+"] "+quips[i].Text+"\n")
    } else {
      fmt.Fprintln(v, "@"+quips[i].User.Screenname+" ["+quips[i].User.Name+"] "+quips[i].Text+"\n")
    }
  }
  }()

	return nil
})


  }
  return nil
}


func refresh(g *gocui.Gui, v *gocui.View) error {
  go ReadPub(g, v)
  go g.SetCurrentView("main")
	return nil
}

func ReadPub(g *gocui.Gui, v *gocui.View) error {
  g.SetCurrentView("main")
  q2 := qw.NewAuth()
  q2.Node = "gs.sdf.org"
  t2 := time.Now()
  quips, err := q2.GetPublic(true)
  t3 := time.Now()
  if err != nil {
    return err
  }
  fmt.Fprintln(v, "Took ", t3.Sub(t2))
  for i := range quips {
    if quips[i].User.Screenname == quips[i].User.Name {
      fmt.Fprintln(v, "[@"+quips[i].User.Screenname+"] "+quips[i].Text+"\n")
    } else {
      fmt.Fprintln(v, "@"+quips[i].User.Screenname+" ["+quips[i].User.Name+"] "+quips[i].Text+"\n")
    }
  }
  g.SetCurrentView("main")
	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "main" {
		return g.SetCurrentView("main")
	}
	return g.SetCurrentView("main")
}



func reload(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "shell" {
		return g.SetCurrentView("main")
	}
	return g.SetCurrentView("shell")
}

func returnStatus() string {

return "go-quitter stand-by"

}
