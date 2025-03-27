package shankskit

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/asticode/go-astilectron"
)

// Creates a new List object with some default values
func NewList(contents ...interface{}) *List {
	return &List{
		Contents:   contents,
		RangeIndex: -1,
		RangeItem:  nil,
	}
}

type List struct {
	Contents   []interface{}
	RangeIndex int
	RangeItem  interface{}
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
		case string:
			representation += fmt.Sprintf(`'%+v'`, typedItem)
		case *List:
			representation += typedItem.GetSyntax()
		case []interface{}:
			tempItem := NewList(typedItem...)
			representation += tempItem.GetSyntax()
		default:
			representation += fmt.Sprintf(`%+v`, typedItem)
		}

		representation += ", "
	}
	if 1 < len(representation) {
		representation = representation[0 : -2+len(representation)]
	}
	representation += `]`
	return representation
}

func NewGrid(rows int, columns int) *Grid {
	placeholderValue := "nil"
	contents := [][]string{}
	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		row_instance := []string{}
		for columnIndex := 0; columnIndex < columns; columnIndex++ {
			row_instance = append(row_instance, placeholderValue)
		}
		contents = append(contents, row_instance)
	}

	gridInstance := Grid{
		Contents: contents,
		Rows:     rows,
		Columns:  columns,
	}

	return &gridInstance
}

type Grid struct {
	Contents [][]string
	Rows     int
	Columns  int
}

func (g *Grid) GetValue(xCoord int, yCoord int) string {
	return g.Contents[yCoord][xCoord]
}

func (g *Grid) GetSyntax() string {
	var syntax string
	var padding_size int
	var value string
	for yCoord := range g.Rows {
		for xCoord := range g.Columns {
			padding_size = max(padding_size, len(g.GetValue(xCoord, yCoord)))
		}
	}

	for yCoord := range g.Rows {
		for xCoord := range g.Columns {
			value = g.GetValue(xCoord, yCoord)
			syntax += PaddRight(value, padding_size, " ") + "|"
		}
		syntax += "\n"
	}

	return syntax
}

func (g *Grid) InsertRow(row []string, index int) {
	// Check for out-of-bounds index
	if index < 0 || index > g.Rows {
		fmt.Println("Error: Invalid index for row insertion.")
		return
	}

	// Ensure the row has the correct number of columns
	if len(row) != g.Columns {
		fmt.Println("Error: Row length must match the number of columns in the grid.")
		return
	}

	var contents [][]string
	rows := g.Rows

	for i := 0; i < rows; i++ {
		if i == index {
			contents = append(contents, row) // Insert the new row
		}
		contents = append(contents, g.Contents[i])
	}

	// Edge case: If inserting at the last index, append at the end
	if index == g.Rows {
		contents = append(contents, row)
	}

	g.Contents = contents
	g.Rows++
}

func (g *Grid) InsertColumn(column []string, index int) {
	if index < 0 || index > g.Columns {
		fmt.Println("Invalid index for column insertion")
		return
	}

	if len(column) != g.Rows {
		fmt.Println("Error: Column length must match the number of rows in the grid.")
		return
	}

	for rowIndex := 0; rowIndex < g.Rows; rowIndex++ {
		newRow := make([]string, 0, g.Columns+1) // Create a new row with increased capacity
		for colIndex := 0; colIndex < g.Columns; colIndex++ {
			if colIndex == index {
				newRow = append(newRow, column[rowIndex]) // Insert before index
			}
			newRow = append(newRow, g.Contents[rowIndex][colIndex])
		}
		if index == g.Columns { // Edge case: inserting at the last index
			newRow = append(newRow, column[rowIndex])
		}
		g.Contents[rowIndex] = newRow
	}
	g.Columns++
}

func (g *Grid) Print() {
	syntax := g.GetSyntax()
	Print(syntax)
}

func PaddLeft(text string, totallength int, paddingcharacter string) string {
	var paddingLength int
	var output string
	paddingLength = totallength - len(text)
	for i := 0; i < paddingLength; i++ {
		output += paddingcharacter
	}
	output += text
	return output
}

func PaddRight(text string, totallength int, paddingcharacter string) string {
	var paddingLength int
	var output string
	output += text
	paddingLength = totallength - len(text)
	for i := 0; i < paddingLength; i++ {
		output += paddingcharacter
	}
	return output
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

func ReadLines(path string) []string {
	var lines []string
	file, err := os.Open(path)
	HandleErrors(err)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
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

func IntToString(i int) string {
	return fmt.Sprintf("%d", i)
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

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting server:", err)
		}
	}()

	url := "http://localhost:" + port

	fmt.Println("Server running on", url)

	logger := log.New(os.Stderr, "", log.LstdFlags)
	a, err := astilectron.New(logger, astilectron.Options{
		AppName:        appName,
		SingleInstance: true,
	})
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
		Width:      &width,
		Height:     &height,
		Frame:      &frameless,
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
