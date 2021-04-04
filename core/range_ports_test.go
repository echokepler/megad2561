package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRangePortsModule_GetRange(t *testing.T) {
	t.Parallel()

	t.Run("Should be return correct range", func(t *testing.T) {
		testCases := []struct {
			condition RangePortsModule
			expected  []int
		}{
			{
				condition: XP1,
				expected:  []int{1, 14},
			},
			{
				condition: XP2,
				expected:  []int{15, 29},
			},
			{
				condition: XT2,
				expected:  []int{30, 35},
			},
			{
				condition: XP5andXP6,
				expected:  []int{36, 37},
			},
			{
				condition: RangePortsModule("45"),
				expected:  []int{45, 45},
			},
		}

		for _, tCase := range testCases {
			pRange, err := tCase.condition.GetRange()
			if err != nil {
				t.Error(err)
			}

			assert.EqualValues(t, tCase.expected, pRange)
		}
	})
}
