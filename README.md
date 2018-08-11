# go-bitflyer

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

go-bitflyer is a Go bindings for [bitFlyer Lightning API](https://lightning.bitflyer.jp/docs?lang=en).

Using fasthttp around network request to speed up.

## Usage

```golang
package main

import (
  "log"

  "github.com/noda-sin/go-bitflyer/pkg/api/auth"
  "github.com/noda-sin/go-bitflyer/pkg/api/v1"
  "github.com/noda-sin/go-bitflyer/pkg/api/v1/markets"
  "github.com/noda-sin/go-bitflyer/pkg/api/v1/permissions"
)

func main() {
  client := v1.NewClient(&v1.ClientOpts{
    AuthConfig: &auth.AuthConfig{
      APIKey:    "**********************",
      APISecret: "********************************************",
    },
  })

  if resp, err := client.Permissions(&permissions.Request{}); err != nil {
    log.Fatalln(err)
  } else {
    log.Println(resp)
  }

  if resp, err := client.Markets(&markets.Request{}); err != nil {
    log.Fatalln(err)
  } else {
    log.Println(resp)
  }
}

```
