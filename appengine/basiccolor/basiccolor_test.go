package basiccolor

import (
	"fmt"
	"testing"
)

func TestClosest_numerical(t *testing.T) {
	type testcase struct {
		h, s, l float64
		expect  Color
	}

	testcases := []testcase{
		{360, 0, 0.93, White},
		{360, 0, 0.95, White},
		{270, 0.1, 0.93, White},
		{270, 0.1, 0.95, White},
		{270, 0.11, 0.938, White},
		{270, 0.15, 0.97, White},
		{180, 0.15, 0.971, White},
		{180, 0.25, 0.984, White},
		{180, 0.4, 1.00, White},
		{180, 0.75, 0.99, White},
		{180, 1, 0.99, White},
		{180, 1, 1, White},
		{0, 0.05, 0.9, White},
		{0, 0, 0.9299, White},
		{0, 0, 0.9098039215686274, White},

		{0, 0, 0.05, Black},
		{0, 0.08, 0.2, Black},
		{270, 0.2, 0.08, Black},
		{270, 0.25, 0.03, Black},
		{90, 0.5, 0.0499999, Black},
		{90, 0.75, 0, Black},
		{360, 1, 0.039999, Black},
		{0, 0.079, 0.201, Black},
		{307, 0.1, 0.2, Black},

		{0, 0.06, 0.3, Gray},
		{0, 0.03, 0.5, Gray},
		{0, 0.04, 0.6, Gray},
		{0, 0.06, 0.7, Gray},
		{0, 0.07, 0.8, Gray},

		{343, 0.4, 0.4, Pink},
		{360, 0.8, 0.8, Pink},
		{0, 0.9, 0.8, Pink},
		{7, .84, 0.48, Red},
		{13.9, 0.9, 0.8, Orange},
		{14, 0.4, 0.4, Brown},
		{14, 0.4, 0.34, Brown},
		{14, 0.4, 0.8, Orange},
		{17, 0.5, 0.4, Brown},
		{17, 0.5, 0.41, Orange},
		{30, 0.5, 0.4, Brown},
		{35, 0.5, 0.7, Orange},

		{46, 0.2, 0.7, Yellow},
		{68.9, 0.8, 0.8, Yellow},

		{70, 0.8, 0.8, Green},

		{175, 0.8, 0.8, Blue},
		{210, 0.15, 0.59, Gray}, // slate gray

		{250, 0.8, 0.8, Violet},
		{287, 0.4, 0.4, Violet},
		{287, 0.8, 0.8, Pink},
		{306.9, 0.3, 0.3, Violet},
		{306.9, 0.7, 0.7, Pink},
	}

	testName := func(tt testcase) string {
		return fmt.Sprintf("%f_%f_%f_%s", tt.h, tt.s, tt.l, tt.expect)
	}

	for _, tt := range testcases {
		t.Run(testName(tt), func(t *testing.T) {
			got := Closest(HSL{tt.h, tt.s, tt.l, 1})
			if tt.expect != got {
				t.Errorf("wrong color: expected: %s, got: %s", tt.expect, got)
			}
		})
	}
}

func TestClosest_handpicked(t *testing.T) {
	type testcase struct {
		h, s, l float64
		expect  Color
		desc    string
	}

	type undesirableResultTestcase struct {
		h, s, l float64
		desired Color
		expect  Color
		desc    string
	}

	testcases := []testcase{
		{46, .96, .49, Yellow, "Buttercup"},
		{0, 0, .99, White, "XX white"},
		{0, 0, 0, Black, "XX black"},
		{19, .96, .54, Orange, "Autumn Orange flame"},
		{7, .71, .17, Brown, "All Melody wood"},
		{342, .97, .86, Pink, "Wolfgang Amadeus Phoenix pink"},
		{347, .97, .86, Pink, "Wolfgang Amadeus Phoenix pink [hand modified]"},
		{179, .73, .53, Blue, "Gorillaz The Now Now cyan"},
		{289, .19, .45, Violet, "alt-J This is All Yours violet"},
		{213, .46, .38, Blue, "Josef Salvat In Your Prime background"},
		{192, .05, .60, Gray, "Dan Deacon Mystic Familiar gray"},
		{43, .75, .87, Yellow, "How To Be a Human Being border"},
		{360, 1, .40, Red, "Feeling Free block letters"},
		{175, .09, .27, Gray, "Feeling Free background"},
		{53, .58, .69, Yellow, "Holy Fire sky yellow"},
		{355, .64, .41, Red, "Currents Tame Impala midstreak stark"},
		{188, .36, .70, Blue, "What Went Down background"}, // looks green wholly, but really is blue
		{187, .71, .77, Blue, "Dan Deacon bicycle soundtrack"},
		{240, .32, .13, Blue, "Marks Beauvois"}, // also Violet would work
		{28, .85, .44, Orange, "Champyons"},
		{173, .46, .18, Green, "Holy fire special edition sky green"},
		{175, .58, .19, Blue, "darkish cyan blue [handwritten]"},
		{0, 0, .73, Gray, "Bone Digger gray"},
		{137, .21, .13, Green, "Tourist Wash clothes green"},
		{0, 0, .06, Black, "Cross of Changes background"},
		{341, .83, .53, Pink, "French 79 pink band"},
		{241, .96, .09, Blue, "Cupa Cupa blue"},
		{0, 0, .95, White, "Tailwhip white"},
		{41, .63, .49, Orange, "All India Radio Mountain Top yellow"},
		{17, .92, .69, Orange, "Exits remix"},
		{45, .68, .82, Yellow, "Foals Antidotes"},
		{248, .35, .12, Violet, "Glass Animals EP"},
		{73, .18, .29, Green, "Costello Music green dress"},
		{7, .72, .58, Red, "Costello Music red dress"},
		{344, .69, .26, Pink, "Everything Not Saved Will Be Lost Pt. 1"},
		{357, .70, .32, Red, "Animal Collective Sung Tongs"},
	}

	undesirables := []undesirableResultTestcase{
		{182, .20, .71, Green, Blue, "Part Time Glory background"},
	}

	for _, tt := range testcases {
		t.Run(tt.desc, func(t *testing.T) {
			got := Closest(HSL{tt.h, tt.s, tt.l, 1})
			if tt.expect != got {
				t.Errorf("wrong color: expected: %s, got: %s", tt.expect, got)
			}
		})
	}

	for _, tt := range undesirables {
		t.Run(tt.desc, func(t *testing.T) {
			got := Closest(HSL{tt.h, tt.s, tt.l, 1})
			if tt.desired == got {
				t.Errorf("now produces desired color: %s", tt.desired)
				return
			}
			if tt.expect != got {
				t.Errorf("wrong color: expected: %s, got: %s", tt.expect, got)
				return
			}
		})
	}
}

func TestClosest_fuzz(t *testing.T) {
	// TODO
}
