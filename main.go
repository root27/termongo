package main

import (
	"flag"
	"github.com/jroimartin/gocui"
	"log"
)

var (
	host   string
	client *mongoClient
	colls  []string
	dbName string
)

func main() {

	flag.StringVar(&host, "host", "", "Host to connect to ( e.g. mongodb://localhost:27017 )")

	flag.StringVar(&dbName, "db", "", "Database to connect to")

	flag.Parse()

	if host == "" {

		log.Println("Please provide a host to connect to")
		return

	}

	c, err := connectToMongoDB()

	client = c

	if err != nil {

		log.Println("Error connecting to MongoDB: ", err)
		return
	}

	collects, err := client.getCollections()

	colls = collects

	if err != nil {

		log.Println("Error getting collections: ", err)
		return
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("collections", gocui.KeyArrowDown, gocui.ModNone, nextCursorLine); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("collections", gocui.KeyArrowUp, gocui.ModNone, prevCursorLine); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("collections", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
