# Local dependency source

`scrapemate/` contains the source of `github.com/gosom/scrapemate` at
`v1.3.0`, copied from the Go module distribution. Its original MIT license
and copyright notice are preserved in `scrapemate/LICENSE`.

The root `go.mod` uses a relative `replace` directive so that regular Go
and Docker builds include our browser lifecycle corrections. No absolute
paths, prebuilt executables, browser installations, or collected leads are
needed to build the fork.

Keep this copy pinned when updating dependencies, and review local changes
before replacing it with a newer upstream release. Run its tests from the
repository root so they use the application's dependency versions:

```sh
go test github.com/gosom/scrapemate/adapters/fetchers/jshttp
```

## Local changes

- `playwrightPage.isClosed()` now reports the actual Playwright page state.
  Previously it inverted the state, recreating healthy pages and reusing
  closed ones. A regression test exercises both states through slot acquisition.
- Chromium no longer receives `--single-process`. Removing this flag resolved
  the browser termination observed during Windows collection tests.

The corrected code completed two small Jundiai searches on Windows with Go
1.27.0 and Playwright 1.61.1: 20 auto repair businesses and 20 dental clinics.
No collected data is included in this repository. This is a targeted runtime
fix, not a resolution of all upstream security or empty-result reporting issues.
