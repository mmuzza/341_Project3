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
	"fmt"
)

func main() {
	fmt.Println("starting ...")
	display.initialize(1024, 1024)

	rect := Rectangle{Point{100, 300}, Point{600, 900}, red}
	err := rect.draw(&display)
	if err != nil {
		fmt.Println("rect: ", err)
	}

	rect2 := Rectangle{Point{0, 0}, Point{100, 1024}, green}
	err = rect2.draw(&display)
	if err != nil {
		fmt.Println("rect2: ", err)
	}

	rect3 := Rectangle{Point{0, 0}, Point{100, 1022}, 102}
	err = rect3.draw(&display)
	if err != nil {
		fmt.Println("rect3: ", err)
	}

	circ := Circle{Point{500, 500}, 200, green}
	err = circ.draw(&display)
	if err != nil {
		fmt.Println("circ: ", err)
	}

	tri := Triangle{Point{100, 100}, Point{600, 300}, Point{859, 850}, yellow}
	err = tri.draw(&display)
	if err != nil {
		fmt.Println("tri: ", err)
	}

	display.screenShot("output")
}
