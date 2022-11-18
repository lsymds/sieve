package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/lsymds/sieve"
	"github.com/lsymds/sieve/http"
	"github.com/rivo/tview"
)

func main() {
	operationsStore := sieve.NewOperationsStore()

	go runProxy(operationsStore)

	app := tview.NewApplication()
	pages := tview.NewPages()

	// Operation page - details about the operation that occurred.
	operationPage := tview.NewFlex()
	operationPage.SetDirection(tview.FlexRowCSS)
	operationPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			pages.SwitchToPage("overview")
			return nil
		}

		return event
	})

	requestDetails := tview.NewTextView()
	requestDetails.SetBorder(true).SetTitle("Request")
	operationPage.AddItem(requestDetails, 0, 1, true)

	responseDetails := tview.NewTextView()
	responseDetails.SetBorder(true).SetTitle("Response")
	operationPage.AddItem(responseDetails, 0, 1, true)

	pages.AddPage("operation", operationPage, true, false)

	// Overview page - essentially a giant dashboard of operations that have occurred.
	overviewPage := tview.NewFlex()
	overviewPage.SetDirection(tview.FlexColumnCSS)

	header := tview.NewTextView()
	header.SetBorder(true)
	fmt.Fprintf(header, `The proxy is running on http://localhost:8080/`)
	overviewPage.AddItem(header, 3, 0, false)

	operationsList := tview.NewList()
	operationsList.SetBorder(true)
	operationsList.SetTitle("Operations")
	operationsList.SetSelectedFunc(func(i int, p string, secondaryText string, r rune) {
		oid := strings.Split(secondaryText, " - ")[0]
		operation := operationsStore.GetOperationById(oid)

		requestDetails.Clear()

		// Populate the request overview.
		fmt.Fprintf(
			requestDetails,
			`URL:
	%s

Host:
	%s

Path:
	%s`,
			operation.Request.FullUrl,
			operation.Request.Host,
			operation.Request.Path,
		)

		// Switch the page to the operation page.
		pages.SwitchToPage("operation")
	})
	overviewPage.AddItem(operationsList, 0, 1, true)

	pages.AddPage("overview", overviewPage, true, true)

	// Subscribe to operations changes.
	operationsStore.AddListener(func(o sieve.Operation) {
		app.QueueUpdateDraw(func() {
			for _, i := range operationsList.FindItems("", o.Id, true, true) {
				operationsList.RemoveItem(i)
			}

			operationsList.InsertItem(
				0,
				o.Request.FullUrl,
				fmt.Sprintf("%s - %s", o.Id, o.CreatedAt.Local().Format("15:04:05")),
				0,
				nil,
			)
		})
	})

	app.SetRoot(pages, true).Run()
}

func runProxy(operationsStore *sieve.OperationsStore) {
	server, err := http.NewHttpServer(operationsStore)
	if err != nil {
		panic(err)
	}

	panic(server.ListenAndServe(":8080"))
}
