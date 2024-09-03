package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

func nextCursorLine(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	if err := v.SetCursor(cx, cy+1); err != nil && oy+cy >= len(v.BufferLines()) {
		v.SetOrigin(ox, oy+1)
	}
	return nil
}

func prevCursorLine(g *gocui.Gui, v *gocui.View) error {

	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	if err := v.SetCursor(cx, cy-1); err != nil && oy+cy > 0 {
		v.SetOrigin(ox, oy-1)
	}
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {

	var l string

	_, cy := v.Cursor()

	l, err := v.Line(cy)

	if err != nil {

		l = "Error getting line"

		return err

	}

	fmt.Println(l)

	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "collections" {
		_, err := g.SetCurrentView("query")
		return err
	}

	if v.Name() == "query" {
		_, err := g.SetCurrentView("results")
		return err
	}

	if v.Name() == "results" {
		_, err := g.SetCurrentView("collections")
		return err
	}

	return nil
}
