package logging

// StandardFields are expected to always have meaningful values.
// They are generated internally according to particular rules, so it's unlikely you'll need to create them yourself.
type StandardFields struct {
	// Type is calculated as "ServiceType-LogName-FormatVersion"
	Type           string      `json:"type"`
	ServiceType    string      `json:"servicetype"`
	LogName        string      `json:"logname"`
	FormatVersion  string      `json:"formatversion"`
	Host           string      `json:"host"`
	Severity       string      `json:"severity"`
	SequenceNumber int64       `json:"sequencenumber"`
	Timestamp      string      `json:"@timestamp"`
	LogMessage     interface{} `json:"logmessage"`
}

func newStandardFields(l *logBase, s Severity, m interface{}) *StandardFields {
	return &StandardFields{
		Type:           l.Type,
		ServiceType:    l.Config.ServiceType,
		LogName:        l.LogName,
		FormatVersion:  l.formatVersionString,
		Host:           l.Config.hostString,
		Severity:       s.String(),
		SequenceNumber: seq(),
		Timestamp:      timestamp(),
		LogMessage:     m,
	}
}
