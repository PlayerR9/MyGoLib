package Document

import (
	"errors"
	"fmt"

	fss "github.com/PlayerR9/MyGoLib/FString/Section"
)

// Page is a type that represents a page of a document.
type Page struct {
	// sections are the sections of the page.
	sections []fss.Sectioner

	// width is the width of the page.
	width int

	// height is the height of the page.
	height int
}

func (p *Page) View() ([]*RenderedPage, error) {
	renders, err := p.sections[0].ApplyRender(p.width, p.height)
	if err != nil {
		return nil, fmt.Errorf("could not apply render: %w", err)
	}

	if len(renders) == 0 {
		// No renders were created.
		return nil, errors.New("no renders were created")
	}

	// FIXME: This is a temporary fix.

	sol := make([]*RenderedPage, 0, len(renders))

	for i, render := range renders {
		sol = append(sol, NewRenderedPage(render, 1, i+1))
	}

	return sol, nil
}
