package telemetry

var gutVersion = "dev"

func GetBuildInfo() string {
	return gutVersion
}

func GetConsentState() bool {
	return false
}

func IsConsentStateKnown() bool {
	return true
}
