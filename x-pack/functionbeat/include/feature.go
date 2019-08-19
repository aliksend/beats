// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package include

import (
	"github.com/aliksend/beats/libbeat/feature"
	"github.com/aliksend/beats/x-pack/functionbeat/manager/aws"
	"github.com/aliksend/beats/x-pack/functionbeat/provider/local/local"
)

// Bundle feature enabled.
var Bundle = feature.MustBundle(
	local.Bundle,
	aws.Bundle,
)

func init() {
	feature.MustRegisterBundle(Bundle)
}
