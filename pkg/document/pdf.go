package document

import (
	"bytes"
	"github.com/signintech/gopdf"
	"github.com/signintech/gopdf/fontmaker/core"
	"github.com/sirupsen/logrus"
	"gitlab.com/72nd/acc/pkg/document/utils"
	"strings"
)

type Doc struct {
	fontSize        int
	defaultFontSize int
	lineSpread      float64
	capValue        float64
	fontStyle       string
	currentX        float64
	currentY        float64
	Pdf             gopdf.GoPdf
}

func NewDoc(fontSize int, lineSpread float64) Doc {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4, Unit: gopdf.Unit_MM})
	if err := pdf.AddTTFFontByReaderWithOption("lato", bytes.NewBuffer(utils.LatoHeavy()), gopdf.TtfOption{Style: gopdf.Bold}); err != nil {
		logrus.Fatal("error adding lato heavy to utils: ", err)
	}
	latoRegular := utils.LatoRegular()
	if err := pdf.AddTTFFontByReaderWithOption("lato", bytes.NewBuffer(latoRegular), gopdf.TtfOption{Style: gopdf.Regular, UseKerning: true}); err != nil {
		logrus.Fatal("error adding lato regular to utils: ", err)
	}
	var parser core.TTFParser
	if err := parser.ParseByReader(bytes.NewBuffer(latoRegular)); err != nil {
		logrus.Fatal("error while parsing font for height calculation: ", err)
	}
	doc := Doc{
		fontSize:        fontSize,
		defaultFontSize: fontSize,
		lineSpread:      lineSpread,
		capValue:        float64(parser.CapHeight()) * 1000.0 / float64(parser.UnitsPerEm()),
		fontStyle:       "",
		currentX:        0,
		currentY:        0,
		Pdf:             pdf,
	}
	doc.SetFontSize(fontSize)
	return doc
}

func (d *Doc) SetPosition(x, y float64) {
	d.Pdf.SetX(x)
	d.currentX = x
	d.Pdf.SetY(y)
	d.currentY = y
}

func (d *Doc) SetFontSize(size int) {
	d.fontSize = size
	if err := d.Pdf.SetFont("lato", "", size); err != nil {
		logrus.Fatal("error while changing Pdf font size: ", err)
	}
}

func (d *Doc) DefaultFontSize() {
	d.SetFontSize(d.defaultFontSize)
}

func (d *Doc) SetFontStyle(style string) {
	if err := d.Pdf.SetFont("lato", style, d.fontSize); err != nil {
		logrus.Fatal("error while changing Pdf font style: ", err)
	}
	d.fontStyle = style
}

func (d *Doc) DefaultFontStyle() {
	d.SetFontStyle("")
}

func (d *Doc) AddText(x, y float64, content string) {
	d.SetPosition(x, y)
	if err := d.Pdf.Cell(nil, content); err != nil {
		logrus.Fatal("error adding text to Pdf: ", err)
	}
}

func (d *Doc) AddFormattedText(x, y float64, content string, size int, style string) {
	d.SetFontSize(size)
	d.SetFontStyle(style)
	d.AddText(x, y, content)
	d.DefaultFontSize()
	d.DefaultFontStyle()
}

func (d *Doc) AddMultilineText(x, y float64, content string) {
	data := strings.Split(content, "\n")
	for i := range data {
		d.AddText(x, y, data[i])
		y += d.LineHeight()
	}
}

func (d *Doc) AddFormattedMultilineText(x, y float64, content string, size int, style string) {
	d.SetFontSize(size)
	d.SetFontStyle(style)
	data := strings.Split(content, "\n")
	for i := range data {
		d.AddText(x, y, data[i])
		y += d.LineHeight()
	}
	d.DefaultFontSize()
	d.DefaultFontStyle()
}

func (d Doc) LineHeight() float64 {
	return d.capValue * float64(d.fontSize) / 2000.0 * d.lineSpread
}

func (d Doc) textLineWidth() float64 {
	return gopdf.PageSizeA4.W - d.Pdf.MarginLeft() - d.Pdf.MarginRight()
}
