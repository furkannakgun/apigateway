package subextractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateImsi(t *testing.T) {
	type args struct {
		imsi string
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name:     "Valid IMSI",
			args:     args{imsi: "208260123456789"},
			expected: "208260123456789",
		},
		{
			name:     "Invalid IMSI",
			args:     args{imsi: "20826012345sg89"},
			expected: "",
		},
		{
			name:     "Short IMSI",
			args:     args{imsi: "208260"},
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, validateImsi(tt.args.imsi))
		})
	}
}

func TestValidateMsisdn(t *testing.T) {
	type args struct {
		msisdn string
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name:     "Valid MSISDN",
			args:     args{msisdn: "442012345678"},
			expected: "+442012345678",
		},
		{
			name:     "Valid MSISDN +",
			args:     args{msisdn: "+442012345678"},
			expected: "+442012345678",
		},
		{
			name:     "Valid MSISDN 00",
			args:     args{msisdn: "00442012345678"},
			expected: "+442012345678",
		},
		{
			name:     "Invalid MSISDN",
			args:     args{msisdn: "+4420123xs678"},
			expected: "",
		},
		{
			name:     "Short MSISDN",
			args:     args{msisdn: "44201"},
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, validateMsisdn(tt.args.msisdn))
		})
	}
}
func TestGetSubscriberIdentityFromPath(t *testing.T) {
	type args struct {
		pathParams map[string]string
	}
	tests := []struct {
		name     string
		args     args
		expected Subscriber
	}{
		{
			name:     "Path contains MSISDN",
			args:     args{pathParams: map[string]string{"subscriberMSISDN": "34600556800"}},
			expected: Subscriber{MSISDN: "+34600556800"},
		},
		{
			name:     "Path does not contain MSISDN",
			args:     args{pathParams: map[string]string{"subscriber": "34600556800"}},
			expected: Subscriber{},
		},
		{
			name:     "Path contains mixed-case MSISDN",
			args:     args{pathParams: map[string]string{"subscriberMsisdn": "34600556800"}},
			expected: Subscriber{MSISDN: "+34600556800"},
		},
		{
			name:     "Path contains IMSI",
			args:     args{pathParams: map[string]string{"IMSI": "310260123456789"}},
			expected: Subscriber{IMSI: "310260123456789"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSubscriberIdentityFromPath(tt.args.pathParams)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetSubscriberIdentityFromBody(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name     string
		args     args
		expected Subscriber
	}{
		{
			name: "Body contains MSISDN",
			args: args{body: []byte(
				`{"duration": "PT20M","maxBandwidth": 10000,"name": "Promotional offer number 5","subscriberMsisdn": "34600556800"}`,
			)},
			expected: Subscriber{MSISDN: "+34600556800"},
		},
		{
			name: "Body contains IMSI",
			args: args{body: []byte(
				`{"duration": "PT20M","maxBandwidth": 10000,"name": "Promotional offer number 5","IMSI": "310260123456789"}`,
			)},
			expected: Subscriber{IMSI: "310260123456789"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSubscriberIdentityFromBody(tt.args.body)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetSubscriberFromMatch(t *testing.T) {
	type args struct {
		matches []string
		value   string
	}
	tests := []struct {
		name     string
		args     args
		expected Subscriber
	}{
		{
			name:     "Matches MSISDN",
			args:     args{matches: []string{"34600556800", "34600556800", ""}, value: "34600556800"},
			expected: Subscriber{MSISDN: "+34600556800"},
		},
		{
			name:     "No match",
			args:     args{matches: []string{"", "", ""}, value: "3460055"},
			expected: Subscriber{},
		},
		{
			name:     "Matches IMSI",
			args:     args{matches: []string{"310260123456789", "", "310260123456789"}, value: "310260123456789"},
			expected: Subscriber{IMSI: "310260123456789"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSubscriberfromMatch(tt.args.matches, tt.args.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}
