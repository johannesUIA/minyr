
package yr_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestYr(t *testing.T) {
	var input string
scanner := bufio.NewScanner(os.Stdin)

for scanner.Scan() {
    input = scanner.Text()
    if input == "q" || input == "exit" {
        fmt.Println("exit")
        os.Exit(0)
    } else if input == "convert" {
        fmt.Println("Konverterer alle målingene gitt i grader Celsius til grader Fahrenheit.")
        // funksjon som gjør åpner fil, leser linjer, gjør endringer og lagrer nye linjer i en ny fil

    // flere else-if setninger     
    } else {
        fmt.Println("Venligst velg convert, average eller exit:")

    }

}
}