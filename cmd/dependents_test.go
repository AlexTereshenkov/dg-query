/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseAdjacencyLists(t *testing.T) {
	dg := AdjacencyList{"foo": {"bar", "baz"}, "spam": {"eggs", "bar"}}
	rdg := reverseAdjacencyLists(dg)
	expected := AdjacencyList{"bar": {"foo", "spam"}, "baz": {"foo"}, "eggs": {"spam"}}
	assert.True(t, reflect.DeepEqual(rdg, expected))
}
