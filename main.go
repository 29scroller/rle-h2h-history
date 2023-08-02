package main

import (
	"fmt"
	"rle-h2h-history/headtohead_history"
)

func main() {
	var choseInput, choseTeams bool
	var firstTeamInfo, secondTeamInfo headtohead_history.Team
	var firstTeamPlayers, secondTeamPlayers []headtohead_history.Player

	//headtohead_history.CollectAllMatchesInCSV()

	for !choseInput {
		userInput := headtohead_history.UserChoosingInput()
		switch userInput {
		case 1:
			firstTeamInfo, secondTeamInfo = headtohead_history.UserEnteringTeams()
			choseInput = true
			choseTeams = true
			fmt.Printf("Finding active players of %s\n", firstTeamInfo.Name)
			firstTeamPlayers = headtohead_history.FindActivePlayersByTeamID(firstTeamInfo.Id)
			fmt.Printf("Finding active players of %s\n", secondTeamInfo.Name)
			secondTeamPlayers = headtohead_history.FindActivePlayersByTeamID(secondTeamInfo.Id)
		case 2:
			firstTeamPlayers, secondTeamPlayers = headtohead_history.UserEnteringAllPlayers()
			choseInput = true
		default:
			fmt.Println("Need either 1 or 2 to continue")
		}
	}

	fmt.Println("Checking if all players have files...")
	headtohead_history.CheckIfTeamPlayersHaveFiles(firstTeamPlayers)
	headtohead_history.CheckIfTeamPlayersHaveFiles(secondTeamPlayers)

	choseInput = false
	for !choseInput {
		userInput := headtohead_history.UserChoosingWhetherToCheckPlayers()
		switch userInput {
		case 1:
			fmt.Println("Checking all players for updates...")
			for i := 0; i < len(firstTeamPlayers); i++ {
				headtohead_history.CheckIfPlayerHasFile(firstTeamPlayers[i])
				headtohead_history.CheckPlayerForNewMatches(firstTeamPlayers[i])
			}
			for i := 0; i < len(secondTeamPlayers); i++ {
				headtohead_history.CheckIfPlayerHasFile(firstTeamPlayers[i])
				headtohead_history.CheckPlayerForNewMatches(secondTeamPlayers[i])
			}
			choseInput = true
		case 0:
			fmt.Println("Got it, won't check")
			choseInput = true
		default:
			fmt.Println("Need either 1 or 0 to continue")
		}
	}

	fmt.Println("Reading data from players' files...")
	var firstTeamMatches, secondTeamMatches []headtohead_history.Match
	headtohead_history.LoadAndAppendTeamMatchesFromFiles(&firstTeamMatches, firstTeamPlayers, firstTeamInfo.Name)
	headtohead_history.LoadAndAppendTeamMatchesFromFiles(&secondTeamMatches, secondTeamPlayers, secondTeamInfo.Name)

	if !choseTeams {
		firstTeamInfo.Name = "Team 1"
		secondTeamInfo.Name = "Team 2"
	}

	fmt.Println("Deleting duplicates...")
	firstTeamMatches = headtohead_history.RemoveDuplicatesinMatches(firstTeamMatches, firstTeamInfo.Name)
	secondTeamMatches = headtohead_history.RemoveDuplicatesinMatches(secondTeamMatches, secondTeamInfo.Name)
	fmt.Println()

	eligibleMatchCounter := 0
	var firstAdjustedSeriesSum, firstAdjustedGamesSum float64
	var secondAdjustedSeriesSum, secondAdjustedGamesSum float64

	/* 	fmt.Println("Starting to check matches...")
	   	if len(firstTeamMatches) < len(secondTeamMatches) {
	   		fmt.Printf("%s has less matches, checking their matches (their results are first in output)\n", firstTeamInfo.Name)
	   		firstAdjustedSeriesSum, secondAdjustedSeriesSum, firstAdjustedGamesSum, secondAdjustedGamesSum, eligibleMatchCounter = headtohead_history.CalculateMatchesForTeam(firstTeamMatches, firstTeamPlayers, secondTeamPlayers)
	   	} else {
	   		fmt.Printf("%s has less matches, checking their matches (their results are first in output)\n", secondTeamInfo.Name)
	   		secondAdjustedSeriesSum, firstAdjustedSeriesSum, secondAdjustedGamesSum, firstAdjustedGamesSum, eligibleMatchCounter = headtohead_history.CalculateMatchesForTeam(secondTeamMatches, secondTeamPlayers, firstTeamPlayers)
	   	}
	*/

	fmt.Println("Starting to check matches differently...")
	if len(firstTeamMatches) < len(secondTeamMatches) {
		fmt.Printf("%s has less matches, checking their matches (their results are first in output)\n", firstTeamInfo.Name)
		firstAdjustedSeriesSum, firstAdjustedGamesSum, secondAdjustedSeriesSum, secondAdjustedGamesSum, eligibleMatchCounter = headtohead_history.CalculateScoreForMatchup(firstTeamMatches, firstTeamPlayers, secondTeamPlayers)
	} else {
		fmt.Printf("%s has less matches, checking their matches (their results are first in output)\n", secondTeamInfo.Name)
		secondAdjustedSeriesSum, secondAdjustedGamesSum, firstAdjustedSeriesSum, firstAdjustedGamesSum, eligibleMatchCounter = headtohead_history.CalculateScoreForMatchup(secondTeamMatches, secondTeamPlayers, firstTeamPlayers)
	}

	fmt.Println()
	fmt.Println("Calculation is complete!")
	fmt.Printf("Matches found: %d\n", eligibleMatchCounter)
	fmt.Printf("%s results: series - %.3f, games - %.3f\n", firstTeamInfo.Name, firstAdjustedSeriesSum, firstAdjustedGamesSum)
	fmt.Printf("%s results: series - %.3f, games - %.3f\n", secondTeamInfo.Name, secondAdjustedSeriesSum, secondAdjustedGamesSum)
}
