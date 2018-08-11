package v1

import (
	"log"

	"github.com/noda-sin/go-bitflyer/pkg/api/v1/ticker"
	"testing"
)

func BenchmarkClient_Ticker(b *testing.B) {
	client := NewClient(&ClientOpts{
		FastHttp: false,
	})
	for i := 0; i < 100; i++ {
		_, err := client.Ticker(&ticker.Request{})
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func BenchmarkClient_TickerFastHttp(b *testing.B) {
	client := NewClient(&ClientOpts{
		FastHttp: true,
	})
	for i := 0; i < 100; i++ {
		_, err := client.Ticker(&ticker.Request{})
		if err != nil {
			log.Fatalln(err)
		}
	}
}
