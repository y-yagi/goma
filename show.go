package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/now"
	"github.com/olekukonko/tablewriter"
)

var nothingMsg = "Tomato is nothing (=･x･=)\n"

func showTodayTomatoes() error {
	tomatoes, err := selectTomatos(now.BeginningOfDay(), now.EndOfDay())
	if err != nil {
		return err
	}

	if len(tomatoes) == 0 {
		fmt.Fprintf(os.Stdout, nothingMsg)
		return nil
	}

	w := os.Stdout
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"id", "Time", "Tag"})
	var values = []string{}

	for i, tomato := range tomatoes {
		values = append(values, strconv.Itoa(i+1))
		values = append(values, tomato.CreatedAt.Format("15:04"))
		values = append(values, tomato.Tag)
		table.Append(values)
		values = nil
	}

	table.Render()
	return nil

}

func showTomatoes(showRange string) error {
	var start time.Time
	var end time.Time

	if showRange == "today" {
		return showTodayTomatoes()
	}

	if showRange == "all" {
		start = time.Date(2000, 01, 01, 00, 00, 00, 0, time.Now().Location())
	} else if showRange == "week" {
		start = now.BeginningOfWeek()
	} else if showRange == "mohth" {
		start = now.BeginningOfMonth()
	}

	end = time.Now()

	tagSummaries, err := selectTagSummary(start, end)
	if err != nil {
		return err
	}

	if len(tagSummaries) == 0 {
		fmt.Fprintf(os.Stdout, nothingMsg)
		return nil
	}

	w := os.Stdout
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Count", "Tag"})
	var values = []string{}

	for _, tagSummary := range tagSummaries {
		values = append(values, strconv.Itoa(tagSummary.Count))
		values = append(values, tagSummary.Tag)
		table.Append(values)
		values = nil
	}

	table.Render()
	return nil
}
