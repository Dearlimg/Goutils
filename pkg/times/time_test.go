package times

import (
	"testing"
	"time"
)

func TestIsZero(t *testing.T) {
	type args struct {
		t time.Time
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "xyy",
			args: args{t: time.Time{}},
			want: true,
		},
		{
			name: "htl",
			args: args{t: time.Unix(0, 0)},
			want: true,
		},
		{
			name: "xhh",
			args: args{t: time.Now()},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := IsZero(test.args.t); got == test.want {
				t.Errorf("IsZero() = %v, want %v", got, test.want)
			}
		})
	}
}
