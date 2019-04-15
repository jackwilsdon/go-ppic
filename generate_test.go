package ppic_test

import (
	"github.com/jackwilsdon/go-ppic/ppictest"
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
		expected [8]string
	}{
		{
			text: "jackwilsdon",
			expected: [8]string{
				"# # ####",
				"# ## ###",
				"    #   ",
				"# #   ##",
				"  #  ## ",
				"     # #",
				"##   #  ",
				"#    ###",
			},
		},
		{
			text: "jackwilsdon",
			mX:   true,
			expected: [8]string{
				"# #  # #",
				"########",
				"# #### #",
				" ###### ",
				"        ",
				"#      #",
				"# #  # #",
				"  ####  ",
			},
		},
		{
			text: "jackwilsdon",
			mY:   true,
			expected: [8]string{
				"# # ####",
				"# ## ###",
				"    #   ",
				"# #   ##",
				"# #   ##",
				"    #   ",
				"# ## ###",
				"# # ####",
			},
		},
		{
			text: "jackwilsdon",
			mX:   true,
			mY:   true,
			expected: [8]string{
				"# #  # #",
				"########",
				"# #### #",
				" ###### ",
				" ###### ",
				"# #### #",
				"########",
				"# #  # #",
			},
		},
	}

	for i, c := range cases {
		grid := ppic.Generate(c.text, c.mX, c.mY)

		err := ppictest.Compare(grid, c.expected)

		if err != nil {
			t.Errorf("%s for case %d", err, i)
		}
	}
}
