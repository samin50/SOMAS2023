package agent5

import (
	utils "SOMAS2023/internal/common/utils"
	voting "SOMAS2023/internal/common/voting"

	objects "SOMAS2023/internal/common/objects"

	"fmt"

	"github.com/google/uuid"
)

type ResourceAllocationParams struct {
	resourceNeed          float64 // 0-1, how much energy the agent needs, could be set to 1 - energyLevel
	resourceDemand        float64 // 0-1, how much energy the agent wants, might differ from resourceNeed
	resourceProvision     float64 // 0-1, how much energy the agent has given to reach a goal (could be either the sum of pedaling forces since last lootbox, or the latest pedalling force, or something else
	resourceAppropriation float64 // 0-1, the proportion of what the server allocates that the agent actually gets, for MVP, set to 1
}

type BikerAction int

type Iagent5 interface {
	objects.IBaseBiker

	DecideForce(direction uuid.UUID)                                // ** defines the vector you pass to the bike: [pedal, brake, turning]
	DecideJoining(pendinAgents []uuid.UUID) map[uuid.UUID]bool      // ** decide whether to accept or not accept bikers, ranks the ones
	ChangeBike() uuid.UUID                                          // ** called when biker wants to change bike, it will choose which bike to try and join
	ProposeDirection() uuid.UUID                                    // ** returns the id of the desired lootbox based on internal strategy
	FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap // ** stage 3 of direction voting
	DecideAllocation() voting.IdVoteMap                             // ** decide the allocation parameters
}

type agent5 struct {
	*objects.BaseBiker
	test string
}

func (b *agent5) DecideAllocation() voting.IdVoteMap {
	fmt.Printf("test \n")
	return calculateResourceAllocation(b.GetGameState(), b, "equal")

}

// func (b *agent5) DecideCurrentBikers(bikers []uuid.UUID) map[uuid.UUID]bool {

// }

// func (b *agent5) FinalDirectionVote(proposals []uuid.UUID) voting.LootboxVoteMap {

// }

// this function is going to be called by the server to instantiate bikers in the MVP
func GetIagent5(totColours utils.Colour, bikeId uuid.UUID) Iagent5 {
	return &agent5{
		BaseBiker: objects.GetBaseBiker(totColours, bikeId),
		test:      "hello",
	}
}

// this function will be used by GetTeamAgent to get the ref to the BaseBiker
func Getagent5(totColours utils.Colour, bikeId uuid.UUID) *agent5 {
	return &agent5{
		BaseBiker: objects.GetBaseBiker(totColours, bikeId),
		test:      "hello",
	}
}
