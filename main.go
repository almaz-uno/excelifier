package main

import (
	"log"
	"os"

	"github.com/tealeg/xlsx"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s %s %s", os.Args[0], "srcFile.xml", "dstFile.xlsx")
	}

	err := doExport(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Done.")

}

func doExport(srcFile, dstFile string) error {
	inFile, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer inFile.Close() // nolint: errcheck

	tt, err := extractTasks(inFile)
	if err != nil {
		return err
	}

	xFile := xlsx.NewFile()
	sheet, err := xFile.AddSheet("Tasks")
	if err != nil {
		return err
	}

	var style = xlsx.NewStyle()
	style.Font.Bold = true
	style.ApplyFont = true

	var cell *xlsx.Cell

	row := sheet.AddRow()
	cell = row.AddCell()
	cell.Value = `WBS`
	cell.SetStyle(style)

	cell = row.AddCell()
	cell.Value = `Name`
	cell.SetStyle(style)
	sheet.Col(1).Width = 50.0

	cell = row.AddCell()
	cell.Value = `Duration, days`
	cell.SetStyle(style)

	for _, t := range tt {

		row := sheet.AddRow()
		row.OutlineLevel = uint8(t.OutlineLevel)

		row.AddCell().SetString(t.WBS)

		cell = row.AddCell()
		cell.SetString(t.Name)

		style := xlsx.NewStyle()
		style.Alignment.Indent = t.OutlineLevel
		style.ApplyAlignment = true

		cell.SetStyle(style)
		row.AddCell().SetFloat(t.Hours() / 8.0)

	}

	err = xFile.Save(dstFile)
	if err == nil {
		log.Printf("%d rows were wrote in file %s", len(tt), dstFile)
	}

	return err
}
