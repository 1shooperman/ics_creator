package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/1shooperman/ics_creator/ics"
)

func main() {

	year := flag.String("year", "", "Year to be processed")
	inFile := flag.String("file", "", "File to be processed")
	label := flag.String("label", "", "Label for the calendar")

	flag.Parse()

	if year == nil || *year == "" {
		fmt.Println("Invalid year:", year)
		os.Exit(1)
	}

	if inFile == nil || *inFile == "" {
		fmt.Println("Invalid file:", inFile)
		os.Exit(1)
	}

	if label == nil || *label == "" {
		fmt.Println("Invalid label:", label)
		os.Exit(1)
	}

	broker := ics.NewICSBroker()

	// Scrape the event list from the URL
	ev, err := broker.ScrapeEventList(*inFile)
	if err != nil {
		fmt.Println("Error scraping event list:", err)
		os.Exit(1)
	}

	var events []ics.Event
	for _, event := range ev {
		events = append(events, ics.Event{
			Summary:     *label,
			StartDate:   event,
			EndDate:     event,
			Description: *label,
		})
	}

	// Create an ICS file handle with the year in its name
	fileName := fmt.Sprintf("%s.ics", *year)
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer outFile.Close()

	// write the calendar file
	if err = broker.Write(outFile, events); err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}

	fmt.Println("ICS file created successfully:", fileName)
}
