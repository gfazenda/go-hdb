// SPDX-FileCopyrightText: 2014-2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Package unicode implements UTF-8 to CESU-8 and vice versa transformations.
package unicode

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/SAP/go-hdb/driver/internal/unicode/cesu8"
	"golang.org/x/text/transform"
)

var (
	// Utf8ToCesu8Transformer implements the golang.org/x/text/transform/Transformer interface for UTF-8 to CESU-8 transformation.
	Utf8ToCesu8Transformer = new(utf8ToCesu8Transformer)
	// Cesu8ToUtf8Transformer implements the golang.org/x/text/transform/Transformer interface for CESU-8 to UTF-8 transformation.
	Cesu8ToUtf8Transformer = new(cesu8ToUtf8Transformer)
	// ErrInvalidUtf8 means that a transformer detected invalid UTF-8 data.
	ErrInvalidUtf8 = errors.New("invalid UTF-8")
)

// InvalidCesu8Error is the error raised when a transformer detects invalid CESU-8 data.
type InvalidCesu8Error struct {
	p int    // position of error in []byte
	v []byte // value where error occurs
}

func (e *InvalidCesu8Error) Error() string {
	return fmt.Sprintf("invalid CESU-8: %x at pos: %d", e.v, e.p)
}

func newInvalidCesu8Error(p int, v []byte) *InvalidCesu8Error {
	// copy value
	cv := make([]byte, len(v))
	copy(cv, v)
	return &InvalidCesu8Error{p: p, v: cv}
}

type utf8ToCesu8Transformer struct{ transform.NopResetter }

func (t *utf8ToCesu8Transformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	i, j := 0, 0
	for i < len(src) {
		if src[i] < utf8.RuneSelf {
			if j < len(dst) {
				dst[j] = src[i]
				i++
				j++
			} else {
				return j, i, transform.ErrShortDst
			}
		} else {
			if !utf8.FullRune(src[i:]) {
				return j, i, transform.ErrShortSrc
			}
			r, n := utf8.DecodeRune(src[i:])
			if r == utf8.RuneError {
				return j, i, ErrInvalidUtf8
			}
			m := cesu8.RuneLen(r)
			if m == -1 {
				panic("internal UTF-8 to CESU-8 transformation error")
			}
			if j+m <= len(dst) {
				cesu8.EncodeRune(dst[j:], r)
				i += n
				j += m
			} else {
				return j, i, transform.ErrShortDst
			}
		}
	}
	return j, i, nil
}

type cesu8ToUtf8Transformer struct{ transform.NopResetter }

func (t *cesu8ToUtf8Transformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	i, j := 0, 0
	for i < len(src) {
		if src[i] < utf8.RuneSelf {
			if j < len(dst) {
				dst[j] = src[i]
				i++
				j++
			} else {
				return j, i, transform.ErrShortDst
			}
		} else {
			if !cesu8.FullRune(src[i:]) {
				return j, i, transform.ErrShortSrc
			}
			r, n := cesu8.DecodeRune(src[i:])
			if r == utf8.RuneError {
				return j, i, newInvalidCesu8Error(i, src[i:i+n])
			}
			m := utf8.RuneLen(r)
			if m == -1 {
				panic("internal CESU-8 to UTF-8 transformation error")
			}
			if j+m <= len(dst) {
				utf8.EncodeRune(dst[j:], r)
				i += n
				j += m
			} else {
				return j, i, transform.ErrShortDst
			}
		}
	}
	return j, i, nil
}
