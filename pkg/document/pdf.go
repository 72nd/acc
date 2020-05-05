package document

import (
	"bytes"
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/phpdave11/gofpdi"
	"github.com/signintech/gopdf"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

type PdfVersion int

const (
	v13 PdfVersion = iota
	v14
	v15
	v16
	v17
)

type Pdf struct {
	SrcPath             string
	DstPath             string
	useConvertedVersion bool
	tmpPath             string
	pdf                 gopdf.GoPdf
}

func NewPdf(srcPath, dstPath string) *Pdf {
	return &Pdf{
		SrcPath: srcPath,
		DstPath: dstPath,
	}
}

func (p *Pdf) Generate(props Properties) {
	if p.getPdfVersion() > v15 {
		p.useConvertedVersion = true
		p.tmpPath = p.downConvert()
	}
	p.initPdf()
	p.processPdf(props)
	p.safeAndCleanup(props.DstName)
}

func (p Pdf) getPdfVersion() PdfVersion {
	out, err := exec.Command("pdfinfo", p.SrcPath).Output()
	if err != nil {
		logrus.Fatalf("running pdfinfo for file «%s»: %s", p.SrcPath, err)
	}
	result := regexp.MustCompile(`PDF version:\s*(.*)`).FindStringSubmatch(string(out))
	if len(result) != 2 {
		logrus.Fatal("version matching failed: ", result)
	}
	switch result[1] {
	case "1.3":
		return v13
	case "1.4":
		return v14
	case "1.5":
		return v15
	case "1.6":
		return v16
	case "1.7":
		return v17
	default:
		logrus.Fatal("version matching failed, given value was ", result[1])
	}
	return v13
}

func (p *Pdf) downConvert() string {
	tmpPs, err := ioutil.TempFile("", "ps.*.ps")
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpPs.Name()) }()
	err = exec.Command("pdftops", p.SrcPath, tmpPs.Name()).Run()
	if err != nil {
		logrus.Fatalf("error running pdftops with %s %s: %s", p.SrcPath, tmpPs.Name(), err)
	}

	tmpPdf, err := ioutil.TempFile("", "pdf.*.pdf")
	if err != nil {
		logrus.Fatal("failed to create tmp file: ",  err)
	}
	err = exec.Command("gs",
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		fmt.Sprintf("-sOutputFile=%s", tmpPdf.Name()),
		tmpPs.Name()).Run()
	if err != nil {
		logrus.Fatal(err)
	}
	return tmpPdf.Name()
}

func (p Pdf) countPages() int {
	imp := gofpdi.NewImporter()
	imp.SetSourceFile(p.getSrcPath())
	return len(imp.GetPageSizes())
}

func (p Pdf) getSrcPath() string {
	if p.useConvertedVersion {
		return p.tmpPath
	}
	return p.SrcPath
}

func (p *Pdf) initPdf() {
	box, err := rice.FindBox("fonts")
	if err != nil {
		logrus.Error("rice find box failed: ", err)
	}
	latoHeavy, err := box.Bytes("Lato-Heavy.ttf")
	if err != nil {
		logrus.Errorf("could not load lato heavy: ", err)
	}
	latoRegular, err := box.Bytes("Lato-Regular.ttf")
	if err != nil {
		logrus.Errorf("could not load lato regular: ", err)
	}
	p.pdf = gopdf.GoPdf{}
	p.pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 595.28, H: 841.89}})
	if err := p.pdf.AddTTFFontByReaderWithOption("lato", bytes.NewBuffer(latoHeavy), gopdf.TtfOption{Style: gopdf.Bold}); err != nil {
		logrus.Fatal(err)
	}
	if err := p.pdf.AddTTFFontByReaderWithOption("lato", bytes.NewBuffer(latoRegular), gopdf.TtfOption{Style: gopdf.Regular}); err != nil {
		logrus.Fatal(err)
	}
}

func (p *Pdf) processPdf(props Properties) {
	cnt := p.countPages()
	p.processFirstPage(props, cnt)
	if cnt > 1 {
		for i := 2; i <= cnt; i++ {
			p.processOtherPage(props, i, cnt)
		}
	}
}

func (p *Pdf) processFirstPage(props Properties, maxPageNr int) {
	p.pdf.AddPage()
	p.pdf.SetLineWidth(0.1)
	p.pdf.SetFillColor(255, 255, 255)
	p.pdf.RectFromUpperLeftWithStyle(40, 130, 500, 680, "FD")
	p.pdf.SetFillColor(0, 0, 0)

	titleLine := fmt.Sprintf("%s: %s  //  page 1 of %d", props.Type, props.Identifier, maxPageNr)
	addText(&p.pdf, titleLine, 40, 40, 20, "Bold")
	addText(&p.pdf, props.Line1, 40, 63, 12, "")
	addText(&p.pdf, props.Line2, 40, 79, 12, "")
	addText(&p.pdf, props.Line3, 40, 95, 12, "")
	addText(&p.pdf, props.Line4, 40, 111, 12, "")

	tpl := p.pdf.ImportPage(p.getSrcPath(), 1, "/MediaBox")
	p.pdf.UseImportedTemplate(tpl, 45, 150, 480, 0)
}

func (p *Pdf) processOtherPage(props Properties, pageNr, maxPageNr int) {
	p.pdf.AddPage()
	p.pdf.SetFillColor(255, 255, 255)
	p.pdf.RectFromUpperLeftWithStyle(40, 100, 500, 700, "FD")
	p.pdf.SetFillColor(0, 0, 0)

	titleLine := fmt.Sprintf("Document: %s  //  page %d of %d", props.Identifier, pageNr, maxPageNr)
	addText(&p.pdf, titleLine, 40, 40, 20, "Bold")

	tpl := p.pdf.ImportPage(p.getSrcPath(), pageNr, "/MediaBox")
	p.pdf.UseImportedTemplate(tpl, 40, 100, 500, 0)
}

func (p *Pdf) safeAndCleanup(fileName string) {
	if p.DstPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			logrus.Fatal(err)
		}
		p.DstPath = wd
	}
	if err := p.pdf.WritePdf(p.DstPath); err != nil {
		logrus.Fatal("error while writing pdf: ", err)
	}
}

func addText(pdf *gopdf.GoPdf, content string, x, y float64, size int, style string) {
	if err := pdf.SetFont("lato", style, size); err != nil {
		logrus.Fatal(err)
	}
	pdf.SetX(x)
	pdf.SetY(y)
	if err := pdf.Cell(nil, content); err != nil {
		logrus.Fatal(err)
	}
}
