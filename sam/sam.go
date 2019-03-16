package sam

import (
	// "fmt"
	"bufio"
	"os"
	"strings"
)

type SamParser struct {
	OnOneLines   map[string]func(line string) error
	OnMultiLines map[string]func(lines []string) error
	OnTable      func(columns map[string]string) error
}

func (sam SamParser) tableEach(line string, scanner *bufio.Scanner) error {
	tokens := strings.Split(line, `|`)
	tableColumns := make([]string, len(tokens))
	// if strings.Trim(tokens[0], ` `) == "" {
	// 	return nil
	// }
	// fmt.Printf(" >> %v\n", tableColumns)
	// fmt.Printf(" >> LINE >> : %v\n",tokens)
	for i, name := range tokens {
		tableColumns[i] = name
	}

	scanner.Scan() // Skip bar line.

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, `|`) {
			return nil
		}

		tokens := strings.Split(line, `|`)
		tableValues := make(map[string]string, len(tokens))
		for i, token := range tokens {
			// fmt.Printf(" >> %v", tableValues)
			tableValues[tableColumns[i]] = strings.Trim(token, ` `)
		}
		// fmt.Printf(" >>>> tableValues: %v\n",tableValues)

		if err := sam.OnTable(tableValues); err != nil {
			return err
		}
	}

	return nil
}

func (sam SamParser) Start(filePath string) error {
	fp, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()

		// fmt.Printf(" >> line: %v\n", line)

		if strings.HasPrefix(line, `|`) {
			if err := sam.tableEach(line, scanner); err != nil {
				return nil
			}

			continue
		}

		for search, onOneLine := range sam.OnOneLines {
			// fmt.Printf(" ?? %v\n", search)

			tokens := strings.Split(line, ` `)
			// fmt.Printf(" ?? %v", tokens)
			if tokens[0] == search {
				onOneLine(tokens[1])
			}
		}
	}

	return nil
}
