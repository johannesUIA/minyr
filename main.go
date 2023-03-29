package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/johannesUIA/funtemps/conv"
)

func calculateAverageTemp(unit string) (float64, error) {
	// Åpne inndatafil for lesing
	inputFileName := "kjevik-temp-celsius-20220318-20230318.csv"
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		return 0, err
	}
	defer inputFile.Close()

	// Les data fra inndatafil og beregn gjennomsnittet
	reader := csv.NewReader(inputFile)
	reader.Comma = ';'

	var sum, count float64
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		temperatureStr := strings.TrimSpace(record[3])
		if temperatureStr == "" {
			continue // hopp over tomme celler
		}

		temperature, err := strconv.ParseFloat(temperatureStr, 64)
		if err != nil {
			continue // hopp over ugyldige temperaturdata
			
		}

		if unit == "f" {
			temperature = conv.CelsiusToFarhenheit(temperature)
		}

		sum += temperature
		count++
	}

	if count == 0 {
		return 0, fmt.Errorf("ingen gyldige temperaturdata funnet")
	}

	averageTemp := sum / count
	return averageTemp, nil
}

func main() {
	// Sjekk om programmet blir kjørt med korrekt kommando
	if len(os.Args) < 2 || (os.Args[1] != "convert" && os.Args[1] != "average") {
		log.Fatal("Bruk kommando 1 eller 2: go run main.go convert\n      go run main.go average [f/c]")
	}
	
	if os.Args[1] == "average" {
		// Sjekk om ønsket enhet er spesifisert
		if len(os.Args) < 3 {
			log.Fatal("Du må spesifisere ønsket enhet for gjennomsnittet: 'f' for Fahrenheit eller 'c' for Celsius")
		}
		unit := os.Args[2]
		if unit != "f" && unit != "c" {
			log.Fatal("Ugyldig enhet. Du må velge 'f' for Fahrenheit eller 'c' for Celsius")
		}
		// Kall funksjonen for å beregne gjennomsnittet
		averageTemp, err := calculateAverageTemp(unit)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Gjennomsnittstemperaturen er %.2f %s\n", averageTemp, unit)
		return
	}
	
	// Sett opp navn på inndata- og utdatafiler
	inputFileName := "kjevik-temp-celsius-20220318-20230318.csv"
	outputFileName := "kjevik-temp-fahr-20220318-20230318.csv"

	// Sjekk om utdatafil allerede eksisterer
	if _, err := os.Stat(outputFileName); !os.IsNotExist(err) {
		// Filen eksisterer
		fmt.Printf("Utskriftsfilen '%s' finnes allerede.\n", outputFileName)
		var answer string
		fmt.Print("Generer igjen? (j/n): ")
		fmt.Scanln(&answer)
		if answer != "j" {
			os.Exit(0)
		}
	}

	// Åpne inndatafil for lesing
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	// Åpne utdatafil for skriving
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Les data fra inndatafil og skriv til utdatafil
	reader := csv.NewReader(inputFile)
	reader.Comma = ';'
	writer := csv.NewWriter(outputFile)
	writer.Comma = ';'

	var record []string
	firstLine := true
	lastLine := false
	for {
		record, err = reader.Read()
		if err == io.EOF {
			if !lastLine {
				// Write last line
				writer.Write([]string{"Data er basert på gyldig data (per 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET); endringen er gjort av Tony Le, Bunyamin og Johannes"})
				lastLine = true
			}
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if firstLine {
			// Skriv første linje uendret
			err := writer.Write(record)
			if err != nil {
				log.Fatal(err)
			}
			firstLine = false
			continue
		}

		temperatureStr := strings.TrimSpace(record[3])
		if temperatureStr == "" {
			continue // hopp over tomme celler
		}
		
		// Konverter Celsius til Fahrenheit og skriv til utdatafil
		temperature, err := strconv.ParseFloat(temperatureStr, 64)
		if err != nil {
    		log.Fatal(err)
		}
		fmt.Printf("Konverterer temperaturen %v til Fahrenheit.\n", temperature)
		fahrenheit := conv.CelsiusToFarhenheit(temperature)
		fmt.Printf("Konverterte Celsius-verdien %v til Fahrenheit-verdien %v.\n", temperature, fahrenheit)
		record[3] = fmt.Sprintf("%0.1f", fahrenheit)

		err = writer.Write(record)
		if err != nil {
    		log.Fatal(err)
		}
	}

	writer.Flush()

	fmt.Printf("Utskriftsfilen '%s' er generert.\n", outputFileName)
}
