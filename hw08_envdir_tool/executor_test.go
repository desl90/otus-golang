package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	testCases := []struct {
		name         string
		command      []string
		env          Environment
		expectedCode int
		callable     func(t *testing.T)
	}{
		{
			name:         "no args",
			command:      nil,
			env:          nil,
			expectedCode: 1,
			callable:     nil,
		}, {
			name:         "empty cmd",
			command:      []string{},
			env:          nil,
			expectedCode: 1,
			callable:     nil,
		}, {
			name:         "invalid cmd",
			command:      []string{"program"},
			env:          nil,
			expectedCode: 1,
			callable:     nil,
		}, {
			name:         "negative",
			command:      []string{"test", "git", "-ls"},
			env:          nil,
			expectedCode: 1,
			callable:     nil,
		}, {
			name:    "positive",
			command: []string{"COMMAND"},
			env: Environment{
				"command": EnvValue{
					Value: "run",
				},
			},
			expectedCode: 1,
			callable: func(t *testing.T) {
				t.Helper()

				expectedEnv := map[string]string{
					"command": "run",
				}

				currentEnv := make(map[string]string, len(expectedEnv))

				for key := range expectedEnv {
					val, exists := os.LookupEnv(key)

					if exists {
						currentEnv[key] = val
					}
				}

				require.Equal(t, expectedEnv, currentEnv)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code := RunCmd(tc.command, tc.env)

			require.Equal(t, tc.expectedCode, code)

			if tc.callable != nil {
				tc.callable(t)
			}
		})
	}
}
