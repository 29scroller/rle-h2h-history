package main

import (
	"fmt"
	"headtohead_history"
)

func main() {
	var teamFound bool
	var firstTeamName, secondTeamName string
	var firstTeamInfo, secondTeamInfo headtohead_history.Team

	for !teamFound {
		firstTeamInfo, teamFound = headtohead_history.UserEnteringTeam("first", firstTeamName)
	}
	fmt.Printf("Finding active players of %s by ID\n", firstTeamInfo.Name)
	firstTeamPlayers := headtohead_history.FindActivePlayersByTeamID(firstTeamInfo.Id)

	teamFound = false
	for !teamFound {
		secondTeamInfo, teamFound = headtohead_history.UserEnteringTeam("second", secondTeamName)
	}
	fmt.Printf("Finding active players of %s by ID\n", secondTeamInfo.Name)
	secondTeamPlayers := headtohead_history.FindActivePlayersByTeamID(secondTeamInfo.Id)

	var firstTeamMatches, secondTeamMatches []headtohead_history.Match
	headtohead_history.LoadAndAppendTeamMatchesFromFiles(&firstTeamMatches, firstTeamPlayers, firstTeamInfo.Name)
	headtohead_history.LoadAndAppendTeamMatchesFromFiles(&secondTeamMatches, secondTeamPlayers, secondTeamInfo.Name)

	fmt.Println("Deleting duplicates...")
	firstTeamMatches = headtohead_history.RemoveDuplicatesinMatches(firstTeamMatches, firstTeamInfo.Name)
	secondTeamMatches = headtohead_history.RemoveDuplicatesinMatches(secondTeamMatches, secondTeamInfo.Name)
	fmt.Println()

	eligibleMatchCounter := 0
	var firstAdjustedSeriesSum, firstAdjustedGamesSum float64
	var secondAdjustedSeriesSum, secondAdjustedGamesSum float64

	fmt.Println("Starting to check matches...")
	if len(firstTeamMatches) < len(secondTeamMatches) {
		fmt.Printf("%s has less matches, checking their matches (their results are first in output)\n", firstTeamInfo.Name)
		firstAdjustedSeriesSum, secondAdjustedSeriesSum, firstAdjustedGamesSum, secondAdjustedGamesSum, eligibleMatchCounter = headtohead_history.CalculateMatchesForTeam(firstTeamMatches, firstTeamPlayers, secondTeamPlayers)
	} else {
		fmt.Printf("%s has less matches, checking their matches (their results are first in output)\n", secondTeamInfo.Name)
		secondAdjustedSeriesSum, firstAdjustedSeriesSum, secondAdjustedGamesSum, firstAdjustedGamesSum, eligibleMatchCounter = headtohead_history.CalculateMatchesForTeam(secondTeamMatches, secondTeamPlayers, firstTeamPlayers)
	}
	fmt.Println()
	fmt.Println("Calculation is complete!")
	fmt.Printf("Matches found: %d\n", eligibleMatchCounter)
	fmt.Printf("%s results: series - %.2f, games - %.2f\n", firstTeamInfo.Name, firstAdjustedSeriesSum, firstAdjustedGamesSum)
	fmt.Printf("%s results: series - %.2f, games - %.2f\n", secondTeamInfo.Name, secondAdjustedSeriesSum, secondAdjustedGamesSum)
}
