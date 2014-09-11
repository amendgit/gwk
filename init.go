// Copyright 2014 By Jshi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gwk

import (
	"gwk/vango"
	"gwk/views"
)

func Init() {
	vango.InitVango()
	// The |InitViews| must after the |InitVango|. Cause the font should init
	// first.
	views.InitViews()
}
