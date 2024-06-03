package SiteNavigator

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
)

// Context is the context of the session.
type Context struct {
	// ctx is the context of the session.
	ctx context.Context

	// ctxCancel is the cancel function of the context.
	ctxCancel context.CancelFunc

	// allocCancel is the cancel function of the allocator.
	allocCancel context.CancelFunc
}

// Close closes the context.
func (c *Context) Close() {
	if c.ctxCancel != nil {
		c.ctxCancel()
	}

	if c.allocCancel != nil {
		c.allocCancel()
	}
}

// Context returns the context of the session.
//
// Returns:
//   - context.Context: The context of the session.
func (c *Context) Context() context.Context {
	return c.ctx
}

// InitializeContext initializes a new context.
//
// Returns:
//   - *Context: The new context.
func InitializeContext() *Context {
	// Create a new browser context
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false), // run in headful mode
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, ctxCancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	return &Context{
		ctx:         ctx,
		ctxCancel:   ctxCancel,
		allocCancel: allocCancel,
	}
}

// NewSubContext creates a new sub context.
//
// Returns:
//   - *Context: The new sub context.
func (c *Context) NewSubContext() *Context {
	var newContext Context

	newContext.ctx, newContext.ctxCancel = chromedp.NewContext(c.ctx)
	newContext.allocCancel = nil

	return &newContext
}

// ParseHTML parses the HTML of the URL.
//
// Parameters:
//   - url: The URL of the HTML.
//   - loadedSignal: The signal that the page has loaded.
//
// Returns:
//   - *html.Node: The HTML node of the URL.
//   - error: The error that occurred while parsing the HTML.
func (c *Context) ParseHTML(url string, loadedSignal ...chromedp.Action) (*html.Node, error) {
	var document string

	tasks := slices.Insert(chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.OuterHTML("html", &document),
	}, 1, loadedSignal...)

	err := chromedp.Run(c.ctx, tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to open the URL %s: %w", url, err)
	}

	return html.Parse(strings.NewReader(document))
}

// GetLastPage gets the last page of the URL.
//
// Parameters:
//   - url: The URL of the page.
//   - waitTask: The task to wait for the page to load.
//   - f: The function to extract the last page from the HTML.
//
// Returns:
//   - int: The last page of the URL.
//   - error: The error that occurred while getting the last page.
func (c *Context) GetLastPage(url string, wait WaitFunc, f ExtractFunc[int]) (int, error) {
	err := chromedp.Run(c.ctx, wait(url))
	if err != nil {
		return 0, fmt.Errorf("failed to open the URL %s: %w", url, err)
	}

	var document string

	err = chromedp.Run(c.ctx, chromedp.OuterHTML("html", &document))
	if err != nil {
		return 0, fmt.Errorf("failed to get the outer HTML: %w", err)
	}

	doc, err := html.Parse(strings.NewReader(document))
	if err != nil {
		return 0, fmt.Errorf("failed to parse html: %w", err)
	}

	last_page, err := f(doc)
	if err != nil {
		return 0, fmt.Errorf("failed to extract the number of pages: %w", err)
	}

	return last_page, nil
}

// GetArticleNodes gets the <article> elements on the page.
//
// Returns:
//   - []*cdp.Node: The <article> elements on the page.
//   - error: The error that occurred while getting the <article> elements.

// GetNodes gets the nodes that match the selector.
//
// Parameters:
//   - sel: The selector of the nodes.
//   - opt: The options of the selector.
//
// Returns:
//   - []*cdp.Node: The nodes that match the selector.
//   - error: The error that occurred while getting the nodes.
func (c *Context) GetNodes(sel any, opt func(*chromedp.Selector)) ([]*cdp.Node, error) {
	var nodes []*cdp.Node

	err := chromedp.Run(c.ctx, chromedp.Nodes(sel, &nodes, opt))
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// RunTasks runs the tasks on the session.
//
// Parameters:
//   - tasks: The tasks to run.
//
// Returns:
//   - error: The error that occurred while running the tasks.
func (c *Context) RunTasks(tasks chromedp.Tasks) error {
	return chromedp.Run(c.ctx, tasks)
}
