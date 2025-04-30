package main

import (
	"reflect"
	"testing"

	"github.com/Adz-ai/AdventOfCode2024/utility"
)

func TestPart1(t *testing.T) {
	input, err := utility.ParseTextFile("test")
	if err != nil {
		t.Fatal(err)
	}
	want := 1928
	got := part1(input[0])
	if !reflect.DeepEqual(got, want) {
		t.Errorf("part1() = %v, want %v", got, want)
	}
}

func TestPart2(t *testing.T) {
	input, err := utility.ParseTextFile("test")
	if err != nil {
		t.Fatal(err)
	}
	want := 2858
	got := part2(input[0])
	if !reflect.DeepEqual(got, want) {
		t.Errorf("part2() = %v, want %v", got, want)
	}
}
