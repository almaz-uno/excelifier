package main

import (
	"log"
	"os"
	"time"

	"github.com/tealeg/xlsx"
)

const timeLayout = "2006-01-02T15:04:05"

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
	//sheet.Col(1).Width = 50.0

	cell = row.AddCell()
	cell.Value = `Duration, days`
	cell.SetStyle(style)

	cell = row.AddCell()
	cell.Value = `Work time, days`
	cell.SetStyle(style)

	cell = row.AddCell()
	cell.Value = `Start`
	cell.SetStyle(style)

	cell = row.AddCell()
	cell.Value = `Finish`
	cell.SetStyle(style)

	cell = row.AddCell()
	cell.Value = `Cost`
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
		row.AddCell().SetFloat(t.DurationHours() / 8.0)

		cell.SetStyle(style)
		row.AddCell().SetFloat(t.WorkHours() / 8.0)

		if t, err := time.Parse(timeLayout, t.Start); err == nil {
			row.AddCell().SetDate(t)
		}

		if t, err := time.Parse(timeLayout, t.Finish); err == nil {
			row.AddCell().SetDate(t)
		}

		row.AddCell().SetFloat(t.Cost / 100.0)

	}

	err = xFile.Save(dstFile)
	if err == nil {
		log.Printf("%d rows were wrote in file %s", len(tt), dstFile)
	}

	return err
}
