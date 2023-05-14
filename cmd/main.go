package main

import (
	"fmt"
	"os"

	"github.com/Angelmaneuver/xlsx2html/internal/kamipro"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "excel2html",
		Usage: "Generates HTML codes from the contents of an Excel sheets.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "`Path` to the Excel file to be used for generate.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output `Path` for HTML to be generate.",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			err := kamipro.Start(ctx.String("i"), ctx.String("o"))
			if err != nil {
				return cli.Exit(err, -1)
			}

			fmt.Println("Process is completed.")
			return nil
		},
	}

	app.Run(os.Args)
}
