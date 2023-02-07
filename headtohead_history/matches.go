package headtohead_history

import (
	"fmt"
	"math"
	"time"
)

// ParseAllMatchesOfPlayer parses all matches of a player by player ID into MultipleMatches struct.
func ParseAllMatchesOfPlayer(playerId string) MultipleMatches {
	url := "https://zsr.octane.gg/matches?mode=3&perPage=500&player=" + playerId
	rawData := UrlToByteSlice(url)
	var playerMatches MultipleMatches
	UnmarshalObject(rawData, &playerMatches)
	fmt.Println("Matches found:", len(playerMatches.Matches))
	return playerMatches
}

// ParseMatchInfo parses match info by match ID into Match struct.
func ParseMatchInfo(matchId string) Match {
	url := "https://zsr.octane.gg/matches/" + matchId
	rawData := UrlToByteSlice(url)
	var matchInfo Match
	UnmarshalObject(rawData, &matchInfo)
	return matchInfo
}

// FindPlayersOfTeamInMatch searches for each player of a team in Blue and Orange sides of match.
func FindPlayersOfTeamInMatch(team []Player, match Match) (blueCount, orangeCount uint8) {
	for teamIndex := 0; teamIndex < 3; teamIndex++ {
		for blueIndex := 0; blueIndex < 3; blueIndex++ {
			if match.Blue.PlayerUp[blueIndex].Player.Slug == team[teamIndex].Slug {
				blueCount++
			}
		}
		for orangeIndex := 0; orangeIndex < 3; orangeIndex++ {
			if match.Orange.PlayerUp[orangeIndex].Player.Slug == team[teamIndex].Slug {
				orangeCount++
			}
		}
	}
	return
}

// LoadAndAppendTeamMatchesFromFiles loads matches from players' files and appends them in one slice of matches.
func LoadAndAppendTeamMatchesFromFiles(matches *[]Match, players []Player, teamName string) {
	for i := 0; i < len(players); i++ {
		*matches = append(*matches, LoadPlayerMatchesFromFile(players[i])...)
	}
	fmt.Printf("Total number of matches for %s = %d\n", teamName, len(*matches))
	fmt.Println()
}

// IsMatchEligible checks whether match is suitable for calculation.
func IsMatchSuitable(blueTeammates, blueOpponents, orangeTeammates, orangeOpponents uint8) (btoo, otbo bool) {
	if blueTeammates > 0 && orangeOpponents > 0 {
		btoo = true
	}
	if orangeTeammates > 0 && blueOpponents > 0 {
		otbo = true
	}
	return
}

// CompletionCoefficient calculates completion coefficient based on number of teammates and opponents on each side.
func CompletionCoefficient(blueCount, orangeCount uint8) (kcomp float64) {
	switch {
	case blueCount == 3 && orangeCount == 3:
		kcomp = 1
	case blueCount == 2 && orangeCount == 3:
		kcomp = float64(5) / float64(6)
	case blueCount == 3 && orangeCount == 2:
		kcomp = float64(5) / float64(6)
	case blueCount == 2 && orangeCount == 2:
		kcomp = float64(2) / float64(3)
	case blueCount == 1 && orangeCount == 3:
		kcomp = float64(1) / float64(2)
	case blueCount == 3 && orangeCount == 1:
		kcomp = float64(1) / float64(2)
	case blueCount == 1 && orangeCount == 2:
		kcomp = float64(1) / float64(3)
	case blueCount == 2 && orangeCount == 1:
		kcomp = float64(1) / float64(3)
	case blueCount == 1 && orangeCount == 1:
		kcomp = float64(1) / float64(6)
	default:
		fmt.Println("Something's wrong with the CompletionCoefficient calculation")
		kcomp = -1
	}
	return
}

// DateCoefficient calculates date coefficient based on how many days passed since match date.
func DateCoefficient(date string) (kdate float64) {
	date = date[:10]
	layout := "2006-01-02"
	matchDate, _ := time.Parse(layout, date)
	daysPassed := (time.Since(matchDate).Hours() / 24)
	daysPassed = math.Round(daysPassed)
	switch {
	case daysPassed < 91:
		kdate = 1
	case daysPassed < 182:
		kdate = 0.8
	case daysPassed < 365:
		kdate = 0.6
	case daysPassed < 730:
		kdate = 0.4
	case daysPassed >= 730:
		kdate = 0.2
	default:
		fmt.Println("Something's wrong with the DateCoefficient calculation")
		kdate = -1
	}
	return
}

// AdjustSeriesAndGames adjusts the result of the match by multiplying it on both coefficients defined earlier
// Result is stored in AdjustedResult struct.
func AdjustSeriesAndGames(match *Match, blueCount, orangeCount uint8) {
	kcomp := CompletionCoefficient(blueCount, orangeCount)
	kdate := DateCoefficient(match.Date)
	match.Blue.AdjustedResult.AdjustedGames += (float64(match.Blue.Score) * kcomp * kdate)
	match.Orange.AdjustedResult.AdjustedGames += (float64(match.Orange.Score) * kcomp * kdate)
	if match.Blue.Winner {
		match.Blue.AdjustedResult.AdjustedSeries = kcomp * kdate
	} else if match.Orange.Winner {
		match.Orange.AdjustedResult.AdjustedSeries = kcomp * kdate
	} else {
		match.Blue.AdjustedResult.AdjustedSeries = kcomp * kdate * 0.5
		match.Orange.AdjustedResult.AdjustedSeries = kcomp * kdate * 0.5
	}
}

// RemoveDuplicatesinMatches removes duplicate matches from slice of Matches struct.
func RemoveDuplicatesinMatches(matches []Match, teamName string) (result []Match) {
	inResult := make(map[string]bool)
	for _, match := range matches {
		if _, ok := inResult[match.Id]; !ok {
			inResult[match.Id] = true
			result = append(result, match)
		}
	}
	fmt.Printf("Deleted %d duplicates, new count of matches for team %s = %d\n", len(matches)-len(result), teamName, len(result))
	return
}

// CalculateMatchesForTeam calculates adjusted series and score for slice of matches, given slices of players of both teams.
func CalculateMatchesForTeam(ourTeamMatches []Match, ourTeamPlayers, theirTeamPlayers []Player) (ourSeries, theirSeries, ourGames, theirGames float64, counter int) {
	for i := 0; i < len(ourTeamMatches); i++ {
		ourBlueCount, ourOrangeCount := FindPlayersOfTeamInMatch(ourTeamPlayers, ourTeamMatches[i])
		theirBlueCount, theirOrangeCount := FindPlayersOfTeamInMatch(theirTeamPlayers, ourTeamMatches[i])
		btoo, otbo := IsMatchSuitable(ourBlueCount, theirBlueCount, ourOrangeCount, theirOrangeCount)
		if btoo {
			AdjustSeriesAndGames(&ourTeamMatches[i], ourBlueCount, theirOrangeCount)
			ourSeries += ourTeamMatches[i].Blue.AdjustedResult.AdjustedSeries
			theirSeries += ourTeamMatches[i].Orange.AdjustedResult.AdjustedSeries
			ourGames += ourTeamMatches[i].Blue.AdjustedResult.AdjustedGames
			theirGames += ourTeamMatches[i].Orange.AdjustedResult.AdjustedGames
			fmt.Printf("Match %s from %s is eligible! Adjusted series: %.2f - %.2f, adjusted games: %.2f - %.2f\n", ourTeamMatches[i].MEvent.Name, ourTeamMatches[i].Date, ourTeamMatches[i].Blue.AdjustedResult.AdjustedSeries, ourTeamMatches[i].Orange.AdjustedResult.AdjustedSeries, ourTeamMatches[i].Blue.AdjustedResult.AdjustedGames, ourTeamMatches[i].Orange.AdjustedResult.AdjustedGames)
		}
		if otbo {
			AdjustSeriesAndGames(&ourTeamMatches[i], theirBlueCount, ourOrangeCount)
			ourSeries += ourTeamMatches[i].Orange.AdjustedResult.AdjustedSeries
			theirSeries += ourTeamMatches[i].Blue.AdjustedResult.AdjustedSeries
			ourGames += ourTeamMatches[i].Orange.AdjustedResult.AdjustedGames
			theirGames += ourTeamMatches[i].Blue.AdjustedResult.AdjustedGames
			fmt.Printf("Match %s from %s is eligible! Adjusted series: %.2f - %.2f, adjusted games: %.2f - %.2f\n", ourTeamMatches[i].MEvent.Name, ourTeamMatches[i].Date, ourTeamMatches[i].Orange.AdjustedResult.AdjustedSeries, ourTeamMatches[i].Blue.AdjustedResult.AdjustedSeries, ourTeamMatches[i].Orange.AdjustedResult.AdjustedGames, ourTeamMatches[i].Blue.AdjustedResult.AdjustedGames)
		}
		if btoo || otbo {
			counter++
		}
	}
	return
}
