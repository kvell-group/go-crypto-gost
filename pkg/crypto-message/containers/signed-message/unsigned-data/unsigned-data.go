package unsigneddata

import (
	"encoding/asn1"

	"github.com/kvell-group/go-crypto-gost/pkg/crypto-message/containers"

	"github.com/nobuenhombre/suikat/pkg/ge"
)

type Container []byte

func DecodeDER(data containers.DER) (*Container, error) {
	var (
		compound asn1.RawValue
		content  Container
	)

	if len(data) > 0 {
		_, err := asn1.Unmarshal(data, &compound)
		if err != nil {
			return nil, ge.Pin(err)
		}
	}

	// Compound octet string
	if compound.IsCompound && compound.Tag == 4 {
		_, err := asn1.Unmarshal(compound.Bytes, &content)
		if err != nil {
			return nil, ge.Pin(err)
		}
	} else {
		// assuming this is tag 04
		content = compound.Bytes
	}

	return &content, nil
}
