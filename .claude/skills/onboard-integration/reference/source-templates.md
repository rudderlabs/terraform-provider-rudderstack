# Source Code Templates

## Modify: `rudderstack/integrations/sources/sources.go`

Add a new `c.Sources.Register(...)` block at the end of the `init()` function. Most sources use `SkipConfig: true` with empty properties:

```go
c.Sources.Register("{name}", c.ConfigMeta{
	APIType:    "{APIType from db-config.json}",
	Properties: []c.ConfigProperty{},
	SkipConfig: true,
})
```

## Modify: `rudderstack/integrations/sources/sources_test.go`

Add TWO new test functions at the end of the file — one unit test and one E2E acceptance test. Note: sources use `configs` (NOT aliased as `c`):

```go
// Unit test
func TestSourceResource{PascalCaseName}(t *testing.T) {
	cmt.AssertSource(t, "{name}", []configs.TestConfig{configs.EmptyTestConfig})
}

// E2E acceptance test
func TestAccSource{PascalCaseName}(t *testing.T) {
	acc.AccAssertSource(t, "{name}", emptyTestConfigs)
}
```

The `acc` import is already present in `sources_test.go`. The `emptyTestConfigs` var is a shared package-level variable.

## Create: `examples/source_{name}.tf`

```hcl
resource "rudderstack_source_{name}" "example" {
  name = "example-{name}"
}
```

## Create: `templates/resources/source_{name}.md.tmpl`

```markdown
---
page_title: "rudderstack_source_{name} Resource - terraform-provider-rudderstack"
subcategory: ""
description: |-

---

# rudderstack_source_{name} (Resource)

This resource represents a {Display Name} event stream source. For more information check https://www.rudderstack.com/docs/sources/{name-with-hyphens}/

## Example Usage

{{tffile "examples/source_{name}.tf"}}

{{ .SchemaMarkdown | trimspace }}
```
