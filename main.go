package main

import (
	"github.com/daehyeok/awsurfline/surfline"
	"github.com/deanishe/awgo"
)

var (
	wf *aw.Workflow // Our Workflow object
)

func init() {
	wf = aw.New() // Initialise workflow
}

func run() {

	var query string
	if args := wf.Args(); len(args) > 0 {
		query = args[0]
	}

	if query != "" {
		wf.Filter(query)
	}

	// We could also set this modifier via Alfred's GUI.

	items, _ := surfline.Query(query)
	for _, it := range items {
		wf.NewItem(it.Title).
			Subtitle(it.SubTitle).
			Arg(it.Url).
			UID(it.Url).
			Valid(true)
	}

	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
