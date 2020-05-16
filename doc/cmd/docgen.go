package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
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
	fBodyClass   = flag.String("bodyClass", "", "class name for body")
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
	Title     string
	Content   template.HTML
	BodyClass string
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

		if filepath.Ext(i.Name()) == ".md" {
			name := strings.TrimSuffix(i.Name(), filepath.Ext(i.Name()))
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// convert markdown to html
			rend := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
				AbsolutePrefix:     *fLinkPrefix,
				HeadingLevelOffset: 1,
				Flags:              blackfriday.CommonHTMLFlags,
			})
			out := blackfriday.Run(b,
				blackfriday.WithExtensions(blackfriday.CommonExtensions | blackfriday.AutoHeadingIDs & ^blackfriday.Autolink),
				blackfriday.WithRenderer(rend))

			// render the html
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, TemplateArgs{
				Title:     convertNameToTitle(name),
				Content:   template.HTML(out),
				BodyClass: *fBodyClass,
			}); err != nil {
				return fmt.Errorf("execute template: %s", err)
			}

			outDir := filepath.Join(*fDst, name)
			if name == "index" {
				outDir = filepath.Join(*fDst)
			}
			if err := os.MkdirAll(outDir, 0755); err != nil {
				return fmt.Errorf("make directory %s: %s", outDir, err)
			}
			if err := ioutil.WriteFile(filepath.Join(outDir, "index.html"), buf.Bytes(), permFile); err != nil {
				return fmt.Errorf("write file: %s", err)
			}
		} else {
			// copy the file directly to dst
			if err := copyFile(path, filepath.Join(*fDst, i.Name())); err != nil {
				return fmt.Errorf("copy file: %s", err)
			}
		}

		return nil
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func convertNameToTitle(f string) string {
	if f == "index" {
		return *fIndexTitle
	}
	return *fTitlePrefix + strings.ReplaceAll(f, "__", "/")
}
