package ccextractor

import (
	"github.com/ttacon/libphonenumber"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/subextractor"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/structlog"
)

// ISO Alpha-2 Country code mappings
var mccToCountryCode = map[string]string{
	"202": "GR", // Greece
	"204": "NL", // Netherlands
	"206": "BE", // Belgium
	"208": "FR", // France
	"212": "MC", // Monaco
	"213": "AD", // Andorra
	"214": "ES", // Spain
	"216": "HU", // Hungary
	"218": "BA", // Bosnia and Herzegovina
	"219": "HR", // Croatia
	"220": "RS", // Serbia
	"222": "IT", // Italy
	"226": "RO", // Romania
	"228": "CH", // Switzerland
	"230": "CZ", // Czech Republic
	"231": "SK", // Slovakia
	"232": "AT", // Austria
	"234": "GB", // United Kingdom
	"235": "GB", // United Kingdom (special use)
	"238": "DK", // Denmark
	"240": "SE", // Sweden
	"242": "NO", // Norway
	"244": "FI", // Finland
	"246": "LT", // Lithuania
	"247": "LV", // Latvia
	"248": "EE", // Estonia
	"250": "RU", // Russia
	"255": "UA", // Ukraine
	"257": "BY", // Belarus
	"259": "MD", // Moldova
	"260": "PL", // Poland
	"262": "DE", // Germany
	"266": "GI", // Gibraltar
	"268": "PT", // Portugal
	"270": "LU", // Luxembourg
	"272": "IE", // Ireland
	"274": "IS", // Iceland
	"276": "AL", // Albania
	"278": "MT", // Malta
	"280": "CY", // Cyprus
	"282": "GE", // Georgia
	"283": "AM", // Armenia
	"284": "BG", // Bulgaria
	"286": "TR", // Turkey
	"288": "FO", // Faroe Islands
	"289": "GE", // Abkhazia (disputed region, using Georgia's code as a placeholder.)
	"290": "GL", // Greenland
	"292": "SM", // San Marino
	"293": "SI", // Slovenia
	"294": "MK", // North Macedonia
	"295": "LI", // Liechtenstein
	"297": "ME", // Montenegro
}

// Returns the Country Code in ISO Alpha format (CC/MCC) for a given MSISDN/IMSI
func GetCountryCodeAsAlpha2(subscriber subextractor.Subscriber) string {
	countryCode := ""
	if subscriber.MSISDN != "" {
		countryCode = GetCountryCodeFromMsisdnAsAlpha2(subscriber.MSISDN)
	} else if subscriber.IMSI != "" {
		countryCode = GetCountryCodeFromImsiAsAlpha2(subscriber.IMSI)
	}
	return countryCode
}

func GetCountryCodeFromMsisdnAsAlpha2(msisdn string) string {
	logger := structlog.GetLogger("nebula")
	num, err := libphonenumber.Parse(msisdn, "")
	if err != nil {
		logger.Errorv(err.Error())
		return ""
	}
	return libphonenumber.GetRegionCodeForNumber(num)
}

func GetCountryCodeFromImsiAsAlpha2(imsi string) string {
	mcc_code := imsi[:3] // Country code always first three digits of an IMSI
	return mccToCountryCode[mcc_code]
}
