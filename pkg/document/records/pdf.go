package records

import (
	"bytes"
	"fmt"
	"github.com/phpdave11/gofpdi"
	"github.com/signintech/gopdf"
	"github.com/sirupsen/logrus"
	"github.com/72nd/acc/pkg/document/utils"
	"image"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
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

type ScrType int

const (
	PDF ScrType = iota
	PNG
)

type Pdf struct {
	SrcPath             string
	ScrType             ScrType
	DstPath             string
	useConvertedVersion bool
	tmpPath             string
	pdf                 gopdf.GoPdf
}

func NewPdf(srcPath, dstPath string) *Pdf {
	scrType := PDF
	if path.Ext(srcPath) == ".png" {
		scrType = PNG
	}
	return &Pdf{
		SrcPath: srcPath,
		ScrType: scrType,
		DstPath: dstPath,
	}
}

func (p *Pdf) Generate(props Properties, downConvert bool) {
	if p.ScrType == PDF && downConvert {
		p.useConvertedVersion = true
		p.tmpPath = p.downConvert()
	}
	p.initPdf()
	p.processPdf(props)
	p.safeAndCleanup()
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

	tmpPdf, err := ioutil.TempFile("", "utils.*.utils")
	if err != nil {
		logrus.Fatal("failed to create tmp file: ", err)
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
	if p.ScrType == PNG {
		return 1
	}
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
	p.pdf = gopdf.GoPdf{}
	p.pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	if err := p.pdf.AddTTFFontByReaderWithOption("lato", bytes.NewBuffer(utils.LatoHeavy()), gopdf.TtfOption{Style: gopdf.Bold}); err != nil {
		logrus.Fatal("error adding lato heavy to utils: ", err)
	}
	if err := p.pdf.AddTTFFontByReaderWithOption("lato", bytes.NewBuffer(utils.LatoRegular()), gopdf.TtfOption{Style: gopdf.Regular}); err != nil {
		logrus.Fatal("error adding lato regular to utils: ", err)
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
	p.pdf.SetLineWidth(0.5)
	p.pdf.SetFillColor(255, 255, 255)
	p.pdf.RectFromUpperLeftWithStyle(40, 146, 500, 680, "FD")
	p.pdf.SetFillColor(0, 0, 0)

	titleLine := fmt.Sprintf("%s: %s  //  page 1 of %d", props.Type, props.Identifier, maxPageNr)
	addText(&p.pdf, titleLine, 40, 40, 20, "Bold")
	addText(&p.pdf, props.Line1, 40, 63, 12, "")
	addText(&p.pdf, props.Line2, 40, 79, 12, "")
	addText(&p.pdf, props.Line3, 40, 95, 12, "")
	addText(&p.pdf, props.Line4, 40, 111, 12, "")
	addText(&p.pdf, props.Line5, 40, 127, 12, "")

	if p.ScrType == PDF {
		tpl := p.pdf.ImportPage(p.getSrcPath(), 1, "/MediaBox")
		p.pdf.UseImportedTemplate(tpl, 45, 150, 480, 0)
	} else {
		p.pdf.SetX(45)
		p.pdf.SetY(150)
		rect := fitImage(p.SrcPath, 480, 660)
		if err := p.pdf.Image(p.SrcPath, 45, 150, &rect); err != nil {
			logrus.Errorf("error while including image %s into pdf: %s", p.SrcPath, err)
		}
	}
}

func (p *Pdf) processOtherPage(props Properties, pageNr, maxPageNr int) {
	p.pdf.AddPage()
	p.pdf.SetLineWidth(0.5)
	p.pdf.SetFillColor(255, 255, 255)
	p.pdf.RectFromUpperLeftWithStyle(40, 100, 500, 700, "FD")
	p.pdf.SetFillColor(0, 0, 0)

	titleLine := fmt.Sprintf("%s: %s  //  page %d of %d", props.Type, props.Identifier, pageNr, maxPageNr)
	addText(&p.pdf, titleLine, 40, 40, 20, "Bold")

	tpl := p.pdf.ImportPage(p.getSrcPath(), pageNr, "/MediaBox")
	p.pdf.UseImportedTemplate(tpl, 45, 105, 480, 0)
}

func (p *Pdf) safeAndCleanup() {
	if p.DstPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			logrus.Fatal(err)
		}
		p.DstPath = wd
	}
	if err := p.pdf.WritePdf(p.DstPath); err != nil {
		logrus.Fatal("error while writing utils: ", err)
	}
}

func addText(pdf *gopdf.GoPdf, content string, x, y float64, size int, style string) {
	if err := pdf.SetFont("lato", style, size); err != nil {
		logrus.Fatal("error while changing Pdf font style: ", err)
	}
	pdf.SetX(x)
	pdf.SetY(y)
	if err := pdf.Cell(nil, content); err != nil {
		logrus.Fatal(err)
	}
}

func fitImage(path string, containerWidth, containerHeight int) gopdf.Rect {
	reader, err := os.Open(path)
	if err != nil {
		logrus.Fatalf("couldn't open image \"%s\": %s", path, err)
	}
	defer reader.Close()
	img, _, err := image.DecodeConfig(reader)
	if err != nil {
		logrus.Fatalf("error while reading image \"%s\": %s", path, err)
	}
	iWidth := float64(img.Width)
	iHeight := float64(img.Height)
	cWidth := float64(containerWidth)
	cHeight := float64(containerHeight)
	if cWidth/cHeight > iWidth/iHeight {
		return gopdf.Rect{W: iWidth * cHeight / iHeight, H: cHeight}
	}
	return gopdf.Rect{W: cWidth, H: iHeight * cWidth / iWidth}
}
