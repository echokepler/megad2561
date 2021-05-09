package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoteCommander_Parse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		expected []ICommand
		cmd      string
	}{
		{
			name: "Should correct parse when all cases",
			expected: []ICommand{
				&BaseCommand{
					Command: Command{
						Value: 8,
					},
					TargetPort: 8,
				},
				&PauseCommand{
					Command: Command{
						Value: 5,
					},
				},
				&RepeatCommand{
					Command: Command{
						Value: 4,
					},
				},
				&GlobalCommand{
					Command: Command{
						Value: 1,
					},
				},
			},
			cmd: "8:8;p5;r4;g1",
		},
		{
			name: "Should skip parse unknown commands",
			expected: []ICommand{
				&BaseCommand{
					Command: Command{
						Value: 2,
					},
					TargetPort: 32,
				},
			},
			cmd: "unknown3;32:2",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			commander, err := NewRemoteCommander(tc.cmd)
			if err != nil {
				t.Error(err)
			}

			assert.EqualValues(t, tc.expected, commander.commands)
		})
	}
}
