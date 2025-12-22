

# 📦 go-shankskit

**go-shankskit** is a small Go library to help you build desktop GUI apps using [Astilectron](https://github.com/asticode/go-astilectron) with an embedded HTTP server.
It’s ideal for quickly wrapping a local web UI into a native-looking desktop window.

---

## 📥 Installation

```bash
go get github.com/shashankraocoding/go-shankskit
```

---

## 🚀 Quick Start

```go
package main

import (
    "net/http"
    "github.com/shashankraocoding/go-shankskit"
)

func main() {
    routes := map[string]http.HandlerFunc{
        "/": index,
    }

    settings := shankskit.AppSettings{
        Routes:      routes,
        AppName:     "MyApp",
        Width:       800,
        Height:      600,
        Fullscreen:  false,
        VisibleUI:   true,
        AlwaysOnTop: false,
        Transparent: false,
        Webview: true, 
    }

    // Start the embedded HTTP server and Astilectron window
    _, a, server := shankskit.StartApp(settings)

    // Wait for the window to close, then shut down the server
    shankskit.HandleShutDown(a, server)
}

func index(w http.ResponseWriter, r *http.Request) {
    shankskit.Respond(w, "index.html", nil)
}
```

Run this and you’ll see `Server running on http://localhost:<random port>` in your terminal.

---

## 🧱 API Reference

### `type AppSettings`

Configuration for your desktop app:

```go
type AppSettings struct {
    Routes      map[string]http.HandlerFunc // Route handlers
    AppName     string                      // Application name
    Width       int                         // Window width
    Height      int                         // Window height
    Fullscreen  bool                        // Launch in fullscreen
    VisibleUI   bool                        // Show window UI/frame
    Transparent bool                        // Transparent window
    AlwaysOnTop bool                        // Keep window on top
}
```

> **Note:** The HTTP server always binds to a random free port.
> The chosen port is printed to stdout.

---

### `func StartApp(settings AppSettings) (*astilectron.Window, *astilectron.Astilectron, *http.Server)`

Starts the HTTP server and creates an Astilectron window pointing to it.

**Returns:**
Window, Astilectron instance, HTTP server instance.

---

### `func HandleShutDown(a *astilectron.Astilectron, server *http.Server)`

Blocks until the window is closed, then gracefully shuts down the embedded HTTP server.

---

### `func Respond(w http.ResponseWriter, filePath string, data interface{})`

Serve a Go template file from disk.

---

### `func RespondWithEmbed(w http.ResponseWriter, filePath string, fs fs.FS, data interface{})`

Serve a Go template embedded with `embed.FS` (via `template.ParseFS`).

---

## 📂 Example HTML Template (`index.html`)

```html
<!DOCTYPE html>
<html>
<head><title>Hello</title></head>
<body>
    <h1>Hello, {{.}}</h1>
</body>
</html>
```

---

## 🛠 Dependencies

* [go-astilectron](https://github.com/asticode/go-astilectron)
* Standard Go libraries (`net/http`, `text/template`, etc.)

---

## 📦 Packaging Your App

To bundle your app into an executable with Electron:

1. Use [astilectron-bundler](https://github.com/asticode/go-astilectron-bundler)
2. Configure `bundler.json`
3. Run `go generate` + `go build`

> [!Remember!]
> To include astilectron in the build, build with the tag: `-tags astilectron`
> 
---

## ⚠️ Disclaimer

This library is provided **as-is**, without warranty of any kind.
Use of this library is at your own risk.
The author assumes **no responsibility** for damages or misuse.

---

