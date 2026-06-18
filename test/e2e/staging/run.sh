#!/usr/bin/env bash
# run.sh — Staging smoke runner for the terraform-provider-rudderstack.
#
# Usage:
#   ./run.sh [--backfill] [path/to/secret.tfvars]
#
# The script:
#   1. Builds the provider locally (into .bin/) and wires a TF dev-override.
#   2. Runs terraform apply.
#   3. Asserts no drift via terraform plan -detailed-exitcode.
#   4. Prints the created resource IDs from terraform output.
#   5. On EXIT (success or failure) runs terraform destroy to clean up staging.
#
# With --backfill the script would trigger and poll a sync on the rETL
# connection; however, no trigger endpoint currently exists in the
# rudder-iac client (v0.17.1).  That branch prints a clear TODO and exits
# non-zero — see the --backfill section below for the tracking comment.
#
# Prerequisites:
#   - Go ≥ 1.21
#   - Terraform ≥ 1.0
#   - secret.tfvars (or the path supplied as $1 / $TFVARS_FILE) with at least:
#       access_token, bq_project, bq_dataset, bq_table, bq_credentials

set -euo pipefail

# ── Argument parsing ────────────────────────────────────────────────────────
BACKFILL=false
TFVARS_FILE="${TFVARS_FILE:-}"   # can also be set in the environment

args=()
for arg in "$@"; do
  case "$arg" in
    --backfill) BACKFILL=true ;;
    *) args+=("$arg") ;;
  esac
done

# First non-flag positional arg overrides the var-file path.
if [[ ${#args[@]} -gt 0 ]]; then
  TFVARS_FILE="${args[0]}"
fi

# Resolve relative paths against the script's own directory.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ -z "${TFVARS_FILE}" ]]; then
  TFVARS_FILE="${SCRIPT_DIR}/secret.tfvars"
fi

if [[ ! -f "${TFVARS_FILE}" ]]; then
  echo "ERROR: var-file not found: ${TFVARS_FILE}"
  echo "       Create it (see README) or pass a path as the first argument."
  exit 1
fi

# ── Repo root (two levels above test/e2e/staging) ───────────────────────────
REPO_ROOT="$(cd "${SCRIPT_DIR}/../../.." && pwd)"

# ── Provider binary paths ────────────────────────────────────────────────────
BIN_DIR="${SCRIPT_DIR}/.bin"
PROVIDER_BINARY="${BIN_DIR}/terraform-provider-rudderstack"

# ── Temp dev-override CLI config ─────────────────────────────────────────────
OVERRIDE_CFG="$(mktemp "${TMPDIR:-/tmp}/tf-dev-override-XXXXXX.tfrc")"
trap 'rm -f "${OVERRIDE_CFG}"' EXIT

# ── Build the provider ───────────────────────────────────────────────────────
echo "==> Building provider binary …"
mkdir -p "${BIN_DIR}"
go build -o "${PROVIDER_BINARY}" "${REPO_ROOT}"
echo "    Provider written to: ${PROVIDER_BINARY}"

# ── Write dev-override config and export for Terraform ───────────────────────
cat > "${OVERRIDE_CFG}" <<HCL
provider_installation {
  dev_overrides {
    "rudderstack.com/rudderlabs/rudderstack" = "${BIN_DIR}"
  }
  direct {}
}
HCL
export TF_CLI_CONFIG_FILE="${OVERRIDE_CFG}"
echo "==> TF_CLI_CONFIG_FILE=${TF_CLI_CONFIG_FILE}"

# ── Destroy trap (always runs on EXIT) ───────────────────────────────────────
# Best-effort: errors here must NOT mask the original script exit code.
_destroy() {
  local original_exit=$?
  echo "==> [trap] Running terraform destroy (cleanup) …"
  terraform -chdir="${SCRIPT_DIR}" destroy -auto-approve \
    -var-file="${TFVARS_FILE}" || true
  echo "==> [trap] Destroy complete."
  rm -f "${OVERRIDE_CFG}"  # this trap replaces the earlier rm-only trap; fold its cleanup in here
  exit "${original_exit}"
}
trap '_destroy' EXIT

# ── Apply ────────────────────────────────────────────────────────────────────
echo "==> Running terraform apply …"
terraform -chdir="${SCRIPT_DIR}" apply -auto-approve \
  -var-file="${TFVARS_FILE}"
echo "==> Apply succeeded."

# ── Assert: plan must show zero drift ────────────────────────────────────────
# terraform plan -detailed-exitcode exit codes:
#   0 = success, no diff (what we want)
#   1 = error
#   2 = success, but diff exists (provider Read is wrong / resource drifted)
echo "==> Asserting no drift (terraform plan -detailed-exitcode) …"
set +e
terraform -chdir="${SCRIPT_DIR}" plan -detailed-exitcode \
  -var-file="${TFVARS_FILE}"
PLAN_EXIT=$?
set -e

case "${PLAN_EXIT}" in
  0)
    echo "==> ASSERT PASSED: plan reports no drift — provider Read path is correct."
    ;;
  1)
    echo "FAIL: terraform plan returned an error (exit 1)."
    echo "      Check the output above for details."
    exit 1
    ;;
  2)
    echo "FAIL: terraform plan detected drift (exit 2)."
    echo "      At least one resource does not match the control-plane state."
    echo "      This indicates a bug in a provider Read/Refresh function."
    exit 2
    ;;
  *)
    echo "FAIL: terraform plan returned unexpected exit code ${PLAN_EXIT}."
    exit "${PLAN_EXIT}"
    ;;
esac

# ── Print created IDs ────────────────────────────────────────────────────────
echo "==> Resource IDs created in staging:"
terraform -chdir="${SCRIPT_DIR}" output

# ── Optional backfill / sync trigger ─────────────────────────────────────────
if [[ "${BACKFILL}" == "true" ]]; then
  # NOTE: The rudder-iac client (github.com/rudderlabs/rudder-iac,
  # currently at v0.17.1-0.20260612051227-31f63ee269cf) does NOT expose a
  # sync-trigger or backfill endpoint.  The RETLConnectionStore interface
  # (api/client/retl/retl.go) only covers CRUD + SetConnectionExternalId;
  # there is no Trigger/Run/StartSync method and no corresponding HTTP path.
  #
  # Until that endpoint is added to the client (or confirmed via the
  # RudderStack public API docs), this branch cannot be wired correctly.
  #
  # What needs confirming:
  #   - The HTTP method and path for triggering a manual rETL sync, e.g.:
  #       POST /v2/retl-connections/{id}/start    (guessed — NOT confirmed)
  #       POST /v2/retl-connections/{id}/trigger  (guessed — NOT confirmed)
  #   - The polling endpoint and its terminal state field, e.g.:
  #       GET  /v2/retl-connections/{id}/sync-status/{syncId}
  #   - The request/response shape.
  #
  # File to watch: api/client/retl/connections.go in the rudder-iac repo.
  # Once that method exists, replace the block below with the real curl.
  echo ""
  echo "ERROR: --backfill is not yet wired."
  echo ""
  echo "  No sync-trigger endpoint was found in the rudder-iac client"
  echo "  (api/client/retl/connections.go, v0.17.1).  The API path and"
  echo "  polling contract need to be confirmed before this can be"
  echo "  implemented correctly."
  echo ""
  echo "  See the comment block in run.sh (--backfill section) for details."
  exit 3
fi

echo "==> Smoke run complete."
