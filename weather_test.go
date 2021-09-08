package weather_test

import (
	"os"
	"testing"
	"weather"

	"github.com/google/go-cmp/cmp"
)

func TestFormatURL(t *testing.T) {
	t.Parallel()
	location := "London"
	token := "dummy_token"
	want := "??? You need to fill this part in! ???"
	got := weather.FormatURL(location, token)
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestParseJSON(t *testing.T) {
	t.Parallel()
	f, err := os.Open("testdata/london.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	want := weather.Conditions{
		Summary:            "Drizzle",
		TemperatureCelsius: 7.17,
	}
	got := weather.ParseJSON(f)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
