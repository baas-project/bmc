package ipmi

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/gopacket"
)

func TestAuthenticationPayload(t *testing.T) {
	table := []struct {
		wire      []byte
		payload   *AuthenticationPayload
		remaining []byte
	}{
		{
			// too short
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			nil,
			nil,
		},
		{
			// not authentication payload
			[]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			nil,
			nil,
		},
		{
			// simultaneously wildcard and not-None
			[]byte{0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00},
			nil,
			nil,
		},
		{
			// wildcard
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1},
			&AuthenticationPayload{
				Wildcard: true,
			},
			[]byte{0x1},
		},
		{
			// RAKP-HMAC-SHA256
			[]byte{0x00, 0x00, 0x00, 0x08, 0x03, 0x00, 0x00, 0x00},
			&AuthenticationPayload{
				Algorithm: AuthenticationAlgorithmHMACSHA256,
			},
			[]byte{},
		},
	}
	payload := &AuthenticationPayload{}
	for _, test := range table {
		if test.payload != nil {
			sb := gopacket.NewSerializeBuffer()
			serialiseErr := test.payload.Serialise(sb)
			got := sb.Bytes()

			switch {
			case serialiseErr != nil:
				t.Errorf("serialise %v failed with %v, wanted %v", test.payload,
					serialiseErr, test.wire)
			case !bytes.Equal(got, test.wire[:8]):
				t.Errorf("serialise %v = %v, want %v", test.payload, got, test.wire)
			}
		}

		remaining, deserialiseErr := payload.Deserialise(test.wire, gopacket.NilDecodeFeedback)
		switch {
		case deserialiseErr == nil && test.payload == nil:
			t.Errorf("deserialise %v succeeded with %v, wanted error", test.wire,
				payload)
		case deserialiseErr != nil && test.payload != nil:
			t.Errorf("deserialise %v failed with %v, wanted %v", test.wire, deserialiseErr,
				test.payload)
		case deserialiseErr == nil && test.payload != nil:
			if !bytes.Equal(remaining, test.remaining) {
				t.Errorf("remaining from deserialising %v = %v, want %v", test.wire,
					remaining, test.remaining)
			}
			if diff := cmp.Diff(test.payload, payload); diff != "" {
				t.Errorf("deserialise %v = %v, want %v: %v", test.wire, payload,
					test.payload, diff)
			}
		}
	}
}
