package inbound

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createRequest(filename string) *http.Request {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	// Build POST request
	req, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(file))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=xYzZY")
	req.Header.Set("User-Agent", "Twilio-SendGrid-Test")
	return req
}

func TestParse(t *testing.T) {
	// Build a table of tests to run with each one having a name, the sample data file to post,
	// and the expected HTTP response from the handler
	tests := []struct {
		name          string
		file          string
		expectedError error
	}{
		{
			name: "DefaultData",
			file: "./sample_data/default_data.txt",
		},
		{
			name:          "BadData",
			file:          "./sample_data/bad_data.txt",
			expectedError: fmt.Errorf("multipart: NextPart: EOF"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(subTest *testing.T) {
			//Load POST body
			req := createRequest(test.file)

			// Invoke callback handler
			email, err := Parse(req)
			if test.expectedError != nil {
				assert.Error(subTest, err, "expected an error to occur")
				return
			}

			assert.NoError(subTest, err, "did NOT expect an error to occur")

			from := "test@example.com"
			assert.Equalf(subTest, email.Envelope.From, from, "Expected From: %s, Got: %s", from, email.Envelope.From)
		})
	}
}
