package pkg

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type SimpleNodeRend struct {
	Writer html.Writer
}

func NewSimpleNodeRend() renderer.NodeRenderer {
	return &SimpleNodeRend{
		Writer: html.DefaultWriter,
	}
}

func (r *SimpleNodeRend) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindString, r.renderText)
}

func (r *SimpleNodeRend) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	r.Writer.Write(w, node.Text(source))
	return ast.WalkContinue, nil
}

func ReadFileComplex(path, fileName string, nameAsDate bool) *FileContent {
	var (
		err           error
		haveDateParam bool
		buf           bytes.Buffer
		b             []byte
		context       = parser.NewContext()
		currentTime   = time.Now()
		isoDate       = currentTime.Format("2006-Jan-02")
		fullDate      = currentTime.Format("January 2, 2006")

		markdown = goldmark.New(
			goldmark.WithExtensions(meta.Meta),
			goldmark.WithRenderer(renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(NewSimpleNodeRend(), 1000)))),
		)
		data = &FileContent{
			Name:    fileName,
			Date:    strings.Split(fileName, ".")[0],
			Content: "",
			Err:     nil,
		}
	)

	file, err := os.Open(path + fileName)
	if err != nil {
		data.Err = err
		return data
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err = ioutil.ReadAll(file)

	if err = markdown.Convert(b, &buf, parser.WithContext(context)); err != nil {
		data.Err = err
		return data
	}
	metaData := meta.Get(context)

	if !nameAsDate {
		if data.Date, haveDateParam = metaData["Date"].(string); !haveDateParam {
			return nil
		}
	}

	if data.Date != isoDate && data.Date != fullDate {
		return nil
	}
	data.Content = buf.String()
	return data
}
