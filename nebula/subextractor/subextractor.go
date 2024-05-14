package subextractor

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/ttacon/libphonenumber"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/structlog"
)

type Subscriber struct {
	IMSI   string
	MSISDN string
}

func GetSubscriberFromRequest(pathParams map[string]string, body []byte) (s Subscriber) {
	s = getSubscriberIdentityFromPath(pathParams)
	if s == (Subscriber{}) {
		// We didn't find a subscriber in the path
		s = getSubscriberIdentityFromBody(body)
	}
	return
}

func getSubscriberfromMatch(matches []string, value string) (s Subscriber) {
	s = Subscriber{}
	if matches[1] != "" {
		// We found an msisdn
		s.MSISDN = validateMsisdn(value)
	}
	if matches[2] != "" {
		// we found an imsi
		s.IMSI = validateImsi(value)
	}
	return
}

// Returns an MSISDN/IMSI from a request path
func getSubscriberIdentityFromPath(pathParams map[string]string) (s Subscriber) {
	s = Subscriber{}
	regex := regexp.MustCompile(`(?i)(msisdn)|(imsi)`)
	for key, value := range pathParams {
		matches := regex.FindStringSubmatch(key)
		if len(matches) != 0 {
			s = getSubscriberfromMatch(matches, value)
			return
		}
	}
	return
}

// Returns an MSISDN/IMSI from a request body
func getSubscriberIdentityFromBody(body []byte) (s Subscriber) {
	logger := structlog.GetLogger("nebula")
	s = Subscriber{}
	var converted_body map[string]interface{}
	regex := regexp.MustCompile(`(?i)(msisdn)|(imsi)`)
	err := json.Unmarshal(body, &converted_body)
	if err != nil {
		logger.Errorv(err)
		return
	}
	for key, value := range converted_body {
		matches := regex.FindStringSubmatch(key)
		if len(matches) != 0 {
			s = getSubscriberfromMatch(matches, value.(string))
			return
		}
	}
	return
}

func validateMsisdn(subscriber string) string {
	logger := structlog.GetLogger("nebula")
	var msisdn strings.Builder
	if subscriber[0:1] == "+" { // if MSISDN already contains leading '+'
		msisdn.WriteString(subscriber)
	} else if subscriber[0:2] == "00" { // if MSISDN contains leading '00'
		msisdn.WriteString("+")
		msisdn.WriteString(subscriber[2:]) // remove leading '00'
	} else {
		msisdn.WriteString("+")
		msisdn.WriteString(subscriber)
	}
	num, err := libphonenumber.Parse(msisdn.String(), "")
	if err != nil {
		logger.Error("Error parsing MSISDN")
		logger.Errorv(err)
		return ""
	}
	if !libphonenumber.IsValidNumber(num) {
		logger.Error("Invalid MSISDN")
		return ""
	}
	// Return the phonenumber without the + as a string
	return libphonenumber.Format(num, libphonenumber.E164)
}

func validateImsi(imsi string) string {
	logger := structlog.GetLogger("nebula")
	imsiRegex := regexp.MustCompile(`\b\d{14,15}\b`)
	if imsiRegex.MatchString(imsi) {
		return imsi
	}
	logger.Error("Invalid IMSI")
	return ""
}
