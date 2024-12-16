// Package verifier contains clients for various attestation verifiers.
// It is meant for launcher use and testing; the API is not stable.
package verifier

import (
	"context"

	attestpb "github.com/google/go-tpm-tools/proto/attest"
	"github.com/google/go-tpm-tools/verifier/oci"
	"google.golang.org/genproto/googleapis/rpc/status"
)

// Client is a common interface to various attestation verifiers.
type Client interface {
	CreateChallenge(ctx context.Context) (*Challenge, error)
	VerifyAttestation(ctx context.Context, request VerifyAttestationRequest) (*VerifyAttestationResponse, error)
}

// Challenge is the response for CreateChallenge. It is used in the
// get challenge part of a remote attestation protocol. The challenge
// will be verified as part of VerifyAttestation.
type Challenge struct {
	Name   string
	Nonce  []byte
	ConnID string
}

// TokenOptions contains fields that will be passed to the Attestation Service TokenOptions field.
// These fields are used to customize several claims in the token from the Attestation service.
type TokenOptions struct {
	CustomAudience string
	CustomNonce    []string
	TokenType      string
}

// VerifyAttestationRequest is passed in on VerifyAttestation. It contains the
// Challenge from CreateChallenge, optional GcpCredentials linked to the
// attestation, the Attestation generated from the TPM, and optional container image signatures associated with the workload.
type VerifyAttestationRequest struct {
	Challenge      *Challenge
	GcpCredentials [][]byte
	// Attestation is for TPM attestation
	Attestation              *attestpb.Attestation
	ContainerImageSignatures []oci.Signature
	TokenOptions             TokenOptions
	// TDCCELAttestation is for TDX CCEL RTMR attestation
	TDCCELAttestation *TDCCELAttestation
}

type TDCCELAttestation struct {
	CcelAcpiTable     []byte
	CcelData          []byte
	CanonicalEventLog []byte
	TdQuote           []byte
	// still needs following two for GCE info
	AkCert            []byte
	IntermediateCerts [][]byte
}

// VerifyAttestationResponse is the response from a successful
// VerifyAttestation call.
type VerifyAttestationResponse struct {
	ClaimsToken []byte
	PartialErrs []*status.Status
}
