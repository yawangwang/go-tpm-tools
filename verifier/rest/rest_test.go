package rest

import (
	"testing"

	confidentialcomputingpb "cloud.google.com/go/confidentialcomputing/apiv1/confidentialcomputingpb"
	"github.com/google/go-cmp/cmp"
	sabi "github.com/google/go-sev-guest/abi"
	spb "github.com/google/go-sev-guest/proto/sevsnp"
	tabi "github.com/google/go-tdx-guest/abi"
	tpb "github.com/google/go-tdx-guest/proto/tdx"
	tgtestdata "github.com/google/go-tdx-guest/testing/testdata"
	"github.com/google/go-tpm-tools/verifier"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/testing/protocmp"
)

// Make sure our conversion function can handle empty values.
func TestConvertEmpty(t *testing.T) {
	if _, err := convertChallengeFromREST(&confidentialcomputingpb.Challenge{}); err != nil {
		t.Errorf("Converting empty challenge: %v", err)
	}
	_ = convertRequestToREST(verifier.VerifyAttestationRequest{})
	if _, err := convertResponseFromREST(&confidentialcomputingpb.VerifyAttestationResponse{}); err != nil {
		t.Errorf("Converting empty challenge: %v", err)
	}
}

const (
	emptyReport = `
	version: 2
	policy: 0xa0000
	signature_algo: 1
	report_data: '\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01'
	family_id: '\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
  image_id: '\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
	measurement: '\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
  host_data: '\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
  id_key_digest: '\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
  author_key_digest: '\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
  report_id: '\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
  report_id_ma: '\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
  chip_id: '\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
 	signature: '\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00'
	`
	extraGUID = "00000000-0000-c0de-0000-000000000000"
)

func TestConvertSEVSNPProtoToREST(t *testing.T) {
	report := &spb.Report{}
	if err := prototext.Unmarshal([]byte(emptyReport), report); err != nil {
		t.Fatalf("Unable to unmarshal SEV-SNP report proto: %v", err)
	}

	rawCertTable := testRawCertTable(t)
	certTable := &sabi.CertTable{}
	if err := certTable.Unmarshal(rawCertTable.table); err != nil {
		t.Fatalf("Failed to unmarshal certTable bytes: %v", err)
	}
	sevsnp := &spb.Attestation{Report: report, CertificateChain: certTable.Proto()}

	got, err := convertSEVSNPProtoToREST(sevsnp)
	if err != nil {
		t.Errorf("failed to convert SEVSNP proto to API proto: %v", err)
	}

	wantReport, err := sabi.ReportToAbiBytes(report)
	if err != nil {
		t.Fatalf("Unable to convert SEV-SNP report proto to ABI bytes: %v", err)
	}

	want := &confidentialcomputingpb.VerifyAttestationRequest_SevSnpAttestation{
		SevSnpAttestation: &confidentialcomputingpb.SevSnpAttestation{
			AuxBlob: rawCertTable.table,
			Report:  wantReport,
		},
	}

	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("SEVSNP API proto mismatch: %s", diff)
	}
}

type testCertTable struct {
	table    []byte
	extraraw []byte
}

func testRawCertTable(t testing.TB) *testCertTable {
	t.Helper()
	headers := make([]sabi.CertTableHeaderEntry, 6) // ARK, ASK, VCEK, VLEK, extra, NULL
	arkraw := []byte("ark")
	askraw := []byte("ask")
	vcekraw := []byte("vcek")
	vlekraw := []byte("vlek")
	extraraw := []byte("extra")

	var err error
	headers[0].GUID, err = uuid.Parse(sabi.ArkGUID)
	if err != nil {
		t.Fatalf("cannot parse uuid: %v", err)
	}
	headers[0].Offset = uint32(len(headers) * sabi.CertTableEntrySize)
	headers[0].Length = uint32(len(arkraw))

	headers[1].GUID, err = uuid.Parse(sabi.AskGUID)

	if err != nil {
		t.Fatalf("cannot parse uuid: %v", err)
	}
	headers[1].Offset = headers[0].Offset + headers[0].Length
	headers[1].Length = uint32(len(askraw))

	headers[2].GUID, err = uuid.Parse(sabi.VcekGUID)
	if err != nil {
		t.Fatalf("cannot parse uuid: %v", err)
	}
	headers[2].Offset = headers[1].Offset + headers[1].Length
	headers[2].Length = uint32(len(vcekraw))

	headers[3].GUID, err = uuid.Parse(sabi.VlekGUID)
	if err != nil {
		t.Fatalf("cannot parse uuid: %v", err)
	}
	headers[3].Offset = headers[2].Offset + headers[2].Length
	headers[3].Length = uint32(len(vlekraw))

	headers[4].GUID, err = uuid.Parse(extraGUID)
	if err != nil {
		t.Fatalf("cannot parse uuid: %v", err)
	}
	headers[4].Offset = headers[3].Offset + headers[3].Length
	headers[4].Length = uint32(len(extraraw))

	result := &testCertTable{
		table:    make([]byte, headers[4].Offset+headers[4].Length),
		extraraw: extraraw,
	}
	for i, cert := range [][]byte{arkraw, askraw, vcekraw, vlekraw, extraraw} {
		if err := (&headers[i]).Write(result.table[i*sabi.CertTableEntrySize:]); err != nil {
			t.Fatalf("could not write header %d: %v", i, err)
		}
		copy(result.table[headers[i].Offset:], cert)
	}
	return result
}

func TestConvertTDXProtoToREST(t *testing.T) {
	testCases := []struct {
		name     string
		quote    func() *tpb.QuoteV4
		wantPass bool
	}{
		{
			name: "successful TD quote conversion",
			quote: func() *tpb.QuoteV4 {
				tdx, err := tabi.QuoteToProto(tgtestdata.RawQuote)
				if err != nil {
					t.Fatalf("Unable to convert Raw TD Quote to TDX V4 quote: %v", err)
				}

				quote, ok := tdx.(*tpb.QuoteV4)
				if !ok {
					t.Fatal("Quote format not supported, want QuoteV4 format")
				}
				return quote
			},
			wantPass: true,
		},
		{
			name:     "nil TD quote conversion",
			quote:    func() *tpb.QuoteV4 { return nil },
			wantPass: false,
		},
	}

	for _, tc := range testCases {
		got, err := convertTDXProtoToREST(tc.quote())
		if err != nil && tc.wantPass {
			t.Errorf("failed to convert TDX proto to API proto: %v", err)
		}

		if tc.wantPass {
			want := &confidentialcomputingpb.VerifyAttestationRequest_TdCcel{
				TdCcel: &confidentialcomputingpb.TdxCcelAttestation{
					TdQuote: tgtestdata.RawQuote,
				},
			}

			if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
				t.Errorf("TDX API proto mismatch: %s", diff)
			}
		}
	}
}
