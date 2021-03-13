package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/purefun/gql-gen-dapr/generator"
	"github.com/urfave/cli/v2"
)

func Execute() {
	app := cli.NewApp()
	app.Name = "gql-gen-dapr"
	app.Usage = "Generate dapr app using GraphQL schema"
	app.UsageText = "gql-gen-dapr graphql-file [flags...]"
	app.Version = generator.Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "package", Aliases: []string{"pkg"}, Usage: "package name", DefaultText: "main"},
		&cli.StringFlag{Name: "service", Aliases: []string{"s"}, Usage: "service name", DefaultText: "The name of --schemaFile"},
		&cli.StringFlag{Name: "out", Aliases: []string{"o"}, Usage: "output dir", DefaultText: "same dir with --schemaFile"},
	}
	app.Action = func(ctx *cli.Context) error {
		schemaFile := ctx.Args().Get(0)
		if schemaFile == "" {
			return fmt.Errorf("graphql file should be set, \n\nfor example: gql-gen-dapr schema.graphql")
		}
		packageName := ctx.String("package")
		if packageName == "" {
			packageName = "main"
		}
		serviceName := ctx.String("service")
		if serviceName == "" {
			serviceName = strings.TrimSuffix(filepath.Base(schemaFile), filepath.Ext(schemaFile))
		}
		out := ctx.String("out")
		if out == "" {
			out = filepath.Dir(schemaFile)
		}

		schemaBytes, err := ioutil.ReadFile(schemaFile)
		if err != nil {
			return err
		}
		g := generator.NewGenerator(generator.Options{
			PackageName:       packageName,
			ServiceName:       serviceName,
			GenHeaderComments: true,
		})

		g.AddSource(schemaFile, string(schemaBytes))

		outString, err := g.Generate()
		if err != nil {
			return err
		}

		outFile := strings.TrimSuffix(out, "/") + "/" + serviceName + ".dapr.go"

		err = ioutil.WriteFile(outFile, []byte(outString), 0644)
		if err != nil {
			return err
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
