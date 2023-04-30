package models


type BandInfo struct {
	Id 				int 					`json:"id"`
	CreationDate	int 					`json:"creationDate"`
	Name 			string 					`json:"name"`
	Image 			string 					`json:"image"`
	FirstAlbum		string 					`json:"firstAlbum"`
	Members			[]string 				`json:"members"`
	Relations		map[string][]string 	`json:"datesLocations"`
}
