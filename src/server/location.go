package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type latLong struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"lon"`
}

type locationPcoordinates struct {
	Name        string  `json:"name"`
	Coordinates latLong `json:"coordinates"`
}

type Response struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}

var cache sync.Map

//Checking cached data to determine if we can populate the coordinates
/*
potential problem is that if alot of locations are uncached then the geoLocation request timing difference
when cache misses can slow down the requests for multiple requests at the same time
since we won't have alot of requests tho it's fine. -> fix undetermined :d
*/
func geoLocation(location string) (latLong, error) {
	if cached, good := cache.Load(location); good { //If good is true for cached data
		return cached.(latLong), nil //if this sends anything but coord data we're ducked :)
	}

	//dividing at "," then cutting spaces and returning the search if 2 parts exist
	parts := strings.Split(location, ",")
	search := strings.TrimSpace(location)
	if len(parts) == 2 {
		city := strings.TrimSpace(parts[0])
		country := strings.TrimSpace(parts[1])
		search = city + ", " + country
	}

	//base endpoint + map for query string building and encoding into url
	baseUrl := "https://nominatim.openstreetmap.org/search"
	parameters := url.Values{}
	parameters.Add("q", search) //q is nominatim base search parameter delimited by q="search/query"
	parameters.Add("format", "json")
	parameters.Add("limit", "1") //getting the first response
	baseRequest := baseUrl + "?" + parameters.Encode()

	//net/http import base request handler for request timeouts and server hangs -> if not server degrades
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	//Making GET request with baseRequest param
	req, err := http.NewRequest("GET", baseRequest, nil)
	if err != nil { //if any error
		return latLong{}, fmt.Errorf("There was an errorW: %v", err) //just error message
	}
	req.Header.Set("User-Agent", "GroupTracker") //if baseRequest isn't errored then pass to necessary User-Agent(nominatim needs this) declaration

	//Basically works with http.Client to get request response
	resp, err := client.Do(req)
	if err != nil {
		return latLong{}, fmt.Errorf("There was an error %v", err) //error handling also returns 0 value for latLong
	}
	defer resp.Body.Close() //closes the http response body to avoid eating resources

	//reads and stocks http data into "body" and turns it indo bytes so we can use it later
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return latLong{}, fmt.Errorf("There was an error %v", err) //error handling !
	}

	var results []Response               //declaring Reponse slice
	err = json.Unmarshal(body, &results) //Fills Response slice with JSON data from http response inside body where JSON data is part of ([]bytes)
	if err != nil {
		return latLong{}, fmt.Errorf("There was an error %v", err) //say it with me : error handling ! :_)
	}
	if len(results) == 0 {
		return latLong{}, fmt.Errorf("No results found %s", location) //if nothing is found well return value 0 from latLong
	}

	//Converting result String to float64 because latLong struct is float64
	var convCords latLong //Declaring latLong variable
	/*
		Sscan is used to transform string to float and & delimiter is used to change the active parameter
		instead of keeping it hardcoded from result
	*/
	_, err = fmt.Sscanf(results[0].Latitude, "%f", &convCords.Lat) //Takes parsed latitude and puts it in convCords.Lat
	if err != nil {
		return latLong{}, fmt.Errorf("Latitude error %v", err)
	}

	//Same as for Latitude
	_, err = fmt.Sscanf(results[0].Longitude, "%f", &convCords.Long)
	if err != nil {
		return latLong{}, fmt.Errorf("Longitude error %v", err)
	}

	cache.Store(location, convCords) //Just stores results into cache to not re-request from API (we do be stealin them duckies)
	time.Sleep(1 * time.Second)      //avoid rate limits from nominatim by limiting request return times

	return convCords, nil //returns the lat and long values :d
}

//Return a string slice map of coordinates to convert location names to map coordinates
/*
The function doesn't really behave as a setter it just gives a map to use foro name searching a location by key
the function name can be confusing :)
*/
func setLocations(markers []string) map[string]latLong {
	result := make(map[string]latLong) //defining the string slice map respecting the latLong structure
	for _, pin := range markers {      //Going through markers -> pin makes up the loaction name
		coordinates, err := geoLocation(pin) //determines coordinate data
		if err != nil {
			fmt.Printf("Cannot get the code location '%s': %v\n", pin, err) //ERROR HANDLINGGG!!!!
			coordinates = latLong{Lat: 0, Long: 0}                          //Fallback to 0 if error (i hope nobody is having a concert in a gulf)
		}
		result[pin] = coordinates //stocks coordinate data with location names in a string slice
	}
	return result //returns the result list
}

//Returning the location names + coordinates via locationPcoordinates structure
/*
As pointed out for setLocations the getLocations doesn't really work with a setter
it's just a slice representation of the same result as the setLocations but it's easier to use
if we are transfering to frontend because the slice is structured as a JSON array = easier/ordered
*/
func getLocations(markers []string) []locationPcoordinates {
	//just a slice with struct, starts empty (0) and determines size as length of markers
	result := make([]locationPcoordinates, 0, len(markers))

	//Checking for location coordinates
	for _, pin := range markers {
		coordinates, err := geoLocation(pin)
		if err != nil {
			fmt.Printf("Cannot get code location '%s': %v\n", pin, err)
			coordinates = latLong{Lat: 0, Long: 0}
		}
		//Adding to the result same as the setLocations but with a different struct format
		result = append(result, locationPcoordinates{
			Name:        pin,
			Coordinates: coordinates,
		})
	}
	return result //Returning result
}
