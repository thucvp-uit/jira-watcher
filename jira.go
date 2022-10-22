package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"github.com/k3a/html2text"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var jiraURL = os.Getenv("J_JIRA_URL")

// check here for the correct value
// jira_url/rest/activity-stream/1.0/config
var excludeConfluence = os.Getenv("J_EXCLUDE_CONFLUENCE")
var token = os.Getenv("J_JIRA_TOKEN")
var users = os.Getenv("J_WATCH_USERS")

// invert time zone here
var timeZone = -7

func main() {
	//get input value from command line
	inUsers := flag.String("u", users, "list of username separated by a comma")
	isVerbose := flag.Bool("v", false, "is verbose mode")
	inDate := flag.String("d", time.Now().Format("02-01"), "date format DD-MM-YYYY or DD-MM - default is current date")
	flag.Parse()

	//prepare parameters
	formattedUser := strings.ReplaceAll(*inUsers, ",", " ")
	formattedUser = strings.ReplaceAll(formattedUser, " ", "+")
	user := fmt.Sprintf("streams=user+IS+%v", formattedUser)
	currentYear := time.Now().Year()
	if len(*inDate) < 10 {
		*inDate = fmt.Sprintf("%v-%v", *inDate, currentYear)
	}
	timeFrom, _ := time.Parse("02-01-2006", *inDate)
	timeFrom = timeFrom.Add(time.Hour * time.Duration(timeZone))
	timeTo := timeFrom.Add(time.Hour * 24)
	dateRange := fmt.Sprintf("streams=update-date+BETWEEN+%v+%v", timeFrom.UnixMilli(), timeTo.UnixMilli())
	maxResult := fmt.Sprintf("maxResults=%v", 1000)
	issueComment := "issues=activity+IS+comment:post"
	url := fmt.Sprintf("%v/activity?%v&%v&%v&%v&%v", jiraURL, user, dateRange, maxResult, excludeConfluence, issueComment)

	fmt.Println(url)
	//fmt.Printf("Username: %v active from %v to %v\n", *inUsers, timeFrom.Format("02-01-2006"), timeTo.Format("02-01-2006"))

	//validate data
	if err := validateData(*inUsers); err != nil {
		log.Fatalln(err)
	}

	//make request data
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	var feeds Feed
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	err = xml.Unmarshal(body, &feeds)
	if err != nil {
		log.Fatalln(err)
	}

	var groupedEntries = make(map[string][]Entry)

	for _, entry := range feeds.Entries {
		author := entry.Author
		key := author.UserName
		groupedEntries[key] = append(groupedEntries[key], entry)
	}
	for _, key := range strings.Split(formattedUser, "+") {
		entries := groupedEntries[key]
		name := getUserName(key)
		fmt.Printf("[%v] make [%v] comments on %v\n", name, len(entries), *inDate)
		if *isVerbose {
			printActionDetail(entries)
		}
	}
}

func getUserName(username string) string {
	url := fmt.Sprintf("%v/rest/api/latest/user?username=%v", jiraURL, username)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	displayName := gjson.Get(string(body), "displayName")
	return displayName.String()
}

func validateData(userNames string) error {
	if len(userNames) == 0 {
		return errors.New("list username can't be empty")
	}

	if len(token) == 0 {
		return errors.New("jira token can't be empty")
	}

	if len(excludeConfluence) == 0 {
		return errors.New("exclude confluence can't be empty")
	}

	if len(jiraURL) == 0 {
		return errors.New("jira URL can't be empty")
	}

	return nil
}

func printActionDetail(entries []Entry) {
	for _, entry := range entries {
		fmt.Println("----------------------------------------------------------------")
		fmt.Println(html2text.HTML2Text(entry.Content))
	}
	fmt.Println("================================================================")
}
