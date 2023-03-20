package upapi

import (
	"context"

	"github.com/hashicorp/go-cleanhttp"
)

const (
	defaultBaseURL   = "https://uptime.com/api/v1/"
	defaultUserAgent = "uptime-client-go"
)

// API manages communication with the Uptime.com API.
type API interface {
	Checks() ChecksEndpoint
	Contacts() ContactsEndpoint
	Integrations() IntegrationsEndpoint
	Tags() TagsEndpoint
	Outages() OutagesEndpoint
	ProbeServers() ProbeServersEndpoint
}

// New returns a new API client instance.
func New(opts ...Option) (api API, err error) {
	var cbd CBD = &struct {
		Doer
		RequestBuilder
		ResponseDecoder
	}{
		cleanhttp.DefaultClient(),
		&requestBuilderImpl{},
		&responseDecoderImpl{},
	}
	cbd, err = applyOptions(cbd, opts...)
	if err != nil {
		return nil, err
	}

	var defs []Option
	if cbd.(Doer) == nil {
		panic("no Doer")
	}

	rq, err := cbd.BuildRequest(context.Background(), "GET", "/", nil, nil)
	if rq.URL.Host == "" {
		defs = append(defs, WithBaseURL(defaultBaseURL))
	}
	if rq.Header.Get("User-Agent") == "" {
		defs = append(defs, WithUserAgent(defaultUserAgent))
	}
	cbd, err = applyOptions(cbd, defs...)
	if err != nil {
		return nil, err
	}

	api = &apiImpl{
		CBD:          cbd,
		checks:       NewChecksEndpoint(cbd),
		contacts:     NewContactsEndpoint(cbd),
		integrations: NewIntegrationsEndpoint(cbd),
		tags:         NewTagsEndpoint(cbd),
		outages:      NewOutagesEndpoint(cbd),
		probeServers: NewProbeServersEndpoint(cbd),
	}
	return api, nil
}

type apiImpl struct {
	CBD
	checks       ChecksEndpoint
	contacts     ContactsEndpoint
	integrations IntegrationsEndpoint
	tags         TagsEndpoint
	outages      OutagesEndpoint
	probeServers ProbeServersEndpoint
}

func (api *apiImpl) Checks() ChecksEndpoint {
	return api.checks
}

func (api *apiImpl) Contacts() ContactsEndpoint {
	return api.contacts
}

func (api *apiImpl) Integrations() IntegrationsEndpoint {
	return api.integrations
}

func (api *apiImpl) Tags() TagsEndpoint {
	return api.tags
}

func (api *apiImpl) Outages() OutagesEndpoint {
	return api.outages
}

func (api *apiImpl) ProbeServers() ProbeServersEndpoint {
	return api.probeServers
}