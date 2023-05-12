package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type TourOperator struct {
	Name          string `json:"name"`
	LicenceNumber string `json:"licence_Number"`
	Email         string `json:"email"`
	Website       string `json:"website"`
}

var stateFileName string = "state.json"
var outputFileName string = "output"

func operatorExists(items []TourOperator, licenceNumber string) bool {
	for _, item := range items {
		if item.LicenceNumber == licenceNumber {
			return true
		}
	}

	return false
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}

	return false
}

func getPage(URL string) int {
	if strings.Contains(URL, "=") {
		parts := strings.Split(URL, "=")
		page, err := strconv.Atoi(parts[1])

		if err != nil {
			return 0
		}

		return page
	} else {
		return 0
	}
}

type State struct {
	FailedAt string `json:"failed_at"`
	LastPage int    `json:"last_page"`
}

func (state State) save() {
	jsonData, _ := json.Marshal(state)
	jsonFile, _ := os.Create(stateFileName)
	defer jsonFile.Close()
	jsonFile.WriteString(string(jsonData))
}

func (state *State) load() {
	bytes, _ := os.ReadFile(stateFileName)

	if len(bytes) > 0 {
		json.Unmarshal(bytes, &state)
	}
}

func save(tourOperators []TourOperator) {
	if len(tourOperators) > 0 {
		jsonData, err := json.Marshal(tourOperators)

		if err != nil {
			panic(err)
		}

		jsonFile, err := os.Create(outputFileName + ".json")

		if err != nil {
			panic(err)
		}

		defer jsonFile.Close()

		if _, err = jsonFile.WriteString(string(jsonData)); err != nil {
			panic(err)
		}

		csvFile, err := os.Create(outputFileName + ".csv")

		if err != nil {
			panic(err)
		}

		defer csvFile.Close()

		writer := csv.NewWriter(csvFile)

		headers := []string{"Name", "Licence Number", "Email", "Website"}
		writer.Write(headers)

		for _, p := range tourOperators {
			record := []string{p.Name, p.LicenceNumber, p.Email, p.Website}
			writer.Write(record)
		}

		writer.Flush()
	}
}

func loadData() []TourOperator {
	var tourOperators []TourOperator
	bytes, _ := os.ReadFile(outputFileName + ".json")

	if len(bytes) > 0 {
		json.Unmarshal(bytes, &tourOperators)
	}

	return tourOperators
}

func main() {
	var state = &State{
		FailedAt: "",
		LastPage: 0,
	}

	// load prev state
	state.load()

	// initializing the slice of structs to store the data to scrape
	tourOperators := loadData()

	// initializing the list of pages to scrape with an empty slice
	var pagesToScrape []string

	// the first pagination URL to scrape
	pageToScrape := "https://utb.go.ug/tour-operators"

	// initializing the list of pages discovered with a pageToScrape
	pagesDiscovered := []string{pageToScrape}

	// creating a new Colly instance
	c := colly.NewCollector()

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	// get number of pages
	c.OnHTML("li.pager-last", func(e *colly.HTMLElement) {
		lastPage := e.ChildAttr("a", "href")

		// Start from page 1
		first := 1

		// if failed at, start from where it failed.
		if state.FailedAt != "" {
			first = getPage(state.FailedAt) + 1
			state.FailedAt = ""
			state.save()
		}

		// last page
		last := getPage(lastPage)

		if last > 0 {
			// save last page
			if state.LastPage == 0 {
				state.LastPage = last
				state.save()
			}

			for i := first; i <= last; i++ {
				// discovering a new page
				newPaginationLink := fmt.Sprintf("%v?page=%v", pageToScrape, i)

				// if the page discovered is new
				if !contains(pagesToScrape, newPaginationLink) {
					// if the page discovered should be scraped
					if !contains(pagesDiscovered, newPaginationLink) {
						pagesToScrape = append(pagesToScrape, newPaginationLink)
					}

					pagesDiscovered = append(pagesDiscovered, newPaginationLink)
				}
			}
		}
	})

	// iterating over the list of HTML product elements
	c.OnHTML("li.views-row", func(e *colly.HTMLElement) {
		// initializing a new TourOperator instance
		tourOperator := TourOperator{}

		// scraping the data of interest
		tourOperator.Name = e.ChildText("h4")
		tourOperator.Website = e.ChildAttr("a", "href")
		tourOperator.Email = e.ChildText(".views-field-field-company-email")
		tourOperator.LicenceNumber = e.ChildText(".views-field-field-licence-number")

		// adding the TourOperator instance with scraped data to the list of tour operators
		if !operatorExists(tourOperators, tourOperator.LicenceNumber) {
			tourOperators = append(tourOperators, tourOperator)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("-----------------------------------------------------")
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		state.FailedAt = fmt.Sprintf("%v", r.Request.URL)
		// save state
		state.save()

		// save data
		save(tourOperators)

		log.Println("Something went wrong: ", err)

		// retry after 5 seconds
		time.Sleep(5 * time.Second)
		log.Println("Retrying...")
		r.Request.Retry()
	})

	c.OnScraped(func(r *colly.Response) {
		URL := r.Request.URL
		fmt.Println("Finished", URL)

		currentPage := getPage(fmt.Sprintf("%v", URL))

		// until there is still a page to scrape
		if len(pagesToScrape) != 0 && currentPage < state.LastPage {
			// getting the current page to scrape and removing it from the list
			pageToScrape = pagesToScrape[0]
			pagesToScrape = pagesToScrape[1:]

			// visiting a new page
			c.Visit(pageToScrape)
		} else {
			save(tourOperators)
		}
	})

	startWithPage := pageToScrape

	// or start from where script failed at.
	if state.FailedAt != "" {
		startWithPage = state.FailedAt
	}

	// visiting the first page
	c.Visit(startWithPage)
}
