package svg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// MISSING is a the string identifier for a path that can not be parsed.
const MISSING string = "X"

// UNKNOWN is a the string identifier for a path type that is not supported.
const UNKNOWN string = ""

// POINT is the string identifier for a `Point`.
const POINT string = "P"

// CUBIC is the string identifier for a `Cubic` path.
const CUBIC string = "C"

// LINE is the string identifier for a `Line` path.
const LINE string = "L"

var (
	_ Coordinates = Point{}
	_ Coordinates = Cubic{}
	_ Coordinates = Line{}
)

// Coordinate is an interface representation the coordinates defined by a `Path` instance.
type Coordinates interface {
	Type() string
}

// Point is a two dimensional X,Y coordinate.
type Point [2]float64

// Type returns the `POINT` string.
func (p Point) Type() string {
	return POINT
}

// Cubic is a list of a tuples of `Point` instances that make up a cubic curve path.
type Cubic [][3]Point

// Type returns the `CUBIC` string.
func (c Cubic) Type() string {
	return CUBIC
}

// Line is a list of pairs of `Point` instances that make up a line path.
type Line []Point

// Type returns the `LINE` string.
func (l Line) Type() string {
	return LINE
}

// DerivePathType returns a string identifier for the SVG path (`d` attribute) defined by 'd'.
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

// DerivePathType returns a `Coordinates` instance derived from the SVG path (`d` attribute) defined by 'd'.
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
