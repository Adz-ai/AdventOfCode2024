package main

import (
	"aoc2024/utility"
	"reflect"
	"testing"
)

func TestPart1(t *testing.T) {
	input, err := utility.ParseTextFile("day3", "test")
	if err != nil {
		t.Fatal(err)
	}
	want := 161
	got := part1(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("part1() = %v, want %v", got, want)
	}
}

func TestPart2(t *testing.T) {
	input := make([]string, 0)
	input = append(input, "xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))")
	want := 48
	got := part2(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("part2() = %v, want %v", got, want)
	}
}
