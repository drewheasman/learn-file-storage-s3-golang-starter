package main

import "testing"

func TestGetVideoAspectRatio(t *testing.T) {
	file := "samples/boots-video-horizontal.mp4"
	expected := "landscape"

	output, err := getVideoAspectRatio(file)
	if err != nil {
		t.Fatalf("getVideoAspectRatio(\"%v\") resulted in error: %v", file, err)
	}
	if output != expected {
		t.Fatalf("getVideoAspectRatio(\"%v\") expected: \"%v\", resulted: in \"%v\"", file, expected, output)
	}

	file = "samples/boots-video-vertical.mp4"
	expected = "portrait"

	output, err = getVideoAspectRatio(file)
	if err != nil {
		t.Fatalf("getVideoAspectRatio(\"%v\") resulted in error: %v", file, err)
	}
	if output != expected {
		t.Fatalf("getVideoAspectRatio(\"%v\") expected: \"%v\", resulted: in \"%v\"", file, expected, output)
	}
}
