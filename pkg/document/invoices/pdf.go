package invoices

import (
	"fmt"
	"github.com/signintech/gopdf"
	"github.com/sirupsen/logrus"
	"github.com/72nd/acc/pkg/document"
	"github.com/72nd/acc/pkg/schema"
	"math"
	"os"
	"time"
)

// InvoiceDocument is a Doc report which generates a invoice-letter.
type InvoiceDocument struct {
	document.Doc
	place string
}

// NewInvoiceDocument returns a new InvoiceDocument.
func NewInvoiceDocument(fontSize int, place string) InvoiceDocument {
	return InvoiceDocument{
		Doc:   document.NewDoc(fontSize, 1.2),
		place: place,
	}
}

// Generate generates a PDF for a given InvoiceDocument and returns it as a gopdf.GoPdf element.
func (d *InvoiceDocument) Generate(company schema.Company, invoice schema.Invoice, customer schema.Party) gopdf.GoPdf {
	d.Doc.Pdf.AddPage()
	d.Doc.Pdf.SetLineWidth(0.1)
	d.Doc.Pdf.SetMargins(20, 10, 20, 10)
	d.Doc.Pdf.SetFillColor(0, 0, 0)
	d.header(company)
	d.address(company, customer)

	return d.Doc.Pdf
}

func (d *InvoiceDocument) header(company schema.Company) {
	d.Doc.AddFormattedMultilineText(20, 20, fmt.Sprintf(
		"%s\n%s %d\n%d %s\nTelefon: %s\nE-Mail: %s\nURL: %s",
		company.Name,
		company.Street,
		company.StreetNr,
		company.PostalCode,
		company.Place,
		company.Phone,
		company.Mail,
		company.Url,
	), 10, "")
}

func (d *InvoiceDocument) address(company schema.Company, customer schema.Party) {
	sender := fmt.Sprintf(
		"%s, %s %d, %d %s",
		company.Name,
		company.Street,
		company.StreetNr,
		company.PostalCode,
		company.Place,
	)
	d.Doc.AddFormattedText(115, 50, sender, 7, "")
	y := 50 + math.Round(d.Doc.LineHeight()/1.3)
	d.Pdf.Line(115, y, 190, y)
	d.Doc.AddFormattedMultilineText(115, y+2*d.Doc.LineHeight(), customer.AddressLines(), 10, "")

	placeDate := fmt.Sprintf("%s, %s", d.place, time.Now().Format("02.01.2006"))
	placeDateWidth, _ := d.Pdf.MeasureTextWidth(placeDate)
	placeDateX := gopdf.PointsToUnits(gopdf.Unit_MM, gopdf.PageSizeA4.W) - d.Pdf.MarginRight() - placeDateWidth
	d.Doc.AddText(placeDateX, 100, placeDate)
}

func save(pdf gopdf.GoPdf, dstPath string) {
	if dstPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			logrus.Fatal(err)
		}
		dstPath = wd
	}
	if err := pdf.WritePdf(dstPath); err != nil {
		logrus.Fatal("error while writing utils: ", err)
	}
}
