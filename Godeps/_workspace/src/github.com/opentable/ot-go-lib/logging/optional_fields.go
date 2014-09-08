package logging

import (
	"net/http"
)

// OptionalFields probably need a better name.
// They represent the set of mainly OpenTable-specific HTTP headers that ought to be set by the Front Door Service.
type OptionalFields struct {
	UserAgent          string `json:"user-agent,omitempty",`
	AcceptLanguage     string `json:"accept-language,omitempty"`
	OTRequestId        string `json:"ot-requestid,omitempty"`
	OTUserId           string `json:"ot-userid,omitempty"`
	OTSessionId        string `json:"ot-sessionid,omitempty"`
	OTReferringHost    string `json:"ot-referringhost,omitempty"`
	OTReferringService string `json:"ot-referringservice,omitempty"`
	OTDomain           string `json:"ot-domain,omitempty"` // One of: com, couk, jp, de, commx
}

func newOptionalFields(r *http.Request) *OptionalFields {
	// Note: Case sensitive, canonicalised header keys.
	return &OptionalFields{
		UserAgent:          r.Header.Get("User-Agent"),
		AcceptLanguage:     r.Header.Get("Accept-Language"),
		OTRequestId:        r.Header.Get("Ot-Requestid"),
		OTUserId:           r.Header.Get("Ot-Userid"),
		OTSessionId:        r.Header.Get("Ot-Sessionid"),
		OTReferringHost:    r.Header.Get("Ot-Referringhost"),
		OTReferringService: r.Header.Get("Ot-Referringservice"),
		OTDomain:           r.Header.Get("Ot-Domain"),
	}
}
