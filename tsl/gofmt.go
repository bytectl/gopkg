// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO(gri): This file and the file src/go/format/internal.go are
// the same (but for this comment and the package name). Do not modify
// one without the other. Determine if we can factor out functionality
// in a public API. See also #11844 for context.

package tsl

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
)

func gofmt(src string) string {

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "codec", src, parser.ParseComments)
	if err != nil {
		fmt.Printf("ParseExpr error: %v\n", err)
		return src
	}
	var buf bytes.Buffer
	err = format.Node(&buf, fset, node)
	if err != nil {
		fmt.Printf("format.Node error: %v\n", err)
		return src
	}

	return buf.String()
}
