package kafka

import (
	"reflect"
	"testing"
)

func TestGeneratePriorityPartitions(t *testing.T) {
	partitions := []int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	tests := []struct {
		name           string
		partitions     []int32
		priorityRatios map[MessagePriority]float64
		expected       map[MessagePriority][]int32
		expectError    bool
	}{
		{
			name:       "Valid ratios summing to 1.0",
			partitions: partitions,
			priorityRatios: map[MessagePriority]float64{
				HighMessagePriority:   0.5,
				MediumMessagePriority: 0.3,
				LowMessagePriority:    0.2,
			},
			expected: map[MessagePriority][]int32{
				HighMessagePriority:   {0, 1, 2, 3, 4},
				MediumMessagePriority: {5, 6, 7},
				LowMessagePriority:    {8, 9},
			},
			expectError: false,
		},
		{
			name:       "Ratios not summing to 1.0",
			partitions: partitions,
			priorityRatios: map[MessagePriority]float64{
				HighMessagePriority:   0.5,
				MediumMessagePriority: 0.4,
			},
			expected:    nil,
			expectError: true,
		},
		{
			name:       "Single priority level",
			partitions: partitions,
			priorityRatios: map[MessagePriority]float64{
				HighMessagePriority: 1.0,
			},
			expected: map[MessagePriority][]int32{
				HighMessagePriority:   partitions,
				MediumMessagePriority: {},
				LowMessagePriority:    {},
			},
			expectError: false,
		},
		{
			name:       "Zero partitions",
			partitions: []int32{},
			priorityRatios: map[MessagePriority]float64{
				HighMessagePriority:   0.5,
				MediumMessagePriority: 0.3,
				LowMessagePriority:    0.2,
			},
			expected: map[MessagePriority][]int32{
				HighMessagePriority:   {},
				MediumMessagePriority: {},
				LowMessagePriority:    {},
			},
			expectError: false,
		},
		{
			name:           "Empty priority ratios",
			partitions:     partitions,
			priorityRatios: map[MessagePriority]float64{},
			expected:       map[MessagePriority][]int32{},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GeneratePriorityPartitions(tt.partitions, tt.priorityRatios)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !comparePriorityPartitions(result, tt.expected) {
				t.Errorf("Result does not match expected output.\nExpected: %v\nGot: %v", tt.expected, result)
			}
		})
	}
}

func comparePriorityPartitions(a, b map[MessagePriority][]int32) bool {
	if len(a) != len(b) {
		return false
	}
	for k, vA := range a {
		vB, ok := b[k]
		if !ok {
			return false
		}
		if !reflect.DeepEqual(vA, vB) {
			return false
		}
	}
	return true
}
