package main

import (
	"flag"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"time"
)

const (
	AwsRegion = "eu-west-1"
)

var (
	TimestreamDatabase *string
	TimestreamTable    *string
	FinestraLocationId *string
)

func main() {
	finestraUsername := flag.String("finestra_username", "", "your finestra username, provided by Salford University")
	finestraPassword := flag.String("finestra_password", "", "your finestra password, provided by Salford University")
	awsKey := flag.String("aws_key", "", "AWS access key")
	awsSecret := flag.String("aws_secret", "", "AWS secret")
	TimestreamDatabase = flag.String("timestream_database", "", "AWS secret")
	TimestreamTable = flag.String("timestream_table", "", "AWS secret")
	FinestraLocationId = flag.String("finestra_location_id", "", "AWS secret")
	flag.Parse()

	if *finestraUsername == "" || *finestraPassword == "" || *awsKey == "" || *awsSecret == "" || *TimestreamDatabase == "" || *TimestreamTable == "" || *FinestraLocationId == "" {
		flag.Usage()
		return
	}

	AwsCustomCreds = aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(*awsKey, *awsSecret, ""),
	)

	setupAwsTimestreamWriteService()

	apiKey, err := fetchApiKey(finestraUsername, finestraPassword)
	handleError(err)

	yesterday := getYesterdayDate()
	locationData := fetchLocationData(apiKey, FinestraLocationId, yesterday)
	parseCSVData(locationData)
}

func fetchApiKey(username *string, password *string) (string, error) {
	apiKey, err := authenticate(username, password)
	handleError(err)
	return apiKey, nil
}

func getYesterdayDate() string {
	return time.Now().Add(-24 * time.Hour).Format("2006-01-02")
}

func fetchLocationData(apiKey string, locationID *string, date string) string {
	locationData, err := exportLocationData(apiKey, locationID, date)
	handleError(err)
	return locationData
}

func parseCSVData(locationData string) {
	records, err := parseCSV(locationData)
	handleError(err)
	parseCSVRow(records)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
