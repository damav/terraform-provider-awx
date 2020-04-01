# awx-go

[![Build Status](https://travis-ci.org/davidfischer-ch/awx-go.svg?branch=master)](https://travis-ci.org/davidfischer-ch/awx-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/davidfischer-ch/awx-go)](https://goreportcard.com/report/github.com/davidfischer-ch/awx-go)
[![codecov](https://codecov.io/gh/davidfischer-ch/awx-go/branch/master/graph/badge.svg)](https://codecov.io/gh/davidfischer-ch/awx-go)

AWX SDK for the Go programming language.

![AWX-GO-ROBOT](images/awx-go-robot.png)

## Installing

If you are using Go 1.5 with the GO15VENDOREXPERIMENT=1 vendoring flag, or 1.6 and higher you can use the following command to retrieve the SDK. The SDK will be included.

```
go get -u github.com/davidfischer-ch/awx-go
```

## Example

We can simply import awx-go and call its services, such as PingService:

```
import (
    "log"
    awxGo "github.com/davidfischer-ch/awx-go"
)

func main() {
    awx := awxGo.NewAWX("http://awx.your_server_host.com", "your_awx_username", "your_awx_passwd", nil)
    result, err := awx.PingService.Ping()
    if err != nil {
        log.Fatalf("Ping awx err: %s", err)
    }

    log.Println("Ping awx: ", result)
}
```

More examples can be found at [here](https://github.com/davidfischer-ch/awx-go/tree/master/examples).

## Roadmap

awx-go is still in development, and its roadmap could be found at [here](https://github.com/davidfischer-ch/awx-go/blob/master/ROADMAP.md).

## Contribute

There are many ways to contribute to awx-go.

* Submit bugs via [Github issues](https://github.com/davidfischer-ch/awx-go/issues);
* Submit a [pull request](https://github.com/davidfischer-ch/awx-go/pulls) for fixes or features;
* Mail [me](mailto:wjx_colstu@hotmail.com)
