package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type ApiCoords struct {
	Locations []Coords
	Name      string
}

func (self *ApiCoords) FeedApiCord(cities []string, name string) { // remplis la struct ApiCoords des coordonnées (lat , lon) des villes rentré en paramètre
	self.Locations = []Coords{}
	self.Name = name
	Coordstemp := Coords{}
	for _, api := range cities {
		req, err := http.Get("https://nominatim.openstreetmap.org/search/" + api + "?format=json&addressdetails=1&limit=1&polygon_svg=1")
		if err != nil {
			log.Printf("[ERROR] - %v\n", err)
			return
		}
		data, err := ioutil.ReadAll(req.Body)
		arbitre := false  // vaut true si on a rencontré lat ou lon dans la chaine de caractère
		compt := 0        // nombre de " rencontré une fois que abitre == true
		latorlon := false // false = on cherche la coord de lat , true = on cherche la coord de lon
		result := ""
		for index, letter := range string(data) { // On parcours les données pour prendre ce que l'on veut car il ne veut pas récup lat et lon via unmarchall
			if arbitre {
				if letter == '"' {
					compt += 1
					if compt == 3 {
						arbitre = false
						compt = 0
						coord, err := strconv.ParseFloat(result, 64)
						if err != nil {
							fmt.Println(err)
							return
						}
						if !latorlon {
							Coordstemp.Lat = coord
							latorlon = true
						} else {
							Coordstemp.Lon = coord
							latorlon = false
						}
						result = ""
					}
				} else if compt == 2 {
					result = result + string(letter)
				}
			}
			if letter == 'l' && (string(data)[index:index+5] == "lat\":" || string(data)[index:index+5] == "lon\":") {
				arbitre = true
			}
		}
		self.Locations = append(self.Locations, Coordstemp)
	}
}
