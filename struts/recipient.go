package struts

import (
	"encoding/csv"
	"io"
	"os"
)

type Recipient struct {
	Email          string
	EmployeeNumber string
}

func (r *Recipient) GetRecipients(path string) ([]Recipient, error) {
	var recipients []Recipient

	f, err := os.Open(path)

	if err != nil {
		return nil, err
	} else {
		defer f.Close()
		csvReader := csv.NewReader(f)
		for {
			r, err := csvReader.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			recipient := Recipient{
				Email:          r[0],
				EmployeeNumber: r[1],
			}
			recipients = append(recipients, recipient)
		}
	}

	return recipients, nil
}
