# JODA Web
This project provides a web interface for the [JODA](https://github.com/JODA-Explore/JODA) semi-structured data processor.
The web interface gives a basic overview of the data and allows the execution of queries.
Results can be viewed and downloaded as JSON files.


Additionally, it includes a web interface for the [BETZE](https://github.com/JODA-Explore/BETZE) benchmark generator.
It can be used to create exploratory benchmark sessions to compare the performance of multiple systems.

To help explore unknown data sets, the web interface also includes a data analysis tool.
It shows a graph of the structure and data types of a set of JSON documents.
It also includes an exploration helper, which can be used to find interesting data in the documents.

## Installation
The web interface is written in [Go](https://golang.org/).
To install it, you need to have Go installed.
Then, you can simply run `go get github.com/JODA-Explore/JODA-Web` to install the web interface.
The web interface can then be started by running `joda-web` in the terminal.

### Docker
The web interface can also be run in a docker container.
To do so, you need to have docker installed.
Then, you can simply run `docker run -p 8080:8080 ghcr.io/joda-explore/joda-web/joda-web:latest` to start the web interface.

We also included a docker-compose file, which can be used to start the web interface together with a JODA server instance.

## Usage
For the web interface to work, you need to have a running server instance of JODA.
The interface will try to automatically detect a running instance.
If it can't be found automatically, you can specify the address of the server in the settings.
