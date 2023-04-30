package models

import (
	"encoding/json"
	"groupie-tracker/controller"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ApiData struct {
	AllBands           []BandInfo
	DisplayBands       []BandInfo
	CacheByFirstAlbum  []BandInfo
	CacheByMembers     map[int][]BandInfo
	CacheByYearStarted map[int][]BandInfo
	CacheByMemberName  map[string][]BandInfo
	CacheByCountry     map[string][]BandInfo
}

func (self *ApiData) FeedApi() {
	/*
		Method of ApiData

		Extract the data from the heroku api to self.AllBands and self.DisplayBands.
		Current content of self.AllBands and self.DisplayBands will be erased and replaced by new data.
	*/
	artistTemp := []BandInfo{}
	locTemp := map[string][]BandInfo{}

	log.Printf("[INFO] - Extracting data from API...\n")
	for _, api := range []string{"artists", "relation"} {
		req, err := http.Get("https://groupietrackers.herokuapp.com/api/" + api)
		if err != nil {
			log.Printf("[ERROR] - While reaching the API.\n%v\n", err.Error())
			return
		}

		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("[ERROR] - While reading response from the API.\n%v\n", err.Error())
			return
		}

		if api == "artists" {
			err = json.Unmarshal(data, &artistTemp)
		} else {
			err = json.Unmarshal(data, &locTemp)
		}

		if err != nil {
			log.Printf("[WARNING] - Ignoring some useless JSON data from the API.\n", err)
		}
	}

	// Deep copy of tmpBandsData to avoid storing it on memory
	// Since it's a local variable it should only exists in the function's scope.
	self.AllBands = make([]BandInfo, len(artistTemp))
	copy(self.AllBands, artistTemp)

	for index, band := range locTemp["index"] {
		formatedRelation := map[string][]string{}
		for place, date := range band.Relations {
			formatedPlace := strings.Title(strings.ReplaceAll(strings.ReplaceAll(place, "-", " - "), "_", " "))
			formatedRelation[formatedPlace] = date
		}
		self.AllBands[index].Relations = formatedRelation
	}

	// Shallow copy of AllBands since it's not a function scoped variable.
	self.DisplayBands = self.AllBands

	log.Printf("[INFO] - Data succesfully extracted!\n")
}

func (self *ApiData) CreateCaches() {
	/*
		Method of ApiData

		Creates all the Caches of the program.
		If there are already exisiting, their content will be erased.
	*/
	self.CacheByFirstAlbum = []BandInfo{}
	self.CacheByMembers = map[int][]BandInfo{}
	self.CacheByCountry = map[string][]BandInfo{}
	self.CacheByYearStarted = map[int][]BandInfo{}
	self.CacheByMemberName = map[string][]BandInfo{}
}

func (self *ApiData) CitiesTab(id int) ([]string, string) { // Renvoie un tableau de string contenant les noms des villes ou le groupe ayant l'id "id" à fais des concerts
	result := []string{}
	temp := ""
	var name string
	flag := true
	for _, band := range self.AllBands {
		if band.Id == id {
			name = band.Name
			for key, _ := range band.Relations {
				temp = strings.Split(key, "-")[0]
				for _, word := range result {
					if word == temp {
						flag = false
						break
					}
				}
				if flag {
					result = append(result, temp)
				}
				flag = true
			}
		}
	}
	return result, name
}

func (self *ApiData) RootHandler(webpage http.ResponseWriter, request *http.Request) {
	/*
		Method of ApiData

		This is the main function of the structure, it handles an HTTP connection on the "/" path.
		If the request method is POST, it will try to perform some sorting based on the content of the
		request. Otherwise, it will simply serve the basic main HTML page with the self.DisplayBands
		displayed.

		:param webpage: the page we're writing on
		:param request: the current request
	*/
	if request.URL.Path != "/" && request.URL.Path != "" {
		controller.ServeFile(webpage, "404.html", nil)
		return
	}
	self.DisplayBands = self.AllBands
	if request.Method == "POST" {
		if self.CacheByYearStarted[2022] == nil {
			self.GetBandsByYearStarted(2022) // This is necessary. See function GetBandsByYearStarted.
		}

		self.DisplayBands = []BandInfo{}
		request.ParseForm()

		if len(request.Form["filter_startingyear"]) != 0 {
			year := controller.AtoiSlice(request.Form["filter_startingyear"])[0]
			self.GetBandsByYearStarted(year)
			if request.Form["filter_startingyear"][0] == "2022" {
				self.GetAllBandsById()
			}
		}

		if len(request.Form["filter_nmembers"]) != 0 {
			sizes := controller.AtoiSlice(request.Form["filter_nmembers"])
			self.GetBandsByMemberSize(sizes)
		}

		if len(request.Form["filter_location"]) != 0 && request.Form["filter_location"][0] != "" {
			self.GetBandsByCountry(strings.Title(request.Form["filter_location"][0]))
		}

		if len(request.Form["filter_firstalbum"]) != 0 {
			self.GetBandsByFirstAlbum()
		}

		if len(request.Form["input-search"]) != 0 {
			self.DisplayBands = []BandInfo{}
			data := strings.Split(request.Form["input-search"][0], " (")
			if len(data) > 1 {
				if data[1] == "band)" {
					self.GetBandByName(data[0])
				} else if data[1] == "artist)" {
					self.GetBandByMemberName(data[0])
				}
			} else {
				number, err := strconv.Atoi(data[0])
				if err == nil {
					self.GetBandsByYearStarted(2022)
					self.GetBandsByMemberSize([]int{number})
					if len(self.DisplayBands) == 0 {
						self.GetBandsByYearStarted(number)
					}
				} else {
					self.GetBandByName(data[0])
					if len(self.DisplayBands) == 0 {
						self.GetBandByMemberName(data[0])
					}
				}
			}
		}

		if len(self.DisplayBands) == 0 {
			self.DisplayBands = self.AllBands
			log.Printf("Nothing to show")
		}
	}

	controller.ServeFile(webpage, "index.html", self)
}

func (self *ApiData) GetBandByName(name string) {
	/*
		Method of ApiData

		Simply adds to self.DisplayBands the band corresponding with the name
		requested.

		:param name: The name of the band
	*/
	for _, band := range self.AllBands {
		if strings.ToLower(band.Name) == strings.ToLower(name) {
			self.DisplayBands = append(self.DisplayBands, band)
			return
		}
	}
}

func (self *ApiData) GetBandByMemberName(name string) {
	/*
		Method of ApiData

		If there is no cache for the requested name, the function will attempt to fill it.
		It will eitherway add to self.DisplayBands the current content of the cache for the requested name (that should now be filled).

		:param name: The name of the member requested
	*/
	name = strings.ToLower(name)
	if self.CacheByMemberName[name] == nil {
		log.Printf("[INFO] - BandByMemberName [%v] not found in cache, creating it...", name)
		flag := false
		for _, band := range self.AllBands {
			for _, member := range band.Members {
				if strings.ToLower(member) == name {
					self.CacheByMemberName[name] = append(self.CacheByMemberName[name], band)
					flag = true
					break
				}
				if flag {
					break
				}
			}
		}
		log.Printf("\t     Cache for BandByMemberName [%v] created succesfully, %v objects added.",
			name, len(self.CacheByMemberName[name]))
	}
	self.DisplayBands = append(self.DisplayBands, self.CacheByMemberName[name]...)
}

func (self *ApiData) GetBandsByYearStarted(year int) {
	/*
		Method of ApiData

		If there is no cache for the year requested, the function will attempt to fill it.
		To do so, the function will seek for a greater year in cache and then use its' content
		to fill the cache of the requested year. If no greater year is found, it will use all of the
		bands to fill the cache. To use this function at its peak efficiency, we need to first store in cache
		the year 2022, so every other year will use 2022 to construct themselves.

		Then adds to self.DisplayBands the current content of the cache for the requested year (that should now be filled).

		:param year: The requested year
	*/
	if self.CacheByYearStarted[year] == nil {
		log.Printf("[INFO] - CreationDate [%v] not found in cache, creating it...", year)
		if len(self.CacheByYearStarted) != 0 {
			keys := GetKeys(self.CacheByYearStarted)
			closest := controller.GetClosestTo(keys, year)
			if closest != -1 {
				tmp := self.CacheByYearStarted[closest]

				for i := len(tmp) - 1; i >= 0; i-- {
					if tmp[i].CreationDate > year {
						tmp = tmp[:i]
					}
				}

				self.CacheByYearStarted[year] = make([]BandInfo, len(tmp))
				copy(self.CacheByYearStarted[year], tmp)

				self.DisplayBands = self.CacheByYearStarted[year]

				log.Printf("\t     Cache for CreationDate [%v] created succesfully, %v objects added.",
					year, len(self.CacheByYearStarted[year]))
				return
			}
		}
		// If nothing is in cache:
		tmp := make([]BandInfo, len(self.AllBands))
		copy(tmp, self.AllBands)
		sort.Slice(tmp, func(i, j int) bool {
			return tmp[i].CreationDate < tmp[j].CreationDate
		})

		self.CacheByYearStarted[year] = []BandInfo{}
		for _, band := range tmp {
			if band.CreationDate > year {
				break
			}
			self.CacheByYearStarted[year] = append(self.CacheByYearStarted[year], band)
		}
		log.Printf("\t     Cache for CreationDate [%v] created succesfully, %v objects added.",
			year, len(self.CacheByYearStarted[year]))

	}
	self.DisplayBands = append(self.DisplayBands, self.CacheByYearStarted[year]...)
}

func (self *ApiData) GetBandsByMemberSize(size []int) {
	/*
		Method of ApiData

		If there is no cache for the requested member size, the function will attempt to fill it.
		It will eitherway add to self.DisplayBands the content of the cache for the member size (that should now be filled).

		:param size: The size of the band
	*/
	tempBands := []BandInfo{}
	for _, number := range size {
		if self.CacheByMembers[number] == nil {
			log.Printf("[INFO] - MemberSize [%v] not found in cache, creating it...", number)
			self.CacheByMembers[number] = []BandInfo{}
			for _, band := range self.DisplayBands {
				if len(band.Members) == number {
					self.CacheByMembers[number] = append(self.CacheByMembers[number], band)
				}
			}
			log.Printf("\t     Cache for MemberSize [%v] created succesfully, %v objects added.",
				number, len(self.CacheByMembers[number]))
		}
		tempBands = append(tempBands, self.CacheByMembers[number]...)
	}
	self.DisplayBands = tempBands
}

func (self *ApiData) GetBandsByCountry(country string) {
	/*
		Method of ApiData

		If there is no cache for the requested country, the function will attempt to fill it.
		It will eitherway add to self.DisplayBands the content of the cache for the country (that should now be filled).

		:param country: The requested country
	*/
	if self.CacheByCountry[country] == nil {
		log.Printf("[INFO] - BandCountry [%v] not found in cache, creating it...", country)
		self.CacheByCountry[country] = []BandInfo{}
		for _, band := range self.DisplayBands {
			for location := range band.Relations {
				location = strings.Split(location, "- ")[1]
				if location == country {
					self.CacheByCountry[country] = append(self.CacheByCountry[country], band)
					break
				}
			}
		}
		log.Printf("\t     Cache for BandCountry [%v] created succesfully, %v objects added.",
			country, len(self.CacheByCountry[country]))
	}
	self.DisplayBands = self.CacheByCountry[country]
}

func (self *ApiData) GetBandsByFirstAlbum() {
	/*
		Method of ApiData

		This sorts self.DisplayBands by the date of the first album of the bands.
	*/
	sort.Slice(self.DisplayBands, func(i, j int) bool {
		date1 := strings.Split(self.DisplayBands[i].FirstAlbum, "-") // ["12", "02", "2020"]
		date2 := strings.Split(self.DisplayBands[j].FirstAlbum, "-")
		if date1[2] < date2[2] {
			return true
		} else if date1[2] > date2[2] {
			return false
		} else {
			if date1[1] < date2[1] {
				return true
			} else if date1[1] > date2[1] {
				return false
			} else {
				if date1[0] < date2[0] {
					return true
				} else {
					return false
				}
			}
		}
	})
}

func (self *ApiData) GetAllBandsById() {
	/*
		Method of ApiData

		This repopulate self.DisplayBands then sort it by the id of the bands
	*/
	self.DisplayBands = self.AllBands
	sort.Slice(self.DisplayBands, func(i, j int) bool {
		return self.DisplayBands[i].Id < self.DisplayBands[j].Id
	})
}

func (self *ApiData) WaitThenRefreshApi() {
	/*
		Method of ApiData

		Loop forever, waits 24 hours then refresh the api and empty the caches
	*/
	for true {
		time.Sleep(24 * time.Hour)
		log.Printf("[INFO] - 30s since last cache update, refreshing the API.")
		self.FeedApi()
		self.CreateCaches()
	}
}

// TODO: found a way to put this in controller/utils.go
func GetKeys(items map[int][]BandInfo) []int {
	result := []int{}
	for key, _ := range items {
		result = append(result, key)
	}
	return result
}

func KeepOnlyDuplicates(longarray []BandInfo, shorterarray []BandInfo) []BandInfo {
	result := []BandInfo{}
	allKeys := make(map[int]bool)

	for _, element := range longarray {
		allKeys[element.Id] = true
	}

	for _, element := range shorterarray {
		if allKeys[element.Id] {
			result = append(result, element)
		}
	}

	return result
}
