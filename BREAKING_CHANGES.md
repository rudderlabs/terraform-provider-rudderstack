# Breaking changes

## Secret (Sensitive) destination fields now show a permanent plan diff

**Applies to:** the release that adds handling for backend secret redaction (the
version containing this document). Affects every `rudderstack_destination_*`
resource that has a secret field — e.g. `api_key`, `api_secret`, `access_key`,
`access_key_id`, `private_key`, `ca_certificate`, `event_key`, webhook `headers`,
and similar credential fields.

### What you'll see

After you upgrade, `terraform plan` will **always report a change** on the secret
fields of these destinations, even when you've changed nothing:

```
  ~ resource "rudderstack_destination_amplitude" "example" {
      ~ config {
          ~ api_secret = (sensitive value)
        }
    }
```

`terraform apply` re-sends the secret (from your configuration) and succeeds. The
next `plan` shows the same diff again. This is **expected and harmless** — it is
not an error, and it does not recreate or disable the destination.

### Why this happens

The RudderStack API was hardened to **stop returning secret values in destination
read responses** (they come back empty/redacted). Terraform reconciles desired
state (your `.tf`) against the *actual* state it reads from the API. Because the
API no longer returns the secret, the provider cannot know the secret's current
value on the backend, so it **cannot detect whether it matches your configuration**.

Faced with an unreadable field, the provider keeps your **configuration
authoritative**: on every apply it re-asserts the secret from your `.tf`. The
visible cost is this permanent diff. The alternative — silently assuming the
secret is unchanged — was rejected because it would let the value drift out of
sync with your configuration without Terraform ever noticing or correcting it.

### What this means for you

- **No action is required.** Applies remain safe and idempotent in effect; the
  secret you declare in configuration is what ends up on the destination.
- **Configuration is the source of truth.** If someone changes a destination's
  secret in the RudderStack dashboard, your next `terraform apply` overwrites it
  back to the value in your `.tf`.
- **Plan review / CI:** pipelines that fail on a non-empty plan, or that gate on
  "no changes", will now see a diff for secret-bearing destinations. Adjust those
  checks to tolerate secret-field changes on these resources.
- **Keep secrets in variables**, as you should already, so the perpetual diff
  renders as `(sensitive value)` and never prints the secret.

### Notes

- Non-secret fields are unaffected — they still read back and reconcile normally,
  and a real drift on them is still detected.
- Importing a destination brings in its non-secret configuration; secret fields
  come in empty (the API can't return them) and are reconciled from your
  configuration on the next apply.
