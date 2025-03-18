package shankskit 

import ( 
	"fmt"; 
	"os"; 
	"bufio"; 
	"strings"; 
	"reflect"; 	
) 

type List struct { 
	Contents []interface{} 
}

func (l *List) Append (object interface{}) { 
	l.Contents = append(l.Contents, object) 
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
