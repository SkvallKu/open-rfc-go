package rfc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/oisee/open-rfc-go/internal/client"
)

func TestLogonErrorRejectsOnlyARefusal(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		want error
	}{
		{"refused", fmt.Errorf("%w: SAP rejected logon with status 1", client.ErrLogonRefused), ErrLogonRejected},
		{"malformed", fmt.Errorf("%w: CPIC logon response is malformed", client.ErrSession), ErrProtocol},
		{"no answer", errors.New("read tcp: i/o timeout"), ErrTransport},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := logonError(tc.err)
			for _, e := range []error{ErrLogonRejected, ErrProtocol, ErrTransport} {
				if errors.Is(got, e) != (e == tc.want) {
					t.Errorf("logonError(%v) = %v; errors.Is(%v) = %v", tc.err, got, e, errors.Is(got, e))
				}
			}
		})
	}
}
