package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Support an interface similar - but not identical - to the IPython (canonical Jupyter kernel).
// See http://ipython.readthedocs.io/en/stable/api/generated/IPython.display.html#IPython.display.display
// for a good overview of the support types.

const (
	MIMETypeHTML       = "text/html"
	MIMETypeJavaScript = "application/javascript"
	MIMETypeJPEG       = "image/jpeg"
	MIMETypeJSON       = "application/json"
	MIMETypeLatex      = "text/latex"
	MIMETypeMarkdown   = "text/markdown"
	MIMETypePNG        = "image/png"
	MIMETypePDF        = "application/pdf"
	MIMETypeSVG        = "image/svg+xml"
	MIMETypeText       = "text/plain"
)

/**
 * general interface, allows libraries to fully specify
 * how their data is displayed by Jupyter.
 * Supports multiple MIME formats.
 *
 * Note that Data defined above is an alias:
 * libraries can implement Renderer without importing gophernotes
 */
type Renderer = interface {
	Render() Data
}

/**
 * simplified interface, allows libraries to specify
 * how their data is displayed by Jupyter.
 * Supports multiple MIME formats.
 *
 * Note that MIMEMap defined above is an alias:
 * libraries can implement SimpleRenderer without importing gophernotes
 */
type SimpleRenderer = interface {
	SimpleRender() MIMEMap
}

/**
 * specialized interfaces, each is dedicated to a specific MIME type.
 *
 * They are type aliases to emphasize that method signatures
 * are the only important thing, not the interface names.
 * Thus libraries can implement them without importing gophernotes
 */
type HTMLer = interface {
	HTML() string
}
type JavaScripter = interface {
	JavaScript() string
}
type JPEGer = interface {
	JPEG() []byte
}
type JSONer = interface {
	JSON() map[string]interface{}
}
type Latexer = interface {
	Latex() string
}
type Markdowner = interface {
	Markdown() string
}
type PNGer = interface {
	PNG() []byte
}
type PDFer = interface {
	PDF() []byte
}
type SVGer = interface {
	SVG() string
}

// injected as placeholder in the interpreter, it's then replaced at runtime
// by a closure that knows how to talk with Jupyter
func stubDisplay(Data) error {
	return errors.New("cannot display: connection with Jupyter not available")
}

// if vals[] contain a single non-nil value which is auto-renderable,
// convert it to Data and return it.
// otherwise return MakeData("text/plain", fmt.Sprint(vals...))
func (kernel *Kernel) autoRenderResults(vals []interface{}) Data {
	return MakeData(MIMETypeText, fmt.Sprint(vals...))
}

var autoRenderers = map[string]func(Data, interface{}) Data{
	"Renderer": func(d Data, i interface{}) Data {
		if r, ok := i.(Renderer); ok {
			x := r.Render()
			d.Data = merge(d.Data, x.Data)
			d.Metadata = merge(d.Metadata, x.Metadata)
			d.Transient = merge(d.Transient, x.Transient)
		}
		return d
	},
	"SimpleRenderer": func(d Data, i interface{}) Data {
		if r, ok := i.(SimpleRenderer); ok {
			x := r.SimpleRender()
			d.Data = merge(d.Data, x)
		}
		return d
	},
	"HTMLer": func(d Data, i interface{}) Data {
		if r, ok := i.(HTMLer); ok {
			d.Data = ensure(d.Data)
			d.Data[MIMETypeHTML] = r.HTML()
		}
		return d
	},
	"JavaScripter": func(d Data, i interface{}) Data {
		if r, ok := i.(JavaScripter); ok {
			d.Data = ensure(d.Data)
			d.Data[MIMETypeJavaScript] = r.JavaScript()
		}
		return d
	},
	"JPEGer": func(d Data, i interface{}) Data {
		if r, ok := i.(JPEGer); ok {
			d.Data = ensure(d.Data)
			d.Data[MIMETypeJPEG] = r.JPEG()
		}
		return d
	},
	"JSONer": func(d Data, i interface{}) Data {
		if r, ok := i.(JSONer); ok {
			d.Data = ensure(d.Data)
			d.Data[MIMETypeJSON] = r.JSON()
		}
		return d
	},
	"Latexer": func(d Data, i interface{}) Data {
		if r, ok := i.(Latexer); ok {
			d.Data = ensure(d.Data)
			d.Data[MIMETypeLatex] = r.Latex()
		}
		return d
	},
	"Markdowner": func(d Data, i interface{}) Data {
		if r, ok := i.(Markdowner); ok {
			d.Data = ensure(d.Data)
			d.Data[MIMETypeMarkdown] = r.Markdown()
		}
		return d
	},
	"PNGer": func(d Data, i interface{}) Data {
		if r, ok := i.(PNGer); ok {
			d.Data = ensure(d.Data)
			d.Data[MIMETypePNG] = r.PNG()
		}
		return d
	},
	"PDFer": func(d Data, i interface{}) Data {
		if r, ok := i.(PDFer); ok {
			d.Data = ensure(d.Data)
			d.Data[MIMETypePDF] = r.PDF()
		}
		return d
	},
	"SVGer": func(d Data, i interface{}) Data {
		if r, ok := i.(SVGer); ok {
			d.Data = ensure(d.Data)
			d.Data[MIMETypeSVG] = r.SVG()
		}
		return d
	},
	"Image": func(d Data, i interface{}) Data {
		if r, ok := i.(image.Image); ok {
			b, mimeType, err := encodePng(r)
			if err != nil {
				d = makeDataErr(err)
			} else {
				d.Data = ensure(d.Data)
				d.Data[mimeType] = b
				d.Metadata = merge(d.Metadata, imageMetadata(r))
			}
		}
		return d
	},
}

// detect and render data types that should be auto-rendered graphically
func fillDefaults(data Data, arg interface{}, s string, b []byte, mimeType string, err error) Data {
	if err != nil {
		return makeDataErr(err)
	}
	if data.Data == nil {
		data.Data = make(MIMEMap)
	}
	// cannot autodetect the mime type of a string
	if len(s) != 0 && len(mimeType) != 0 {
		data.Data[mimeType] = s
	}
	// ensure plain text is set
	if data.Data[MIMETypeText] == "" {
		if len(s) == 0 {
			s = fmt.Sprint(arg)
		}
		data.Data[MIMETypeText] = s
	}
	// if []byte is available, use it
	if len(b) != 0 {
		if len(mimeType) == 0 {
			mimeType = http.DetectContentType(b)
		}
		if len(mimeType) != 0 && mimeType != MIMETypeText {
			data.Data[mimeType] = b
		}
	}
	return data
}

// do our best to render data graphically
func render(mimeType string, data interface{}) Data {
	var s string
	var b []byte
	var err error
	switch data := data.(type) {
	case string:
		s = data
	case []byte:
		b = data
	case io.Reader:
		b, err = ioutil.ReadAll(data)
	case io.WriterTo:
		var buf bytes.Buffer
		data.WriteTo(&buf)
		b = buf.Bytes()
	default:
		panic(fmt.Errorf("unsupported type, cannot render: %T", data))
	}
	return fillDefaults(Data{}, data, s, b, mimeType, err)
}

func makeDataErr(err error) Data {
	return Data{
		Data: MIMEMap{
			"ename":     "ERROR",
			"evalue":    err.Error(),
			"traceback": nil,
			"status":    "error",
		},
	}
}

func Any(mimeType string, data interface{}) Data {
	return render(mimeType, data)
}

// same as Any("", data), autodetects MIME type
func Auto(data interface{}) Data {
	return render("", data)
}

func MakeData(mimeType string, data interface{}) Data {
	d := Data{
		Data: MIMEMap{
			mimeType: data,
		},
	}
	if mimeType != MIMETypeText {
		d.Data[MIMETypeText] = fmt.Sprint(data)
	}
	return d
}

func MakeData3(mimeType string, plaintext string, data interface{}) Data {
	return Data{
		Data: MIMEMap{
			MIMETypeText: plaintext,
			mimeType:     data,
		},
	}
}

func File(mimeType string, path string) Data {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return Any(mimeType, bytes)
}

func HTML(html string) Data {
	return MakeData(MIMETypeHTML, html)
}

func JavaScript(javascript string) Data {
	return MakeData(MIMETypeJavaScript, javascript)
}

func JPEG(jpeg []byte) Data {
	return MakeData(MIMETypeJPEG, jpeg)
}

func JSON(json map[string]interface{}) Data {
	return MakeData(MIMETypeJSON, json)
}

func Latex(latex string) Data {
	return MakeData3(MIMETypeLatex, latex, "$"+strings.Trim(latex, "$")+"$")
}

func Markdown(markdown string) Data {
	return MakeData(MIMETypeMarkdown, markdown)
}

func Math(latex string) Data {
	return MakeData3(MIMETypeLatex, latex, "$$"+strings.Trim(latex, "$")+"$$")
}

func PDF(pdf []byte) Data {
	return MakeData(MIMETypePDF, pdf)
}

func PNG(png []byte) Data {
	return MakeData(MIMETypePNG, png)
}

func SVG(svg string) Data {
	return MakeData(MIMETypeSVG, svg)
}

// MIME encapsulates the data and metadata into a Data.
// The 'data' map is expected to contain at least one {key,value} pair,
// with value being a string, []byte or some other JSON serializable representation,
// and key equal to the MIME type of such value.
// The exact structure of value is determined by what the frontend expects.
// Some easier-to-use functions for common formats supported by the Jupyter frontend
// are provided by the various functions above.
func MIME(data, metadata MIMEMap) Data {
	return Data{data, metadata, nil}
}
