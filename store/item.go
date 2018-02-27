package store

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Variable is a piece of configuration
type Variable struct {
	Name  string `dynamodbav:"name" json:"name"`
	Value string `dynamodbav:"value" json:"value"`
}

// Item is the format of the configuratoin stored in dynamodb
type Item struct {
	ID          string     `dynamodbav:"id" json:"id"`
	Application string     `dynamodbav:"application" json:"application"`
	Environment string     `dynamodbav:"environment" json:"environment"`
	Variables   []Variable `dynamodbav:"variables" json:"variables"`
}

// PrintVars prints the variables in the item
func (item *Item) PrintVars(format string) {
	format = strings.ToLower(format)
	if format == "json" {
		item.printJSON()
	} else {
		item.printPlain()
	}
}

func (item *Item) printPlain() {
	for i := range item.Variables {
		fmt.Printf("%s=%s\n", item.Variables[i].Name, item.Variables[i].Value)
	}
}

func (item *Item) printJSON() {
	// The default json.unmarshal HTML escapes
	// We create a custom encoder so we don't have to HTML escape
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "   ")
	err := encoder.Encode(item.Variables)
	if err != nil {
		// TODO debug print error or something
		fmt.Println("ERROR: ", err)
	}
}

// TODO this is pretty darn primitive so make it more robust
// support other formats and what not
func parseVariables(variables string) []Variable {
	split := strings.Split(variables, ",")
	vars := make([]Variable, len(split))
	for i := range split {
		j := strings.Index(split[i], "=")
		vars[i] = Variable{
			Name:  string(split[i][:j]),
			Value: string(split[i][j+1:]),
		}
	}
	return vars
}

func parseVariablesFromFile(fileName string) ([]Variable, error) {
	variables := make([]Variable, 0)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Parse line
		line := scanner.Text()
		line = strings.TrimPrefix(line, "export") // remove export if exists
		line = strings.TrimLeft(line, " \t")      // remove all spaces on left
		words := strings.SplitN(line, "=", 2)     // split into no more than 2 strings
		if len(words) != 2 {                      // should be two strings
			// Skip. Something went wrong
			// TODO do we want to print an error?
			continue
		}
		variables = append(variables, Variable{Name: words[0], Value: words[1]})
	}

	if err = scanner.Err(); err != nil {
		return variables, err
	}
	return variables, nil
}

// CreateItem creates an item
func CreateItem(id, application, environment string, variables []Variable) Item {
	return Item{
		ID:          id,
		Application: application,
		Environment: environment,
		Variables:   variables,
	}
}

func (i *Item) String() string {
	b, _ := json.MarshalIndent(i, "", "\t")
	return string(b)
}

// attempt to decode the Variable values from base64
func (item *Item) decode() {
	// TODO add debug maybe?
	for j := range item.Variables {
		tmp, err := base64.StdEncoding.DecodeString(item.Variables[j].Value)
		if err == nil {
			item.Variables[j].Value = string(tmp)
		}
	}
}

func (item *Item) encode() {
	// TODO add debug maybe?
	for j := range item.Variables {
		tmp := base64.StdEncoding.EncodeToString([]byte(item.Variables[j].Value))
		item.Variables[j].Value = tmp
	}
}
