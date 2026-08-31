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
