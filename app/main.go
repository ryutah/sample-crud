package main

import (
	"net/http"
	"sample"
)

func init() {
	http.Handle("/foo", sample.GetCRUDSampleHandler())
}
