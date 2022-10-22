package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awscw "github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	awscwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	flag "github.com/spf13/pflag"
)

const (
	TARGET                 = "http://127.0.0.1:80/status"
	TARGET_EXPECTED_STATUS = 200
)

var (
	interval                          = flag.DurationP("interval", "i", time.Minute, "Interval of checks (use golang time.Duration format. cf: https://pkg.go.dev/time#ParseDuration)")
	namespace                         = flag.StringP("cwnamespace", "n", "", "Namespace of metrics to put to CloudWatch")
	metric                            = flag.StringP("cwmetricname", "m", "", "Metric name of metrics to put to CloudWatch")
	val        float64                = 1.0
	dims                              = flag.StringToStringP("dimensions", "d", make(map[string]string, 0), "Dimensions of metrics to put to CloudWatch")
	dimensions []awscwtypes.Dimension = nil
)

func init() {
	flag.Parse()

	dimensions = make([]awscwtypes.Dimension, 0, len(*dims))
	for k, v := range *dims {
		dimName := k
		dimValue := v
		dimensions = append(dimensions, awscwtypes.Dimension{
			Name:  &dimName,
			Value: &dimValue,
		})
	}
}

func main() {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion("ap-northeast-1"),
	)
	if err != nil {
		panic(err)
	}

	cwClient := awscw.NewFromConfig(awsCfg)

	log("HealthCheck loop starts")
	ticker := time.NewTicker(*interval)
	for {

		<-ticker.C

		resp, err := http.Get(TARGET)
		if err != nil {
			reportError(cwClient, err)
			continue
		}

		if resp.StatusCode != TARGET_EXPECTED_STATUS {
			reportError(cwClient, nil)
			continue
		}

		reportSucceed(cwClient)

	}
}

func reportError(client *awscw.Client, checkError error) {
	log(fmt.Sprintf("check error happens: %#v", checkError))
}

func reportSucceed(client *awscw.Client) {
	log("check succeeded")
	metrics := &awscw.PutMetricDataInput{
		MetricData: []awscwtypes.MetricDatum{
			{
				MetricName: metric,
				Value:      &val,
				Dimensions: dimensions,
			},
		},
		Namespace: namespace,
	}
	_, err := client.PutMetricData(context.Background(), metrics)
	if err != nil {
		log(fmt.Sprintf("failed to put metrics: %s: %#v", err.Error(), err))
		return
	}

	log("succeed to report")
}

func log(msg string) {
	fmt.Printf("[%s] %s\n", time.Now().Format(time.RFC3339Nano), msg)
}
