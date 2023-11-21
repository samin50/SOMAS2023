
import (
	"...\internal\common\objects\BaseBiker.go"
)


func getReputation(agent_id uuid.UUID) float64 
{
	// Returns the reputation of a biker as a float64
	return 0.0
}

func getBikeIds() []uuid.UUID 
{
	// Returns all bike ids in the game as a list of uuid.UUIDs
	bike_ids := []uuid.UUID {.....}
	return bike_ids
}

func getBikerListFromBike(bike_id uuid.UUID) []uuid.UUID 
{
	// Returns all biker ids on a given bike as a list of uuid.UUIDs
	biker_ids := []uuid.UUID {.....}
	return biker_ids
}

func calculateAverageReputationOfEachBike(bike_ids []uuid.UUID) map[uuid.UUID]float64 
{
	// Returns a map of bike ids to average reputation of bikers on that bike
	avg_rep := make(map[uuid.UUID]float64)

	// For each bike in the bike_ids, find the average reputation of the bikers on that bike.
	for bike_id := range bike_ids
	{
		// Get list of bikers on bike
		biker_ids := getBikerListFromBike(bike_id)
		
		// Calculate average reputation of bikers on bike
		sum_reps := 0.0
		for biker_id := range biker_ids
		{
			sum_reps += getReputation(biker_id)
		}
		
		avg_rep[bike_id] = sum_reps / len(biker_ids)
	}

	return avg_rep
}

func (bb *BaseBiker) DecideAction() BikerAction {
	// Get Bike IDs
	bike_ids := []uuid.UUID {getBikeIds()}

	// Get number of bikes
	num_of_bikes := len(bike_ids)

	// For each bike, calculate the average reputation of the bikers on that bike.
	avg_reps_per_bike := calculateAverageReputationOfBike(bike_ids)	

	// Calculate the average reputation of all bikes
	sum_reps := 0.0

	for _, value := range avg_reps_per_bike 
	{
		sum_reps += value
	}

	averageReputation = sum_reps / num_of_bikes

	// Calculate our bike reputation
	ourBikeReputation := calculateAverageReputationOfBike(bb.megaBikeId)

	// If our bike reputation is less than the average reputation of all bikes, then we should change bike.
	if (ourBikeReputation < averageReputation) 
	{
		return ChangeBike;
	}
	else
	{
		return Pedal;
	}
}

func (bb *BaseBiker) ChangeBike() uuid.UUID 
{
	// For each bike, calculate the average reputation of the bikers on that bike.
	avg_reps_per_bike := calculateAverageReputationOfBike(bike_ids)	

	// Change to bike with highest average reputation:
	maxValue := float64(-1)
	maxKey := uuid.Nil
	
	for key, value := range avg_reps_per_bike {
		if value > maxValue {
			maxValue = value
			maxKey = key
		}
	}

	return maxKey
}