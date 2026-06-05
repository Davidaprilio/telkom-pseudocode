The job is failing because the build step tries to compile the package `./cmd/telupscThe job is failing because the build step is trying to compile the package path `./cmd/telupsc`, but that directory doesn’t exist in the repository at the commit used by the run (`ref: d3e5`, but that path doesn’t exist in the repository at the commit used by the run (`d3e5aa51975b9d54a43501ba350cadd9b8b90aea`).

From the logs:
- `stat .../cmd/telupsc: directory not found`

And from the workflow definition (`.github/workflows/release.yml` @ `d3e5aa51975b9d54a43501ba350cadd9b8b90aea`),aa51975b9d54a43501ba350cadd9b8b90aea`).

From the logs:
- `stat .../cmd/telupsc: directory not found`
- The workflow is explicitly running:  
  `go build ... -o "$OUT" ./cmd/telupsc` (in `.github/workflows/release.yml` at `ref d3e5aa5...`, line 58)

## Fix (choose one)

### Option A (recommended): Point the the failing line is:

```yaml
CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build ... -o "$OUT" ./cmd/telupsc
```

## Fix options (pick one)

### Option A (recommended): Build the correct main package path
Update the workflow to point at the actual Go “main” package for your CLI. Common alternatives are `.` (repo root), `./cmd/telupsc` (current expectation), or `./ workflow at the correct Go package
Update the build command to target the actual module/package that contains your `main` package.

In `.github/workflows/release.yml` (`ref d3e5aa5...`), change line 58 from:

```bash
CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build -trimpath -ldflags "-s -w -X main.Version=${VERSION}" -o "$OUT" ./cmd/telupsc
``cmd/<actual-name>`.

**Change** the build target in `.github/workflows/release.yml`:

```diff
- CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build -trimpath -ldflags "-s -w -X main.Version=${VERSION}" -o "$OUT" ./cmd/telupsc
+ CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build -trimpath -ldflags "-s -w`

to one of the following (whichever matches your repo layout):

**If your `main` package is at repo root:**
```bash
CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build -trimpath -ldflags "-s -w -X main.Version=${VERSION}" -o "$OUT" .
```

**If the command folder exists but has a different name (example `cmd/telupsc-cli`):**
```bash
CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build -trimpath -ldflags "-s -w -X main.Version=${VERSION}" -o "$OUT" ./cmd/telupsc-cli
```

**If `main.go` is somewhere else (example `./app`):**
```bash
CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build -trimpath -ldflags "-s -w -X main.Version=${VERSION}" - -X main.Version=${VERSION}" -o "$OUT" .
```

Use `.` only if your `package main` (and `func main()`) is at the repository root. If your main package is elsewhere, replace `.` with that path (for example `./cmd/telucode`).

### Option B: Add the missing `cmd/telupsc` command directory
If the intended structure is a Cobra CLI under `cmd/telupsc`, create it and ensure it containso "$OUT" ./app
```

After this change, the build will stop failing on “directory not found”.

### Option B: Add the missing directory `cmd/telupsc`
If the intended structure really is `cmd/telupsc`, then create it and place the entrypoint there:

**`cmd/telupsc/main.go`**
```go
package main

import "fmt"

var Version = "dev"

func main() {
	fmt.Println("telupsc version:", Version a `main.go` entrypoint.

Example `cmd/telupsc/main.go`:

```go
package main

import "fmt"

var Version = "dev"

func main() {
	fmt.Println("telupsc version:", Version)
	// TODO: call your cobra root command Execute() here
}
```

This makes the workflow’s existing `go build ... ./cmd/telupsc` succeed.

## Why this resolves it
The failure is not about dependencies or Go version (Go is)
	// call your real app entrypoint here
}
```

This makes the existing workflow path (`./cmd/telupsc`) valid.

## Notes
- The failure happens before compilation (it’s a filesystem path error), so fixing the package path (Option A) or creating the expected directory (Option B) resolves the job.
- The workflow file that needs editing is: `.github/workflows/release.yml` at `ref d3e5aa51975b9d54a435 set from `go.mod`, which is `go 1.22`). It’s a path mismatch: `go build` can’t compile a package directory that doesn’t exist. Pointing the workflow at the correct package (Option A) or adding the expected directory (Option B) removes the `stat ... directory not found` error and allows cross-platform builds to proceed.