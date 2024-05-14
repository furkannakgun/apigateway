package ccextractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/subextractor"
)

func TestGetCountryCode(t *testing.T) {
	type args struct {
		subscriber subextractor.Subscriber
	}
	tests := []struct {
		name         string
		args         args
		expected_str string
	}{
		{
			name:         "Get CC with Spanish MSISDN",
			args:         args{subscriber: subextractor.Subscriber{MSISDN: "+34600556800"}},
			expected_str: "ES",
		},
		{
			name:         "Get CC with UK MSISDN",
			args:         args{subscriber: subextractor.Subscriber{MSISDN: "+442012345678"}},
			expected_str: "GB",
		},
		{
			name:         "Get CC with French IMSI",
			args:         args{subscriber: subextractor.Subscriber{IMSI: "208260123456789"}},
			expected_str: "FR",
		},
		{
			name:         "Get CC with Italian IMSI",
			args:         args{subscriber: subextractor.Subscriber{IMSI: "222260123456789"}},
			expected_str: "IT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCountryCodeAsAlpha2(tt.args.subscriber)
			assert.Equal(t, tt.expected_str, result)
		})
	}
}
