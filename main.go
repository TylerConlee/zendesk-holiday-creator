package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	// Establish variables to store data from flags in.
	var (
		schedule  int
		subdomain string
		user      string
		key       string
	)

	// Open a locally sourced file with holidays listed in it
	// This file should be comma separated with one holiday per line
	file, err := os.Open("holidays.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Define what each of the flags do
	// -id: The schedule ID in Zendesk, can be found in the URL:
	// https://<subdomain>.zendesk.com/agent/admin/schedules/397207/holidays
	// where 397207 is the schedule ID
	flag.IntVar(&schedule, "id", 0, "ID for the Zendesk schedule")

	// -user: The username that will access the Zendesk API. Can be used in the
	// form user@company.com/token
	flag.StringVar(&user, "user", "test@circleci.com/token", "Zendesk username for API")

	// -url: The subdomain for your Zendesk instance:
	// https://<subdomain>.zendesk.com
	flag.StringVar(&subdomain, "url", "circleci", "Zendesk subdomain (i.e. <companyname>.zendesk.com)")

	// -key: The API key that will be used to interact with Zendesk
	// API keys can be generated at https://<subdomain>.zendesk.com/agent/admin/api/settings
	flag.StringVar(&key, "key", "XXXXX", "Zendesk API key")

	// Parse the inputted flags into usable variables
	flag.Parse()

	// Using the schedule ID and the subdomain, compile a URL to send the API
	//request to
	url := fmt.Sprintf("https://%s.zendesk.com/api/v2/business_hours/schedules/%d/holidays.json", subdomain, schedule)

	// Create a new CSV parser using the contents of the local file imported
	// above
	r := csv.NewReader(file)

	// Loop through each line of the file and call the sendHoliday function
	for {
		// Read the line into the variable "record"
		record, err := r.Read()
		// Stop looping if this is the end of the file
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// Send the line, the API url and the username and API key to the
		// SendHoliday function. Record represents each component of the comma
		// separated text file, where record[0] is the name of the holiday and
		// record[1] is the date
		sendHoliday(record, url, user, key)
	}
}

// SendHoliday takes a slice of strings representing a given holiday, and,
// using the username and API key provided, sends it through an API request to
// Zendesk
func sendHoliday(holiday []string, url string, user string, key string) {

	// Create a JSON string formatted in a way that Zendesk expects based on
	// https://developer.zendesk.com/rest_api/docs/support/schedules#create-a-holiday
	payload := fmt.Sprintf(`{"holiday": {"name": "%s", "start_date": "%s", "end_date": "%s"}}`, holiday[0], holiday[1], holiday[1])

	// Output this string to the console for clarity
	fmt.Println(payload)

	// Turn the string into a []byte to allow for it to be used in the HTTP
	// request
	jsonStr := []byte(payload)

	// Set up a new post request to the Zendesk URL created above, and include
	// the formatted JSON as the body of the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("Error creating new HTTP request")
	}

	// Log the request to the console
	fmt.Println(req)

	// Include the username and API key in the request
	req.SetBasicAuth(user, key)

	// Set the content type to JSON to be accepted by Zendesk
	req.Header.Set("Content-Type", "application/json")

	// create custom http.Client to manually set timeout and disable keepalive
	// in an attempt to avoid EOF errors
	var netClient = &http.Client{
		Timeout: time.Second * 480,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}

	// Execute the request and capture the response in "resp"
	resp, err := netClient.Do(req)

	if err != nil {
		fmt.Printf("Error creating new HTTP request: %s", err.Error())

	}
	defer resp.Body.Close()

	// Output the status of the completed request
	fmt.Printf("%s", resp.Status)
}
