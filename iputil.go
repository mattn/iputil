package iputil

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/domainr/whois"
)

var (
	toUint32 = binary.BigEndian.Uint32
)

type Range struct {
	Begin net.IP
	End   net.IP
}

func (r *Range) CIDR() string {
	dist := toUint32(r.Begin) ^ toUint32(r.End)

	bits, mask := 32, 0
	for {
		if dist == 0 {
			break
		}
		dist >>= 1
		bits -= 1
		mask = (mask << 1) | 1
	}
	return fmt.Sprintf("%s/%d", r.Begin, bits)
}

func Ranges(query string) ([]Range, error) {
	req := &whois.Request{
		Query: query,
		Host:  "whois.apnic.net",
	}
	if err := req.Prepare(); err != nil {
		return nil, err
	}

	resp, err := whois.DefaultClient.Fetch(req)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(strings.NewReader(resp.String()))

	var rngs []Range

	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "inetnum:") {
			continue
		}
		token := strings.SplitN(scanner.Text(), ":", 2)
		if len(token) != 2 {
			continue
		}
		token = strings.SplitN(token[1], "-", 2)
		if len(token) != 2 {
			continue
		}
		rngs = append(rngs, Range{
			Begin: net.ParseIP(strings.TrimSpace(token[0])).To4(),
			End:   net.ParseIP(strings.TrimSpace(token[1])).To4(),
		})
	}

	// sort by begin address
	sort.Slice(rngs, func(i, j int) bool {
		return toUint32(rngs[i].Begin) < toUint32(rngs[j].Begin)
	})

	// join ranges
	i := 0
	for i < len(rngs)-1 {
		if toUint32(rngs[i].End)+1 == toUint32(rngs[i+1].Begin) {
			rngs = append(rngs[:i], rngs[i+2:]...)
		} else {
			i++
		}
	}

	return rngs, nil
}
