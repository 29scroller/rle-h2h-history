package headtohead_history

import (
	"fmt"
	"strings"
)

// FindTeamByName finds team in zsr.octane.gg by its name.
func FindTeamByName(teamName string) (Team, bool) {
	//Writing team name in lowercase, replacing space with %20, trimming whitespace
	teamName = strings.TrimSpace(strings.ToLower(strings.Replace(teamName, " ", "%20", 1)))
	url := "https://zsr.octane.gg/teams?name=" + teamName
	rawData := UrlToByteSlice(url)
	var teamInfo Player
	UnmarshalObject(rawData, &teamInfo)
	if len(teamInfo.Teams) == 0 { //separate function
		fmt.Println("Could not find any teams with that name, please try again")
		return Team{}, false
	} else if len(teamInfo.Teams) > 1 { //separate function
		fmt.Println("Found multiple teams??? Choose the right team and enter its number")
		fmt.Print("Teams found: \n")
		for i := 0; i < len(teamInfo.Teams); i++ {
			fmt.Printf("%d - %s \n", i+1, teamInfo.Teams[i].Name)
		}
		fmt.Println()
		fmt.Println("Enter a number of a needed team and press Enter (0 for searching another team)")
		var res int
		fmt.Scanln(&res)
		switch {
		case res == 0:
			return Team{}, false
		case res > len(teamInfo.Teams):
			return Team{}, false
		default:
			fmt.Println("You chose", teamInfo.Teams[res-1].Name)
			return teamInfo.Teams[res-1], true
		}
	} else { //separate function
		fmt.Printf("Found team %s; confirm it's the right team by typing 1 or type something else to search different team\n", teamInfo.Teams[0].Name)
		var res int
		fmt.Scanln(&res)
		if res == 1 {
			return teamInfo.Teams[0], true
		} else {
			return Team{}, false
		}
	}
}

// FindActivePlayersByTeamID finds team by team ID and returns all players of a team who are neither substitute nor coach.
func FindActivePlayersByTeamID(teamId string) []Player {
	url := "https://zsr.octane.gg/players?team=" + teamId
	rawData := UrlToByteSlice(url)
	var playersInfo MultiplePlayers
	UnmarshalObject(rawData, &playersInfo)
	fmt.Print("Active players are: \n")
	var activePlayers []Player
	for i := 0; i < len(playersInfo.Players); i++ {
		if !playersInfo.Players[i].Substitute && !playersInfo.Players[i].Coach {
			activePlayers = append(activePlayers, playersInfo.Players[i])
		}
	}
	for i := 0; i < len(activePlayers); i++ {
		fmt.Print(activePlayers[i].Tag + " | ")
	}
	fmt.Println()
	return activePlayers
}
