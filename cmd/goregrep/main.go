package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/goregrep/goregrep/goregrep"
	"golang.org/x/tools/imports"
)

func main() {
	cmd := kong.Parse(&CLI)

	switch cmd.Command() {
	case "regenerate":
		err := grep(CLI.Regenerate.File, CLI.Regenerate.References)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprint(os.Stderr, cmd.Command())
		os.Exit(1)
	}
}

var CLI struct {
	Regenerate struct {
		File       string   `default:"goregrep.yaml" help:"Specify an alternate YAML file."`
		References []string `help:"Specify YAML references."`
	} `cmd:"" help:"Replace generated code."`
}

func grep(pth string, refs []string) error {
	yml, err := os.Open(pth)
	if err != nil {
		return err
	}

	var opts []goregrep.Option

	for _, pth := range refs {
		yml, err := os.Open(pth)
		if err != nil {
			return err
		}

		opts = append(opts, goregrep.WithReferences(yml))
	}

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os: get current/working directory: %w", err)
	}

	opts = append(opts, goregrep.WithDirectory(dir))

	gofmt := imports.Options{
		Fragment:  true,
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
	}

	opts = append(opts, goregrep.WithGofmt(&gofmt))

	err = goregrep.New(yml, opts...)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	return nil
}
