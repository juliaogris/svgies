package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	str := "M-2.18557e-06 50C-9.78513e-07 77.6142 22.3858 100 50 100L50 -2.18557e-06C22.3858 -9.78513e-07 -3.39263e-06 22.3858 -2.18557e-06 50Z"
	if len(os.Args) >= 2 {
		str = os.Args[1]
	}
	r := regexp.MustCompile(`[MZCLVH][^MZCLVH]*`)
	match := r.FindAllStringSubmatch(str, -1)
	bailIfErr(validate(match))
	segments, err := convert(match)
	bailIfErr(err)
	round(segments)
	fmt.Println(toString(segments))
}

func validate(match [][]string) error {
	if len(match) < 3 {
		return fmt.Errorf("error: expected at least 3 segments M1 2 C1 2 3 4 5 6Z, got %d", len(match))
	}
	for i, m := range match {
		if len(m) != 1 || len(m[0]) < 1 {
			return fmt.Errorf("submatch %d not 1", i)
		}
	}
	if match[0][0][0] != 'M' || match[len(match)-1][0][0] != 'Z' {
		return fmt.Errorf("error: m start or end")
	}
	for i, m := range match[1 : len(match)-1] {
		ch := m[0][0]
		if ch != 'C' && ch != 'L' && ch != 'V' && ch != 'H' {
			return fmt.Errorf("error: missing C|L|V|H prefix in %d", i)
		}
	}
	return nil
}

func convert(match [][]string) ([][]float64, error) {
	var err error
	segments := make([][]float64, len(match)-1)
	for i, m := range match[:len(match)-1] { // skip the last Z
		if segments[i], err = convertSegment(m[0]); err != nil {
			return nil, err
		}
		command := m[0][0]
		if command == 'V' {
			prevX := segments[i-1][len(segments[i-1])-2]
			segments[i] = append([]float64{prevX}, segments[i][0])
		} else if command == 'H' {
			prevY := segments[i-1][len(segments[i-1])-1]
			segments[i] = append(segments[i], prevY)
		}
	}
	return segments, nil
}

var expectedLengths = map[byte]int{
	'M': 2,
	'C': 6,
	'Z': 0,
	'L': 2,
	'V': 1,
	'H': 1,
}

func convertSegment(segment string) ([]float64, error) {
	expectedLength := expectedLengths[segment[0]]
	nums := strings.Fields(segment[1:])
	if len(nums) != expectedLength {
		return nil, fmt.Errorf("expected %d numbers, got %d", expectedLength, len(nums))
	}
	result := make([]float64, len(nums))
	var err error
	for i, num := range nums {
		result[i], err = strconv.ParseFloat(num, 64)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func round(segments [][]float64) {
	for _, segment := range segments {
		for i, seg := range segment {
			segment[i] = math.Round(seg)
		}
	}
}

func bailIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func toString(segments [][]float64) string {
	s := make([]string, len(segments))
	for i, segment := range segments {
		s[i] = toStringSegment(segment)
	}
	return "[" + strings.Join(s, ",\n ") + "]"
}

func toStringSegment(segment []float64) string {
	vals := make([]string, len(segment))
	for i, val := range segment {
		vals[i] = fmt.Sprintf("%.0f", val)
	}
	return "[" + strings.Join(vals, ", ") + "]"
}
