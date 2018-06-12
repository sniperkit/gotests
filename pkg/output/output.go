package output

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/tools/imports"

	models "github.com/sniperkit/gotests/pkg/models"
	render "github.com/sniperkit/gotests/pkg/render"
)

type Options struct {
	PrintInputs bool
	Subtests    bool
	TemplateDir string
}

func Process(head *models.Header, funcs []*models.Function, opt *Options) ([]byte, error) {
	if opt != nil && opt.TemplateDir != "" {
		err := render.LoadCustomTemplates(opt.TemplateDir)
		if err != nil {
			return nil, fmt.Errorf("loading custom templates: %v", err)
		}
	}

	tf, err := ioutil.TempFile("", "gotests_")
	if err != nil {
		return nil, fmt.Errorf("ioutil.TempFile: %v", err)
	}
	defer tf.Close()
	defer os.Remove(tf.Name())
	b := &bytes.Buffer{}
	if err := writeTests(b, head, funcs, opt); err != nil {
		return nil, err
	}

	return ProcessImports(tf.Name(), b.Bytes())
}

func ProcessImports(filename string, src []byte) ([]byte, error) {
	out, err := imports.Process(filename, src, nil)
	if err != nil {
		return nil, fmt.Errorf("imports.Process: %v", err)
	}
	return out, nil
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func writeTests(w io.Writer, head *models.Header, funcs []*models.Function, opt *Options) error {
	b := bufio.NewWriter(w)
	if err := render.Header(b, head); err != nil {
		return fmt.Errorf("render.Header: %v", err)
	}
	for _, fun := range funcs {
		if err := render.TestFunction(b, fun, opt.PrintInputs, opt.Subtests); err != nil {
			return fmt.Errorf("render.TestFunction: %v", err)
		}
	}
	return b.Flush()
}
