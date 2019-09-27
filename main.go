package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/cppforlife/go-patch/patch"
	flag "github.com/spf13/pflag"
	"golang.org/x/xerrors"
	yaml "gopkg.in/yaml.v2"
)

var outPath = flag.StringP("output", "o", "-",
	"Name of file to write to")

func render(args []string, output io.Writer) error {
	if len(args) < 1 {
		return xerrors.New("no YAML file given")
	}

	docBytes, err := ioutil.ReadFile(args[0])
	if err != nil {
		return xerrors.Errorf("could not read YAML document %s: %w", args[0], err)
	}
	var doc interface{}
	err = yaml.Unmarshal(docBytes, &doc)
	if err != nil {
		return xerrors.Errorf("could not unmarshal YAML document %s: %w", args[0], err)
	}

	for _, opPath := range args[1:] {
		opBytes, err := ioutil.ReadFile(opPath)
		if err != nil {
			return xerrors.Errorf("failed to read ops file %s: %w", opPath, err)
		}
		var opDefs []patch.OpDefinition
		err = yaml.Unmarshal(opBytes, &opDefs)
		if err != nil {
			return xerrors.Errorf("could not unmarshal ops file %s: %w", opPath, err)
		}
		ops, err := patch.NewOpsFromDefinitions(opDefs)
		if err != nil {
			return xerrors.Errorf("could not create ops from ops file %s: %w", opPath, err)
		}
		result, err := ops.Apply(doc)
		if err != nil {
			if docBytes, marshalErr := yaml.Marshal(doc); marshalErr == nil {
				fmt.Fprintf(os.Stderr, "%s\n", docBytes)
			}
			return xerrors.Errorf("could not apply ops file %s: %w", opPath, err)
		}
		doc = result
	}

	docBytes, err = yaml.Marshal(doc)
	if err != nil {
		return xerrors.Errorf("could not marshal result: %w", err)
	}

	fmt.Fprintf(output, "%s\n", string(docBytes))

	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:    %s file.yaml [op.yaml...]\n", os.Args[0])
		flag.CommandLine.PrintDefaults()
	}
	flag.Parse()

	var err error
	var output = os.Stdout
	if *outPath != "-" {
		output, err = os.Create(*outPath)
	}
	if err == nil {
		err = render(flag.Args(), output)
	}

	if err != nil {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "\n%+v\n", err)
		os.Exit(2)
	}
}
