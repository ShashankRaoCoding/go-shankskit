
---

# 📦 go-shankskit

**go-shankskit** is a minimal personal use Go library that lets you build desktop GUI apps using [Astilectron](https://github.com/asticode/go-astilectron) with an embedded HTTP server. Ideal for quick tools, dashboards, and web-based UIs wrapped into native desktop apps.

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
        Port:        "8080",
        Routes:      routes,
        AppName:     "MyApp",
        Width:       800,
        Height:      600,
        Fullscreen:  false,
        VisibleUI:   true,
        AlwaysOnTop: false,
        Transparent: false,
    }

    _, a, server := shankskit.StartApp(settings)
    shankskit.HandleShutDown(a, server)
}

func index(w http.ResponseWriter, r *http.Request) {
    shankskit.Respond(w, "index.html", nil)
}
```

---

## 🧱 API Reference

### `type AppSettings`

Defines settings for your desktop app.

```go
type AppSettings struct {
    Port        string                          // HTTP server port
    Routes      map[string]http.HandlerFunc     // Route handlers
    AppName     string                          // Application name
    Width       int                             // Window width
    Height      int                             // Window height
    Fullscreen  bool                            // Launch in fullscreen
    VisibleUI   bool                            // Show window UI/frame
    Transparent bool                            // Transparent window
    AlwaysOnTop bool                            // Keep window on top
}
```

---

### `func StartApp(settings AppSettings) (*astilectron.Window, *astilectron.Astilectron, *http.Server)`

Starts the HTTP server and creates an Astilectron window pointing to it.

**Arguments:**

* `settings` – App configuration.

**Returns:**

* Window, Astilectron instance, HTTP server instance.

---

### `func HandleShutDown(a *astilectron.Astilectron, server *http.Server)`

Blocks until the window is closed and then gracefully shuts down the HTTP server.

---

### `func Respond(w http.ResponseWriter, filePath string, data interface{})`

Utility function for serving Go templates.

**Arguments:**

* `w` – `http.ResponseWriter`
* `filePath` – path to `.html` file
* `data` – any data passed to the template

---

### `func ListSubPaths(dir string) []string`

Recursively lists all files in a directory (including nested).

**Example:**

```go
paths := shankskit.ListSubPaths("./assets")
```

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
* Standard Go libraries

---

## 📦 Packaging Your App

To bundle your app into an executable with Electron:

1. Use [astilectron-bundler](https://github.com/asticode/go-astilectron-bundler)
2. Configure `bundler.json`
3. Run `go generate` + `go build`

Let me know if you want a bundler config template!

---

## ⚠️ Disclaimer

This library is provided **as-is**, without warranty of any kind, express or implied. Use of this library is at your own risk. In no event shall the author or contributors be held liable for any damages or losses, including but not limited to direct, indirect, incidental, or consequential damages, loss of data, or business interruption arising out of the use of or inability to use this software, even if advised of the possibility of such damages.

By using this package, you agree that the author assumes **no responsibility for negligence, misuse, or misinterpretation** of the software.

