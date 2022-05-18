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

func TestSkipZeroValue(t *testing.T) {
	assert.True(t, configs.SkipZeroValue(""))
	assert.True(t, configs.SkipZeroValue(0))
	assert.True(t, configs.SkipZeroValue(false))
	assert.True(t, configs.SkipZeroValue([]interface{}{}))
	assert.False(t, configs.SkipZeroValue("123"))
	assert.False(t, configs.SkipZeroValue(123))
	assert.False(t, configs.SkipZeroValue(true))
	assert.False(t, configs.SkipZeroValue([]interface{}{1, 2, 3}))
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
	p := configs.Discriminator("f", map[string]interface{}{
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

func TestArrayWithStrings(t *testing.T) {
	p := configs.ArrayWithStrings("oneTrustCookieCategories.web", "oneTrustCookieCategory", "onetrust_cookie_categories.0.web")

	a, err := p.FromStateFunc(`{}`, `{ "onetrust_cookie_categories": [ { "web": [ "a", "b" ] } ]}`)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"oneTrustCookieCategories": {
			"web": [
				{ "oneTrustCookieCategory": "a" },
				{ "oneTrustCookieCategory": "b" }
			]
		}
	}`, a)

	s, err := p.ToStateFunc(`{}`, `{
		"oneTrustCookieCategories": {
			"web": [
				{ "oneTrustCookieCategory": "a" },
				{ "oneTrustCookieCategory": "b" }
			]
		}
	}`)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"onetrust_cookie_categories": [{
			"web": [ "a", "b" ]
		}]
	}`, s)
}

func TestArrayWithObjects(t *testing.T) {
	p := configs.ArrayWithObjects("eventChannelSettings", "event_channel_settings", map[string]string{
		"eventName":    "name",
		"eventChannel": "channel",
		"eventRegex":   "regex",
	})

	a, err := p.FromStateFunc(`{}`, `{
		"event_channel_settings": [
			{ "name": "n1", "channel": "c1", "regex": "r1" },
			{ "name": "n2", "channel": "c2", "regex": "r2" }
		]
	}`)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"eventChannelSettings": [
			{ "eventName": "n1", "eventChannel": "c1", "eventRegex": "r1" },
			{ "eventName": "n2", "eventChannel": "c2", "eventRegex": "r2" }
		]
	}`, a)

	s, err := p.ToStateFunc(`{}`, `{
		"eventChannelSettings": [
			{ "eventName": "n1", "eventChannel": "c1", "eventRegex": "r1", "extra": "e1" },
			{ "eventName": "n2", "eventChannel": "c2", "eventRegex": "r2", "extra": "e2" }
		]
	}`)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"event_channel_settings": [
			{ "name": "n1", "channel": "c1", "regex": "r1" },
			{ "name": "n2", "channel": "c2", "regex": "r2" }
		]
	}`, s)
}
