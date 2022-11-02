package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/tools"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// IncludeActivations contains operations available on IncludeVersion resource
	IncludeActivations interface {
		// ActivateInclude creates a new include activation, which deactivates any current activation
		ActivateInclude(context.Context, ActivateIncludeRequest) (*ActivationIncludeResponse, error)

		// DeactivateInclude deactivates the include activation
		DeactivateInclude(context.Context, DeactivateIncludeRequest) (*DeactivationIncludeResponse, error)

		// GetIncludeActivation gets details about an activation
		GetIncludeActivation(context.Context, GetIncludeActivationRequest) (*IncludeActivationResponse, error)

		// ListIncludeActivations lists all activations for all versions of the include, on both production and staging networks
		ListIncludeActivations(context.Context, ListIncludeActivationsRequest) (*IncludeActivationsResponse, error)
	}

	// ActivateIncludeRequest contains parameters used to activate include
	ActivateIncludeRequest ActivateOrDeactivateIncludeRequest

	// DeactivateIncludeRequest contains parameters used to deactivate include
	DeactivateIncludeRequest ActivateOrDeactivateIncludeRequest

	// ActivateOrDeactivateIncludeRequest contains parameters used to activate or deactivate include
	ActivateOrDeactivateIncludeRequest struct {
		IncludeID              string            `json:"-"`
		Version                int               `json:"includeVersion"`
		Network                ActivationNetwork `json:"network"`
		Note                   string            `json:"note"`
		NotifyEmails           []string          `json:"notifyEmails"`
		AcknowledgeWarnings    []string          `json:"acknowledgeWarnings,omitempty"`
		AcknowledgeAllWarnings bool              `json:"acknowledgeAllWarnings"`
		IgnoreHTTPErrors       *bool             `json:"ignoreHttpErrors,omitempty"`
	}

	// ActivationIncludeResponse represents a response object returned by ActivateInclude operation
	ActivationIncludeResponse struct {
		ActivationID   string `json:"-"`
		ActivationLink string `json:"activationLink"`
	}

	// DeactivationIncludeResponse represents a response object returned by DeactivateInclude operation
	DeactivationIncludeResponse struct {
		ActivationID   string `json:"-"`
		ActivationLink string `json:"activationLink"`
	}

	// GetIncludeActivationRequest contains parameters used to get the include activation
	GetIncludeActivationRequest struct {
		IncludeID    string
		ActivationID string
	}

	// IncludeActivationResponse represents a response object returned by GetIncludeActivation
	IncludeActivationResponse struct {
		AccountID   string                `json:"accountId"`
		ContractID  string                `json:"contractId"`
		GroupID     string                `json:"groupId"`
		Activations IncludeActivationsRes `json:"activations"`
		Validations *Validations          `json:"validations,omitempty"`
	}

	// Validations represent include activation validation object
	Validations struct {
		ValidationSummary          ValidationSummary  `json:"validationSummary"`
		ValidationProgressItemList ValidationProgress `json:"validationProgressItemList"`
		Network                    ActivationNetwork  `json:"network"`
	}

	// ValidationSummary represent include activation validation summary object
	ValidationSummary struct {
		CompletePercent      float64 `json:"completePercent"`
		HasValidationError   bool    `json:"hasValidationError"`
		HasValidationWarning bool    `json:"hasValidationWarning"`
		HasSystemError       bool    `json:"hasSystemError"`
		HasClientError       bool    `json:"hasClientError"`
		MessageState         string  `json:"messageState"`
	}

	// ValidationProgress represents include activation validation progress object
	ValidationProgress struct {
		ErrorItems []ErrorItem `json:"errorItemsList"`
	}

	// ErrorItem represents validation progress error item object
	ErrorItem struct {
		VersionID             int    `json:"versionId"`
		PropertyName          string `json:"propertyName"`
		VersionNumber         int    `json:"versionNumber"`
		HasValidationError    bool   `json:"hasValidationError"`
		HasValidationWarning  bool   `json:"hasValidationWarning"`
		ValidationResultsLink string `json:"validationResultsLink"`
	}

	// ListIncludeActivationsRequest contains parameters used to list the include activations
	ListIncludeActivationsRequest struct {
		IncludeID  string
		ContractID string
		GroupID    string
	}

	// IncludeActivationsResponse represents a response object returned by ListIncludeActivations
	IncludeActivationsResponse struct {
		AccountID   string                `json:"accountId"`
		ContractID  string                `json:"contractId"`
		GroupID     string                `json:"groupId"`
		Activations IncludeActivationsRes `json:"activations"`
	}

	// IncludeActivationsRes represents Activations object
	IncludeActivationsRes struct {
		Activations []IncludeActivation `json:"items"`
	}

	// IncludeActivation represents an include activation object
	IncludeActivation struct {
		ActivationID        string                  `json:"activationId"`
		Network             ActivationNetwork       `json:"network"`
		ActivationType      ActivationType          `json:"activationType"`
		Status              ActivationStatus        `json:"status"`
		SubmitDate          string                  `json:"submitDate"`
		UpdateDate          string                  `json:"updateDate"`
		Note                string                  `json:"note"`
		NotifyEmails        []string                `json:"notifyEmails"`
		FMAActivationState  string                  `json:"fmaActivationState"`
		FallbackInfo        *ActivationFallbackInfo `json:"fallbackInfo"`
		IncludeID           string                  `json:"includeId"`
		IncludeName         string                  `json:"includeName"`
		IncludeType         IncludeType             `json:"includeType"`
		IncludeVersion      int                     `json:"includeVersion"`
		IncludeActivationID string                  `json:"includeActivationId"`
	}
)

// Validate validates ActivateIncludeRequest
func (i ActivateIncludeRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IncludeID":    validation.Validate(i.IncludeID, validation.Required),
		"Version":      validation.Validate(i.Version, validation.Required),
		"Network":      validation.Validate(i.Network, validation.Required),
		"NotifyEmails": validation.Validate(i.NotifyEmails, validation.Required),
	})
}

// Validate validates DeactivateIncludeRequest
func (i DeactivateIncludeRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IncludeID":    validation.Validate(i.IncludeID, validation.Required),
		"Version":      validation.Validate(i.Version, validation.Required),
		"Network":      validation.Validate(i.Network, validation.Required),
		"NotifyEmails": validation.Validate(i.NotifyEmails, validation.Required),
	})
}

// Validate validates GetIncludeActivationRequest
func (i GetIncludeActivationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IncludeID":    validation.Validate(i.IncludeID, validation.Required),
		"ActivationID": validation.Validate(i.ActivationID, validation.Required),
	})
}

// Validate validates ListIncludeActivationsRequest
func (i ListIncludeActivationsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"IncludeID":  validation.Validate(i.IncludeID, validation.Required),
		"ContractID": validation.Validate(i.ContractID, validation.Required),
		"GroupID":    validation.Validate(i.GroupID, validation.Required),
	})
}

var (
	// ErrActivateInclude is returned in case an error occurs on ActivateInclude operation
	ErrActivateInclude = errors.New("activate include")
	// ErrDeactivateInclude is returned in case an error occurs on DeactivateInclude operation
	ErrDeactivateInclude = errors.New("deactivate include")
	// ErrGetIncludeActivation is returned in case an error occurs on GetIncludeActivation operation
	ErrGetIncludeActivation = errors.New("get include activation")
	// ErrListIncludeActivations is returned in case an error occurs on ListIncludeActivations operation
	ErrListIncludeActivations = errors.New("list include activations")
)

func (p *papi) ActivateInclude(ctx context.Context, params ActivateIncludeRequest) (*ActivationIncludeResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ActivateInclude")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrActivateInclude, ErrStructValidation, err)
	}

	if params.IgnoreHTTPErrors == nil {
		params.IgnoreHTTPErrors = tools.BoolPtr(true)
	}

	requestBody := struct {
		ActivateIncludeRequest
		ActivationType ActivationType `json:"activationType"`
	}{
		params,
		ActivationTypeActivate,
	}

	uri := fmt.Sprintf("/papi/v1/includes/%s/activations", params.IncludeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrActivateInclude, err)
	}

	var result ActivationIncludeResponse
	resp, err := p.Exec(req, &result, requestBody)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrActivateInclude, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrActivateInclude, p.Error(resp))
	}

	id, err := ResponseLinkParse(result.ActivationLink)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrActivateInclude, ErrInvalidResponseLink, err)
	}
	result.ActivationID = id

	return &result, nil
}

func (p *papi) DeactivateInclude(ctx context.Context, params DeactivateIncludeRequest) (*DeactivationIncludeResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("DeactivateInclude")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeactivateInclude, ErrStructValidation, err)
	}

	if params.IgnoreHTTPErrors == nil {
		params.IgnoreHTTPErrors = tools.BoolPtr(true)
	}

	requestBody := struct {
		DeactivateIncludeRequest
		ActivationType ActivationType `json:"activationType"`
	}{
		params,
		ActivationTypeDeactivate,
	}

	uri := fmt.Sprintf("/papi/v1/includes/%s/activations", params.IncludeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeactivateInclude, err)
	}

	var result DeactivationIncludeResponse
	resp, err := p.Exec(req, &result, requestBody)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeactivateInclude, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrDeactivateInclude, p.Error(resp))
	}

	id, err := ResponseLinkParse(result.ActivationLink)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeactivateInclude, ErrInvalidResponseLink, err)
	}
	result.ActivationID = id

	return &result, nil
}

func (p *papi) GetIncludeActivation(ctx context.Context, params GetIncludeActivationRequest) (*IncludeActivationResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetIncludeActivation")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetIncludeActivation, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/papi/v1/includes/%s/activations/%s", params.IncludeID, params.ActivationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetIncludeActivation, err)
	}

	var result IncludeActivationResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetIncludeActivation, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetIncludeActivation, p.Error(resp))
	}

	return &result, nil
}

func (p *papi) ListIncludeActivations(ctx context.Context, params ListIncludeActivationsRequest) (*IncludeActivationsResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("ListIncludeActivations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListIncludeActivations, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/papi/v1/includes/%s/activations", params.IncludeID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListIncludeActivations, err)
	}

	q := uri.Query()
	q.Add("contractId", params.ContractID)
	q.Add("groupId", params.GroupID)
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListIncludeActivations, err)
	}

	var result IncludeActivationsResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListIncludeActivations, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListIncludeActivations, p.Error(resp))
	}

	return &result, nil
}
