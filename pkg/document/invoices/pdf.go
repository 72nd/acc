package invoices

import (
	"bytes"
	rice "github.com/GeertJohan/go.rice"
	"github.com/signintech/gopdf"
	"github.com/sirupsen/logrus"
	"gitlab.com/72th/acc/pkg/schema"
	"os"
)

func initPdf() gopdf.GoPdf {
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
	pdf := gopdf.GoPdf{}
	if err := pdf.AddTTFFontByReaderWithOption("lato", bytes.NewBuffer(latoHeavy), gopdf.TtfOption{Style: gopdf.Bold}); err != nil {
		logrus.Fatal("error adding lato heavy to pdf: ", err)
	}
	if err := pdf.AddTTFFontByReaderWithOption("lato", bytes.NewBuffer(latoRegular), gopdf.TtfOption{Style: gopdf.Regular}); err != nil {
		logrus.Fatal("error adding lato regular to pdf: ", err)
	}
	return pdf
}

func page(pdf gopdf.GoPdf, inv schema.Invoice) gopdf.GoPdf {
	pdf.AddPage()
	pdf.SetLineWidth(0.1)
	pdf.SetFillColor(255, 255, 255)
	pdf.RectFromUpperLeftWithStyle(40, 130, 500, 680, "FD")
	pdf.SetFillColor(0, 0, 0)
	return pdf
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
		logrus.Fatal("error while writing pdf: ", err)
	}
}
