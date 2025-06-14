// GoGOST -- Pure Go GOST cryptographic functions library
// Copyright (C) 2015-2022 Sergey Matveev <stargrave@stargrave.org>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 3 of the License.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package gost3410

import (
	"errors"
	"math/big"

	"github.com/kvell-group/go-crypto-gost/pkg/gost28147"
	"github.com/kvell-group/go-crypto-gost/pkg/gost341194"
)

// RFC 4357 VKO GOST R 34.10-2001 key agreement function.
// UKM is user keying material, also called VKO-factor.
func (prv *PrivateKey) KEK2001(pub *PublicKey, ukm *big.Int) ([]byte, error) {
	if prv.C.PointSize() != 32 {
		return nil, errors.New("gogost/gost3410: KEK2001 is only for 256-bit curves")
	}
	key, err := prv.KEK(pub, ukm)
	if err != nil {
		return nil, err
	}
	h := gost341194.New(&gost28147.SboxIdGostR341194CryptoProParamSet)
	if _, err = h.Write(key); err != nil {
		return nil, err
	}
	return h.Sum(key[:0]), nil
}
