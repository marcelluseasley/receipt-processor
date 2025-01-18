package pointsprocessor

import "testing"

func Test_processRetailerName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Target should be 6",
			args: args{
				name: "Target",
			},
			want: 6,
		},
		{
			name: "M&M Corner Market should be 14",
			args: args{
				name: "M&M Corner Market",
			},
			want: 14,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processRetailerName(tt.args.name); got != tt.want {
				t.Errorf("retailerName() = %v, want %v", got, tt.want)
			}
		})
	}
}
