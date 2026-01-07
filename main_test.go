package main

import "testing"

func parseEnv(t *testing.T, contents string) ParseResult {
	t.Helper()
	return ParseDotenv([]byte(contents))
}

func TestMerge_DefaultsBlankAndPreservesComments(t *testing.T) {
	source := "A=1\n# c\nB=2\nC=3"
	destination := "A=9\nB=\nD=4"

	got := string(Merge(parseEnv(t, source), parseEnv(t, destination), MergeOptions{}))
	want := "A=9\n# c\nB=\nC=\n"

	if got != want {
		t.Fatalf("unexpected merge output\nwant:\n%q\ngot:\n%q", want, got)
	}
}

func TestMerge_UseSourceDefaultsAndFillEmpty(t *testing.T) {
	source := "A=1\nB=2\nC=3"
	destination := "A=9\nB="

	got := string(Merge(parseEnv(t, source), parseEnv(t, destination), MergeOptions{
		UseSourceDefaults: true,
		FillEmpty:         true,
	}))
	want := "A=9\nB=2\nC=3\n"

	if got != want {
		t.Fatalf("unexpected merge output\nwant:\n%q\ngot:\n%q", want, got)
	}
}

func TestMerge_KeepUnusedAndAdoptDestinationTail(t *testing.T) {
	source := "A=1\nB=2"
	destination := "A=9  # from dest\nB=2\nD=4\nE=5"

	got := string(Merge(parseEnv(t, source), parseEnv(t, destination), MergeOptions{
		KeepUnused: true,
	}))
	want := "A=9  # from dest\nB=2\n\nD=4\nE=5\n"

	if got != want {
		t.Fatalf("unexpected merge output\nwant:\n%q\ngot:\n%q", want, got)
	}
}

func TestMerge_SourceTailWins(t *testing.T) {
	source := "A=1  # src"
	destination := "A=9  # dest"

	got := string(Merge(parseEnv(t, source), parseEnv(t, destination), MergeOptions{}))
	want := "A=9  # src\n"

	if got != want {
		t.Fatalf("unexpected merge output\nwant:\n%q\ngot:\n%q", want, got)
	}
}
