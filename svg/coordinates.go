package svg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const MISSING string = "X"
const UNKNOWN string = ""
const POINT string = "P"
const CUBIC string = "C"
const LINE string = "L"

var (
	_ Coordinates = Point{}
	_ Coordinates = Cubic{}
	_ Coordinates = Line{}
)

type Coordinates interface {
	Type() string
}

type Point [2]float64

func (p Point) Type() string {
	return POINT
}

type Cubic [][3]Point

func (c Cubic) Type() string {
	return CUBIC
}

type Line []Point

func (l Line) Type() string {
	return LINE
}

func DerivePathType(d string) string {

	path_re, err := regexp.Compile(`^M\s{0,}\d+\,\d+\s{0,}(L|C)`)

	if err != nil {
		return MISSING
	}

	if !path_re.MatchString(d) {
		return UNKNOWN
	}

	m := path_re.FindStringSubmatch(d)
	return m[1]
}

func DeriveCoordinates(d string) (Coordinates, error) {

	t := DerivePathType(d)

	switch t {
	case LINE:
		return deriveCoordinatesForLine(d)
	case CUBIC:
		return deriveCoordinatesForCurve(d)
	default:
		return nil, fmt.Errorf("Unsuported coordinate type, '%s' (%s)", t, d)
	}
}

func deriveCoordinatesForLine(d string) (Coordinates, error) {

	d = strings.TrimLeft(d, "M")
	d = strings.TrimRight(d, "Z")

	pairs := strings.Split(d, "L")
	coords := make([]Point, len(pairs))

	for i, pair := range pairs {

		xy := strings.Split(pair, ",")

		if len(xy) != 2 {
			return nil, fmt.Errorf("Invalid coordinate (%s) at offset %d", pair, i)
		}

		x, err := strconv.ParseFloat(xy[0], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[0], i, err)
		}

		y, err := strconv.ParseFloat(xy[1], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[1], i, err)
		}

		coords[i] = Point{x, y}
	}

	return Line(coords), nil
}

func deriveCoordinatesForCurve(d string) (Coordinates, error) {

	re, err := regexp.Compile(`^M\s{0,}\d+\,\d+\s{0,}C`)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile pattern, %w", err)
	}

	if !re.MatchString(d) {
		return nil, fmt.Errorf("Invalid preamble, %w", err)
	}

	m := re.FindStringSubmatch(d)

	prefix := m[0]
	suffix := "Z"

	d = strings.Replace(d, prefix, "", 1) // trying to TrimLeft w/ prefix results in weirdness
	d = strings.TrimRight(d, suffix)

	curves := strings.Split(d, "C")
	coords := make([][3]Point, len(curves))

	for i, curve := range curves {

		xy := strings.Split(curve, ",")

		if len(xy) != 6 {
			return nil, fmt.Errorf("Invalid coordinate (%s) at offset %d", curve, i)
		}

		x1, err := strconv.ParseFloat(xy[0], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[0], i, err)
		}

		y1, err := strconv.ParseFloat(xy[1], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[1], i, err)
		}

		x2, err := strconv.ParseFloat(xy[2], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[2], i, err)
		}

		y2, err := strconv.ParseFloat(xy[3], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[3], i, err)
		}

		x, err := strconv.ParseFloat(xy[4], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[4], i, err)
		}

		y, err := strconv.ParseFloat(xy[5], 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse x value '%s' at offset %d, %w", xy[5], i, err)
		}

		coords[i] = [3]Point{
			Point{x1, y1},
			Point{x2, y2},
			Point{x, y},
		}
	}

	return Cubic(coords), nil
}
