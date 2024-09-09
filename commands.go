package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"strings"
)

func nextCursorLine(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()
		_, oy := v.Origin()

		// Prevent cursor from moving beyond the collection length
		if cy+oy < len(colls)-1 {
			if err := v.SetCursor(0, cy+1); err != nil {
				if oy < len(colls)-1 {
					if err := v.SetOrigin(0, oy+1); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func prevCursorLine(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, oy := v.Origin()
		_, cy := v.Cursor()

		// Prevent cursor from moving beyond the first line
		if cy+oy > 0 {
			if err := v.SetCursor(0, cy-1); err != nil {
				if oy > 0 {
					if err := v.SetOrigin(0, oy-1); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {

	_, cy := v.Cursor()

	l, err := v.Line(cy)

	collName = l

	if err != nil {

		collName = "Error getting line"

		return err

	}

	view, _ := g.View("results")

	view.Clear()

	fmt.Fprintf(view, "Selected collection: %s\n", collName)

	docs, _ := client.numberOfDocs(collName)

	fmt.Fprintf(view, "Number of documents: %d\n", docs)

	stats, _ := client.getCollSize(collName)

	fmt.Fprintf(view, "Collection size: %v Bytes\n", stats["size"])

	return nil
}

func execute(g *gocui.Gui, v *gocui.View) error {

	query, _ := v.Line(0)

	command := strings.Split(query, "(")[0]

	switch command {
	case "find":
		documents, _ := client.findAll(collName)

		view, _ := g.View("results")
		view.Clear()
		fmt.Fprintf(view, "Query: %s\n", query)

		fmt.Fprintf(view, "Command: %s\n", command)

		fmt.Fprintf(view, "Collection: %s\n", collName)

		fmt.Fprintf(view, "Results: %v\n", string(documents))

	case "findOne":

		termWidth, termHeight := g.Size()

		if v, err := g.SetView("filter", 0, 0, termWidth-1, termHeight-1); err != nil {

			if err != gocui.ErrUnknownView {
				return err
			}

			v.Title = `Filter ( e.g. {"name":"joe"} ) `
			v.Editable = true

			_, err = g.SetCurrentView("filter")

		}

	case "findOneAndUpdate":
		termwidth, termheight := g.Size()

		if v, err := g.SetView("update", 0, 0, termwidth-1, termheight-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = `Update (e.g. {"$set":{"name":"joe"}}`
			v.Editable = true
		}
		if v, err := g.SetView("filter", 0, 0, termwidth-1, termheight-1); err != nil {

			if err != gocui.ErrUnknownView {
				return err

			}

			v.Title = `Filter ( e.g. {"name":"joe"} ) `
			v.Editable = true

			_, _ = g.SetCurrentView("filter")

		}

	case "insertOne":

		termwidth, termheight := g.Size()

		if v, err := g.SetView("insertOne", 0, 0, termwidth-1, termheight-1); err != nil {
			if err != gocui.ErrUnknownView {

				return err

			}

			v.Title = "Insert Document"
			v.Editable = true
			_, _ = g.SetCurrentView("insertOne")
		}

	}
	return nil

}
func readFilter(g *gocui.Gui, v *gocui.View) error {
	filt, _ := v.Line(0)

	filter = filt

	if _, err := g.View("update"); err != nil {

		_ = g.DeleteView("filter")

		_, _ = g.SetCurrentView("query")
		documents, _ := client.findOne(collName, filter)

		view, _ := g.View("results")
		view.Clear()
		fmt.Fprintf(view, "Filter: %s\n", filter)
		fmt.Fprintf(view, "Collection: %s\n", collName)
		fmt.Fprintf(view, "Results: %v\n", string(documents))
		return nil

	}

	_ = g.DeleteView("filter")
	_, _ = g.SetCurrentView("update")

	return nil
}

func insertOne(g *gocui.Gui, v *gocui.View) error {

	document := v.Buffer()

	response, _ := client.insertOne(collName, document)

	_ = g.DeleteView("insertOne")

	_, _ = g.SetCurrentView("query")

	resultView, _ := g.View("results")

	resultView.Clear()

	fmt.Fprintf(resultView, "Result: %v\n", string(response))

	return nil

}

func readUpdate(g *gocui.Gui, v *gocui.View) error {
	update, _ := v.Line(0)
	updateRes, err := client.findOneAndUpdate(collName, filter, update)

	_ = g.DeleteView("update")

	_, _ = g.SetCurrentView("query")
	if err != nil {

		resultView, _ := g.View("results")

		resultView.Clear()

		fmt.Fprintf(resultView, "Error: %v\n", err)

		return nil
	}

	resultView, _ := g.View("results")

	resultView.Clear()

	fmt.Fprintf(resultView, "Update Result: %s\n", string(updateRes))

	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "collections" {
		_, err := g.SetCurrentView("query")
		return err
	}

	if v.Name() == "query" {

		resultView, _ := g.View("results")

		if resultView.Buffer() != "" {
			_, err := g.SetCurrentView("results")
			return err
		}

		_, err := g.SetCurrentView("collections")
		return err
	}

	if v.Name() == "results" {

		_, err := g.SetCurrentView("collections")

		return err

	}

	return nil
}

func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}
