package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/lsymds/sieve"
	"github.com/lsymds/sieve/http"
	"github.com/rivo/tview"
)

var (
	app             = tview.NewApplication()
	pages           = tview.NewPages()
	operationsStore = sieve.NewOperationsStore()
)

// main runs the application - booting the proxy in the background and rendering the initial TUI pages.
func main() {
	go runProxy()

	// Build the initial dashboard/overview page.
	overviewPage := buildOverviewPage()
	pages.AddPage("overview", overviewPage, true, true)

	app.SetRoot(pages, true).Run()
}

// buildOverviewPage builds the dashboard or overview page. It subscribes to proxy operations and adds them to a
// list where the consumer can then opt to view more information about it.
func buildOverviewPage() *tview.Flex {
	overviewPage := tview.NewFlex()
	overviewPage.SetDirection(tview.FlexColumnCSS)
	overviewPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.Stop()
		}

		return event
	})

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

		buildAndAddOperationPage(operation)
	})
	overviewPage.AddItem(operationsList, 0, 1, true)

	// Subscribe to operation changes and add them to the overview list.
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

	pages.AddPage("overview", overviewPage, true, true)

	return overviewPage
}

// buildAndAddOperationPage builds an operation detail page, populates it from the provided operation, and then
// switches it to be the currently rendered view in the TUI.
func buildAndAddOperationPage(operation sieve.Operation) {
	// Build the page.
	operationPage := tview.NewFlex()
	operationPage.SetDirection(tview.FlexRowCSS)
	operationPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			pages.SwitchToPage("overview")
			pages.RemovePage("operation")
			return nil
		}

		return event
	})

	requestDetails := tview.NewTextView()
	requestDetails.SetBorder(true).SetTitle("Request")
	operationPage.AddItem(requestDetails, 0, 1, true)

	responseDetails := tview.NewTextView()
	responseDetails.SetBorder(true).SetTitle("Response")
	operationPage.AddItem(responseDetails, 0, 1, false)

	pages.AddPage("operation", operationPage, true, false)

	// Populate the content from the operation.
	fmt.Fprintf(
		requestDetails,
		`URL:
	%s`,
		operation.Request.FullUrl,
	)

	// Switch to the built page.
	pages.SwitchToPage("operation")
}

// runProxy boots the proxy API and runs it on the configured port. It panics if there is an error listening.
func runProxy() {
	server, err := http.NewHttpServer(operationsStore)
	if err != nil {
		panic(err)
	}

	panic(server.ListenAndServe(":8080"))
}
