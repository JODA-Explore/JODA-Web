package main

import (
	"flag"

	"github.com/JODA-Explore/JODA-Web/internal/server"
)

func main() {

	hostPtr := flag.String("host", "http://localhost:5632", "The address of the JODA server")
	portPtr := flag.Uint("port", 8080, "The port used by JODA Web")

	flag.Parse()

	server.Start(*hostPtr, *portPtr)
}
