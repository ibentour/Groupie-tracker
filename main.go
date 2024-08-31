package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Data structures for JSON parsing
type (
	// Data holds the URLs for artist, location, date, and relation data
	Data struct {
		Artists   string `json:"artists"`
		Locations string `json:"locations"`
		Dates     string `json:"dates"`
		Relation  string `json:"relation"`
	}

	// Artist represents an individual artist's information
	Artist struct {
		ID           int      `json:"id"`
		Image        string   `json:"image"`
		Name         string   `json:"name"`
		Members      []string `json:"members"`
		CreationDate int      `json:"creationDate"`
		FirstAlbum   string   `json:"firstAlbum"`
		Locations    string   `json:"locations"`
		ConcertDates string   `json:"concertDates"`
		Relations    string   `json:"relations"`
	}

	// Locations holds location data for all artists
	Locations struct {
		Index []struct {
			ID        int      `json:"id"`
			Locations []string `json:"locations"`
			Dates     string   `json:"dates"`
		} `json:"index"`
	}

	// Dates holds concert date data for all artists
	Dates struct {
		Index []struct {
			ID    int      `json:"id"`
			Dates []string `json:"dates"`
		} `json:"index"`
	}

	// Relation holds relation data (dates and locations) for all artists
	Relation struct {
		Index []struct {
			ID             int                 `json:"id"`
			DatesLocations map[string][]string `json:"datesLocations"`
		} `json:"index"`
	}
)

// Global variables
var (
	artistData []map[string]interface{} // Holds processed artist data
	templates  *template.Template       // Holds parsed HTML templates
)

// fetchData retrieves data from a given URL and decodes it into the target structure.
func fetchData(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

// gatherData fetches and organizes data into a slice of maps.
func gatherData() error {
	var (
		data      Data
		artists   []Artist
		locations Locations
		relation  Relation
		dates     Dates
	)

	// Fetch the main API data
	if err := fetchData("https://groupietrackers.herokuapp.com/api", &data); err != nil {
		return err
	}

	// Use a WaitGroup to fetch data concurrently
	var wg sync.WaitGroup
	var errs []error
	wg.Add(4)

	// Fetch artists data
	go func() {
		defer wg.Done()
		if err := fetchData(data.Artists, &artists); err != nil {
			errs = append(errs, err)
		}
	}()

	// Fetch locations data
	go func() {
		defer wg.Done()
		if err := fetchData(data.Locations, &locations); err != nil {
			errs = append(errs, err)
		}
	}()

	// Fetch relation data
	go func() {
		defer wg.Done()
		if err := fetchData(data.Relation, &relation); err != nil {
			errs = append(errs, err)
		}
	}()

	// Fetch dates data
	go func() {
		defer wg.Done()
		if err := fetchData(data.Dates, &dates); err != nil {
			errs = append(errs, err)
		}
	}()

	wg.Wait()

	// Check if any errors occurred during fetching
	if len(errs) > 0 {
		return errs[0] // Return the first error encountered
	}

	// Process and combine all data into artistData
	for i := 0; i < len(artists); i++ {
		artistData = append(artistData, map[string]interface{}{
			"Id":           artists[i].ID,
			"Name":         artists[i].Name,
			"Image":        artists[i].Image,
			"Members":      artists[i].Members,
			"CreationDate": artists[i].CreationDate,
			"FirstAlbum":   artists[i].FirstAlbum,
			"Locations":    locations.Index[i].Locations,
			"Dates":        dates.Index[i].Dates,
			"Relation":     relation.Index[i].DatesLocations,
		})
	}

	return nil
}

// indexHandler handles HTTP requests for the home page ("/").
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ensure we're on the root path
	if r.URL.Path != "/" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	// Execute the index template with the artist data
	if err := templates.ExecuteTemplate(w, "index.html", artistData); err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
	}
}

// singleArtistHandler handles HTTP requests for individual artist details.
func singleArtistHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse and validate the artist ID from the URL query
	artistID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || artistID < 1 || artistID > len(artistData) {
		http.Error(w, "Invalid Artist ID", http.StatusBadRequest)
		return
	}

	// Execute the artist template with the specific artist's data
	if err := templates.ExecuteTemplate(w, "artist.html", artistData[artistID-1]); err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
	}
}

// init is called before main() to set up the application
func init() {
	// Fetch all data at startup
	if err := gatherData(); err != nil {
		log.Fatalf("Failed to gather data: %v", err)
	}

	// Set up template functions
	funcMap := template.FuncMap{
		"toString": func(v interface{}) string {
			return fmt.Sprintf("%v", v)
		},
	}

	// Parse HTML templates
	var err error
	templates, err = template.New("").Funcs(funcMap).ParseGlob("templet/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
}

func main() {
	// Serve static files (CSS)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	// Set up route handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/Artist", singleArtistHandler)

	// Start the HTTP server
	fmt.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}