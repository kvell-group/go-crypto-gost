package signaturealgorithm

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kvell-group/go-crypto-gost/pkg/crypto-message/oids"
	publickeyalgorithm "github.com/kvell-group/go-crypto-gost/pkg/crypto-message/oids/algorithm/public-key-algorithm"
	"github.com/kvell-group/go-crypto-gost/pkg/crypto-message/oids/hash"
	"github.com/nobuenhombre/suikat/pkg/ge"
)

type getTest struct {
	in  oids.ID
	out *SignatureAlgorithm
	err error
}

func getTestsGet() []getTest {
	return []getTest{
		{
			in: oids.Tc26Gost34102012256,
			out: &SignatureAlgorithm{
				"GOST-3410_12_256",
				publickeyalgorithm.GostR34102012256,
				hash.UnknownHashFunction,
			},
			err: nil,
		},
		{
			in:  "GJFge7f3u3y6",
			out: nil,
			err: &ge.NotFoundError{
				Key: "GJFge7f3u3y6",
			},
		},
	}
}

func TestGet(t *testing.T) {
	getTests := getTestsGet()

	for i := 0; i < len(getTests); i++ {
		test := &getTests[i]
		out, err := Get(test.in)

		outEqual := reflect.DeepEqual(out, test.out)

		errEqual := err == nil
		if test.err != nil {
			errEqual = errors.Is(err, test.err)
		}

		if !(outEqual && errEqual) {
			t.Errorf(
				"[i=%v], Get(%v), Expected (%v, %v) Actual (%v, %v)\n",
				i, test.in, test.out, test.err, out, err,
			)
		}
	}
}
