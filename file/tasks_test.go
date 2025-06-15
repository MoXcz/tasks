package file

import (
	"testing"
	"time"
)

func Test_newTask(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		record  []string
		want    Task
		wantErr bool
	}{
		{
			name:   "valid task",
			record: []string{"0", "do the dishes", "Sun, 15 Jun 2025 04:07:04 UTC", "true"},
			want: Task{
				ID:         0,
				Task:       "do the dishes",
				CreatedAt:  time.Date(2025, 6, 15, 4, 7, 4, 0, time.UTC),
				IsComplete: true,
			},
		},
		{
			name:    "invalid ID",
			record:  []string{"water", "do the dishes", "Sun, 15 Jun 2025 04:07:04 UTC", "false"},
			wantErr: true,
		},
		{
			name:    "invalid time format",
			record:  []string{"0", "do the dishes", "15 Jun 2025 04:07:04 UTC", "false"},
			wantErr: true,
		},
		{
			name:   "invalid isComplete value defaults to false",
			record: []string{"0", "do the dishes", "Sun, 15 Jun 2025 04:07:04 UTC", "42"},
			want: Task{
				ID:         0,
				Task:       "do the dishes",
				CreatedAt:  time.Date(2025, 6, 15, 4, 7, 4, 0, time.UTC),
				IsComplete: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTask(tt.record)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("newTask() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("newTask() succeeded unexpectedly")
			}
			if tt.want != got {
				t.Errorf("newTask() = %v, want %v", got, tt.want)
			}
		})
	}
}
