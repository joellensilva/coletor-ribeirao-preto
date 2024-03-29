package main

import (
	"log"
	"os"
	"time"
)

type confSpec struct {
	Month string
}

const (
	defaultFileDownloadTimeout = 20 * time.Second // Duração que o coletor deve esperar até que o download de cada um dos arquivos seja concluído
	defaultGeneralTimeout      = 6 * time.Minute  // Duração máxima total da coleta de todos os arquivos. Valor padrão calculado a partir de uma média de execuções ~4.5min
	defaulTimeBetweenSteps     = 5 * time.Second  //Tempo de espera entre passos do coletor."
)

func main() {
	outputFolder := os.Getenv("OUTPUT_FOLDER")
	if outputFolder == "" {
		outputFolder = "./output"
	}

	if err := os.Mkdir(outputFolder, os.ModePerm); err != nil && !os.IsExist(err) {
		log.Fatalf("Error creating output folder(%s): %w", outputFolder, err)
	}

	downloadTimeout := defaultFileDownloadTimeout
	generalTimeout := defaultGeneralTimeout
	timeBetweenSteps := defaulTimeBetweenSteps

	c := crawler{
		downloadTimeout:   downloadTimeout,
		collectionTimeout: generalTimeout,
		timeBetweenSteps:  timeBetweenSteps,
		output:            outputFolder,
	}
	err := c.crawl()
	if err != nil {
		log.Fatalf("Error crawling")
	}

	// O parser do CNJ espera os arquivos separados por \n. Mudanças aqui tem
	// refletir as expectativas lá.
	// fmt.Println(strings.Join(downloads, "\n"))
}
