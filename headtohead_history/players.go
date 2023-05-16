package headtohead_history

import (
	"fmt"
	"strings"
)

// FindPlayerByTag finds player by tag in zsr.octane.gg database and returns player's info
func FindPlayerByTag(playerTag string) (Player, bool) {
	//Writing team name in lowercase, replacing space with %20, trimming whitespace
	playerTag = strings.TrimSpace(strings.ToLower(strings.Replace(playerTag, " ", "%20", 1)))
	url := "https://zsr.octane.gg/players?tag=" + playerTag
	rawData := UrlToByteSlice(url)
	var playerInfo MultiplePlayers
	UnmarshalObject(rawData, &playerInfo)
	if len(playerInfo.Players) == 0 {
		fmt.Println("Could not find any players with that name, please try again")
		return Player{}, false
	} else if len(playerInfo.Players) > 1 {
		fmt.Println("Found multiple players??? Choose the right player and enter its number")
		fmt.Print("PLayers found: \n")
		for i := 0; i < len(playerInfo.Players); i++ {
			fmt.Printf("%d - %s, %s\n", i+1, playerInfo.Players[i].Tag, playerInfo.Players[i].PTeam.Name)
		}
		fmt.Println("Enter a number of a needed player and press Enter (0 for searching another player)")
		var res int
		fmt.Scanln(&res)
		switch {
		case res == 0:
			return Player{}, false
		case res > len(playerInfo.Players):
			return Player{}, false
		default:
			fmt.Println("You chose", playerInfo.Players[res-1].Tag)
			return playerInfo.Players[res-1], true
		}
	} else {
		fmt.Printf("Found player %s of team %s; confirm it's the right player by typing 1 or type something else to search different team\n", playerInfo.Players[0].Tag, playerInfo.Players[0].PTeam.Name)
		var userInput int
		fmt.Scanln(&userInput)
		if userInput == 1 {
			return playerInfo.Players[0], true
		} else {
			return Player{}, false
		}
	}
}
