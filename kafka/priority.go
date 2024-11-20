package kafka

import "fmt"

type MessagePriority string

const (
	HighMessagePriority   MessagePriority = "high"
	MediumMessagePriority MessagePriority = "medium"
	LowMessagePriority    MessagePriority = "low"
)

var SortedPriorities = []MessagePriority{HighMessagePriority, MediumMessagePriority, LowMessagePriority}

var PrioritySet = map[MessagePriority]struct{}{
	HighMessagePriority: {}, MediumMessagePriority: {}, LowMessagePriority: {},
}

var PriorityToRatio = map[MessagePriority]float64{
	HighMessagePriority:   0.5,
	MediumMessagePriority: 0.3,
	LowMessagePriority:    0.2,
}

func GeneratePriorityPartitions(partitions []int32, priorityToRatio map[MessagePriority]float64) (map[MessagePriority][]int32, error) {
	totalPartitions := len(partitions)
	priorityPartitions := make(map[MessagePriority][]int32)
	partitionIndex := 0

	totalRatio := 0.0
	for _, ratio := range priorityToRatio {
		totalRatio += ratio
	}
	if totalRatio < 0.9999 || totalRatio > 1.0001 {
		return nil, fmt.Errorf("total of priority ratios must be 1.0, got %f", totalRatio)
	}

	cumulativePartitions := 0
	for i, priority := range SortedPriorities {
		ratio := priorityToRatio[priority]
		numPartitions := int(float64(totalPartitions) * ratio)
		cumulativePartitions += numPartitions

		if i == len(SortedPriorities)-1 {
			numPartitions += totalPartitions - cumulativePartitions
		}

		assignedPartitions := partitions[partitionIndex : partitionIndex+numPartitions]
		priorityPartitions[priority] = assignedPartitions
		partitionIndex += numPartitions
	}

	return priorityPartitions, nil
}
