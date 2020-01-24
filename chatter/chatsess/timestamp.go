package chatsess

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

const (
	// DATE_FMT - date format
	DATE_FMT = "02-01-2006"
)

// TimetoDB - converts time to *string
func TimetoDB(t time.Time) *string {
	tn := t.Unix()
	return aws.String(strconv.FormatInt(tn, 10))
}

// DBtoTime - converts dbTime to time
func DBtoTime(s *string) time.Time {
	n, _ := strconv.ParseInt(*s, 10, 64)
	return time.Unix(n, 0)
}
