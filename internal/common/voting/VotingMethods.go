package voting

import (
	"math"
	"sort"

	"github.com/google/uuid"
)

// Auxiliary Structure for Sorting
type kv struct {
	Key   uuid.UUID
	Value float64
}

func Plurality(voteMap map[uuid.UUID]map[uuid.UUID]float64, voteWeight map[uuid.UUID]float64) uuid.UUID {
	/*
		Plurality:
			Each voter selects one candidate and the candidate with the most first-placed votes is the winner.
	*/

	//initialise the votes with weights
	var voteList []map[uuid.UUID]float64
	for agent, votes := range voteMap {
		weight := voteWeight[agent]
		weightedvotes := make(map[uuid.UUID]float64)
		for key, value := range votes {
			weightedvotes[key] = value * weight
		}
		voteList = append(voteList, weightedvotes)
	}

	// start
	voteCount := make(map[uuid.UUID]float64)
	var winner uuid.UUID

	for _, preference := range voteList {
		var maxPreference float64
		var firstLootBoxChoice uuid.UUID
		for lootBox, value := range preference {
			if value > maxPreference {
				firstLootBoxChoice = lootBox
				maxPreference = value
			}
		}
		voteCount[firstLootBoxChoice] += maxPreference
	}

	// final step: we need to find the winner with highest count number in map.
	var maxVotes float64

	for lootBox, votes := range voteCount {
		if votes > maxVotes {
			maxVotes = votes
			winner = lootBox
		}
	}

	// return the final winner
	return winner
}

func Runoff(voteMap map[uuid.UUID]map[uuid.UUID]float64, voteWeight map[uuid.UUID]float64) uuid.UUID {
	/*
		Runoff:
			1st round: 	each voter selects one candidate, and the two candidates with most first-placed votes are identified.
						If either already has a majority, this candidate is declared the winner.
			2nd round: 	each voter selects one candidate, the candidate with most votes now is the winner.
	*/
	//initialise the votes with weights
	var voteList []map[uuid.UUID]float64
	for agent, votes := range voteMap {
		weight := voteWeight[agent]
		weightedvotes := make(map[uuid.UUID]float64)
		for key, value := range votes {
			weightedvotes[key] = value * weight
		}
		voteList = append(voteList, weightedvotes)
	}

	// start
	voteCount := make(map[uuid.UUID]float64)
	var winner uuid.UUID

	// ----- first round -----
	// find the count number of each lootbox
	for _, preference := range voteList {
		var maxPreference float64
		var firstLootBoxChoice uuid.UUID
		for lootBox, value := range preference {
			if value > maxPreference {
				firstLootBoxChoice = lootBox
				maxPreference = value
			}
		}
		voteCount[firstLootBoxChoice] += maxPreference
	}

	// find the two candidates with most first-placed votes
	var maxVotes1, maxVotes2 float64
	var winner1, winner2 uuid.UUID
	for lootBox, votes := range voteCount {
		if votes > maxVotes1 {
			winner2 = winner1
			maxVotes2 = maxVotes1
			winner1 = lootBox
			maxVotes1 = votes
		} else if votes > maxVotes2 {
			winner2 = lootBox
			maxVotes2 = votes
		}
	}

	// check if either already has a majority or we need the second round
	if maxVotes1 >= (maxVotes2 * 2) {
		// return the majority lootbox
		return winner1
	} else {
		// ----- second round -----
		voteCount := make(map[uuid.UUID]float64)
		for _, preference := range voteList {
			if preference[winner1] > preference[winner2] {
				voteCount[winner1] += preference[winner1]
			} else {
				voteCount[winner2] += preference[winner2]
			}
		}
		if voteCount[winner1] > voteCount[winner2] {
			winner = winner1
		} else {
			winner = winner2
		}
	}

	return winner
}

func BordaCount(voteMap map[uuid.UUID]map[uuid.UUID]float64, voteWeight map[uuid.UUID]float64) uuid.UUID {
	/*
		BordaCount:
			Each voter rank order all the candidates. With n candidates being ranked k scores (n-k)+1 Borda points.
			The candidate with the highest Borda Score is the winner
	*/
	//initialise the votes with weights
	voteListMap := make(map[uuid.UUID]map[uuid.UUID]float64)
	for agent, votes := range voteMap {
		weight := voteWeight[agent]
		weightedvotes := make(map[uuid.UUID]float64)
		for key, value := range votes {
			weightedvotes[key] = value * weight
		}
		voteListMap[agent] = weightedvotes
	}

	// start
	voteCount := make(map[uuid.UUID]float64)
	var winner uuid.UUID

	// initialise the map with all candidates
	for _, preference := range voteListMap {
		for key := range preference {
			voteCount[key] = 0
		}
	}

	// covert the unodered map into ordered list
	ss := make(map[uuid.UUID][]kv)
	for agent, preference := range voteListMap {
		var s []kv
		for k, v := range preference {
			// ignore the lootbox if value is 0
			if v != 0 {
				s = append(s, kv{k, v})
			}
		}
		// sort the list using preference value of each lootbox
		sort.Slice(s, func(i, j int) bool {
			// in the order from large to small
			return s[i].Value > s[j].Value
		})
		ss[agent] = s
	}

	// calculate the Borda score for each candidates
	for agent, sortedList := range ss {
		usedKeys := make(map[uuid.UUID]bool)
		for i, kv := range sortedList {
			score := float64(len(voteCount)) - float64(i) + 1
			voteCount[kv.Key] += score * voteWeight[agent]
			usedKeys[kv.Key] = true
		}

		// points shared if not explicity ranked
		remainingKeyNumber := float64(len(voteCount)) - float64(len(sortedList))
		remainingScore := (1 + remainingKeyNumber) * remainingKeyNumber / 2
		for key := range voteCount {
			if !usedKeys[key] {
				voteCount[key] += remainingScore / remainingKeyNumber
			}
		}
	}

	// find the winner with highest score
	var maxScore float64
	for key, value := range voteCount {
		if value > maxScore {
			winner = key
			maxScore = value
		}
	}

	return winner
}

func InstantRunoff(voteMap map[uuid.UUID]map[uuid.UUID]float64, voteWeight map[uuid.UUID]float64) uuid.UUID {
	/*
		InstantRunoff:
			Each voter rank orders all candidates, and the candidate with the least number of first-place votes is eliminate.
			This is repeated until only one candidate remains
	*/
	//initialise the votes with weights
	var voteList []map[uuid.UUID]float64
	for agent, votes := range voteMap {
		weight := voteWeight[agent]
		weightedvotes := make(map[uuid.UUID]float64)
		for key, value := range votes {
			weightedvotes[key] = value * weight
		}
		voteList = append(voteList, weightedvotes)
	}

	// start
	voteCount := make(map[uuid.UUID]float64)
	eliminateVote := make(map[uuid.UUID]bool)
	var winner uuid.UUID

	// initialise the map with all candidates
	for _, preference := range voteList {
		for key := range preference {
			voteCount[key] = 0
		}
	}

	// loop to eliminate the least number of first-place votes
	for len(voteCount) > 1 {
		// reset map with value = 0
		for key := range voteCount {
			voteCount[key] = 0
		}

		// count the number of first-place votes for each lootbox
		for _, preference := range voteList {
			var maxScore float64
			var firstLootBoxChoice uuid.UUID
			for key, value := range preference {
				if (value > maxScore) && !eliminateVote[key] {
					maxScore = value
					firstLootBoxChoice = key
				}
			}
			voteCount[firstLootBoxChoice] += maxScore
		}

		// eliminate the lootbox with least votes
		var minVotes float64 = math.MaxFloat64
		var candidateToEliminate uuid.UUID
		for key, value := range voteCount {
			if value < minVotes {
				minVotes = value
				candidateToEliminate = key
			}
		}
		eliminateVote[candidateToEliminate] = true
		delete(voteCount, candidateToEliminate)
	}

	// get the final winner
	for key := range voteCount {
		winner = key
	}

	return winner
}

func Approval(voteMap map[uuid.UUID]map[uuid.UUID]float64, voteWeight map[uuid.UUID]float64) uuid.UUID {
	/*
		Approval:
			A ballot represents not a linear rank order of decreasing preference,
			but rather represents the set of candidates who are 'equally acceptable' to the voter
	*/
	//initialise the votes with weights
	var voteList []map[uuid.UUID]float64
	for agent, votes := range voteMap {
		weight := voteWeight[agent]
		weightedvotes := make(map[uuid.UUID]float64)
		for key, value := range votes {
			weightedvotes[key] = value * weight
		}
		voteList = append(voteList, weightedvotes)
	}

	// start
	voteCount := make(map[uuid.UUID]float64)
	var winner uuid.UUID

	for _, preference := range voteList {
		for key, value := range preference {
			if value > 0 {
				voteCount[key] += value
			}
		}
	}

	// find the lootbox with max score
	var maxVotes float64

	for lootBox, votes := range voteCount {
		if votes > maxVotes {
			maxVotes = votes
			winner = lootBox
		}
	}

	return winner
}

func CopelandScoring(voteMap map[uuid.UUID]map[uuid.UUID]float64, voteWeight map[uuid.UUID]float64) uuid.UUID {
	/*
		CopelandScoring:
			Each voter submits a ballot with a linear rank order.
			A win-loss record, the Copeland Score, is calculated for each candidate.
	*/
	//initialise the votes with weights
	voteListMap := make(map[uuid.UUID]map[uuid.UUID]float64)
	for agent, votes := range voteMap {
		weight := voteWeight[agent]
		weightedvotes := make(map[uuid.UUID]float64)
		for key, value := range votes {
			weightedvotes[key] = value * weight
		}
		voteListMap[agent] = weightedvotes
	}

	// start
	// the map to store the winning score for each lootbox
	scores := make(map[uuid.UUID]float64)

	// iterate the voting
	for agent, vote := range voteListMap {
		for candidate1, score1 := range vote {
			for candidate2, score2 := range vote {
				// do not compare with itself
				if candidate1 == candidate2 {
					continue
				}

				// update the score of each lootbox
				if score1 > score2 {
					scores[candidate1] += voteWeight[agent]
					scores[candidate2] -= voteWeight[agent]
				} else if score1 < score2 {
					scores[candidate1] -= voteWeight[agent]
					scores[candidate2] += voteWeight[agent]
				}
			}
		}
	}

	// find the lootbox with the highest score
	var maxScore float64
	var maxCandidate uuid.UUID
	for candidate, score := range scores {
		if score > maxScore || maxCandidate == uuid.Nil {
			maxScore = score
			maxCandidate = candidate
		}
	}

	return maxCandidate
}
