package measure

import (
	"blockEmulator/params"
	"encoding/csv"
	"log"
	"os"
)

func WriteMetricsToCSV(fileName string, colName []string, colVals [][]string) {
	// Construct directory path
	dirpath := params.DataWrite_path + "supervisor_measureOutput/"
	if err := os.MkdirAll(dirpath, os.ModePerm); err != nil {
		log.Panic(err)
	}

	// Construct target file path
	targetPath := dirpath + fileName + ".csv"

	// Open file, create if it does not exist
	file, err := os.OpenFile(targetPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)

	// Write header if the file is newly created
	fileInfo, err := file.Stat()
	if err != nil {
		log.Panic(err)
	}
	if fileInfo.Size() == 0 {
		if err := writer.Write(colName); err != nil {
			log.Panic(err)
		}
		writer.Flush()
	}

	// Write data
	for _, metricVal := range colVals {
		if err := writer.Write(metricVal); err != nil {
			log.Panic(err)
		}
		writer.Flush()
	}
}
