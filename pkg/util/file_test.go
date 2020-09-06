package util

import (
	"os"
	"testing"
)

func TestRelativeAssetPath(t *testing.T) {
	paths := [][3]string{
		{
			"/tmp/acc/projects/bernhard-beratung-realisierung/scrapers",
			"/tmp/acc/projects/bernhard-beratung-realisierung/scrapers/invoice_01.pdf",
			"invoice_01.pdf",
		},
		{
			"/tmp/acc/projects/bernhard-beratung-realisierung/ausbildungsnachweise",
			"acc/projects/bernhard-beratung-realisierung/ausbildungsnachweise/invoice_01.pdf",
			"invoice_01.pdf",
		},
		{
			"/tmp/acc/projects/politforum-kaefigturm/staat-religion/",
			"/tmp/acc/projects/politforum-kaefigturm/staat-religion/invoice1.pdf",
			"invoice1.pdf",
		},
	}

	crtWd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	if err := os.Chdir("/tmp/"); err != nil {
		t.Error(err)
	}

	for i := range paths {
		rsl := RelativeAssetPath(paths[i][0], paths[i][1])
		if rsl != paths[i][2] {
			t.Errorf("got \"%s\" but expected \"%s\"", rsl, paths[i][2])
		}
	}

	if err := os.Chdir(crtWd); err != nil {
		t.Error(err)
	}
}
