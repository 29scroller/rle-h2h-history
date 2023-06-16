package headtohead_history

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
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
func CountPlayersOfTeamInMatch(team []Player, match Match) (blueCount, orangeCount uint8) {
	for teamIndex := 0; teamIndex < len(team); teamIndex++ {
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

// IsMatchSuitable checks whether match is suitable for calculation.
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
	case blueCount == 2 && orangeCount == 3,
		blueCount == 3 && orangeCount == 2:
		kcomp = 0.75
	case blueCount == 2 && orangeCount == 2:
		kcomp = 0.5
	case blueCount == 1 && orangeCount == 3,
		blueCount == 3 && orangeCount == 1:
		kcomp = 0.5
	case blueCount == 1 && orangeCount == 2,
		blueCount == 2 && orangeCount == 1:
		kcomp = 0.25
	case blueCount == 1 && orangeCount == 1:
		kcomp = 0.1
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
	case daysPassed < 45:
		kdate = 1
	case daysPassed < 91:
		kdate = 0.8
	case daysPassed < 182:
		kdate = 0.5
	case daysPassed < 365:
		kdate = 0.25
	case daysPassed < 730:
		kdate = 0.1	
	case daysPassed >= 730:
		kdate = 0.01
	default:
		fmt.Println("Something's wrong with the DateCoefficient calculation")
		kdate = -1
	}
	return
}

// AdjustSeriesAndGames adjusts the result of the match by multiplying it on both coefficients defined earlier
func AdjustSeriesAndGames(match Match, blueCount, orangeCount uint8) (blueSeries, blueGames, orangeSeries, orangeGames float64) {
	kcomp := CompletionCoefficient(blueCount, orangeCount)
	kdate := DateCoefficient(match.Date)
	blueGames += (float64(match.Blue.Score) * kcomp * kdate)
	orangeGames += (float64(match.Orange.Score) * kcomp * kdate)
	if match.Blue.Winner {
		blueSeries = kcomp * kdate
	} else if match.Orange.Winner {
		orangeSeries = kcomp * kdate
	} else {
		blueSeries = kcomp * kdate * 0.5
		orangeSeries = kcomp * kdate * 0.5
	}
	return
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
		ourBlueCount, ourOrangeCount := CountPlayersOfTeamInMatch(ourTeamPlayers, ourTeamMatches[i])
		theirBlueCount, theirOrangeCount := CountPlayersOfTeamInMatch(theirTeamPlayers, ourTeamMatches[i])
		btoo, otbo := IsMatchSuitable(ourBlueCount, theirBlueCount, ourOrangeCount, theirOrangeCount)
		if btoo {
			ourSeries, ourGames, theirSeries, theirGames = AdjustSeriesAndGames(ourTeamMatches[i], ourBlueCount, theirOrangeCount)
			fmt.Printf("Match %s from %s is eligible! Adjusted series: %.2f - %.2f, adjusted games: %.2f - %.2f\n", ourTeamMatches[i].MEvent.Name, ourTeamMatches[i].Date, ourSeries, theirSeries, ourGames, theirGames)
		}
		if otbo {
			theirSeries, theirGames, ourSeries, ourGames = AdjustSeriesAndGames(ourTeamMatches[i], ourBlueCount, theirOrangeCount)
			fmt.Printf("Match %s from %s is eligible! Adjusted series: %.2f - %.2f, adjusted games: %.2f - %.2f\n", ourTeamMatches[i].MEvent.Name, ourTeamMatches[i].Date, ourSeries, theirSeries, ourGames, theirGames)
		}
		if btoo || otbo {
			counter++
		}
	}
	return
}

// CollectSuitableMatches checks whether matches are suitable for set conditions and collects them into slice of matches
func CollectSuitableMatches(ourTeamMatches []Match, ourTeamPlayers, theirTeamPlayers []Player) (suitableMatches []Match, counter int) {
	for i := 0; i < len(ourTeamMatches); i++ {
		ourBlueCount, ourOrangeCount := CountPlayersOfTeamInMatch(ourTeamPlayers, ourTeamMatches[i])
		theirBlueCount, theirOrangeCount := CountPlayersOfTeamInMatch(theirTeamPlayers, ourTeamMatches[i])
		btoo, otbo := IsMatchSuitable(ourBlueCount, theirBlueCount, ourOrangeCount, theirOrangeCount)
		if btoo {
			ourTeamMatches[i].btoo = true
		}
		if otbo {
			ourTeamMatches[i].otbo = true
		}
		if btoo || otbo {
			counter++
			date := ourTeamMatches[i].Date[:10]
			layout := "2006-01-02"
			ourTeamMatches[i].DateTime, _ = time.Parse(layout, date)
			ourTeamMatches[i].BlueT = ourBlueCount
			ourTeamMatches[i].BlueO = theirBlueCount
			ourTeamMatches[i].OrangeT = ourOrangeCount
			ourTeamMatches[i].OrangeO = theirOrangeCount
			suitableMatches = append(suitableMatches, ourTeamMatches[i])
		}
	}
	return
}

// DistributeSuitableMatchesToSlices writes each suitable match to corresponding slice and retuns slice of slices of matches
func DistributeSuitableMatchesToSlices(suitableMatches []Match) (allSlices [][]Match) {
	var m3v3, m2v3, m2v2, m1v3, m1v2, m1v1 []Match
	for i := 0; i < len(suitableMatches); i++ {
		switch {
		case suitableMatches[i].BlueT == 3 && suitableMatches[i].OrangeO == 3,
			suitableMatches[i].OrangeT == 3 && suitableMatches[i].BlueO == 3:
			m3v3 = append(m3v3, suitableMatches[i])
		case suitableMatches[i].BlueT == 2 && suitableMatches[i].OrangeO == 3,
			suitableMatches[i].BlueT == 3 && suitableMatches[i].OrangeO == 2,
			suitableMatches[i].OrangeT == 2 && suitableMatches[i].BlueO == 3,
			suitableMatches[i].OrangeT == 3 && suitableMatches[i].BlueO == 2:
			m2v3 = append(m2v3, suitableMatches[i])
		case suitableMatches[i].BlueT == 2 && suitableMatches[i].OrangeO == 2,
			suitableMatches[i].OrangeT == 2 && suitableMatches[i].BlueO == 2:
			m2v2 = append(m2v2, suitableMatches[i])
		case suitableMatches[i].BlueT == 1 && suitableMatches[i].OrangeO == 3,
			suitableMatches[i].BlueT == 3 && suitableMatches[i].OrangeO == 1,
			suitableMatches[i].OrangeT == 1 && suitableMatches[i].BlueO == 3,
			suitableMatches[i].OrangeT == 3 && suitableMatches[i].BlueO == 1:
			m1v3 = append(m1v3, suitableMatches[i])
		case suitableMatches[i].BlueT == 2 && suitableMatches[i].OrangeO == 1,
			suitableMatches[i].BlueT == 1 && suitableMatches[i].OrangeO == 2,
			suitableMatches[i].OrangeT == 2 && suitableMatches[i].BlueO == 1,
			suitableMatches[i].OrangeT == 1 && suitableMatches[i].BlueO == 2:
			m1v2 = append(m1v2, suitableMatches[i])
		case suitableMatches[i].BlueT == 1 && suitableMatches[i].OrangeO == 1,
			suitableMatches[i].OrangeT == 1 && suitableMatches[i].BlueO == 1:
			m1v1 = append(m1v1, suitableMatches[i])
		}
	}
	allSlices = [][]Match{m3v3, m2v3, m2v2, m1v3, m1v2, m1v1}
	return
}

// OperateMatchesWithinSlice processes each match in slice and adds results to the table, then prints it
func OperateMatchesWithinSlice(mxvx []Match, ourPlayers, theirPlayers []Player) (sumOurSeries, sumOurGames, sumTheirSeries, sumTheirGames float64) {
	sort.Slice(mxvx, func(i, j int) bool { return mxvx[i].DateTime.After(mxvx[j].DateTime) })
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Event", "Date", "Team 1", "Players", "Score", "Adjusted Score", "Team 2", "Players")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for i := 0; i < len(mxvx); i++ {
		bluePlayers := []string{mxvx[i].Blue.PlayerUp[0].Player.Tag, mxvx[i].Blue.PlayerUp[1].Player.Tag, mxvx[i].Blue.PlayerUp[2].Player.Tag}
		orangePlayers := []string{mxvx[i].Orange.PlayerUp[0].Player.Tag, mxvx[i].Orange.PlayerUp[1].Player.Tag, mxvx[i].Orange.PlayerUp[2].Player.Tag}
		if mxvx[i].btoo {
			HighlightPlayersOfTeamInBlue(ourPlayers, &mxvx[i])
			HighlightPlayersOfTeamInOrange(theirPlayers, &mxvx[i])
			ourSeries, ourGames, theirSeries, theirGames := AdjustSeriesAndGames(mxvx[i], mxvx[i].BlueT, mxvx[i].OrangeO)
			score := strconv.Itoa(int(mxvx[i].Blue.Score)) + "-" + strconv.Itoa(int(mxvx[i].Orange.Score))
			adjScore := fmt.Sprintf("%.2f", ourGames) + "-" + fmt.Sprintf("%.2f", theirGames)
			tbl.AddRow(mxvx[i].MEvent.Name, mxvx[i].Date, mxvx[i].Blue.TeamUp.Team.Name, bluePlayers, score, adjScore, mxvx[i].Orange.TeamUp.Team.Name, orangePlayers)
			sumOurSeries += ourSeries
			sumOurGames += ourGames
			sumTheirSeries += theirSeries
			sumTheirGames += theirGames
		}
		if mxvx[i].otbo {
			HighlightPlayersOfTeamInBlue(theirPlayers, &mxvx[i])
			HighlightPlayersOfTeamInOrange(ourPlayers, &mxvx[i])
			theirSeries, theirGames, ourSeries, ourGames := AdjustSeriesAndGames(mxvx[i], mxvx[i].OrangeT, mxvx[i].BlueO)
			score := strconv.Itoa(int(mxvx[i].Orange.Score)) + "-" + strconv.Itoa(int(mxvx[i].Blue.Score))
			adjScore := fmt.Sprintf("%.2f", ourGames) + "-" + fmt.Sprintf("%.2f", theirGames)
			tbl.AddRow(mxvx[i].MEvent.Name, mxvx[i].Date, mxvx[i].Orange.TeamUp.Team.Name, orangePlayers, score, adjScore, mxvx[i].Blue.TeamUp.Team.Name, bluePlayers)
			sumOurSeries += ourSeries
			sumOurGames += ourGames
			sumTheirSeries += theirSeries
			sumTheirGames += theirGames
		}
	}
	tbl.Print()
	return
}

// OperateAllSlices handles processing all slices and sums the data between slices
func OperateAllSlices(allSlices [][]Match, ourPlayers, theirPlayers []Player) (sumOurSeries, sumOurGames, sumTheirSeries, sumTheirGames float64) {
	var sliceName string
	for i := 0; i < len(allSlices); i++ {
		switch i {
		case 0:
			sliceName = "3v3"
		case 1:
			sliceName = "2v3"
		case 2:
			sliceName = "2v2"
		case 3:
			sliceName = "1v3"
		case 4:
			sliceName = "1v2"
		case 5:
			sliceName = "1v1"
		}
		if len(allSlices[i]) == 0 {
			fmt.Println("No matches found with", sliceName, "completeness")
		} else {
			fmt.Println(sliceName, "matches:")
			ourSeries, ourGames, theirSeries, theirGames := OperateMatchesWithinSlice(allSlices[i], ourPlayers, theirPlayers)
			sumOurSeries += ourSeries
			sumOurGames += ourGames
			sumTheirSeries += theirSeries
			sumTheirGames += theirGames
		}
	}
	return
}

// CalculateScoreForMatchup calculates score for matchup and returns it
func CalculateScoreForMatchup(ourTeamMatches []Match, ourTeamPlayers, theirTeamPlayers []Player) (sumOurSeries, sumOurGames, sumTheirSeries, sumTheirGames float64, counter int) {
	suitableMatches, counter := CollectSuitableMatches(ourTeamMatches, ourTeamPlayers, theirTeamPlayers)
	allSlices := DistributeSuitableMatchesToSlices(suitableMatches)
	sumOurSeries, sumOurGames, sumTheirSeries, sumTheirGames = OperateAllSlices(allSlices, ourTeamPlayers, theirTeamPlayers)
	return
}

// HighlightPlayersOfTeamInBlue adds InTeam marker to each player in blue team that is found in initial request
func HighlightPlayersOfTeamInBlue(team []Player, match *Match) {
	for i := 0; i < len(team); i++ {
		for j := 0; j < 3; j++ {
			if team[i].Tag == match.Blue.PlayerUp[j].Player.Tag {
				match.Blue.PlayerUp[j].Player.InTeam = true
			}
		}
	}
}

// HighlightPlayersOfTeamInBlue adds InTeam marker to each player in orange team that is found in initial request
func HighlightPlayersOfTeamInOrange(team []Player, match *Match) {
	for i := 0; i < len(team); i++ {
		for j := 0; j < 3; j++ {
			if team[i].Tag == match.Orange.PlayerUp[j].Player.Tag {
				match.Orange.PlayerUp[j].Player.InTeam = true
			}
		}
	}
}
