package tunnel

import (
	"errors"
	"net/netip"
	"testing"

	C "github.com/metacubex/mihomo/constant"
)

func TestShouldFindProcessWithLocalIP(t *testing.T) {
	tests := []struct {
		name      string
		srcIP     netip.Addr
		local     bool
		localErr  error
		want      bool
		wantCheck bool
	}{
		{
			name:  "loopback",
			srcIP: netip.MustParseAddr("127.0.0.1"),
			want:  true,
		},
		{
			name:      "local interface",
			srcIP:     netip.MustParseAddr("192.0.2.1"),
			local:     true,
			want:      true,
			wantCheck: true,
		},
		{
			name:      "lan client",
			srcIP:     netip.MustParseAddr("10.0.0.9"),
			want:      false,
			wantCheck: true,
		},
		{
			name:      "local check error preserves old behavior",
			srcIP:     netip.MustParseAddr("10.0.0.9"),
			localErr:  errors.New("boom"),
			want:      true,
			wantCheck: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checked := false
			got := shouldFindProcessWithLocalIP(&C.Metadata{SrcIP: tt.srcIP}, func(addr netip.Addr) (bool, error) {
				checked = true
				if addr != tt.srcIP {
					t.Fatalf("expected addr %s, got %s", tt.srcIP, addr)
				}
				return tt.local, tt.localErr
			})
			if got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
			if checked != tt.wantCheck {
				t.Fatalf("expected local check %v, got %v", tt.wantCheck, checked)
			}
		})
	}
}
