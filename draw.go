/*
Muhammad Muzzammil
CS 341
12/01/2023
Project 3
Professor: Jon Soloworth
Language: Go
UIN: 661 657 007
*/

package main

import (
	"errors"
	"fmt"
	"os"
)

// Declaration of vars we will be using
// Both erros that will be used for out of bounds and unknown color
// All colors assigned to an integer so later can be assigned to matrix
// cmap holding r,g,b values
var (
	outOfBoundsErr  = errors.New("geometry out of bounds")
	colorUnknownErr = errors.New("color unknown")

	display Display

	red    Color = 0
	green  Color = 1
	blue   Color = 2
	yellow Color = 3
	orange Color = 4
	purple Color = 5
	brown  Color = 6
	black  Color = 7
	white  Color = 8

	cmap = [][]int{
		{255, 0, 0},
		{0, 255, 0},
		{0, 0, 255},
		{255, 255, 0},
		{255, 164, 0},
		{128, 0, 128},
		{165, 42, 42},
		{0, 0, 0},
		{255, 255, 255},
	}
)

// Defining Color as type int so it can accept integer values
// All interfaces are declared as given in the document:
//
//	geometry
//	Point
//	Rectangle
//	Circle
//	Triangle
//	screen
//	Display
type (
	Color int

	geometry interface {
		draw(scn screen) error
		shape() string
	}

	Point struct {
		x, y int
	}

	Rectangle struct {
		ll, ur Point
		c      Color
	}

	Circle struct {
		cp    Point
		r     int
		color Color
	}

	Triangle struct {
		pt0, pt1, pt2 Point
		c             Color
	}

	screen interface {
		initialize(maxX, maxY int)
		getMaxXY() (maxX, maxY int)
		drawPixel(x, y int, c Color) error
		getPixel(x, y int) (Color, error)
		clearScreen()
		screenShot(f string) error
	}

	Display struct {
		maxX, maxY int
		matrix     [][]Color
	}
)

func (r Rectangle) draw(scn screen) error {
	if outOfBounds(r.ll, scn) || outOfBounds(r.ur, scn) {
		return outOfBoundsErr
	}
	if colorUnknown(r.c) {
		return colorUnknownErr
	}

	x0, y0, x1, y1 := r.ll.x, r.ll.y, r.ur.x, r.ur.y

	if x1 < x0 {
		x0, x1 = x1, x0
	}
	if y1 < y0 {
		y0, y1 = y1, y0
	}

	for y := y0; y <= y1; y++ {
		drawLine(x0, y, x1, y, r.c, scn)
	}

	return nil
}

// Helper function for drawing a line using Bresenham's algorithm
func drawLine(x0, y0, x1, y1 int, c Color, scn screen) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	var sx, sy int

	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}

	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		scn.drawPixel(x0, y0, c)

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (r Rectangle) shape() string {
	return "Rectangle"
}
func (c Circle) shape() string {
	return "Circle"
}
func (tri Triangle) shape() string {
	return "Triangle"
}

func (c Circle) draw(scn screen) error {
	if outOfBounds(c.cp, scn) {
		return outOfBoundsErr
	}
	if colorUnknown(c.color) {
		return colorUnknownErr
	}

	centerX, centerY, radius := c.cp.x, c.cp.y, float64(c.r)

	x := int(radius)
	y := 0
	err := 0

	for x >= y {
		drawLine(centerX-x, centerY+y, centerX+x, centerY+y, c.color, scn)
		drawLine(centerX-y, centerY+x, centerX+y, centerY+x, c.color, scn)
		drawLine(centerX-y, centerY-x, centerX+y, centerY-x, c.color, scn)
		drawLine(centerX-x, centerY-y, centerX+x, centerY-y, c.color, scn)

		y++

		if err <= 0 {
			err += 2*y + 1
		}

		if err > 0 {
			x--
			err -= 2*x + 1
		}
	}

	return nil
}

func (display *Display) initialize(maxX, maxY int) {
	display.maxX, display.maxY = maxX, maxY
	display.clearScreen()
}

func (display *Display) getMaxXY() (maxX, maxY int) {
	return display.maxX, display.maxY
}

func (display *Display) drawPixel(x, y int, c Color) error {
	display.matrix[x][y] = c
	return nil
}

func (display *Display) getPixel(x, y int) (Color, error) {
	return display.matrix[x][y], nil
}

func (display *Display) clearScreen() {
	display.matrix = make([][]Color, display.maxY)
	for i := range display.matrix {
		display.matrix[i] = display.initRow(display.maxX)
	}
}

func (display *Display) initRow(length int) []Color {
	row := make([]Color, length)
	display.fillRow(row, white)
	return row
}

func (display *Display) fillRow(row []Color, value Color) {
	for i := range row {
		row[i] = value
	}
}

func (display Display) screenShot(f string) error {
	file, err := os.Create(f + ".ppm")
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "P3\n%d %d\n255\n", display.maxX, display.maxY)

	for y := 0; y < display.maxY; y++ {
		for x := 0; x < display.maxX; x++ {
			color := display.matrix[y][x]
			fmt.Fprintf(file, "%d %d %d ", cmap[color][0], cmap[color][1], cmap[color][2])
		}
		fmt.Fprintln(file)
	}

	return nil
}

func outOfBounds(point Point, scr screen) bool {
	maxWidth, maxHeight := scr.getMaxXY()
	return point.x < 0 || point.x >= maxWidth || point.y < 0 || point.y >= maxHeight
}

func colorUnknown(c Color) bool {
	return c < 0 || int(c) >= len(cmap)
}

// https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func interpolate(l0, d0, l1, d1 int) (values []int) {
	a := float64(d1-d0) / float64(l1-l0)
	d := float64(d0)

	count := l1 - l0 + 1
	for ; count > 0; count-- {
		values = append(values, int(d))
		d = d + a
	}
	return
}

//  https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html

func (tri Triangle) draw(scn screen) (err error) {
	if outOfBounds(tri.pt0, scn) || outOfBounds(tri.pt1, scn) || outOfBounds(tri.pt2, scn) {
		return outOfBoundsErr
	}
	if colorUnknown(tri.c) {
		return colorUnknownErr
	}

	y0 := tri.pt0.y
	y1 := tri.pt1.y
	y2 := tri.pt2.y

	// Sort the points so that y0 <= y1 <= y2
	if y1 < y0 {
		tri.pt1, tri.pt0 = tri.pt0, tri.pt1
	}
	if y2 < y0 {
		tri.pt2, tri.pt0 = tri.pt0, tri.pt2
	}
	if y2 < y1 {
		tri.pt2, tri.pt1 = tri.pt1, tri.pt2
	}

	x0, y0, x1, y1, x2, y2 := tri.pt0.x, tri.pt0.y, tri.pt1.x, tri.pt1.y, tri.pt2.x, tri.pt2.y

	x01 := interpolate(y0, x0, y1, x1)
	x12 := interpolate(y1, x1, y2, x2)
	x02 := interpolate(y0, x0, y2, x2)

	// Concatenate the short sides

	x012 := append(x01[:len(x01)-1], x12...)

	// Determine which is left and which is right
	var x_left, x_right []int
	m := len(x012) / 2
	if x02[m] < x012[m] {
		x_left = x02
		x_right = x012
	} else {
		x_left = x012
		x_right = x02
	}

	// Draw the horizontal segments
	for y := y0; y <= y2; y++ {
		for x := x_left[y-y0]; x <= x_right[y-y0]; x++ {
			scn.drawPixel(x, y, tri.c)
		}
	}
	return
}
