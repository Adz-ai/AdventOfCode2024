package main

import (
	"aoc2024/utility"
	"reflect"
	"testing"
)

func TestPart1(t *testing.T) {
	input, err := utility.ParseTextFile("day2", "test")
	if err != nil {
		t.Fatal(err)
	}
	want := 2
	got := Part1(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Part1() = %v, want %v", got, want)
	}
}

func TestPart2(t *testing.T) {
	input, err := utility.ParseTextFile("day2", "test")
	if err != nil {
		t.Fatal(err)
	}
	want := 4
	got := Part2(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Part2() = %v, want %v", got, want)
	}
}
