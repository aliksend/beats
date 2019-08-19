// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aliksend/beats/metricbeat/scripts/msetlists"

	"github.com/aliksend/beats/libbeat/paths"
	"github.com/aliksend/beats/metricbeat/mb"
	_ "github.com/aliksend/beats/x-pack/metricbeat/include"
	xpackmb "github.com/aliksend/beats/x-pack/metricbeat/mb"
)

func main() {
	path := paths.Resolve(paths.Home, "../x-pack/metricbeat/module")
	lm := xpackmb.NewLightModulesSource(path)
	mb.Registry.SetSecondarySource(lm)

	msList := msetlists.DefaultMetricsets()

	raw, err := json.MarshalIndent(msList, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Marshalling json: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", string(raw))
}
