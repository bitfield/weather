# Weather

This is a Go project template, intended for students at the [Bitfield Institute of Technology](https://bitfieldconsulting.com/golang/bit)—but everyone's welcome to try it if you're interested in learning or practicing Go.

## Description

The aim of this project is to produce a Go library package and accompanying command-line tool that will (briefly) report the current weather conditions for a given location. For example, if the location were 'London, UK', your interaction with the tool might go like this:

```
./weather London, UK
Cloudy 15.2ºC
```

This could also be used to display the weather in your menu bar or terminal status bar, or any kind of dashboard that supports running external programs to get data.

## Getting started

You can use whatever public weather API you like for this; [OpenWeather](https://openweathermap.org/api) is a good choice. You'll need to sign in to the OpenWeather site to create an API token, which in turn you'll need to be able to supply to your program to make the API requests.

## Choosing an endpoint

You'll need to work out what API endpoint your program is going to call to get the weather. OpenWeather has many, but I recommend the [Current Weather](https://openweathermap.org/current) endpoint:

* https://api.openweathermap.org/data/2.5/weather

Your program, like all HTTP clients, should use `https` URLs, so that your network traffic is encrypted with TLS.

## Composing a URL

Before trying to call the API from Go, see if you can call it from a standard command-line HTTP client such as `curl`. Experiment a little to see how to specify the location you want the weather for, and how to include your API token in the request.

Once you know how the URL format works, you can start on the first programming goal: generating the required URL in Go.

## Goal 1: formatting the URL

The first test is partly done for you:

```go
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
```

Set the contents of your `want` variable to what you previously determined the correct URL should be, given the location "London" and the token "dummy_token" (don't use your real token).

This value is what your `FormatURL` function needs to return. Write an appropriate function; when the test passes, you're done!

## Goal 2: parsing the JSON response

OpenWeather responds to requests with a body containing JSON data representing the weather conditions. Your challenge here is to work out how to parse that JSON data to extract the information you need (the weather summary text, such as "Cloudy", and the temperature (which will be in Kelvin, and you need to convert it to Celsius).

This test is also written for you:

```go
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
```

We have some test data in the form of JSON returned by the OpenWeather API. The test opens this file, and passes the file description `f` to the function `ParseJSON`. This function's parameter should be of type `io.Reader`, rather than `*os.File`, because in the real program, it will be called with the `Body` field of the HTTP response object.

It should return a `Conditions` struct (you'll need to define it) matching the `want` variable.

Define `ParseJSON` and `Conditions`, write the JSON parsing code, and when the test passes, you're done!

## Goal 3: an end-to-end integration test

You now have the two pieces of machinery you need to get the current weather conditions for a location: constructing the URL to call, and decoding the response. Now it's time to put them together.

The first integration test is also written for you:

```go
func TestConditionsIntegration(t *testing.T) {
	t.Parallel()
	token := os.Getenv("OPENWEATHER_API_TOKEN")
	if token == "" {
		t.Skip("Please set a valid API key in the environment variable OPENWEATHER_API_TOKEN")
	}
	cond, err := weather.Current("London", token)
	if err != nil {
		t.Fatal(err)
	}
	if cond.Summary == "" {
		t.Errorf("empty summary: %#v", cond)
	}
}
```

It is protected by a build tag on the file [`weather_integration_test.go`](weather_integration_test.go), so this test will not normally be run when you run `go test`. To run it, supply the `integration` tag on the command line:

```
go test -tags=integration
```

The test calls the function `Current` with a location and a token value. Because it will really call the OpenWeather API, you'll need to set an environment variable containing your real token (don't hard-code this in the test).

This function should return some non-nil error if there is any problem getting or decoding the weather data. Otherwise, the real current weather for London should be present in the `Conditions` struct.

We can't test what it actually is, since we don't know it in advance. But we can test that at least the `Summary` field is not empty.

When you write the `Current` function, bear in mind that you already have code to generate the URL to call (in `FormatURL`) and code to decode the response to a `Conditions` struct (`ParseJSON`). All `Current` needs to do is orchestrate these two functions and make the actual HTTP call.

When the test passes, you're done!

## Goal 4: getting the location as input

We're aiming to produce a CLI tool with an interface like this:

```
./weather London, UK
```

That means we'll need some way to parse the command-line arguments to the program, and create a location string from them. Let's call this function `LocationFromArgs`, and try to write a test for it before we actually implement it.

What should it take? Well, it can't look directly at `os.Args`; that would be difficult to test. Instead, we'll have to pass a `[]string` to it (in the real program, this will actually be `os.Args`, but we can use any arbitrary slice of strings in the test).

When you run an executable, the first element of its `os.Args` is the path to the program binary itself (for example, `/usr/bin/weather`). The remaining elements contain everything else specified on the command line, split by whitespace. For example, the command line:

```
./weather London, UK
```

might produce the following `os.Args` slice:

```go
[]string{"/usr/bin/weather", "London,", "UK"}
```

The job of `LocationFromArgs` is to take a slice like this and turn it into the required single string representing the location "London,UK" (URLs can't contain spaces, so OpenWeather takes the location in this format).

Write the test for `LocationFromArgs`, and make sure you test a few different cases, including multi-word locations, single-word locations, and no location at all. (If no location is specified, `LocationFromArgs` should return an error.)

Once the test is complete, write the function so that it passes the test.

## Goal 5: formatting the weather as a string

The machinery we have produces a Go struct `Conditions`, containing the weather summary and temperature. For example:

```go
weather.Conditions{
	Summary:            "Drizzle",
	TemperatureCelsius: 7.17,
}
```

But users would like to see a single string in this format:

```
Drizzle 7.2ºC
```

To do that, we can imagine a `String` method on the `Conditions` type that formats it appropriately (including rounding the temperature to one decimal place). Write a test for such a method, and implement it.

## Goal 6: a `RunCLI` function

We would like to keep our `main.go` as simple as possible, because anything in the `main` package can't be tested, or imported by other programs. In fact, the `main.go` has been written for you:

```go
package main

import "weather"

func main() {
	weather.RunCLI()
}
```

To make this work, you'll need a `RunCLI` function which does the following:

* Gets the location from `os.Args`
* Displays an error if the location isn't set
* Gets the token from the environment
* Displays an error if the token variable isn't set
* Calls OpenWeather API to get the weather for the location
* Decodes the response to a `Conditions` struct
* Formats the result as a string and displays it on the terminal

You needn't write a test for this, since it's simply gluing together all the parts you already built test-first. Implement `RunCLI` and make sure it works when you run `main.go`.

## Goal 7: a client object

We can imagine programs that might want to get the weather for several locations at once. We wouldn't want to have to repeatedly pass the token in for every call to `Current`, so it might make sense to create some `Client` object that can store the token. `Current` would then become a method on the client.

Make this change, test-first. Update your `RunCLI` function to use the client object.

## Goal 8: a substitute TLS server for testing

It will be convenient to test methods like `Current` against a local TLS server that supplies canned data, instead of the real OpenWeather API. To do this, write a test that creates a test TLS server using `httptest.NewTLSServer`.

Write a handler that opens the data file `testdata/london.json` and returns its contents as the response body.

Have the test create a client object with a dummy token (since we won't need a real one) that calls the test TLS server instead of OpenWeather. You'll need to figure out how to inject this server's URL into the client, overriding the OpenWeather URL.

Check the result of calling `Current` against what you know it should be, based on the canned test data.

## Stretch goals

Some more refinements to add to your program if you like:

* Add a command-line flag to display different temperature units (for example Fahrenheit)
* Add a command-line flag for more detailed weather information
* Specifying the location in different ways (for example, latitude and longitude co-ordinates)
* Automatically determine the user's location (for example, by GeoIP query)
* Allow using different weather APIs (so that someone can use your library to write a tool that queries some other API, for example)
* Cache the weather data for a configurable length of time, so that if you run the program regularly on a timer, it doesn't hammer the API. For example, you might cache the data for 15 minutes, so that no matter how often you run the program within a 15-minute period, it won't need to make a new API call.

