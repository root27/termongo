package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

const (
	minWidth  = 80
	minHeight = 24
)

func layout(g *gocui.Gui) error {

	termWidth, termHeight := g.Size()

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

	// Check if terminal size is smaller than minimum requirements
	if termWidth < minWidth || termHeight < minHeight {
		if v, err := g.SetView("error", 0, 0, termWidth-1, termHeight-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Error"
			v.Wrap = true
			v.Frame = true
			v.BgColor = gocui.ColorRed
			v.FgColor = gocui.ColorWhite
			v.Clear()
			v.Write([]byte("Error: Terminal size is too small. Please resize to at least 80x24."))
		}
		return nil
	}

	// Close the error view if the terminal is resized to an acceptable size
	if err := g.DeleteView("error"); err != nil && err != gocui.ErrUnknownView {
		return err
	}

	// Header view
	if v, err := g.SetView("header", 0, 0, termWidth-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Connection"
		v.Wrap = true
		v.Frame = true
		fmt.Fprintln(v, "Connected Host: ", host)
	}

	// Collections view
	if v, err := g.SetView("collections", 0, 3, termWidth/4, termHeight-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Collections"
		v.Wrap = true
		v.Frame = true

		v.Highlight = true
		v.SelBgColor = gocui.ColorRed
		v.SelFgColor = gocui.ColorBlack

		for _, c := range colls {
			fmt.Fprintln(v, c)
		}

		if _, err := g.SetCurrentView("collections"); err != nil {

			return err
		}

	}

	// Query view
	if v, err := g.SetView("query", termWidth/4+1, 3, termWidth-1, termHeight/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Query"
		v.Wrap = true
		v.Frame = true
		v.Editable = true
	}

	// Results view
	if v, err := g.SetView("results", termWidth/4+1, termHeight/2, termWidth-1, termHeight-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Results"
		v.Wrap = true
		v.Frame = true
	}

	return nil

}
