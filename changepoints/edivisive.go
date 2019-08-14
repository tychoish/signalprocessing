package change_points

import (
	"math"
	"math/rand"
	"sort"
)

func Init(seeds ...int64) {
	var seed int64 = 1234
	if seeds != nil {
		seed = seeds[0]
	}
	rand.Seed(seed)
}

var debug = false

func init() {
	Init()
}

func CalculateDiffs(series []float64) (diffs []float64) {
	length := len(series)
	diffs = make([]float64, length*length, length*length)
	for row := 0; row < length; row++ {
		for column := row; column < length; column++ {
			delta := math.Abs(series[row] - series[column])
			diffs[row*length+column] = delta
			diffs[column*length+row] = delta
		}
	}
	return
}

func CalculateQ(term1 float64, term2 float64, term3 float64, suffix int, prefix int) (newq float64) {
	m := float64(suffix)
	n := float64(prefix)

	term1Reg := term1 * (2.0 / (m * n))
	term2Reg := term2 * (2.0 / (n * (n - 1)))
	term3Reg := term3 * (2.0 / (m * (m - 1)))
	newq = float64(int((m * n) / (m + n)))
	newq = newq * (term1Reg - term2Reg - term3Reg)
	return newq
}

func QHat(series []float64) (qhatValues []float64, err error) {
	length := len(series)
	qhatValues = make([]float64, length)
	err = nil

	if length >= 5 {
		diffs := CalculateDiffs(series)

		n := 2
		m := length - n

		term1 := 0.0
		for i := 0; i < n; i++ {
			for j := n; j < length; j++ {
				term1 += diffs[i*length+j]
			}
		}
		term2 := 0.0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				term2 += diffs[i*length+j]
			}
		}

		term3 := 0.0
		for i := n; i < length; i++ {
			for j := i + 1; j < length; j++ {
				term3 += diffs[i*length+j]
			}

		}

		qhatValues[n] = CalculateQ(term1, term2, term3, m, n)

		for n := 3; n < (length - 2); n++ {
			m = length - n
			rowDelta := 0.0

			for j := 0; j < n-1; j++ {
				rowDelta = rowDelta + diffs[(n-1)*length+j]
			}

			columnDelta := 0.0
			for j := n - 1; j < length; j++ {
				columnDelta = columnDelta + diffs[j*length+n-1]
			}

			term1 = term1 - rowDelta + columnDelta
			term2 = term2 + rowDelta
			term3 = term3 - columnDelta

			qhatValues[n] = CalculateQ(term1, term2, term3, m, n)
		}
	}
	return
}

func ExtractQ(qhatValues []float64) (index int, value float64) {
	index = 0
	value = 0
	length := len(qhatValues)
	for i := 0; i < length; i++ {
		if qhatValues[i] > value {
			index = i
			value = qhatValues[i]
		}
	}
	return
}

type Window struct {
	Start int
	End   int
}

type ChangePoint struct {
	Index       int
	Q           float64
	Probability float64
	Order       int
	Window
}

// Shuffle a copy of each window, get the highest Q.
// Return true if ay of the permutations are greater than q.
func shuffleWindows(series []float64, windows []Window, q float64) bool {
	series = append([]float64{}, series...)
	maxQ := -1.0
	for _, w := range windows {
		window := series[w.Start:w.End]
		rand.Shuffle(len(window), func(i, j int) { window[i], window[j] = window[j], window[i] })

		winQs, _ := QHat(window)
		if _, winMaxQs := ExtractQ(winQs); winMaxQs > maxQ {
			maxQ = winMaxQs
		}
	}
	return maxQ >= q
}

// Create Windows, if there are none create a window that encompasses the whole range. If there are
// some windows then split the window that encompasses index and resort.
func createWindows(windows []Window, index int, length int) ([]Window, int) {
	found := -1

	if len(windows) == 0 {
		windows = append(windows, Window{Start: 0, End: length})
	} else {
		for i, current := range windows {
			if current.Start <= index && index <= current.End {
				found = i
			}
		}
		window := windows[found]
		windows[found].End = index

		windows = append(windows, Window{Start: index, End: window.End})
		sort.Slice(windows[:], func(i, j int) bool {
			return windows[i].Start < windows[j].Start
		})
	}
	return windows, index
}

// Create Candidate change points. If there are no change points then, calculate the required values
// for the whole series. If there are change points, then find the new candidate and split it.
func createCandidates(series []float64, candidates []ChangePoint, index int) []ChangePoint {
	if len(candidates) == 0 {
		winQs, _ := QHat(series)
		winIndexMax, winMaxQs := ExtractQ(winQs)

		candidate := ChangePoint{
			Index:  winIndexMax,
			Q:      winMaxQs,
			Order:  0,
			Window: Window{Start: 0, End: len(series)},
		}
		candidates = append(candidates, candidate)
	} else {
		found := 0
		for i, current := range candidates {
			if current.Index == index {
				found = i
			}
		}

		start := candidates[found].Start
		end := candidates[found].End

		buffer := series[start:index]

		winQs, _ := QHat(buffer)
		winIndexMax, winMaxQs := ExtractQ(winQs)
		candidates[found].Index = winIndexMax + start
		candidates[found].Q = winMaxQs
		candidates[found].Start = start
		candidates[found].End = index

		buffer = series[index:end]
		winQs, _ = QHat(buffer)
		winIndexMax, winMaxQs = ExtractQ(winQs)
		candidate := ChangePoint{
			Index:  winIndexMax + index,
			Q:      winMaxQs,
			Window: Window{Start: index, End: end},
		}
		candidates = append(candidates, candidate)

		sort.Slice(candidates[:], func(i, j int) bool {
			return candidates[i].Q < candidates[j].Q
		})
	}
	return candidates
}

func ChangePoints(series []float64, pvalue float64, permutations int) (changePoints []ChangePoint, err error) {
	changePoints = make([]ChangePoint, 0, 10)
	length := len(series)
	windows := make([]Window, 0)
	candidates := make([]ChangePoint, 0)

	probability := 0.0
	index := 0

	for probability <= pvalue {
		windows, index = createWindows(windows, index, length)
		candidates = createCandidates(series, candidates, index)

		candidateQ := candidates[len(candidates)-1]
		countAbove := 0.0 // results from permuted test >= candidateQ
		for i := 0; i < permutations; i++ {
			if shuffleWindows(series, windows, candidateQ.Q) {
				countAbove += 1
			}
		}

		probability = (1.0 + countAbove) / float64(permutations+1)

		// If we are not done, add the change point and setup the next iteration.
		if probability <= pvalue {
			changePoint := candidates[len(candidates)-1]
			changePoint.Probability = probability
			changePoints = append(changePoints, changePoint)
			index = changePoint.Index
		}
	}

	sort.Slice(changePoints[:], func(i, j int) bool {
		return changePoints[i].Order < changePoints[j].Order
	})

	return
}

// An FloatHeap is a min-sorted-list of floats.
type SortedList []float64

// Data interface for sort.
//// Get the length of the Heap
//func (s SortedList) Len() int {
//	return len(s)
//}
//
//// True id i is less than j.
//func (s SortedList) Less(i, j int) bool {
//	return s[i] < s[j]
//}
//
//// Swap i and j.
//func (s SortedList) Swap(i, j int) {
//	s[i], s[j] = s[j], s[i]
//}
//

// Clear the list.
func (s *SortedList) Clear() {
	*s = (*s)[:0]
}

// Insert the list of floats maintaining the sort order.
func (s *SortedList) Insert(floats ...float64) {
	length := len(*s)
	for _, f := range floats {
		// Assume f is appended.
		*s = append(*s, f)
		if length != 0 && f < (*s)[length-1] {
			// Deal with the case where we are inserting before the end.
			index := sort.SearchFloat64s(*s, f)
			copy((*s)[index+1:], (*s)[index:])
			(*s)[index] = f
		}
		length += 1
	}
	// note: the lines above are about 3 times faster than inserting and then sorting
	//*s = append(*s, floats...)
	//sort.Sort((*s))
}

// Remove f, maintaining the sort order.
func (s *SortedList) Remove(f float64) {
	index := sort.Search(len(*s), func(i int) bool { return (*s)[i] > f })
	length := len(*s)
	if index == 0 {
		if (*s)[index] == f {
			*s = (*s)[1:]
		}
	} else if index == length {
		*s = (*s)[0 : length-1]
	} else {
		copy((*s)[index-1:], (*s)[index:])
		*s = (*s)[:length-1]
	}
}

// Calculate the median, assuming the list is sorted.
func (s *SortedList) Median() float64 {
	length := len(*s)
	center := int(length / 2)
	if length%2 != 0 {
		return (*s)[center]
	}
	median := ((*s)[center] + (*s)[center-1]) / 2.0
	return median
}

func EDivisiveWithMedians(series []float64, minSize int) []int {
	n := len(series)
	prev := make([]int, n+1, n+1)
	number := make([]int, n+1, n+1)
	F := make([]float64, n+1, n+1)
	for i := range F {
		F[i] = -3.0
	}

	right := &SortedList{}
	left := &SortedList{}
	for s := 2 * minSize; s < n+1; s += 1 {
		right.Clear()
		left.Clear()
		left.Insert(series[prev[minSize-1] : minSize-1]...)
		right.Insert(series[minSize-1 : s]...)
		for t := minSize; t < s-minSize+1; t += 1 { //modify limits to deal with minSize
			left.Insert(series[t-1])
			right.Remove(series[t-1])

			if prev[t] > prev[t-1] {
				for i := prev[t-1]; i < prev[t]; i += 1 {
					left.Remove(series[i])
				}
			} else if prev[t] < prev[t-1] {
				for i := prev[t]; i < prev[t-1]; i += 1 {
					left.Insert(series[i])
				}
			}

			//calculate statistic value
			leftMedian := left.Median()
			rightMedian := right.Median()

			normalize := float64((t-prev[t])*(s-t)) / math.Pow(float64(s-prev[t]), 2.0)
			tmp := F[t] + normalize*math.Pow(leftMedian-rightMedian, 2.0)

			//check for improved optimal statistic value
			if tmp > F[s] {
				number[s] = number[t] + 1
				F[s] = tmp
				prev[s] = t
			}
		}
	}

	loc := make([]int, 0, len(prev))

	for at := n; at > 0; at = prev[at] {
		value := prev[at]
		if value != 0 {
			loc = append(loc, value)
		}
	}
	sort.Slice(loc[:], func(i, j int) bool {
		return loc[i] < loc[j]
	})
	return loc
}
