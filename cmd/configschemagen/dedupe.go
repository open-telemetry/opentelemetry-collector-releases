// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"crypto/sha1" //nolint:gosec // used only to derive a short, stable $defs key, not for anything security-sensitive
	"encoding/hex"
	"encoding/json"
)

// minDedupSize is the minimum encoded size (bytes, compact JSON) a schema
// fragment must have before it's worth hoisting into a shared $defs entry;
// below this the "$ref" itself would cost more than inlining the fragment.
const minDedupSize = 60

// dedupeSchemas finds schema fragments that are byte-for-byte identical and
// occur more than once across the root skeleton and the $defs pool, and
// factors each one out into a single shared pool entry referenced by every
// occurrence instead. This is a pure size optimization with no loss of
// information: components commonly repeat the same sub-schema verbatim on
// many sibling fields (for example the same "metrics_exclude"/
// "metrics_include" filter description on every attribute toggle a
// receiver's metadata.yaml generates).
func dedupeSchemas(skeleton map[string]any, pool map[string]any) {
	counts := map[string]int{}
	countOccurrences(skeleton, counts)
	for _, def := range pool {
		countOccurrences(def, counts)
	}

	extract := map[string]bool{}
	for key, n := range counts {
		if n >= 2 && len(key) >= minDedupSize {
			extract[key] = true
		}
	}

	defined := map[string]string{} // canonical content -> pool key
	substitute(skeleton, pool, extract, defined)
	for _, def := range pool {
		substitute(def, pool, extract, defined)
	}
}

// countOccurrences walks every schema-valid position reachable from node and
// tallies how many times each exact fragment (by canonical JSON content)
// appears.
func countOccurrences(node any, counts map[string]int) {
	m, ok := node.(map[string]any)
	if !ok {
		return
	}
	counts[canonicalKey(m)]++
	for _, child := range schemaChildren(m) {
		countOccurrences(child, counts)
	}
}

// substitute walks every schema-valid position reachable from node in
// place, replacing any child fragment selected for extraction with a $ref
// into pool.
func substitute(node any, pool map[string]any, extract map[string]bool, defined map[string]string) {
	m, ok := node.(map[string]any)
	if !ok {
		return
	}
	if props, ok := m["properties"].(map[string]any); ok {
		for k, v := range props {
			props[k] = resolveChild(v, pool, extract, defined)
		}
	}
	if pp, ok := m["patternProperties"].(map[string]any); ok {
		for k, v := range pp {
			pp[k] = resolveChild(v, pool, extract, defined)
		}
	}
	switch items := m["items"].(type) {
	case map[string]any:
		m["items"] = resolveChild(items, pool, extract, defined)
	case []any:
		for i, v := range items {
			items[i] = resolveChild(v, pool, extract, defined)
		}
	}
	if ap, ok := m["additionalProperties"].(map[string]any); ok {
		m["additionalProperties"] = resolveChild(ap, pool, extract, defined)
	}
	for _, key := range [...]string{"anyOf", "oneOf", "allOf"} {
		if list, ok := m[key].([]any); ok {
			for i, v := range list {
				list[i] = resolveChild(v, pool, extract, defined)
			}
		}
	}
	if not, ok := m["not"].(map[string]any); ok {
		m["not"] = resolveChild(not, pool, extract, defined)
	}
}

// resolveChild recurses into child first (so nested duplicates are hoisted
// too), then, if child's canonical content was selected for extraction,
// hoists it into pool (exactly once per distinct content) and returns a
// $ref to it in place of the inline fragment.
func resolveChild(child any, pool map[string]any, extract map[string]bool, defined map[string]string) any {
	m, ok := child.(map[string]any)
	if !ok {
		return child
	}
	if _, isRef := m["$ref"]; isRef {
		return m
	}
	key := canonicalKey(m)
	if !extract[key] {
		substitute(m, pool, extract, defined)
		return m
	}
	if poolKey, ok := defined[key]; ok {
		return map[string]any{"$ref": "#/$defs/" + poolKey}
	}
	substitute(m, pool, extract, defined)
	poolKey := "_shared_" + shortHash(key)
	defined[key] = poolKey
	pool[poolKey] = m
	return map[string]any{"$ref": "#/$defs/" + poolKey}
}

// schemaChildren returns the direct schema-valued positions of a schema
// fragment: the keys where the JSON Schema spec itself requires the value
// to be another schema (or list of schemas), as opposed to plain data like
// "description", "default", "enum" or "required".
func schemaChildren(m map[string]any) []any {
	var out []any
	if props, ok := m["properties"].(map[string]any); ok {
		for _, v := range props {
			out = append(out, v)
		}
	}
	if pp, ok := m["patternProperties"].(map[string]any); ok {
		for _, v := range pp {
			out = append(out, v)
		}
	}
	switch items := m["items"].(type) {
	case map[string]any:
		out = append(out, items)
	case []any:
		out = append(out, items...)
	}
	if ap, ok := m["additionalProperties"].(map[string]any); ok {
		out = append(out, ap)
	}
	for _, key := range [...]string{"anyOf", "oneOf", "allOf"} {
		if list, ok := m[key].([]any); ok {
			out = append(out, list...)
		}
	}
	if not, ok := m["not"].(map[string]any); ok {
		out = append(out, not)
	}
	return out
}

// canonicalKey returns a deterministic, compact JSON encoding of m, used
// both to detect byte-identical fragments and as the dedup map key.
// encoding/json sorts map[string]any keys alphabetically, so this is stable
// regardless of the fragment's original key order.
func canonicalKey(m map[string]any) string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

func shortHash(s string) string {
	sum := sha1.Sum([]byte(s)) //nolint:gosec // short stable label only, not a security boundary
	return hex.EncodeToString(sum[:])[:12]
}
