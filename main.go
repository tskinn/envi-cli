package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/tskinn/envi-cli/store"
	"github.com/urfave/cli"
)

func main() {
	var tableName, awsRegion, application, environment, id, variables, filePath string
	app := cli.NewApp()

	app.Description = "A simple key-value store cli for dynamodb"
	app.Name = "envi"
	app.Usage = ""
	app.UsageText = "envi set --application myapp --environment dev --variables=one=eno,two=owt,three=eerht\n   envi get --key <key>\n   envi --getall --key <key>"

	globalFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "table, t",
			Value:       "envi",
			Usage:       "name of the dynamodb to store values",
			Destination: &tableName,
		},
		cli.StringFlag{
			Name:        "region, r",
			Value:       "us-east-1",
			Usage:       "name of the aws region in which dynamodb table resides",
			Destination: &awsRegion,
		},
		cli.StringFlag{
			Name:        "id, i",
			Value:       "",
			Usage:       "id of the application environment combo; if id is not provided then application__environment is used as the id",
			Destination: &id,
		},
		cli.StringFlag{
			Name:        "application, a",
			Value:       "",
			Usage:       "name of the application",
			Destination: &application,
		},
		cli.StringFlag{
			Name:        "environment, e",
			Value:       "",
			Usage:       "name of the environment",
			Destination: &environment,
		},
	}

	setCommand := cli.Command{
		Name:    "set",
		Aliases: []string{"s"},
		Usage:   "save application configuraton in dynamodb",
		Action: func(c *cli.Context) error {
			store.Init(c.String("region"), c.String("table"))
			if c.IsSet("application") && c.IsSet("environment") && (c.IsSet("variables") || c.IsSet("file")) {
				tID := c.String("id")
				if !c.IsSet("id") {
					tID = application + "__" + environment
				}
				if filePath != "" && variables != "" {
					fmt.Println("Cannot do that")
					// TODO print error not supposed to have both set at same time
				} else if filePath != "" {
					return store.SaveFromFile(tID, application, environment, filePath)
				} else if variables != "" {
					return store.Save(tID, application, environment, variables)
				}
				return fmt.Errorf("Bad form")
				// TODO return error because neither was set
			}
			return nil
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "variables, v",
				Value:       "",
				Usage:       "env variables to store in the form of key=value,key2=value2,key3=value3",
				Destination: &variables,
			},
			cli.StringFlag{
				Name:        "file, f",
				Value:       "",
				Usage:       "path to a shell file that exports env vars",
				Destination: &filePath,
			},
		},
	}
	setCommand.Flags = append(setCommand.Flags, globalFlags...)

	getCommand := cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get the application configuration for a particular application",
		Action: func(c *cli.Context) error {
			store.Init(c.String("region"), c.String("table"))
			var item store.Item
			var err error
			if c.IsSet("id") {
				item, err = store.Get(c.String("id"))
			} else if c.IsSet("application") && c.IsSet("environment") {
				item, err = store.Get(c.String("application") + "__" + c.String("environment"))
			} else {
				return errors.New("incorrect flags")
			}
			if err != nil {
				return err
			}
			fmt.Println(item.String())
			return nil
		},
		Flags: []cli.Flag{},
	}
	getCommand.Flags = append(getCommand.Flags, globalFlags...)

	app.Commands = []cli.Command{
		setCommand,
		getCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

// func print(output string, keys, values []string) {
// 	switch output {
// 	case "json":
// 		printJSON(keys, values)
// 	case "text":
// 		printText(keys, values)
// 	default:
// 		printText(keys, values)
// 	}
// }

// func printJSON(keys, values []string) {
// 	list := make([]map[string]string, 0)
// 	for i := 0; i < len(keys); i++ {
// 		tmp := map[string]string{"key": keys[i], "value": values[i]}
// 		list = append(list, tmp)
// 	}
// 	raw, err := json.Marshal(list)
// 	if err != nil {
// 		// print something
// 		return
// 	}
// 	fmt.Printf("%s", raw)
// }

// func printText(keys, values []string) {
// 	for i := 0; i < len(keys); i++ {
// 		fmt.Printf("%s=%s\n", keys[i], values[i])
// 	}
// }
