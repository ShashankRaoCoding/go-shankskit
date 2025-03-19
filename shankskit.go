package shankskit 

import ( 
	"fmt"; 
	"os"; 
	"bufio"; 
	"strings"; 
	"reflect"; 	
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
