package pie

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/elliotchance/testify-stats/assert"
)

var carsContainsTests = []struct {
	ss       cars
	contains car
	expected bool
}{
	{nil, car{"a", "green"}, false},
	{nil, car{}, false},
	{cars{car{"a", "green"}, car{"b", "blue"}, car{"c", "gray"}}, car{"a", "green"}, true},
	{cars{car{"a", "green"}, car{"b", "blue"}, car{"c", "gray"}}, car{"b", "blue"}, true},
	{cars{car{"a", "green"}, car{"b", "blue"}, car{"c", "gray"}}, car{"c", "gray"}, true},
	{cars{car{"a", "green"}, car{"b", "blue"}, car{"c", "gray"}}, car{"A", ""}, false},
	{cars{car{"a", "green"}, car{"b", "blue"}, car{"c", "gray"}}, car{}, false},
	{cars{car{"a", "green"}, car{"b", "blue"}, car{"c", "gray"}}, car{"d", ""}, false},
	{cars{car{"a", "green"}, car{}, car{"c", "gray"}}, car{}, true},
}

func TestCars_Contains(t *testing.T) {
	for _, test := range carsContainsTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.expected, test.ss.Contains(test.contains))
		})
	}
}

var carsFilterTests = []struct {
	ss                cars
	condition         func(car) bool
	expectedFilter    cars
	expectedFilterNot cars
	expectedMap       cars
}{
	{
		nil,
		func(s car) bool {
			return s.Name == ""
		},
		nil,
		nil,
		nil,
	},
	{
		cars{car{"a", "green"}, car{"b", "blue"}, car{"c", "gray"}},
		func(s car) bool {
			return s.Name != "b"
		},
		cars{car{"a", "green"}, car{"c", "gray"}},
		cars{car{"b", "blue"}},
		cars{car{"A", "green"}, car{"B", "blue"}, car{"C", "gray"}},
	},
}

func TestCars_Filter(t *testing.T) {
	for _, test := range carsFilterTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.expectedFilter, test.ss.Filter(test.condition))
		})
	}
}

func TestCars_FilterNot(t *testing.T) {
	for _, test := range carsFilterTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.expectedFilterNot, test.ss.FilterNot(test.condition))
		})
	}
}

func TestCars_Map(t *testing.T) {
	for _, test := range carsFilterTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.expectedMap, test.ss.Map(func(car car) car {
				car.Name = strings.ToUpper(car.Name)

				return car
			}))
		})
	}
}

var carsFirstAndLastTests = []struct {
	ss             cars
	first, firstOr car
	last, lastOr   car
}{
	{
		nil,
		car{},
		car{"default1", "unknown"},
		car{},
		car{"default2", "unknown"},
	},
	{
		cars{car{"foo", "red"}},
		car{"foo", "red"},
		car{"foo", "red"},
		car{"foo", "red"},
		car{"foo", "red"},
	},
	{
		cars{car{"a", "green"}, car{"b", "blue"}},
		car{"a", "green"},
		car{"a", "green"},
		car{"b", "blue"},
		car{"b", "blue"},
	},
	{
		cars{car{"a", "green"}, car{"b", "blue"}, car{"c", "gray"}},
		car{"a", "green"},
		car{"a", "green"},
		car{"c", "gray"},
		car{"c", "gray"},
	},
}

func TestCars_FirstOr(t *testing.T) {
	for _, test := range carsFirstAndLastTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.firstOr, test.ss.FirstOr(car{"default1", "unknown"}))
		})
	}
}

func TestCars_LastOr(t *testing.T) {
	for _, test := range carsFirstAndLastTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.lastOr, test.ss.LastOr(car{"default2", "unknown"}))
		})
	}
}

func TestCars_First(t *testing.T) {
	for _, test := range carsFirstAndLastTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.first, test.ss.First())
		})
	}
}

func TestCars_Last(t *testing.T) {
	for _, test := range carsFirstAndLastTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.last, test.ss.Last())
		})
	}
}

var carsStatsTests = []struct {
	ss       cars
	min, max car
	mode     cars
	len      int
}{
	{
		nil,
		car{},
		car{},
		cars{car{}},
		0,
	},
	{
		cars{},
		car{},
		car{},
		cars{car{}},
		0,
	},
	{
		cars{car{"foo", "red"}},
		car{"foo", "red"},
		car{"foo", "red"},
		cars{car{"foo", "red"}},
		1,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"qux", "cyan"}, car{"foo", "red"}},
		car{"Baz", "black"},
		car{"qux", "cyan"},
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"qux", "cyan"}, car{"foo", "red"}},
		4,
	},
}

func TestCars_Mode(t *testing.T) {
	cmp := func(a, b cars) bool {
		m := make(map[car]struct{})
		for _, i := range a {
			m[i] = struct{}{}
		}
		for _, i := range b {
			if _, ok := m[i]; !ok {
				return false
			}
		}
		return true
	}
	for _, test := range carsStatsTests {
		t.Run("", func(t *testing.T) {
			assert.True(t, cmp(test.mode, cars(test.ss).Mode()))
		})
	}
}

func TestCars_Len(t *testing.T) {
	for _, test := range carsStatsTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.len, cars(test.ss).Len())
		})
	}
}

var carsJSONTests = []struct {
	ss         cars
	jsonString string
}{
	{
		nil,
		`[]`, // Instead of null.
	},
	{
		cars{},
		`[]`,
	},
	{
		cars{car{"foo", "red"}},
		`[{"Name":"foo","Color":"red"}]`,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"qux", "cyan"}, car{"foo", "red"}},
		`[{"Name":"bar","Color":"yellow"},{"Name":"Baz","Color":"black"},{"Name":"qux","Color":"cyan"},{"Name":"foo","Color":"red"}]`,
	},
}

func TestCars_JSONBytes(t *testing.T) {
	for _, test := range carsJSONTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, []byte(test.jsonString), test.ss.JSONBytes())
		})
	}
}
func TestCars_JSONString(t *testing.T) {
	for _, test := range carsJSONTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.jsonString, test.ss.JSONString())
		})
	}
}

var carsJSONIndentTests = []struct {
	ss         cars
	jsonString string
}{
	{
		nil,
		`[]`, // Instead of null.
	},
	{
		cars{},
		`[]`,
	},
	{
		cars{car{"foo", "red"}},
		`[
  {
    "Name": "foo",
    "Color": "red"
  }
]`,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"qux", "cyan"}, car{"foo", "red"}},
		`[
  {
    "Name": "bar",
    "Color": "yellow"
  },
  {
    "Name": "Baz",
    "Color": "black"
  },
  {
    "Name": "qux",
    "Color": "cyan"
  },
  {
    "Name": "foo",
    "Color": "red"
  }
]`,
	},
}

func TestCars_JSONBytesIndent(t *testing.T) {
	for _, test := range carsJSONIndentTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, []byte(test.jsonString), test.ss.JSONBytesIndent("", "  "))
		})
	}
}

func TestCars_JSONStringIndent(t *testing.T) {
	for _, test := range carsJSONIndentTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.jsonString, test.ss.JSONStringIndent("", "  "))
		})
	}
}

var carsSortTests = []struct {
	ss        cars
	sorted    cars
	reversed  cars
	areSorted bool
}{
	{
		nil,
		nil,
		nil,
		true,
	},
	{
		cars{},
		cars{},
		cars{},
		true,
	},
	{
		cars{car{"foo", "red"}},
		cars{car{"foo", "red"}},
		cars{car{"foo", "red"}},
		true,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"foo", "red"}},
		cars{car{"Baz", "black"}, car{"bar", "yellow"}, car{"foo", "red"}},
		cars{car{"foo", "red"}, car{"Baz", "black"}, car{"bar", "yellow"}},
		false,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"qux", "cyan"}, car{"foo", "red"}},
		cars{car{"Baz", "black"}, car{"bar", "yellow"}, car{"foo", "red"}, car{"qux", "cyan"}},
		cars{car{"foo", "red"}, car{"qux", "cyan"}, car{"Baz", "black"}, car{"bar", "yellow"}},
		false,
	},
	{
		cars{car{"Baz", "black"}, car{"bar", "yellow"}},
		cars{car{"Baz", "black"}, car{"bar", "yellow"}},
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		true,
	},
}

func TestCars_Reverse(t *testing.T) {
	for _, test := range carsSortTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.reversed, test.ss.Reverse())
		})
	}
}

var carsSortCustomTests = []struct {
	ss                  cars
	sortedStableByName  cars
	sortedStableByColor cars
}{
	{
		nil,
		nil,
		nil,
	},
	{
		cars{},
		cars{},
		cars{},
	},
	{
		cars{car{"foo", "red"}},
		cars{car{"foo", "red"}},
		cars{car{"foo", "red"}},
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"foo", "red"}},
		cars{car{"Baz", "black"}, car{"bar", "yellow"}, car{"foo", "red"}},
		cars{car{"Baz", "black"}, car{"foo", "red"}, car{"bar", "yellow"}},
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"qux", "cyan"}, car{"foo", "red"}},
		cars{car{"Baz", "black"}, car{"bar", "yellow"}, car{"foo", "red"}, car{"qux", "cyan"}},
		cars{car{"Baz", "black"}, car{"qux", "cyan"}, car{"foo", "red"}, car{"bar", "yellow"}},
	},
	{
		cars{car{"aaa", "yellow"}, car{"aaa", "black"}, car{"bbb", "yellow"}, car{"bbb", "black"}},
		cars{car{"aaa", "yellow"}, car{"aaa", "black"}, car{"bbb", "yellow"}, car{"bbb", "black"}},
		cars{car{"aaa", "black"}, car{"bbb", "black"}, car{"aaa", "yellow"}, car{"bbb", "yellow"}},
	},
}

func carNameLess(a, b car) bool {
	return a.Name < b.Name
}

func carColorLess(a, b car) bool {
	return a.Color < b.Color
}

func TestCars_SortUsing(t *testing.T) {
	isSortedUsing := func(ss cars, less func(a, b car) bool) bool {
		for i := 1; i < len(ss); i++ {
			if less(ss[i], ss[i-1]) {
				return false
			}
		}
		return true
	}

	for _, test := range carsSortCustomTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()

			sortedByName := test.ss.SortUsing(carNameLess)
			assert.True(t, isSortedUsing(sortedByName, carNameLess))
			sortedStableByName := test.ss.SortStableUsing(carNameLess)
			assert.Equal(t, test.sortedStableByName, sortedStableByName)

			sortedByColor := test.ss.SortUsing(carColorLess)
			assert.True(t, isSortedUsing(sortedByColor, carColorLess))
			sortedStableByColor := test.ss.SortStableUsing(carColorLess)
			assert.Equal(t, test.sortedStableByColor, sortedStableByColor)
		})
	}
}

var carsStringsUsingTests = []struct {
	ss        cars
	transform func(car) string
	expected  Strings
}{
	{
		nil,
		func(s car) string {
			return "foo"
		},
		nil,
	},
	{
		cars{},
		func(s car) string {
			return fmt.Sprintf("%s!", s.Name)
		},
		nil,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"foo", "red"}},
		func(s car) string {
			return fmt.Sprintf("%s!", s.Color)
		},
		Strings{"yellow!", "black!", "red!"},
	},
}

func TestCars_StringsUsing(t *testing.T) {
	for _, test := range carsStringsUsingTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.expected, test.ss.StringsUsing(test.transform))
		})
	}
}

func TestCars_Append(t *testing.T) {
	assert.Equal(t,
		len(cars{}.Append()),
		0,
	)

	assert.Equal(t,
		cars{}.Append(car{"bar", "yellow"}),
		cars{car{"bar", "yellow"}},
	)

	assert.Equal(t,
		cars{}.Append(car{"bar", "yellow"}, car{"Baz", "black"}),
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
	)

	assert.Equal(t,
		cars{car{"bar", "yellow"}}.Append(car{"Baz", "black"}),
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
	)

	assert.Equal(t,
		cars{car{"bar", "yellow"}}.Append(car{"Baz", "black"}, car{"foo", "red"}),
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"foo", "red"}},
	)
}

func TestCars_Extend(t *testing.T) {
	assert.Equal(t,
		cars{}.Extend(),
		cars{},
	)

	assert.Equal(t,
		cars{}.Extend(cars{car{"bar", "yellow"}}),
		cars{car{"bar", "yellow"}},
	)

	assert.Equal(t,
		cars{}.Extend(cars{car{"bar", "yellow"}}, cars{car{"Baz", "black"}}),
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
	)

	assert.Equal(t,
		cars{car{"bar", "yellow"}}.Extend(cars{car{"Baz", "black"}}),
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
	)

	assert.Equal(t,
		cars{car{"bar", "yellow"}}.Extend(cars{car{"Baz", "black"}, car{"foo", "red"}}),
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"foo", "red"}},
	)
}

func TestCars_All(t *testing.T) {
	assert.True(t,
		(cars)(nil).All(func(value car) bool {
			return false
		}),
	)

	assert.True(t,
		(cars)(nil).All(func(value car) bool {
			return false
		}),
	)

	// None
	assert.False(t,
		cars{car{"bar", "yellow"}, car{"Baz", "black"}}.All(func(value car) bool {
			return value.Color == "green"
		}),
	)

	// Some
	assert.False(t,
		cars{car{"bar", "yellow"}, car{"Baz", "black"}}.All(func(value car) bool {
			return value.Color == "yellow"
		}),
	)

	// All
	assert.True(t,
		cars{car{"bar", "yellow"}, car{"Baz", "black"}}.All(func(value car) bool {
			return len(value.Name) > 0
		}),
	)
}

func TestCars_Any(t *testing.T) {
	assert.False(t,
		(cars)(nil).Any(func(value car) bool {
			return true
		}),
	)

	assert.False(t,
		(cars)(nil).Any(func(value car) bool {
			return true
		}),
	)

	// None
	assert.False(t,
		cars{car{"bar", "yellow"}, car{"Baz", "black"}}.Any(func(value car) bool {
			return value.Color == "green"
		}),
	)

	// Some
	assert.True(t,
		cars{car{"bar", "yellow"}, car{"Baz", "black"}}.Any(func(value car) bool {
			return value.Color == "yellow"
		}),
	)

	// All
	assert.True(t,
		cars{car{"bar", "yellow"}, car{"Baz", "black"}}.Any(func(value car) bool {
			return len(value.Name) > 0
		}),
	)
}

var carsShuffleTests = []struct {
	ss       cars
	expected cars
	source   rand.Source
}{
	{
		nil,
		nil,
		nil,
	},
	{
		nil,
		nil,
		rand.NewSource(0),
	},
	{
		cars{},
		cars{},
		rand.NewSource(0),
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"foo", "red"}},
		cars{car{"Baz", "black"}, car{"bar", "yellow"}, car{"foo", "red"}},
		rand.NewSource(0),
	},
	{
		cars{car{"bar", "yellow"}},
		cars{car{"bar", "yellow"}},
		rand.NewSource(0),
	},
}

func TestCars_Shuffle(t *testing.T) {
	for _, test := range carsShuffleTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.expected, test.ss.Shuffle(test.source))
		})
	}
}

var carsTopAndBottomTests = []struct {
	ss     cars
	n      int
	top    cars
	bottom cars
}{
	{
		nil,
		1,
		nil,
		nil,
	},
	{
		cars{},
		1,
		nil,
		nil,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		1,
		cars{car{"bar", "yellow"}},
		cars{car{"Baz", "black"}},
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		3,
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		cars{car{"Baz", "black"}, car{"bar", "yellow"}},
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		0,
		nil,
		nil,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		-1,
		nil,
		nil,
	},
}

func TestCars_Top(t *testing.T) {
	for _, test := range carsTopAndBottomTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.top, test.ss.Top(test.n))
		})
	}
}

func TestCars_Bottom(t *testing.T) {
	for _, test := range carsTopAndBottomTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.bottom, test.ss.Bottom(test.n))
		})
	}
}

func TestCars_Each(t *testing.T) {
	var names []string

	names = []string{}
	cars{}.Each(func(car car) {
		names = append(names, car.Name)
	})
	assert.Equal(t, []string{}, names)

	names = []string{}
	cars{car{"bar", "yellow"}, car{"Baz", "black"}}.Each(func(car car) {
		names = append(names, car.Name)
	})
	assert.Equal(t, []string{"bar", "Baz"}, names)
}

var carsRandomTests = []struct {
	ss       cars
	expected car
	source   rand.Source
}{
	{
		nil,
		car{},
		nil,
	},
	{
		nil,
		car{},
		rand.NewSource(0),
	},
	{
		cars{},
		car{},
		rand.NewSource(0),
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"foo", "red"}},
		car{"bar", "yellow"},
		rand.NewSource(0),
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"foo", "red"}},
		car{"foo", "red"},
		rand.NewSource(1),
	},
	{
		cars{car{"bar", "yellow"}},
		car{"bar", "yellow"},
		rand.NewSource(0),
	},
}

func TestCars_Random(t *testing.T) {
	for _, test := range carsRandomTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.expected, test.ss.Random(test.source))
		})
	}
}

var carsSendTests = []struct {
	ss            cars
	recieveDelay  time.Duration
	canceledDelay time.Duration
	expected      cars
}{
	{
		nil,
		0,
		0,
		nil,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		0,
		0,
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		time.Millisecond * 30,
		time.Millisecond * 10,
		cars{car{"bar", "yellow"}},
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		time.Millisecond * 3,
		time.Millisecond * 10,
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
	},
}

func TestCar_Send(t *testing.T) {
	for _, test := range carsSendTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			ch := make(chan car)
			actual := getCarsFromChan(ch, test.recieveDelay)
			ctx := createContextByDelay(test.canceledDelay)

			actualSended := test.ss.Send(ctx, ch)
			close(ch)

			assert.Equal(t, test.expected, actualSended)
			assert.Equal(t, test.expected, actual())
		})
	}
}

var carsDiffTests = map[string]struct {
	ss1     cars
	ss2     cars
	added   cars
	removed cars
}{
	"BothEmpty": {
		nil,
		nil,
		nil,
		nil,
	},
	"OnlyRemovedUnique": {
		cars{car{"a", "green"}, car{"bar", "yellow"}},
		nil,
		nil,
		cars{car{"a", "green"}, car{"bar", "yellow"}},
	},
	"OnlyRemovedDuplicates": {
		cars{car{"a", "green"}, car{"Baz", "black"}, car{"a", "green"}},
		nil,
		nil,
		cars{car{"a", "green"}, car{"Baz", "black"}, car{"a", "green"}},
	},
	"OnlyAddedUnique": {
		nil,
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		nil,
	},
	"OnlyAddedDuplicates": {
		nil,
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"Baz", "black"}, car{"bar", "yellow"}},
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{"Baz", "black"}, car{"bar", "yellow"}},
		nil,
	},
	"AddedAndRemovedUnique": {
		cars{car{"a", "green"}, car{"bar", "yellow"}, car{"Baz", "black"}, car{"qux", "grey"}},
		cars{car{"Baz", "black"}, car{"qux", "grey"}, car{"quux", "red"}, car{"Baz", "magenta"}},
		cars{car{"quux", "red"}, car{"Baz", "magenta"}},
		cars{car{"a", "green"}, car{"bar", "yellow"}},
	},
	"AddedAndRemovedDuplicates": {
		cars{car{"a", "green"}, car{"bar", "yellow"}, car{"Baz", "black"}, car{"Baz", "black"}, car{"qux", "grey"}},
		cars{car{"Baz", "black"}, car{"qux", "grey"}, car{"quux", "red"}, car{"qux", "grey"}, car{"Baz", "magenta"}},
		cars{car{"quux", "red"}, car{"qux", "grey"}, car{"Baz", "magenta"}},
		cars{car{"a", "green"}, car{"bar", "yellow"}, car{"Baz", "black"}},
	},
}

func TestCars_Diff(t *testing.T) {
	for testName, test := range carsDiffTests {
		t.Run(testName, func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss1)()
			defer assertImmutableCars(t, &test.ss2)()

			added, removed := test.ss1.Diff(test.ss2)
			assert.Equal(t, test.added, added)
			assert.Equal(t, test.removed, removed)
		})
	}
}

func TestCars_Strings(t *testing.T) {
	assert.Equal(t, Strings(nil), cars{}.Strings())

	assert.Equal(t,
		Strings{"{a green}", "{bar yellow}", "{Baz black}"},
		cars{car{"a", "green"}, car{"bar", "yellow"}, car{"Baz", "black"}}.Strings())
}

func TestCars_Ints(t *testing.T) {
	assert.Equal(t, Ints(nil), cars{}.Ints())

	assert.Equal(t,
		Ints{0, 0, 0},
		cars{car{"a", "green"}, car{"bar", "yellow"}, car{"Baz", "black"}}.Ints())
}

func TestCars_Float64s(t *testing.T) {
	assert.Equal(t, Float64s(nil), cars{}.Float64s())

	assert.Equal(t,
		Float64s{0, 0, 0},
		cars{car{"a", "green"}, car{"bar", "yellow"}, car{"Baz", "black"}}.Float64s())
}

var carsSequenceTests = []struct {
	ss       cars
	creator  func(int) car
	params   []int
	expected cars
}{
	// n
	{
		nil,
		nil,
		nil,
		nil,
	},
	{
		nil,
		nil,
		[]int{0},
		nil,
	},
	{
		nil,
		nil,
		[]int{-1},
		nil,
	},
	{
		nil,
		func(i int) car { return car{Name: strconv.Itoa(i)} },
		[]int{3},
		cars{{Name: "0"}, {Name: "1"}, {Name: "2"}},
	},
	// range
	{
		nil,
		func(i int) car { return car{Name: strconv.Itoa(i)} },
		[]int{6, 6},
		nil,
	},
	{
		nil,
		func(i int) car { return car{Name: strconv.Itoa(i)} },
		[]int{8, 6},
		nil,
	},
	{
		nil,
		func(i int) car { return car{Name: strconv.Itoa(i)} },
		[]int{3, 6},
		cars{{Name: "3"}, {Name: "4"}, {Name: "5"}},
	},
	{
		nil,
		func(i int) car { return car{Name: strconv.Itoa(i)} },
		[]int{-6, -3},
		cars{{Name: "-6"}, {Name: "-5"}, {Name: "-4"}},
	},
	{
		nil,
		func(i int) car { return car{Name: strconv.Itoa(i)} },
		[]int{-3, -6},
		nil,
	},
	// range with step
	{
		nil,
		func(i int) car { return car{Name: strconv.Itoa(i)} },
		[]int{3, 7, 2},
		cars{{Name: "3"}, {Name: "5"}},
	},
	{
		nil,
		func(i int) car { return car{Name: strconv.Itoa(i)} },
		[]int{-3, -6, -2},
		cars{{Name: "-3"}, {Name: "-5"}},
	},
	{
		nil,
		func(i int) car { return car{Name: strconv.Itoa(i)} },
		[]int{3, 7, 10},
		nil,
	},
}

func TestCars_SequenceUsing(t *testing.T) {
	for _, test := range carsSequenceTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.expected, test.ss.SequenceUsing(test.creator, test.params...))
		})
	}
}

var carsDropTopTests = []struct {
	ss      cars
	n       int
	dropTop cars
}{
	{
		nil,
		1,
		nil,
	},
	{
		cars{},
		1,
		nil,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		-1,
		nil,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		0,
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		1,
		cars{car{"Baz", "black"}},
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		2,
		nil,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		3,
		nil,
	},
}

func TestCars_DropTop(t *testing.T) {
	for _, test := range carsDropTopTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.dropTop, test.ss.DropTop(test.n))
		})
	}
}

var carsDropWhileTests = []struct {
	ss        cars
	f         func(s car) bool
	dropWhile cars
}{
	{
		ss:        nil,
		f:         func(s car) bool { return s.Color == "Blue" },
		dropWhile: cars{},
	},
	{
		ss:        cars{{"Bar", "blue"}, {"Foo", "blue"}, {"Baz", "blue"}, {"Bit", "pink"}, {"Baz", "red"}},
		f:         func(s car) bool { return s.Color == "blue" },
		dropWhile: cars{{"Bit", "pink"}, {"Baz", "red"}},
	},
	{
		ss:        cars{{"Bar", "pink"}, {"Foo", "pink"}, {"Baz", "pink"}},
		f:         func(s car) bool { return s.Color == "pink" },
		dropWhile: cars{},
	},
	{
		ss:        cars{{"Bar", "blue"}, {"Bar", "blue"}, {"Bar", "yellow"}, {"Bar", "black"}, {"Bar", "blue"}},
		f:         func(s car) bool { return s.Color == "red" },
		dropWhile: cars{{"Bar", "blue"}, {"Bar", "blue"}, {"Bar", "yellow"}, {"Bar", "black"}, {"Bar", "blue"}},
	},
}

func TestCars_DropWhile(t *testing.T) {
	for _, test := range carsDropWhileTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.dropWhile, test.ss.DropWhile(test.f))
		})
	}
}

var carsSubSliceTests = []struct {
	ss       cars
	start    int
	end      int
	subSlice cars
}{
	{
		nil,
		1,
		1,
		nil,
	},
	{
		nil,
		1,
		2,
		cars{car{}},
	},
	{
		cars{},
		1,
		1,
		nil,
	},
	{
		cars{},
		1,
		2,
		cars{car{}},
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		-1,
		-1,
		nil,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		-1,
		1,
		nil,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		1,
		-1,
		nil,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		2,
		0,
		nil,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		1,
		1,
		nil,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		1,
		2,
		cars{car{"Baz", "black"}},
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		1,
		3,
		cars{car{"Baz", "black"}, car{}},
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		2,
		2,
		nil,
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		2,
		3,
		cars{car{}},
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}, car{}},
		2,
		3,
		cars{car{}},
	},
}

func TestCars_SubSlice(t *testing.T) {
	for _, test := range carsSubSliceTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.subSlice, test.ss.SubSlice(test.start, test.end))
		})
	}
}

var carsFindFirstUsingTests = []struct {
	ss         cars
	expression func(value car) bool
	expected   int
}{
	{
		nil,
		func(value car) bool { return value.Color == "red" },
		-1,
	},
	{
		cars{},
		func(value car) bool { return value.Name == "ferrari" },
		-1,
	},
	{
		cars{car{Name: "volvo", Color: "brown"}},
		func(value car) bool { return value.Name == "eclipse" },
		-1,
	},
	{
		cars{car{Name: "maverick", Color: "red"}, car{Name: "ferrari", Color: "red"}, car{Name: "polo", Color: "white"}},
		func(value car) bool { return value.Name == "polo" && value.Color == "white" },
		2,
	},
	{
		cars{car{Name: "maverick", Color: "red"}, car{Name: "ferrari", Color: "red"}, car{Name: "polo", Color: "white"}},
		func(value car) bool { return value.Name == "maverick" && value.Color == "white" },
		-1,
	},
}

func TestCars_FindFirstUsing(t *testing.T) {
	for _, test := range carsFindFirstUsingTests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, test.expected, test.ss.FindFirstUsing(test.expression))
		})
	}
}

var carsEqualsTests = []struct {
	ss       cars
	rhs      cars
	expected bool
}{
	{nil, nil, true},
	{cars{}, cars{}, true},
	{nil, cars{}, true},
	{cars{{Name: "1"}, {Name: "2"}}, cars{{Name: "1"}, {Name: "2"}}, true},
	{cars{{Name: "1"}, {Name: "2"}}, cars{{Name: "1"}, {Name: "3"}}, false},
	{cars{{Name: "1"}, {Name: "2"}}, cars{{Name: "1"}}, false},
	{cars{{Name: "2"}}, cars{{Name: "1"}}, false},
	{cars{{Name: "2"}}, nil, false},
}

func TestCars_Equals(t *testing.T) {
	for _, test := range carsEqualsTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			assert.Equal(t, test.expected, test.ss.Equals(test.rhs))
		})
	}
}

var carsShiftAndUnshiftTests = []struct {
	ss      cars
	shifted car
	shift   cars
	params  cars
	unshift cars
}{
	{
		nil,
		car{},
		nil,
		nil,
		cars{},
	},
	{
		nil,
		car{},
		nil,
		cars{},
		cars{},
	},
	{
		nil,
		car{},
		nil,
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
	},
	{
		cars{},
		car{},
		nil,
		nil,
		cars{},
	},
	{
		cars{},
		car{},
		nil,
		cars{},
		cars{},
	},
	{
		cars{},
		car{},
		nil,
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
	},
	{
		cars{car{"bar", "yellow"}},
		car{"bar", "yellow"},
		nil,
		cars{car{"Baz", "black"}},
		cars{car{"Baz", "black"}, car{"bar", "yellow"}},
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		car{"bar", "yellow"},
		cars{car{"Baz", "black"}},
		cars{car{}},
		cars{car{}, car{"bar", "yellow"}, car{"Baz", "black"}},
	},
	{
		cars{car{"bar", "yellow"}, car{"Baz", "black"}},
		car{"bar", "yellow"},
		cars{car{"Baz", "black"}},
		cars{car{}, car{"zzz", "blue"}},
		cars{car{}, car{"zzz", "blue"}, car{"bar", "yellow"}, car{"Baz", "black"}},
	},
}

func TestCars_ShiftAndUnshift(t *testing.T) {
	for _, test := range carsShiftAndUnshiftTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCars(t, &test.ss)()
			shifted, shift := test.ss.Shift()
			assert.Equal(t, test.shifted, shifted)
			assert.Equal(t, test.shift, shift)
			assert.Equal(t, test.unshift, test.ss.Unshift(test.params...))
		})
	}
}

func TestCars_Join(t *testing.T) {
	assert.Equal(t, "", cars(nil).Join("a"))
	assert.Equal(t, "", cars{}.Join("a"))
	car1, car2 := car{Name: "maverick", Color: "red"}, car{Name: "ferrari", Color: "red"}
	assert.Equal(t, "{maverick red}-{ferrari red}", cars{car1, car2}.Join("-"))
}

func TestCars_Pop(t *testing.T) {
	car1, car2 := car{Name: "maverick", Color: "red"}, car{Name: "ferrari", Color: "red"}

	garage := cars{car1, car2}

	assert.Equal(t, &car1, garage.Pop())
	assert.Equal(t, cars{car2}, garage)

	assert.Equal(t, &car2, garage.Pop())
	assert.Equal(t, cars{}, garage)

	assert.Equal(t, (*car)(nil), garage.Pop())
	assert.Equal(t, cars{}, garage)

}

func TestCars_Insert(t *testing.T) {
	car1 := car{Name: "maverick", Color: "red"}
	car2 := car{Name: "ferrari", Color: "red"}
	car3 := car{Name: "panda", Color: "white"}

	assert.Equal(t, cars{}, cars(nil).Insert(0))
	assert.Equal(t, cars{car2, car1}, cars{car1}.Insert(0, car2))
	assert.Equal(t, cars{car1, car2}, cars{car1}.Insert(1, car2))
	assert.Equal(t, cars{car1, car2, car3}, cars{car1, car3}.Insert(1, car2))
}
