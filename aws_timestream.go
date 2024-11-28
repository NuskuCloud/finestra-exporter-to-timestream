package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite/types"
	"strconv"
)

func insertIntoTimestream(databaseName string, tableName string, measureName string, measureValue string, timeInSeconds int64) error {
	measureValueFloat, err := strconv.ParseFloat(measureValue, 64)
	if err != nil {
		return fmt.Errorf("failed to parse measureValue to float64, %v", err)
	}
	record := types.Record{
		Dimensions: []types.Dimension{
			{
				Name:  aws.String("example-dimension"),
				Value: aws.String("example-value"),
			},
		},
		MeasureName:      aws.String(measureName),
		MeasureValue:     aws.String(fmt.Sprintf("%f", measureValueFloat)),
		MeasureValueType: types.MeasureValueTypeDouble,
		Time:             aws.String(strconv.FormatInt(timeInSeconds, 10)),
		TimeUnit:         types.TimeUnitSeconds,
	}

	// Write the record to Timestream
	_, err = TimestreamClient.WriteRecords(context.TODO(), &timestreamwrite.WriteRecordsInput{
		DatabaseName: aws.String(databaseName),
		TableName:    aws.String(tableName),
		Records:      []types.Record{record},
	})
	if err != nil {
		return fmt.Errorf("failed to write record to Timestream, %v", err)
	}

	return nil
}
