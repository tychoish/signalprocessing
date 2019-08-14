package change_points

import (
	"fmt"
	"log"
	"time"

	//"sort"

	"github.com/jimoleary/signalprocessing/util"
	"github.com/mongodb/grip"
	"github.com/mongodb/grip/level"
	"github.com/mongodb/grip/send"
	"github.com/stretchr/testify/assert"

	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/convey/reporting"
)

type ChangePointsFixture struct {
	Series       []float64 `json:"series"`
	Expected     []int     `json:"expected"`
	PValue       float64   `json:"p"`
	Permutations int       `json:"permutations"`
}

type EDMFixture struct {
	Series   []float64 `json:"series"`
	Expected []int     `json:"expected"`
	MinSize  int       `json:"minSize"`
	Beta     float64   `json:"beta"`
	Degree   string    `json:"degree"`
}

type QFixture struct {
	Series   []float64 `json:"series"`
	Expected []float64 `json:"expected"`
}

type ExtractQFixture struct {
	QValues  []float64 `json:"qs"`
	Expected struct {
		Index int     `json:"index"`
		Value float64 `json:"value"`
	} `json:"expected"`
}

func init() {
	reporting.QuietMode()
	err := grip.GetSender().SetLevel(send.LevelInfo{level.Error, level.Error})
	if err != nil {
		fmt.Println(err)
	}
}

func EdmTestHelper(t *testing.T) {
	Convey("Find Change Points", t, func() {
		asserter := assert.New(t)
		// Create fixture with defaults
		fixture := EDMFixture{MinSize: 15, Beta: 2.0, Degree: "Quadratic"}
		if err := util.LoadFixture(t.Name(), &fixture); err != nil {
			asserter.Fail("Error loading fixture")
		}
		series := fixture.Series

		expected := fixture.Expected
		var changePointIndexes []int
		start := time.Now()
		changePointIndexes = EDivisiveWithMedians(series, fixture.MinSize)
		elapsed := time.Since(start)
		log.Printf("\nSort took %s", elapsed)
		asserter.Equal(expected, changePointIndexes)
	})
}


func ChangePointTestHelper(t *testing.T) {
	Convey("Find Change Points", t, func() {
		asserter := assert.New(t)
		// Create fixture with defaults
		fixture := ChangePointsFixture{PValue: .05, Permutations: 100}
		if err := util.LoadFixture(t.Name(), &fixture); err != nil {
			asserter.Fail("Error loading fixture")
		}
		series := fixture.Series
		pValue := fixture.PValue

		permutations := fixture.Permutations
		expected := fixture.Expected
		start := time.Now()
		changePoints, _ := ChangePoints(series, pValue, permutations)
		elapsed := time.Since(start)
		log.Printf("\nEDivisive took %s", elapsed)
		changePointIndexes := make([]int, len(changePoints))
		for i, cp := range changePoints {
			changePointIndexes[i] = cp.Index
		}
		asserter.Equal(expected, changePointIndexes)
	})

}

func TestQHat(t *testing.T) {
	Convey("Calculate QHat Values", t, func() {
		asserter := assert.New(t)
		fixture := QFixture{}
		if err := util.LoadFixture(t.Name(), &fixture); err != nil {
			asserter.Fail("Error loading fixtures")
		}
		series := fixture.Series
		expected := fixture.Expected

		values, _ := QHat(series)
		asserter.Equal(expected, values)
	})
}

func TestExtractQ(t *testing.T) {
	Convey("Extract Max Q", t, func() {
		asserter := assert.New(t)
		fixture := ExtractQFixture{}
		if err := util.LoadFixture(t.Name(), &fixture); err != nil {
			asserter.Fail("Error loading fixtures")
		}

		qValues := fixture.QValues
		expectedIndex := fixture.Expected.Index
		expectedValue := fixture.Expected.Value

		index, value := ExtractQ(qValues)
		asserter.Equal(expectedValue, value)
		asserter.Equal(expectedIndex, index)
	})
}

func TestChangePoints(t *testing.T) {
	ChangePointTestHelper(t)
}

func TestShort(t *testing.T) {
	ChangePointTestHelper(t)
}

func TestShortEDM(t *testing.T) {
	EdmTestHelper(t)
}

func TestMedium(t *testing.T) {
	ChangePointTestHelper(t)
}

func TestMediumEDM(t *testing.T) {
	EdmTestHelper(t)
}

//func TestLarge(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping slow test")
//		return
//	}
//	ChangePointTestHelper(t)
//}

func TestLargeEDM(t *testing.T) {
	EdmTestHelper(t)
}

//func TestHuge(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping slow test")
//	}
//	ChangePointTestHelper(t)
//}
//
//func TestHugeEDM(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping slow test")
//		return
//	}
//	edmTestHelper(t)
//}
//
//func TestHumungousPlusEDM(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping slow test")
//		return
//	}
//	edmTestHelper(t)
//}
//
//func TestHumungousEDM(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping slow test")
//		return
//	}
//	edmTestHelper(t)
//}

func TestInsert(t *testing.T) {
	f := &SortedList{}
	f.Insert(1.0, 2.0, 3.0)
	asserter := assert.New(t)

	asserter.Equal(f, &SortedList{1.0, 2.0, 3.0})

	f.Insert(2.0)
	asserter.Equal(f, &SortedList{1.0, 2.0, 2.0, 3.0})

	f.Insert(0.0)
	asserter.Equal(f, &SortedList{0.0, 1.0, 2.0, 2.0, 3.0})

	f.Insert(4.0)
	asserter.Equal(f, &SortedList{0.0, 1.0, 2.0, 2.0, 3.0, 4.0})
}

func TestRemove(t *testing.T) {
	f := &SortedList{0.0, 1.0, 2.0, 2.0, 3.0, 4.0}
	f.Remove(0.0)
	asserter := assert.New(t)

	asserter.Equal(f, &SortedList{1.0, 2.0, 2.0, 3.0, 4.0})

	f.Remove(2.0)
	asserter.Equal(f, &SortedList{1.0, 2.0, 3.0, 4.0})

	f.Remove(0.0)
	asserter.Equal(f, &SortedList{1.0, 2.0, 3.0, 4.0})

	f.Remove(4.0)
	asserter.Equal(f, &SortedList{1.0, 2.0, 3.0})
}
