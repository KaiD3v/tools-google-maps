package jshttp

import (
	"context"
	"testing"

	"github.com/mxschmitt/playwright-go"
)

// Exercise the actual Playwright wrapper through slot acquisition. The older
// fakePage tests bypassed this wrapper and missed its inverted closed state.
func TestSessionSlotWithPlaywrightPageState(t *testing.T) {
	t.Parallel()

	for _, closed := range []bool{false, true} {
		name := "reuse open page"
		if closed {
			name = "replace closed page"
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			original := &playwrightPage{p: &pageStateStub{closed: closed}}
			runtime := &pageStateRuntime{current: original}
			slot := &sessionSlot{runtime: runtime}

			got, err := slot.acquirePage(context.Background())
			if err != nil {
				t.Fatalf("acquirePage: %v", err)
			}
			if got == nil || got.isClosed() {
				t.Fatal("acquirePage must return a usable, open page")
			}
			if closed && (got == original || runtime.recreations != 1) {
				t.Fatal("a closed page must be replaced exactly once")
			}
			if !closed && (got != original || runtime.recreations != 0) {
				t.Fatal("an open page must be reused without recreation")
			}
		})
	}
}

type pageStateStub struct {
	playwright.Page
	closed bool
}

func (p *pageStateStub) IsClosed() bool { return p.closed }

type pageStateRuntime struct {
	fakeRuntime
	current     *playwrightPage
	recreations int
}

func (r *pageStateRuntime) primaryPage() (page, error) { return r.current, nil }

func (r *pageStateRuntime) recreatePage() error {
	r.recreations++
	r.current = &playwrightPage{p: &pageStateStub{}}
	return nil
}
