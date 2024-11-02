package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type errorData struct {
	Num  int
	Text string
}

type artistData struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	// Index []struct {
	// 	ID             int                 `json:"id"`
	// 	DatesLocations map[string][]string `json:"datesLocations"`
	// } `json:"index`
	Index []struct {
		ID             int                 "json:\"id\""
		DatesLocations map[string][]string "json:\"datesLocations\""
	}
}

// type Location struct {
// 	Index []struct {
// 		ID        int      `json:"id"`
// 		Locations []string `json:"locations"`
// 	} `json:"index"`
// }

// type Dates struct {
// 	Index []struct {
// 		ID    int      `json:"id"`
// 		Dates []string `json:"dates"`
// 	} `json:"index`
// }

type Relation struct {
	// Index []struct {
	// 	ID             int                 `json:"id"`
	// 	DatesLocations map[string][]string `json:"datesLocations"`
	// } `json:"index`
	Index []struct {
		ID             int                 "json:\"id\""
		DatesLocations map[string][]string "json:\"datesLocations\""
	}
}

var allData []artistData

// var allLocations []Location
// var allDates []Dates
var allRelations []Relation

func main() {
	fmt.Println("...Hos geldiniz ¯\\(ツ)/¯ ... Give me just a minute to gather data...")

	FileServer := http.FileServer(http.Dir("docs"))
	http.Handle("/docs/", http.StripPrefix("/docs/", FileServer))

	allData = gatherDataUp("https://groupietrackers.herokuapp.com/api/artists")
	// allLocations = gatherDataUp2("https://groupietrackers.herokuapp.com/api/locations")
	// allDates = gatherDataUp3("https://groupietrackers.herokuapp.com/api/dates")
	
	allRelations = gatherDataUp4("https://groupietrackers.herokuapp.com/api/relation")

	//for _, a := range allLocations {
	//fmt.Println(a)
	//}

	//|| allLocations == nil
	if allData == nil || allRelations == nil {
		fmt.Println("Failed to gather Data from API")
		os.Exit(1)
	}

	http.HandleFunc("/", mainPage)
	http.HandleFunc("/response", response)
	http.HandleFunc("/search", search)
	fmt.Println()
	fmt.Println("Thanks, man (ಥ﹏ಥ) Now Server is listening to port #8080   ᕦ(ò_óˇ)ᕤ")
	http.ListenAndServe(":8080", nil)
}

func mainPage(res http.ResponseWriter, req *http.Request) {
	temp, er := template.ParseFiles("docs/index.html")
	if er != nil {
		err(res, req, http.StatusInternalServerError)
		return
	}
	if req.URL.Path != "/" {
		err(res, req, http.StatusNotFound)
		return
	}
	temp.Execute(res, allData)
}

func response(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("burdaaaaa")
	temp, er := template.ParseFiles("docs/response.html")
	if er != nil {
		log.Fatal(er)
		err(res, req, http.StatusInternalServerError)
		return
	}
	ID1 := req.FormValue("ID")
	req.ParseForm()
	num, err := strconv.Atoi(ID1)
	if err != nil {
		fmt.Println("Dizeyi int'e dönüştürme hatasi:", err)
		return
	}
	// İlgili sanatçının verilerini bulma
	var selectedArtist artistData
	for _, v := range allData {
		if v.ID == num {
			//temp.Execute(res, v)
			selectedArtist = v
			//break
		}
	}

	// İlgili sanatçının lokasyon verilerini bulma ve HTML şablonuna eklemek
	// for _, a := range allLocations {
	// 	for _, loc := range a.Index {
	// 		if selectedArtist.ID == loc.ID {
	// 			selectedArtist.Index = append(selectedArtist.Index, a.Index...)

	// 			temp.Execute(res, selectedArtist)
	// 			//fmt.Fprintln(res, loc, i+1)

	// 		}
	// 	}
	// }

	for _, a := range allRelations {
		for _, loc := range a.Index {
			if selectedArtist.ID == loc.ID {
				selectedArtist.Index = append(selectedArtist.Index, a.Index...)
				temp.Execute(res, selectedArtist)
				//fmt.Fprintln(res, loc)
			}
		}
	}
}
func search(res http.ResponseWriter, req *http.Request) {
	temp, e1 := template.ParseFiles("docs/search.html")
	if e1 != nil {
		err(res, req, http.StatusInternalServerError)
		return
	}
	temp.Execute(res, allData)
}
func err(res http.ResponseWriter, req *http.Request, err int) {
	temp, er := template.ParseFiles("docs/error.html")
	if er != nil {
		log.Fatal(er)
		return
	}
	res.WriteHeader(err)
	errData := errorData{Num: err}
	if err == 404 {
		errData.Text = "Page Not Found"
	} else if err == 400 {
		errData.Text = "Bad Request"
	} else if err == 500 {
		errData.Text = "Internal Server Error"
	}
	fmt.Println(errData)
	temp.Execute(res, errData)
}

func gatherDataUp(link string) []artistData {
	client := &http.Client{
		Timeout: time.Second * 10, // Örnek olarak, 10 saniye olarak ayarlandı
	}
	data1, e1 := client.Get(link)
	//data1, e1 := http.Get(link)
	if e1 != nil {
		log.Fatal(e1)
		return nil
	}
	data2, e2 := ioutil.ReadAll(data1.Body)
	if e2 != nil {
		log.Fatal(e2)
		return nil
	}
	Artists := []artistData{}
	e := json.Unmarshal(data2, &Artists)
	if e != nil {
		log.Fatal(e)
		return nil
	}

	return Artists
}

// func gatherDataUp2(link string) []Location {
// 	data1, e1 := http.Get(link)
// 	if e1 != nil {
// 		log.Fatal(e1)
// 	}
// 	var index Location
// 	err := json.NewDecoder(data1.Body).Decode(&index)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	//fmt.Print(index)

// 	var location []Location
// 	for _, item := range index.Index {
// 		location = append(location, Location{Index: []struct {
// 			ID        int      `json:"id"`
// 			Locations []string `json:"locations"`
// 		}{item}})

// 	}
// 	return location
// }

// func gatherDataUp3(link string) []Dates {
// 	data1, e1 := http.Get(link)
// 	if e1 != nil {
// 		log.Fatal(e1)
// 	}
// 	var index Dates
// 	err := json.NewDecoder(data1.Body).Decode(&index)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	//fmt.Print(index)

// 	var dates []Dates
// 	for _, item := range index.Index {
// 		dates = append(dates, Dates{Index: []struct {
// 			ID    int      `json:"id"`
// 			Dates []string `json:"dates"`
// 		}{item}})

// 	}
// 	return dates
// }

func gatherDataUp4(link string) []Relation {
	client := &http.Client{
		Timeout: time.Second * 10, // Örnek olarak, 10 saniye olarak ayarlandı
	}
	data1, e1 := client.Get(link)
	//data1, e1 := http.Get(link)
	if e1 != nil {
		log.Fatal(e1)
	}

	var index Relation
	err := json.NewDecoder(data1.Body).Decode(&index)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Print(index)

	var relation []Relation
	for _, item := range index.Index {
		relation = append(relation, Relation{Index: []struct {
			ID             int                 `json:"id"`
			DatesLocations map[string][]string `json:"datesLocations"`
		}{item}})

	}

	return relation
}
