package main

import (
	"errors"
	"fmt"
	"github.com/hairyhenderson/gomplate/data"
	"github.com/hairyhenderson/gomplate/env"
	"github.com/hairyhenderson/gomplate/funcs"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
	"text/template"
)

// GomplateOpts -
type GomplateOpts struct {
	version           bool
	dataSources       []string
	dataSourceHeaders []string
	lDelim            string
	rDelim            string

	input       string
	inputFiles  []string
	inputDir    string
	outputFiles []string
	outputDir   string
	excludeGlob []string
}

var opts GomplateOpts

var cleanupHooks = make([]func(), 0)

// Gomplate -
type Gomplate struct {
	funcMap    template.FuncMap
	leftDelim  string
	rightDelim string
}

// Context for templates
type Context struct {
}

// RunTemplate -
func (g *Gomplate) RunTemplate(t *tplate) error {
	context := &Context{}
	tmpl, err := t.toGoTemplate(g)
	if err != nil {
		return err
	}

	switch t.target.(type) {
	case io.Closer:
		// nolint: errcheck
		defer t.target.(io.Closer).Close()
	}
	err = tmpl.Execute(t.target, context)
	return err
}

// NewGomplate -
func NewGomplate(d *data.Data, leftDelim, rightDelim string) *Gomplate {
	return &Gomplate{
		leftDelim:  leftDelim,
		rightDelim: rightDelim,
		funcMap:    initFuncs(d),
	}
}

func runTemplate(o *GomplateOpts) error {
	defer runCleanupHooks()
	d := data.NewData(o.dataSources, o.dataSourceHeaders)
	addCleanupHook(d.Cleanup)

	g := NewGomplate(d, o.lDelim, o.rDelim)

	tmpl, err := gatherTemplates(o)
	if err != nil {
		return err
	}
	for _, t := range tmpl {
		if err := g.RunTemplate(t); err != nil {
			return err
		}
	}
	return nil
}

func addCleanupHook(hook func()) {
	cleanupHooks = append(cleanupHooks, hook)
}

func runCleanupHooks() {
	for _, hook := range cleanupHooks {
		hook()
	}
}

// initFuncs - The function mappings are defined here!
func initFuncs(d *data.Data) template.FuncMap {
	f := template.FuncMap{}
	funcs.AddDataFuncs(f, d)
	funcs.AWSFuncs(f)
	funcs.AddBase64Funcs(f)
	funcs.AddNetFuncs(f)
	funcs.AddReFuncs(f)
	funcs.AddStringFuncs(f)
	funcs.AddEnvFuncs(f)
	funcs.AddConvFuncs(f)
	funcs.AddTimeFuncs(f)
	funcs.AddMathFuncs(f)
	funcs.AddCryptoFuncs(f)
	funcs.AddFileFuncs(f)
	funcs.AddSockaddrFuncs(f)
	return f
}

// Env - Map environment variables for use in a template
func (c *Context) Env() map[string]string {
	env := make(map[string]string)
	for _, i := range os.Environ() {
		sep := strings.Index(i, "=")
		env[i[0:sep]] = i[sep+1:]
	}
	return env
}

func validateOpts(cmd *cobra.Command, args []string) error {
	if cmd.Flag("in").Changed && cmd.Flag("file").Changed {
		return errors.New("--in and --file may not be used together")
	}

	if len(opts.inputFiles) != len(opts.outputFiles) {
		return fmt.Errorf("must provide same number of --out (%d) as --file (%d) options", len(opts.outputFiles), len(opts.inputFiles))
	}

	if cmd.Flag("input-dir").Changed && (cmd.Flag("in").Changed || cmd.Flag("file").Changed) {
		return errors.New("--input-dir can not be used together with --in or --file")
	}

	if cmd.Flag("output-dir").Changed {
		if cmd.Flag("out").Changed {
			return errors.New("--output-dir can not be used together with --out")
		}
		if !cmd.Flag("input-dir").Changed {
			return errors.New("--input-dir must be set when --output-dir is set")
		}
	}
	return nil
}

func newGomplateCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "gomplate",
		Short:   "Process text files with Go templates",
		PreRunE: validateOpts,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTemplate(&opts)
		},
		Args: cobra.NoArgs,
	}
	return rootCmd
}

func initFlags(command *cobra.Command) {
	command.Flags().BoolVarP(&opts.version, "version", "v", false, "print the version")

	command.Flags().StringArrayVarP(&opts.inputFiles, "file", "f", []string{"-"}, "Template `file` to process. Omit to use standard input, or use --in or --input-dir")
	command.Flags().StringVarP(&opts.input, "in", "i", "", "Template `string` to process (alternative to --file and --input-dir)")
	command.Flags().StringVar(&opts.inputDir, "input-dir", "", "`directory` which is examined recursively for templates (alternative to --file and --in)")
	command.Flags().StringArrayVar(&opts.excludeGlob, "exclude", []string{}, "glob of files to not parse")
	command.Flags().StringArrayVarP(&opts.outputFiles, "out", "o", []string{"-"}, "output `file` name. Omit to use standard output.")
	command.Flags().StringVar(&opts.outputDir, "output-dir", ".", "`directory` to store the processed templates. Only used for --input-dir")

	command.Flags().StringArrayVarP(&opts.dataSources, "datasource", "d", nil, "`datasource` in alias=URL form. Specify multiple times to add multiple sources.")
	command.Flags().StringArrayVarP(&opts.dataSourceHeaders, "datasource-header", "H", nil, "HTTP `header` field in 'alias=Name: value' form to be provided on HTTP-based data sources. Multiples can be set.")

	ldDefault := env.Getenv("GOMPLATE_LEFT_DELIM", "{{")
	rdDefault := env.Getenv("GOMPLATE_RIGHT_DELIM", "}}")
	command.Flags().StringVar(&opts.lDelim, "left-delim", ldDefault, "override the default left-`delimiter` [$GOMPLATE_LEFT_DELIM]")
	command.Flags().StringVar(&opts.rDelim, "right-delim", rdDefault, "override the default right-`delimiter` [$GOMPLATE_RIGHT_DELIM]")
}
