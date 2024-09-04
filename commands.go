package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"regexp"
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

		re := regexp.MustCompile(`findOne\((.*)\)`)

		// Find the matching part
		matches := re.FindStringSubmatch(query)
		documents, _ := client.findOne(collName, matches[1])

		view, _ := g.View("results")
		view.Clear()
		fmt.Fprintf(view, "Query: %s\n", query)

		fmt.Fprintf(view, "Command: %s\n", command)

		fmt.Fprintf(view, "Collection: %s\n", collName)

		fmt.Fprintf(view, "Results: %v\n", string(documents))

	}
	return nil

}

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "collections" {
		_, err := g.SetCurrentView("query")
		return err
	}

	if v.Name() == "query" {
		_, err := g.SetCurrentView("collections")
		return err
	}

	return nil
}
