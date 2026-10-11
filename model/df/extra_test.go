package df

import "testing"

func TestConsecutivePattern(t *testing.T) {
	cases := map[string][]int{
		"none": {1, 3, 10, 15, 20},
		"2":    {1, 2, 10, 15, 20},
		"3":    {1, 2, 3, 10, 20},
		"22":   {1, 2, 10, 11, 20},
		"32":   {1, 2, 4, 5, 6},
		"5":    {1, 2, 3, 4, 5},
		"222":  {1, 2, 10, 11, 20, 21},
		"33":   {1, 2, 3, 10, 11, 12},
		"42":   {1, 2, 3, 4, 10, 11},
	}
	for want, balls := range cases {
		f := NewFeature(balls, 49)
		if got := f.ConsecutivePattern(); got != want {
			t.Errorf("%v: got %s want %s", balls, got, want)
		}
	}
}

func TestExtraFeatures(t *testing.T) {
	f := NewFeature([]int{3, 7, 15, 17, 30}, 39)
	if f.AC() != 6 || f.Span() != 27 || f.CountAtLeast(20) != 1 || f.TailSum() != 22 ||
		f.SameTailGroups() != 1 || f.MaxGap() != 13 || f.MinGap() != 2 || f.Mod3Counts() != [3]int{3, 1, 1} {
		t.Errorf("unexpected features: ac=%d span=%d big=%d tail=%d same=%d max=%d min=%d mod3=%v",
			f.AC(), f.Span(), f.CountAtLeast(20), f.TailSum(), f.SameTailGroups(), f.MaxGap(), f.MinGap(), f.Mod3Counts())
	}
}
