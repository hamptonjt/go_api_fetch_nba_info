package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

// Teams - Top-level XML structure for teams endpoint
type Teams struct {
	XMLName xml.Name `xml:"FantasyBasketballNerd"`
	Service string   `xml:"service,attr"`
	Comment string   `xml:",comment"`
	Teams   []Team   `xml:"Team"`
}

// Team - Team details
type Team struct {
	Code       string `xml:"code"`
	Name       string `xml:"name"`
	Conference string `xml:"conference"`
	Division   string `xml:"division"`
}

// Players - Top-Level XML Structure for players endpoint
type Players struct {
	XMLName xml.Name `xml:"FantasyBasketballNerd"`
	Service string   `xml:"service,attr"`
	Comment string   `xml:",comment"`
	Players []Player `xml:"Player"`
}

// Player - Player details
type Player struct {
	PlayerID  int    `xml:"playerId"`
	Name      string `xml:"name"`
	Team      string `xml:"team"`
	Position  string `xml:"position"`
	Height    string `xml:"height"`
	Weight    string `xml:"weight"`
	BirthDate string `xml:"dob"`
	School    string `xml:"school"`
}

// Roster - Team and Players associated to the given team.
type Roster struct {
	Team    Team
	Players []Player
}

func main() {
	now := time.Now()
	fmt.Println("NBA Data Collection - NBA API")
	var teams []Team
	var players []Player

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		teams = getTeams()
		wg.Done()
	}()
	// fmt.Println(teams)

	go func() {
		players = getPlayers()
		wg.Done()
	}()

	wg.Wait()

	var rosters []Roster
	for _, team := range teams {
		// fmt.Println(team.Name, team.Conference, team.Division)
		var roster Roster
		roster.Team = team
		roster.Players = filter(players, team.Code)
		rosters = append(rosters, roster)
	}

	f, err := os.Create("nba_teams.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	for _, r := range rosters {
		fmt.Fprintf(f, "%s\n", r.Team.Name)
		fmt.Fprintf(f, "\t%-25s%-4s%s\n", "Name", "Pos", "College")
		for _, p := range r.Players {
			fmt.Fprintf(f, "\t%-25s%-4s%s\n", p.Name, p.Position, p.School)
		}
		fmt.Fprintf(f, "\n\n")
	}

	fmt.Printf("Took %v\n", time.Now().Sub(now).String())
}

func filter(players []Player, teamCode string) (ret []Player) {
	for _, p := range players {
		if p.Team == teamCode {
			ret = append(ret, p)
		}
	}
	return
}

func getTeams() []Team {
	teamResp, err := http.Get("https://www.fantasybasketballnerd.com/service/teams")
	if err != nil {
		fmt.Println(err)
	}

	bytes, err := ioutil.ReadAll(teamResp.Body)
	if err != nil {
		fmt.Println(err)
	}
	teamResp.Body.Close()

	var teams Teams
	err = xml.Unmarshal(bytes, &teams)
	if err != nil {
		fmt.Println(err)
	}
	return teams.Teams
}

func getPlayers() []Player {
	playerResp, err := http.Get("https://www.fantasybasketballnerd.com/service/players")
	if err != nil {
		fmt.Println(err)
	}
	playerBytes, err := ioutil.ReadAll(playerResp.Body)
	if err != nil {
		fmt.Println(err)
	}
	playerResp.Body.Close()

	var players Players
	err = xml.Unmarshal(playerBytes, &players)
	if err != nil {
		fmt.Println(err)
	}
	return players.Players
}
