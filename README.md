# Golang API

A golang service that exposes an API REST. This service allows a user to create and manage answers, an answer consists of a pair of strings representing a key and a value to the answer. When retrieving answers, the user will recieve the lastest value for the given key.

e.g. of answer in JSON:
```
{
	"key" : "name",
	"value" : "John"
}
```

Answers will then be stored in a .json file for persistance. An answer with same key cannot be created twice nor you can update a non existing answer.

There is also a historic functionality added to see the timeline of an answer in chronological order. Create, edit and delete operations are considered events worth registering, get operations will not be recordered. The history of all answers will also be stored in a .json file for persistance purposes. 

e.g. of event in JSON:
```
{
	"event" : "create",
	"data" : {
		"key": "name",
		"value": "John"
	}
}
```

## Installation

You just need [Go](https://go.dev/) and [Gorilla](https://github.com/gorilla/mux)!
Alongside the code in this repository:

```console
git clone https://github.com/alexrodfe/hacker-news-golang-scraper.git
```

```console
go get github.com/gorilla/mux
```

The repo also includes some unitary testing, if you wish to run the tests on your own you will need to install [testify](https://github.com/stretchr/testify).

```console
go get github.com/stretchr/testify
```

## Usage

To start your own service you can do so by running the following command inside the project's root folder:

```console
go run main.go
```

Now that you have your own instance running the exposed endpoints with their corresponding REST method are the following:

```console
"/answer" POST
"/answer/{key}" GET
"/answer/{key}" DELETE
"/answer/{key}" POST
"/answer/history/{key}" GET
```

POST methods will look for an answer compliant json file in the request's body. This repo's default service address is `localhost:8080`, however you can change this to desire. An example for a GET request could be:

```console
curl http://localhost:8080/answer/example
```

## Testing

To run the tests you may use the conveniently placed makefile:

```console
make test
```
