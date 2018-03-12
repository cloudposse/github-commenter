package main

import (
	"io"
	"strings"
	"text/template"

	"github.com/hairyhenderson/gomplate/data"
	"github.com/hairyhenderson/gomplate/funcs"
	"os"
)

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
