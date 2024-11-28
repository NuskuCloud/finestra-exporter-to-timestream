package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"time"
)

var (
	AwsCustomCreds   *aws.CredentialsCache
	AwsConfig        aws.Config
	TimestreamClient *timestreamwrite.Client
	TimestreamErrCh  = make(chan error, 50)
)

func setupAwsTimestreamWriteService() {
	AwsConfig, _ = config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(AwsCustomCreds))
	tr := loadHttpSettings()

	TimestreamClient = timestreamwrite.NewFromConfig(AwsConfig, func(o *timestreamwrite.Options) {
		o.Region = AwsRegion
		o.HTTPClient = &http.Client{Transport: tr}
	})
}

func loadHttpSettings() *http.Transport {
	/**
	* Recommended Timestream write client SDK configuration:
	*  - Set SDK retry count to 10.
	*  - Use SDK DEFAULT_BACKOFF_STRATEGY
	*  - Request timeout of 20 seconds
	 */

	// Setting 20 seconds for timeout
	tr := &http.Transport{
		ResponseHeaderTimeout: 20 * time.Second,
		// Using DefaultTransport values for other parameters: https://golang.org/pkg/net/http/#RoundTripper
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			KeepAlive: 60 * time.Second,
			Timeout:   60 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// So client makes HTTP/2 requests
	http2.ConfigureTransport(tr)
	return tr
}
