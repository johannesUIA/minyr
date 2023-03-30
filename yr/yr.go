package yr

import (
    	"bufio"
		"encoding/csv"
    	"fmt"
    	"os"
		"errors"
		"log"
    	"strings"
    	"strconv"
    	"github.com/johannesUIA/funtemps/conv"
)

func ConvTemperature() {

	//Setter input og output filnavn.

	//inputFilename := "kjevik-temp-celsius-20220318-20230318.csv"
	outputFilename := "kjevik-temp-fahr-20220318-20230318.csv"

	//Sjekker om kjevik fahr versjon av filen allerede eksisterer.

	if _, err := os.Stat(outputFilename); err == nil {
		//Om filen eksisterer spor den om vi vil genere filen paa nytt.
		fmt.Printf("Fil '%s' eksisterer allerede. Vil du generere filen paa nytt? (j/n): ", outputFilename)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			answer := scanner.Text()
			if strings.ToLower(answer) == "j" || strings.ToLower(answer) == "ja" {
				//Om brukeren vil genere filen aa nytt gaar den ut av loopen.
				break
			} else if strings.ToLower(answer) == "n" || strings.ToLower(answer) == "nei" {
				//Om brukeren ikke onsker aa genere filen paa nytt returner den og gaar ut av funksjonen.
				fmt.Println("Avslutter...")
				return
			} else {
				//Om brukeren ikke gir gyldig input i  Scanneren spor den paa nytt.
				fmt.Print("Invalid answer. Do you want to regenerate the file? (j/n): ")
			}
		}
	}

	inputFile := openInputFile()
	defer inputFile.Close()

	outputFile, err := createOutputFile()
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	//Denne blokken skriver forste linje for den begynner paa loop.

	scanner := bufio.NewScanner(inputFile)
    	if scanner.Scan() {
        	fmt.Fprintln(outputFile, scanner.Text())
    	}


	for scanner.Scan() {
    		line := scanner.Text()
    		fields := strings.Split(line, ";")

		if fields[3] == "" {
    		continue //Skipper om den ikke finner temperatur.
		}


    		celsius, err := strconv.ParseFloat(fields[3], 64)

		if err != nil {
        	log.Fatal(err)
    		}

    		fahrenheit := conv.CelsiusToFarhenheit(celsius)
    		fields[3] = fmt.Sprintf("%.2f", fahrenheit)
    		line = strings.Join(fields, ";")
    		fmt.Fprintln(outputFile, line)
		}

		if err := scanner.Err(); err != nil {
    		log.Fatal(err)
	}
	footer := []string{"Data er basert paa gyldig data (per 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET);endringen er gjort av Tony Le, Bunyamin og Johannes;;;"}
	writer := csv.NewWriter(outputFile)
	err = writer.Write(footer)
	if err != nil {
		fmt.Println("Kunne ikke skrive endelig tekst:", err)
	}
	writer.Flush()



}

func openInputFile() *os.File {
	file, err := os.Open("kjevik-temp-celsius-20220318-20230318.csv")
	if err != nil {
	log.Fatal(err)
	}
return file
}

func createOutputFile() (*os.File, error) {
	outputFilePath := "kjevik-temp-fahr-20220318-20230318.csv"
	if _, err := os.Stat(outputFilePath); err == nil {
		fmt.Printf("File %s already exists. Deleting...\n", outputFilePath)
		err := os.Remove(outputFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not delete file: %v", err)
	}
}

outputFile, err := os.Create(outputFilePath)
if err != nil {
	return nil, fmt.Errorf("could not create file: %v", err)
}
return outputFile, nil
}

func AverageTemp() {
	//Aapner og leser linjene fra filen
	file, err := os.Open("kjevik-temp-celsius-20220318-20230318.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	//Setter variablene sum og count til 0 for loopen begynner.
	sum := 0.0
	count := 0.0

	//Loop som deler opp fields, og fortsetter om fields er under 4.
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ";")
		if len(fields) < 4 {
			continue
		}

		temperature, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			continue
		}

		//Legger alle temperaturverdiene i sammen i sum variablen. Plusser ogsaa paa 1 i count variablen.
		sum += temperature
		count++
	}

	if count > 0 {
		var unit string
		fmt.Println("Vil du ha gjennomsnittstemperaturen i Celsius eller Fahrenheit? (celsius/fahrenheit)")
		fmt.Scanln(&unit)

		//Fahrenheit case
		if strings.ToLower(unit) == "fahrenheit" {
			average := (sum/float64(count))*1.8 + 32
			fmt.Printf("Gjennomsnittstemperaturen i Fahrenheit er: %.2f\n", average)
		} else {
			average := sum / float64(count)
			fmt.Printf("Gjennomsnittstemperaturen i celsius er: %.2f\n", average)
		}
	}
}

func ProcessLine(line string) string {
	if line == "" {
		return ""
	}
	fields := strings.Split(line, ";")
	lastField := ""
	if len(fields) > 0 {
		lastField = fields[len(fields)-1]
	}
	convertedField := ""
	if lastField != "" {
		var err error
		convertedField, err = convertLastField(lastField)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return ""
		}
	}
	if convertedField != "" {
		fields[len(fields)-1] = convertedField
	}
	if line[0:7] == "Data er" {
		return "Data er basert paa gyldig data (per 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET);endringen er gjort av Tony Le, Bunyamin og Johannes;;;"
	} else {
		return strings.Join(fields, ";")
	}
}

func convertLastField(lastField string) (string, error) {
	celsius, err := strconv.ParseFloat(lastField, 64)
	if err != nil {
		return "", err
	}


	fahrenheit := conv.CelsiusToFarhenheit(celsius)


	return fmt.Sprintf("%.1f", fahrenheit), nil
}

func CountLines(inputFile string) int {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	countedLines := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			countedLines++
		}
	}
	return countedLines
}

func AverageTemp1(fileName string) (float64, error) {
    //Aapner og leser linjene fra filen
    file, err := os.Open(fileName)
    if err != nil {
        return 0, err
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)

    //Setter variablene sum og count til 0 for loopen begynner.
    sum := 0.0
    count := 0.0

    //Loop som deler opp fields, og fortsetter om fields er under 4.
    for scanner.Scan() {
        fields := strings.Split(scanner.Text(), ";")
        if len(fields) < 4 {
            continue
        }

        temperature, err := strconv.ParseFloat(fields[3], 64)
        if err != nil {
            continue
        }

        //Legger alle temperaturverdiene i sammen i sum variablen. Plusser ogsaa paa 1 i count variablen.
        sum += temperature
        count++
    }

    if count > 0 {
        average := sum / count
        return average, nil
    }

    return 0, errors.New("no temperature data found")
}