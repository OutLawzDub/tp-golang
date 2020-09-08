package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Templates

var homepageTpl *template.Template

// Weapons

type Map struct {
	Id       string `json:"Id"`
	Name     string `json:"Name"`
	Gamers   []Player
	Selected bool
}

type Weapon struct {
	Id     string `json:"Id,omitempty"`
	Name   string `json:"Name"`
	Damage int    `json:"Damage"`
}

// Player

type Player struct {
	Id       string `json:"Id,omitempty"`
	Username string `json:"Username"`
	PV       int    `json:"PV"`
	Armes    Weapon `json:"Armes"`
}

func (p Player) getInfo() {
	fmt.Println(p.Username, p.PV)
}

var Players []Player
var Maps []Map
var Weapons []Weapon

// Home

func HomeGame(w http.ResponseWriter, r *http.Request) {

}

// Web Player

func returnAllPlayers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllPlayers")
	json.NewEncoder(w).Encode(Players)
}

func returnOnePlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["idPlayer"]

	for _, player := range Players {
		if player.Id == key {
			json.NewEncoder(w).Encode(player)
		}
	}
}

func deleteOnePlayer(idPlayer string) {
	for index, player := range Players {
		if player.Id == idPlayer {
			Players = append(Players[:index], Players[index+1:]...)
		}
	}
}

func addNewPlayer(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	println(reqBody)

	var player Player

	json.Unmarshal(reqBody, &player)

	Players = append(Players, player)

	json.NewEncoder(w).Encode(player)
}

func addPlayerWeapon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idPlayer := vars["idPlayer"]
	idWeapon := vars["idWeapon"]

	var theWeapon Weapon

	for _, weapon := range Weapons {
		if weapon.Id == idWeapon {
			theWeapon = weapon
		}
	}

	var thePlayer Player

	for _, player := range Players {
		if player.Id == idPlayer {
			thePlayer = player
		}
	}

	var player Player

	deleteOnePlayer(thePlayer.Id)

	player = Player{
		Id:       thePlayer.Id,
		Username: thePlayer.Username,
		PV:       thePlayer.PV,
		Armes:    theWeapon,
	}

	Players = append(Players, player)

}

//

// Web Weapons

func returnAllWeapons(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Weapons)
}

func returnOneWeapon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["idWeapon"]

	for _, weapon := range Weapons {
		if weapon.Id == key {
			json.NewEncoder(w).Encode(weapon)
		}
	}
}

func addNewWeapon(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var weapon Weapon

	json.Unmarshal(reqBody, &weapon)

	Weapons = append(Weapons, weapon)

	json.NewEncoder(w).Encode(weapon)
}

//

// Web Maps

func deleteMap(idMap string) {
	for index, carte := range Maps {
		if carte.Id == idMap {
			Maps = append(Maps[:index], Maps[index+1:]...)
		}
	}
}

func selectMap(idMap string) {

	var carto Map

	for _, carte := range Maps {
		if carte.Id == idMap {
			carto = carte
		}
	}

	deleteMap(idMap)

	var newCarte Map

	newCarte = Map{
		Id:       carto.Id,
		Name:     carto.Name,
		Gamers:   carto.Gamers,
		Selected: true,
	}

	Maps = append(Maps, newCarte)
}

func selectOneMap(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idMap := vars["idMap"]

	selectMap(idMap)
}

func returnAllMaps(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Maps)
}

func returnOneMap(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["idMap"]

	for _, carte := range Maps {
		if carte.Id == key {
			json.NewEncoder(w).Encode(carte)
		}
	}
}


//

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	// Game

	myRouter.HandleFunc("/", HomeGame)

	// Player route

	myRouter.HandleFunc("/player/", returnAllPlayers)
	myRouter.HandleFunc("/player/{idPlayer}", returnOnePlayer)
	myRouter.HandleFunc("/playerAdd", addNewPlayer).Methods("POST")
	myRouter.HandleFunc("/player/addWeapon/{idPlayer}/{idWeapon}", addPlayerWeapon)

	// Weapons route

	myRouter.HandleFunc("/weapon/", returnAllWeapons)
	myRouter.HandleFunc("/weapon/{idWeapon}", returnOneWeapon)
	myRouter.HandleFunc("/weaponAdd", addNewWeapon)

	// Weapons Map

	myRouter.HandleFunc("/map/", returnAllMaps)
	myRouter.HandleFunc("/map/{idMap}", returnOneMap)
	myRouter.HandleFunc("/map/select/{idMap}", selectOneMap)

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

//

func main() {
	// Init Maps

	Maps = []Map{
		Map{Id: "1", Name: "de_dust2", Selected: false},
		Map{Id: "2", Name: "de_bhop_easy", Selected: false},
	}

	// Init Weapons

	Weapons = []Weapon{
		Weapon{Id: "1", Name: "AK-47", Damage: 50},
		Weapon{Id: "2", Name: "USP-S", Damage: 25},
		Weapon{Id: "3", Name: "AWP", Damage: 100},
	}

	// Init HTML (?)

	homepageHTML := assets.MustAssetString("templates/index.html")
	homepageTpl = template.Must(template.New("homepage_view").Parse(homepageHTML))

	handleRequests()
}
