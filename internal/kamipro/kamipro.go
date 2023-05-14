package kamipro

import (
	"fmt"
	"path/filepath"

	"github.com/Angelmaneuver/xlsx2html/internal/kamipro/application"
	"github.com/Angelmaneuver/xlsx2html/internal/kamipro/generate"
	"github.com/xuri/excelize/v2"
	"golang.org/x/sync/errgroup"
)

func Start(input string, output string) error {
	if len(output) == 0 {
		output = filepath.Dir(input)
	}

	application, err := application.New()
	if err != nil {
		return err
	}

	f, err := excelize.OpenFile(input)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	eg := errgroup.Group{}

	for _, dataset := range application.Excel.Dataset {
		rows, err := f.GetRows(dataset.Sheet)
		if err != nil {
			return err
		}

		setting := generate.Setting{
			Rarity: dataset.Rarity,
			Icon:   dataset.Icon,
			Output: filepath.Join(output, dataset.Output),
		}
		rows = rows[application.Excel.Skip.Row:]

		eg.Go(func() error {
			return generate.Start(
				&setting,
				&application.Excel.Key,
				&application.Excel.Sort,
				&application.Html,
				&rows,
			)
		})
	}

	return eg.Wait()
}
