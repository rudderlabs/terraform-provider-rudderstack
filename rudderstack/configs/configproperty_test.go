package configs_test

import (
	"testing"

	"github.com/rudderlabs/terraform-provider-rudderstack/rudderstack/configs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleConfigProperty(t *testing.T) {
	p := configs.Simple("a.b", "t.s")

	a, err := p.FromStateFunc(`{ "p": true }`, `{ "t": { "s": "123" } }`)
	require.NoError(t, err)
	assert.JSONEq(t, `{ "p": true, "a": { "b": "123" } }`, a)

	s, err := p.ToStateFunc(`{ "p": true }`, `{ "a": { "b": "123" } }`)
	require.NoError(t, err)
	assert.JSONEq(t, `{ "p": true, "t": { "s": "123" } }`, s)
}

func TestConditionalTrue(t *testing.T) {
	p := configs.Conditional("a.b", "t.s", func(state string) bool {
		return true
	})

	// FromStateFunc writes to api regardless of condition result
	a, err := p.FromStateFunc(`{ "p": true }`, `{ "t": { "s": "123" } }`)
	require.NoError(t, err)
	assert.JSONEq(t, `{ "p": true, "a": { "b": "123" } }`, a)

	// ToStateFunc writes to state since condition returns true
	s, err := p.ToStateFunc(`{ "p": true }`, `{ "a": { "b": "123" } }`)
	require.NoError(t, err)
	assert.JSONEq(t, `{ "p": true, "t": { "s": "123" } }`, s)
}

func TestConditionalFalse(t *testing.T) {
	p := configs.Conditional("a.b", "t.s", func(state string) bool {
		return false
	})

	// FromStateFunc writes to api regardless of condition result
	a, err := p.FromStateFunc(`{ "p": true }`, `{ { "t": { "s": "123" } }`)
	require.NoError(t, err)
	assert.JSONEq(t, `{ "p": true, "a": { "b": "123" } }`, a)

	// ToStateFunc does not write to state since condition returns false
	s, err := p.ToStateFunc(`{ "p": true }`, `{ "a": { "b": "123" } }`)
	require.NoError(t, err)
	assert.JSONEq(t, `{ "p": true }`, s)

}

func TestDiscriminator(t *testing.T) {
	p := configs.Discriminator("f", map[string]string{
		"foo": "FOO",
		"bar": "BAR",
	})

	// foo exists in state so f will be 'FOO' in api
	a, err := p.FromStateFunc(`{ "p": true }`, `{ "foo": true }`)
	require.NoError(t, err)
	assert.JSONEq(t, `{ "p": true, "f": "FOO" }`, a)

	// neither foo or bar exist in state so f will be empty in api
	a, err = p.FromStateFunc(`{ "p": true }`, `{ "notfoo": true }`)
	require.NoError(t, err)
	assert.JSONEq(t, `{ "p": true }`, a)

	// Discirminator ToStateFunc does not alter state
	s, err := p.ToStateFunc(`{ "p": true }`, `{ "f": "FOO" }`)
	require.NoError(t, err)
	assert.JSONEq(t, `{ "p": true }`, s)

}

func TestEquals(t *testing.T) {
	f := configs.Equals("a", "VALUE")
	assert.True(t, f(`{"a":"VALUE"}`))
	assert.False(t, f(`{"a":"NOT VALUE"}`))
	assert.False(t, f(`{"b":"VALUE"}`))
}
