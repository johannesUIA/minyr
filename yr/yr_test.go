
package yr_test

import (
	"testing"
	"bufio"
	"math"
	"os"
	"github.com/johannesUIA/minyr/yr"
)

func TestGetLineCount(t *testing.T) {
	expected := 16756

	file, err := os.Open("../kjevik-temp-fahr-20220318-20230318.csv")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()


	scanner := bufio.NewScanner(file)

	count := 0
	for scanner.Scan() {
		count++
	}

	if count != expected {
		t.Errorf("Unexpected line count: got %d, want %d", count, expected)
	}

}

func TestConvTemperature(t *testing.T) {

		type test struct {
			input string
			want string
		}

tests := []test{
		{input: "Kjevik;SN39040;18.03.2022 01:50;6", want: "Kjevik;SN39040;18.03.2022 01:50;42.8"},
		{input: "Kjevik;SN39040;07.03.2023 18:20;0", want: "Kjevik;SN39040;07.03.2023 18:20;32.0"},
		{input: "Kjevik;SN39040;08.03.2023 02:20;-11", want: "Kjevik;SN39040;08.03.2023 02:20;12.2"},
		{input: "Data er gyldig per 18.03.2023 (CC BY 4.0), Meteorologisk institutt (MET);;;",
			want: "Data er basert paa gyldig data (per 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET);endringen er gjort av Tony Le, Bunyamin og Johannes;;;"},
	}

for _, tc := range tests {
	got := yr.ProcessLine(tc.input)
		if !(tc.want == got) {
			t.Errorf("expected: %v, got: %v", tc.want, got)
		}
	}

}

func TestAverageTemp(t *testing.T) {
    // Call the AverageTemp function with the file name
    result, err := yr.AverageTemp1("../kjevik-temp-celsius-20220318-20230318.csv")

    // Verify that there was no error
    if err != nil {
        t.Errorf("Error while calculating average temperature: %v", err)
    }

    // Verify that the result matches the expected result
    expected := 8.56

tolerance := 0.01
    if math.Abs(result-expected) > tolerance {
        t.Errorf("Expected average temperature of %.2f, but got %.2f", expected, result)
    }

}