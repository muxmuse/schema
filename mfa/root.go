package mfa

import (
	"log"
	"bytes"
	"text/template"
)

// sorry -- still learning

func CatchFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Format(name string, templateStr string, value interface{}) string {
	var buffer bytes.Buffer;
	CatchFatal(
		template.Must(
			template.New(
				name).Parse(
				templateStr)).Execute(
			&buffer, value))

	return buffer.String()
}