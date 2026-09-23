package main

import (
	"flag"
	"fmt"
	"os"

	"git.annabunches.net/annabunches/joyful/internal/gremlin"
)

func main() {
	input := flag.String("input", "joystick_gremlin_profile.xml", "Joystick Gremlin v14 XML profile")
	output := flag.String("output", "joystick_gremlin_profile", "output directory")
	leftPath := flag.String("left-device-path", gremlin.DefaultLeftPath, "left physical joystick path")
	rightPath := flag.String("right-device-path", gremlin.DefaultRightPath, "right physical joystick path")
	flag.Parse()

	result, err := gremlin.Generate(*input, *output, gremlin.Options{LeftPath: *leftPath, RightPath: *rightPath})
	if err != nil {
		fmt.Fprintln(os.Stderr, "gremlin-convert:", err)
		os.Exit(1)
	}
	fmt.Printf("generated %d rules with %d warnings in %s\n", total(result.Counts), len(result.Warnings), *output)
}

func total(counts map[string]int) int {
	n := 0
	for _, count := range counts {
		n += count
	}
	return n
}
