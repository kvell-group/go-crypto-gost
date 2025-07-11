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

package mgm

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"testing"
	"testing/quick"

	"github.com/kvell-group/go-crypto-gost/pkg/gost3412128"
	"github.com/kvell-group/go-crypto-gost/pkg/gost341264"
)

func TestVector(t *testing.T) {
	key := []byte{
		0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
		0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10,
		0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF,
	}
	additionalData := []byte{
		0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04,
		0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03,
		0xEA, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05,
		0x05,
	}
	plaintext := []byte{
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x00,
		0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA, 0x99, 0x88,
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
		0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xEE, 0xFF, 0x0A,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0x99, 0xAA, 0xBB, 0xCC, 0xEE, 0xFF, 0x0A, 0x00,
		0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99,
		0xAA, 0xBB, 0xCC, 0xEE, 0xFF, 0x0A, 0x00, 0x11,
		0xAA, 0xBB, 0xCC,
	}
	c := gost3412128.NewCipher(key)
	nonce := plaintext[:16]
	aead, _ := NewMGM(c, 16)
	sealed := aead.Seal(nil, nonce, plaintext, additionalData)
	if bytes.Compare(sealed[:len(plaintext)], []byte{
		0xA9, 0x75, 0x7B, 0x81, 0x47, 0x95, 0x6E, 0x90,
		0x55, 0xB8, 0xA3, 0x3D, 0xE8, 0x9F, 0x42, 0xFC,
		0x80, 0x75, 0xD2, 0x21, 0x2B, 0xF9, 0xFD, 0x5B,
		0xD3, 0xF7, 0x06, 0x9A, 0xAD, 0xC1, 0x6B, 0x39,
		0x49, 0x7A, 0xB1, 0x59, 0x15, 0xA6, 0xBA, 0x85,
		0x93, 0x6B, 0x5D, 0x0E, 0xA9, 0xF6, 0x85, 0x1C,
		0xC6, 0x0C, 0x14, 0xD4, 0xD3, 0xF8, 0x83, 0xD0,
		0xAB, 0x94, 0x42, 0x06, 0x95, 0xC7, 0x6D, 0xEB,
		0x2C, 0x75, 0x52,
	}) != 0 {
		t.FailNow()
	}
	if bytes.Compare(sealed[len(plaintext):], []byte{
		0xCF, 0x5D, 0x65, 0x6F, 0x40, 0xC3, 0x4F, 0x5C,
		0x46, 0xE8, 0xBB, 0x0E, 0x29, 0xFC, 0xDB, 0x4C,
	}) != 0 {
		t.FailNow()
	}
	_, err := aead.Open(sealed[:0], nonce, sealed, additionalData)
	if err != nil {
		t.FailNow()
	}
	if bytes.Compare(sealed[:len(plaintext)], plaintext) != 0 {
		t.FailNow()
	}
}

func TestSymmetric(t *testing.T) {
	sym := func(keySize, blockSize int, c cipher.Block, nonce []byte) {
		f := func(
			plaintext, additionalData []byte,
			initials [][]byte,
			tagSize uint8,
		) bool {
			if len(plaintext) == 0 && len(additionalData) == 0 {
				return true
			}
			tagSize = 4 + tagSize%uint8(blockSize-4)
			aead, err := NewMGM(c, int(tagSize))
			if err != nil {
				return false
			}
			for _, initial := range initials {
				sealed := aead.Seal(initial, nonce, plaintext, additionalData)
				if bytes.Compare(sealed[:len(initial)], initial) != 0 {
					return false
				}
				pt, err := aead.Open(
					sealed[:0],
					nonce,
					sealed[len(initial):],
					additionalData,
				)
				if err != nil || bytes.Compare(pt, plaintext) != 0 {
					return false
				}
			}
			return true
		}
		if err := quick.Check(f, nil); err != nil {
			t.Error(err)
		}
	}

	key128 := new([gost3412128.KeySize]byte)
	if _, err := rand.Read(key128[:]); err != nil {
		panic(err)
	}
	nonce := make([]byte, gost3412128.BlockSize)
	if _, err := rand.Read(key128[1:]); err != nil {
		panic(err)
	}
	sym(
		gost3412128.KeySize,
		gost3412128.BlockSize,
		gost3412128.NewCipher(key128[:]),
		nonce[:gost3412128.BlockSize],
	)

	key64 := new([gost341264.KeySize]byte)
	copy(key64[:], key128[:])
	sym(
		gost341264.KeySize,
		gost341264.BlockSize,
		gost341264.NewCipher(key64[:]),
		nonce[:gost341264.BlockSize],
	)
}

func BenchmarkMGM64(b *testing.B) {
	key := make([]byte, gost341264.KeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}
	nonce := make([]byte, gost341264.BlockSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}
	nonce[0] &= 0x7F
	pt := make([]byte, 1280+3)
	if _, err := io.ReadFull(rand.Reader, pt); err != nil {
		panic(err)
	}
	c := gost341264.NewCipher(key)
	aead, err := NewMGM(c, gost341264.BlockSize)
	if err != nil {
		panic(err)
	}
	ct := make([]byte, len(pt)+aead.Overhead())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aead.Seal(ct[:0], nonce, pt, nil)
	}
}

func BenchmarkMGM128(b *testing.B) {
	key := make([]byte, gost3412128.KeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}
	nonce := make([]byte, gost3412128.BlockSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}
	nonce[0] &= 0x7F
	pt := make([]byte, 1280+3)
	if _, err := io.ReadFull(rand.Reader, pt); err != nil {
		panic(err)
	}
	c := gost3412128.NewCipher(key)
	aead, err := NewMGM(c, gost3412128.BlockSize)
	if err != nil {
		panic(err)
	}
	ct := make([]byte, len(pt)+aead.Overhead())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aead.Seal(ct[:0], nonce, pt, nil)
	}
}
