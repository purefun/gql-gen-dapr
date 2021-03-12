package cmd

import (
	"io/ioutil"
	"os"

	"github.com/purefun/gql-gen-dapr/generator"
	"github.com/urfave/cli/v2"
)

func Execute() {
	app := cli.NewApp()
	app.Name = "gql-gen-dapr"
	app.Usage = "Generate dapr app using GraphQL schema"
	app.Version = generator.Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "schemaFile", Aliases: []string{"f"}, Usage: "graphql schema file path", Required: true},
		&cli.StringFlag{Name: "package", Aliases: []string{"pkg"}, Usage: "package name", Required: true},
		&cli.StringFlag{Name: "service", Aliases: []string{"s"}, Usage: "service name", Required: true},
		&cli.StringFlag{Name: "out", Aliases: []string{"o"}, Usage: "output dir", Required: true},
	}
	app.Action = func(ctx *cli.Context) error {
		schemaFile := ctx.String("schemaFile")
		packageName := ctx.String("package")
		serviceName := ctx.String("service")
		out := ctx.String("out")

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

		err = ioutil.WriteFile(out+"/generated.go", []byte(outString), 0644)
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
