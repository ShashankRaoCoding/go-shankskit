package shankskit 

import ( 
	"github.com/asticode/go-astilectron" 
	"log" 
	"net/http" 
	"fmt"; 
	"os"; 
	"bufio"; 
	"strings"; 
	"reflect"; 	
	"context"; 
	"text/template"
) 

// Creates a new List object with some default values 
func NewList(contents ...interface{}) *List { 
	return &List{ 
		Contents: contents, 
		RangeIndex: -1 , 
		RangeItem: nil , 
	}
}

type List struct { 
	Contents []interface{} 
	RangeIndex int 
	RangeItem interface{} 
}

	// adds the object to the List's Contents 
	func (l *List) Append(object interface{}) { 
		l.Contents = append(l.Contents, object) 
	}

	// returns the length of the List as an integer 
	func (l *List) GetLength() int { 
		return len(l.Contents) 
	}

	// removes the first instance of the object and does nothing if not found 
	func (l *List) Remove(object interface{}) { 
		objectIndex := l.Find(object) 
		if objectIndex == -1 { 
			// do nothign 
		} else { 
			l.Contents = append(l.Contents[:objectIndex], l.Contents[objectIndex+1:]...) 
			if l.RangeIndex != -1 { 
				l.RangeIndex -= 1 
			}
		}
	}
	
	// returns the index of the first instance of the object 
	// returns -1 if not found 
	func (l *List) Find(object interface{}) int { 
		for index, item := range l.Contents { 
			if reflect.DeepEqual(item, object) { 
				return index 
			}
		}
		return -1 
	}
	
	// Range 
	func (l *List) Range() bool { 
		l.RangeIndex += 1 
		if l.RangeIndex == len(l.Contents) { 
			l.RangeIndex = -1 
			l.RangeItem = nil 
			return false 
		} else { 
			l.RangeItem = l.Contents[l.RangeIndex]
			return true 
		}
	}

	func (l *List) GetSyntax() string { 
		representation := `[` 
		for l.Range() { 
			item := l.RangeItem 
			switch typedItem := item.(type) { 
				case string : 
					representation += fmt.Sprintf(`'%+v'`, typedItem) 
				case *List : 
					representation += typedItem.GetSyntax() 
				case []interface{}: 
					tempItem := NewList(typedItem...) 
					representation += tempItem.GetSyntax() 
				default : 
					representation += fmt.Sprintf(`%+v`, typedItem) 
			}
			
			representation += ", " 
		}
		if 1 < len(representation) { 
			representation = representation[0:-2+len(representation)] 
		}
		representation += `]` 
		return representation 
	}


	type Grid struct { 
		Rows int 
		Columns int 
		Dimensions []int 
		Table *List 
		RangeIndex int 
		RangeRow *List 
	}
		func (g *Grid) Range() bool { 
			g.RangeIndex += 1 
			if g.RangeIndex == g.Rows { 
				g.RangeIndex = -1 
				g.RangeRow = NewList() 
				return false 
			} else { 
				row := (*g.Table).Contents[g.RangeIndex] 
				typedRow, ok := row.(List) 
				if ok { 
					g.RangeRow = &typedRow 
				}
				return true 
			}
			return false 
		}
	
		func (g *Grid) ValueAt(column_index int, row_index int) interface{} { 
			if g.Rows <= row_index { 
				// do nothing 
			} else if g.Columns <= column_index { 
				// do nothing 
			} else { 
				row := (*g.Table).Contents[row_index] 
				typedRow, ok := row.(List) 
				if ok {
					value := typedRow.Contents[column_index] 
					typedValue, ok := value.(interface{}) 
					if ok {
						return typedValue 
					}	
				}
			}
			return nil 
		}
	func (g *Grid) GetSyntax() string { 
		output := "" 
		for g.Range() { 
			row := g.RangeRow 
			output += row.GetSyntax() + "\n" 
		}
		
		return output 

	}

func NewGrid(rows int, columns int) *Grid { 
	table := NewList() 
	for s := 0; s < rows; s++ { 
		row := NewList() 
		for i := 0; i < columns; i++ { 
			row.Append(nil) 
		}
		table.Append(*row)
	}
	grid := &Grid{ 
		Rows: rows, 
		Columns: columns, 
		Dimensions: []int{ 
			rows, 
			columns, 
		}, 
		Table: table, 
		RangeIndex: -1, 
		RangeRow: NewList(), 
	}
	return grid 
}

func GetValue(item interface{}) string { 
	return fmt.Sprintf(`%+v`, item) 
}

	
func Input(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return input
}

func HandleErrors(err error) { 
	if err != nil { 
		fmt.Println(err) 
	} else { 
		// do nothing 
	}
}

func Print(plaintext interface{}) { 
	fmt.Println(plaintext)
}

func SplitString(plaintext string, seperator string) []string { 
	return strings.Split(plaintext, seperator)
}

func readLines(path string) []string { 
	var lines []string 
	file, err := os.Open(path) 
	HandleErrors(err)
	scanner := bufio.NewScanner(file) 
	for scanner.Scan() { 
		lines = append(lines, scanner.Text() ) 
	} 
	return lines 
}

func ReadFile(path string) string { 
	file, err := os.Open(path) 
	HandleErrors(err)
	scanner := bufio.NewScanner(file) 
	output := "" 
	for scanner.Scan() { 
		output = output + scanner.Text() 
	} 
	return output 
}

func FileExists(path string) bool { 
	_, err := os.Stat(path) 
	if err != nil && os.IsNotExist(err) { 
		return false 
	} else { 
		return true 
	}
}

func Type(object interface{}) reflect.Type { 
	return reflect.TypeOf(object) 
}

func StartApp(appName string, port string, routes map[string]http.HandlerFunc) { // Create a new Astilectron instance 

	server := &http.Server{
		Addr: ":" + port,
	}
	
	for url, handlerfunc := range routes { 
		http.HandleFunc(url, handlerfunc) 
	}

	// Run the server in a separate goroutine so we can monitor it
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting server:", err)
		}
	}()
		
	
	url := "http://localhost:" + port

	fmt.Println("Server running on", url)

	logger := log.New(os.Stderr, "", log.LstdFlags) 
	a, err := astilectron.New(logger, astilectron.Options{AppName:appName, SingleInstance: true, }) 
	if err != nil { 
		log.Fatal(err) 
	} 
	defer a.Close()
	if err := a.Start(); err != nil {
		log.Fatal(err)
	}

	frameless := false
	width := 600
	height := 800
	fullscreen := true 
	// Create a new window (frameless window, no UI elements)
	w, err := a.NewWindow(url, &astilectron.WindowOptions{
		Fullscreen: &fullscreen, 
		Width:  &width,
		Height: &height,
		Frame:  &frameless,
	})
	
	if err != nil {
		log.Fatal(err)
	}

	// Show the window
	if err := w.Create(); err != nil {
		log.Fatal(err)
	}

	// Wait for the window to be closed
	a.Wait() // Blocks here until the window is closed

	// Gracefully shut down the server after window is closed
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Error shutting down server:", err)
	}
	fmt.Println("Server stopped")
}

func Respond(filePath string, w http.ResponseWriter, data interface{}) { 
	tmpl, err := template.ParseFiles(filePath)
	HandleErrors(err) 
	tmpl.Execute(w, data)
}
