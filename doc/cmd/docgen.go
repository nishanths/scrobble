package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/russross/blackfriday.v2"
)

var (
	fSrc         = flag.String("src", "", "source directory")
	fDst         = flag.String("dst", "", "output directory")
	fTitlePrefix = flag.String("titlePrefix", "", "title prefix")
	fLinkPrefix  = flag.String("linkPrefix", "", "title prefix")
	fIndexTitle  = flag.String("indexTitle", "", "title for index page")
	fTemplate    = flag.String("template", "", "path to template file")
)

const (
	permFile = 0644
	permDir  = 0755
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

type TemplateArgs struct {
	Title   string
	Content template.HTML
}

func run(ctx context.Context) error {
	flag.Parse()

	if *fSrc == "" || *fDst == "" {
		return errors.New("both -src and -dst are required")
	}
	if *fTemplate == "" {
		return errors.New("must specify -template")
	}

	tmpl := template.Must(template.ParseFiles(*fTemplate))

	if err := os.MkdirAll(*fDst, 0755); err != nil {
		return fmt.Errorf("make directory %s: %s", *fDst, err)
	}
	return filepath.Walk(*fSrc, func(path string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == *fSrc {
			return nil // skip root
		}
		if i.IsDir() {
			return nil // skip directories; only consider top-level files
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(i.Name(), filepath.Ext(i.Name()))
		rend := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			AbsolutePrefix:     *fLinkPrefix,
			HeadingLevelOffset: 1,
			Flags:              blackfriday.CommonHTMLFlags,
		})
		out := blackfriday.Run(b, blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.AutoHeadingIDs), blackfriday.WithRenderer(rend))

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, TemplateArgs{
			Title:   convertNameToTitle(name),
			Content: template.HTML(out),
		}); err != nil {
			return fmt.Errorf("execute template: %s", err)
		}

		outPath := filepath.Join(*fDst, name+".html")
		if err := ioutil.WriteFile(outPath, buf.Bytes(), permFile); err != nil {
			return fmt.Errorf("write file: %s", err)
		}
		return nil
	})
}

func convertNameToTitle(f string) string {
	if f == "index" {
		return *fIndexTitle
	}
	return *fTitlePrefix + strings.ReplaceAll(f, "__", "/")
}
