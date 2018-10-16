package ppic_test

import (
	"reflect"
	"testing"

	"github.com/jackwilsdon/go-ppic"
)

func BenchmarkGenerate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ppic.Generate("jackwilsdon", false, false)
	}
}

func TestGenerate(t *testing.T) {
	cases := []struct {
		text     string
		mX       bool
		mY       bool
		expected [8][8]bool
	}{
		{
			text: "jackwilsdon",
			expected: [8][8]bool{
				{true, false, true, false, true, true, true, true},
				{true, false, true, true, false, true, true, true},
				{false, false, false, false, true, false, false, false},
				{true, false, true, false, false, false, true, true},
				{false, false, true, false, false, true, true, false},
				{false, false, false, false, false, true, false, true},
				{true, true, false, false, false, true, false, false},
				{true, false, false, false, false, true, true, true},
			},
		},
		{
			text: "jackwilsdon",
			mX:   true,
			expected: [8][8]bool{
				{true, false, true, false, false, true, false, true},
				{true, true, true, true, true, true, true, true},
				{true, false, true, true, true, true, false, true},
				{false, true, true, true, true, true, true, false},
				{false, false, false, false, false, false, false, false},
				{true, false, false, false, false, false, false, true},
				{true, false, true, false, false, true, false, true},
				{false, false, true, true, true, true, false, false},
			},
		},
		{
			text: "jackwilsdon",
			mY:   true,
			expected: [8][8]bool{
				{true, false, true, false, true, true, true, true},
				{true, false, true, true, false, true, true, true},
				{false, false, false, false, true, false, false, false},
				{true, false, true, false, false, false, true, true},
				{true, false, true, false, false, false, true, true},
				{false, false, false, false, true, false, false, false},
				{true, false, true, true, false, true, true, true},
				{true, false, true, false, true, true, true, true},
			},
		},
		{
			text: "jackwilsdon",
			mX:   true,
			mY:   true,
			expected: [8][8]bool{
				{true, false, true, false, false, true, false, true},
				{true, true, true, true, true, true, true, true},
				{true, false, true, true, true, true, false, true},
				{false, true, true, true, true, true, true, false},
				{false, true, true, true, true, true, true, false},
				{true, false, true, true, true, true, false, true},
				{true, true, true, true, true, true, true, true},
				{true, false, true, false, false, true, false, true},
			},
		},
	}

	for i, c := range cases {
		grid := ppic.Generate(c.text, c.mX, c.mY)

		if !reflect.DeepEqual(grid, c.expected) {
			t.Errorf("generated grid does not match test data for case %d", i)
		}
	}
}
