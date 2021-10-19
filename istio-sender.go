package main

import (
	"fmt"
        "context"
	"io/ioutil"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
        "go.opentelemetry.io/contrib/propagators/b3"
	"github.com/signalfx/splunk-otel-go/distro"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
        "github.com/signalfx/splunk-otel-go/instrumentation/net/http/splunkhttp"
)

func main() {
        sdk, err := distro.Run()
	if err != nil {
		panic(err)
	}
	// Ensure all spans are flushed before the application exits.
	defer func() {
		if err := sdk.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	otel.SetTextMapPropagator(b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader)))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var handler http.Handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("received request")
			w.Write([]byte("Hello"))

			client := http.Client{
				Transport: otelhttp.NewTransport(http.DefaultTransport),
			}
			sendMessage(client,ctx)
		},
	)

	handler = splunkhttp.NewHandler(handler)
	handler = otelhttp.NewHandler(handler, "my-go-client-app")
	
	http.ListenAndServe(":9090", handler)
}

func sendMessage(client http.Client, ctx context.Context) {

	req, err := http.NewRequestWithContext(ctx,"GET", "http://<address>", nil)

	if err != nil {
	 	log.Fatal(err)
	}
        
        resp, err2 := client.Do(req) 
        if err2 != nil{
		log.Fatal(err2)
	}
        fmt.Println("status code ->",resp.StatusCode)
        if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
    		if err != nil {
        		log.Fatal(err)
    		}
    		bodyString := string(bodyBytes)
        	fmt.Println("response text ->",bodyString)
	}
	defer resp.Body.Close()
}

