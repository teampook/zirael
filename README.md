# Zirael

* [Description](#description)
* [Installation](#installation)
* [Usage](#usage)
* [Thanks & Author](#Thanks)

## Description
Zirael is HTTP Client that wrap your request with AMX authorization header

## Installation
```
go get -u github.com/teampook/zirael
``` 

## Usage

### Importing the package
```go
import "github.com/teampook/zirael"
```
### Sample `GET` request

```go
client := zirael.NewClient(
		"YOUR_API_KEY",
		"YOUR_API_ID",
		"YOUR_NONCE", zirael.WithHTTPTimeout(10 * time.Second))

	request, err := client.Get("https://test-apisbn.kemenkeu.go.id/v1/bank", nil)

	if err != nil {
		panic(err)
	}

	fmt.Println(request)

	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
```

## Thanks
Inspiration from [Heimdal](https://github.com/gojek/heimdall) An enhanced HTTP client for Go by [Gojek](http://gojek.tech/)

Authors: Arif Rakhman -- [https://github.com/arieefrachman]
