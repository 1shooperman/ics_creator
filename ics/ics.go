package ics

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ICSBroker struct {
	shouldAlarm bool
}

// Event represents a calendar event
type Event struct {
	Summary     string
	StartDate   time.Time
	EndDate     time.Time
	Description string
}

func NewICSBroker() *ICSBroker {
	return &ICSBroker{}
}

const format = "20060102"

func (w ICSBroker) Write(file *os.File, events []Event) error {
	// Write the header
	fmt.Fprintln(file, "BEGIN:VCALENDAR")
	fmt.Fprintln(file, "VERSION:2.0")
	fmt.Fprintln(file, "PRODID:-//BrandonShoop//ICSCreator//EN")

	// Write each event
	currentTime := time.Now().UTC()
	for _, event := range events {
		fmt.Fprintln(file, "BEGIN:VEVENT")
		fmt.Fprintf(file, "UID:%s\n", calUID())
		fmt.Fprintf(file, "DTSTAMP;VALUE=DATE:%s\n", currentTime.Format(format))
		fmt.Fprintf(file, "DTSTART;VALUE=DATE:%s\n", formatDateTime(event.StartDate))
		fmt.Fprintf(file, "DTEND;VALUE=DATE:%s\n", formatDateTime(event.EndDate))
		fmt.Fprintf(file, "SUMMARY:%s\n", event.Summary)
		fmt.Fprintf(file, "DESCRIPTION:%s\n", event.Description)

		if !alarm {
			fmt.Fprintln(file, "BEGIN:VALARM")
			fmt.Fprintf(file, "X-WR-ALARMUID:%s\n", calUID())
			fmt.Fprintln(file, "TRIGGER;VALUE=DATE-TIME:19760401T005545Z")
			fmt.Fprintln(file, "ACTION:NONE")
			fmt.Fprintln(file, "X-APPLE-DEFAULT-ALARM:TRUE")
			fmt.Fprintln(file, "END:VALARM")
		}

		fmt.Fprintln(file, "END:VEVENT")
	}

	// Write the footer
	fmt.Fprintln(file, "END:VCALENDAR")

	return nil
}

// formatDateTime formats the date and time in the specific format for ICS files
func formatDateTime(t time.Time) string {
	return t.Format(format)
}

func (w ICSBroker) ScrapeEventList(f string) ([]time.Time, error) {

	// Open the file
	file, err := os.Open(f)
	if err != nil {
		return []time.Time{}, err
	}
	defer file.Close()

	// Compile a regular expression to match the date format
	datePattern := regexp.MustCompile(`\b[A-Z][a-z]+, [A-Z][a-z]+ \d{1,2}, \d{4}\b`)

	// Define the date layout (reference time format)
	const layout = "Monday, January 2, 2006"

	// Create a slice to hold the parsed dates
	var dates []time.Time

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Find all date strings in the line
		dateStrings := datePattern.FindAllString(line, -1)

		// Parse each date string and add to the slice
		for _, dateString := range dateStrings {
			date, err := time.Parse(layout, dateString)
			if err != nil {
				return []time.Time{}, err
			} else {
				dates = append(dates, date)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return []time.Time{}, err
	}

	return dates, nil
}

func calUID() string {
	uid := uuid.New().String()[:14]
	return strings.ReplaceAll(uid, "-", "0")
}
