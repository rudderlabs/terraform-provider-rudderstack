#!/usr/bin/env bash
# run.sh — Per-scenario staging smoke runner for terraform-provider-rudderstack.
#
# Usage:
#   ./run.sh [--backfill] [--tfvars PATH] [scenario ...]
#   TFVARS_FILE=path/to.tfvars ./run.sh            # tfvars via env
#   PAUSE=true ./run.sh customerio                 # hold open before destroy
#
# Each scenario is an independent Terraform root module under scenarios/<name>/
# (see README). With NO scenario args, runs ALL scenarios (skipping _-prefixed
# helper dirs) sequentially: build provider once → per scenario
# init → apply → verify outputs → assert no drift → destroy.
#
# Keep-going: a failed scenario is destroyed and the run continues; the script
# exits non-zero at the end if ANY scenario failed. At most one scenario is live
# in staging at a time.
#
# The script:
#   1. Builds the provider locally (into .bin/) and wires a TF dev-override.
#   2. For each selected scenario, applies, asserts no drift, then destroys.
#   3. On EXIT (incl. failure) destroys any in-flight scenario to clean up staging.
#
# With --backfill the script would trigger and poll a sync on an rETL connection;
# no trigger endpoint exists in the rudder-iac client yet, so that branch prints a
# TODO and exits non-zero (see the --backfill section below).
#
# Prerequisites:
#   - Go ≥ 1.21
#   - Terraform ≥ 1.0
#   - secret.tfvars (or --tfvars / $TFVARS_FILE) with at least:
#       access_token, bq_project, bq_dataset, bq_table, bq_credentials
#     plus the creds each selected scenario needs (see README).

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
SCENARIOS_DIR="${SCRIPT_DIR}/scenarios"
BIN_DIR="${SCRIPT_DIR}/.bin"
PROVIDER_BINARY="${BIN_DIR}/terraform-provider-rudderstack"

# ── Argument parsing ────────────────────────────────────────────────────────
BACKFILL=false
TFVARS_FILE="${TFVARS_FILE:-}"   # can also be set in the environment
REQUESTED=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --backfill) BACKFILL=true; shift ;;
    --tfvars)   [[ $# -ge 2 ]] || { echo "ERROR: --tfvars requires a path argument"; exit 1; }
                TFVARS_FILE="$2"; shift 2 ;;
    --tfvars=*) TFVARS_FILE="${1#*=}"; shift ;;
    -*) echo "ERROR: unknown flag: $1"; exit 1 ;;
    *)  REQUESTED+=("$1"); shift ;;
  esac
done

if [[ -z "${TFVARS_FILE}" ]]; then
  TFVARS_FILE="${SCRIPT_DIR}/secret.tfvars"
fi
if [[ ! -f "${TFVARS_FILE}" ]]; then
  echo "ERROR: var-file not found: ${TFVARS_FILE}"
  echo "       Create it (see README) or pass --tfvars PATH."
  exit 1
fi
# Absolute path — terraform -chdir resolves -var-file relative to the scenario dir.
TFVARS_FILE="$(cd "$(dirname "${TFVARS_FILE}")" && pwd)/$(basename "${TFVARS_FILE}")"

# ── Discover scenarios (skip _-prefixed helper dirs) ────────────────────────
ALL_SCENARIOS=()
for d in "${SCENARIOS_DIR}"/*/; do
  [[ -d "$d" ]] || continue    # no match → glob stays literal; skip it
  name="$(basename "$d")"
  [[ "$name" == _* ]] && continue
  ALL_SCENARIOS+=("$name")
done
if [[ ${#ALL_SCENARIOS[@]} -eq 0 ]]; then
  echo "ERROR: no scenarios found under ${SCENARIOS_DIR}"
  exit 1
fi

# ── Resolve requested set (validate names) ──────────────────────────────────
SCENARIOS=()
if [[ ${#REQUESTED[@]} -eq 0 ]]; then
  SCENARIOS=("${ALL_SCENARIOS[@]}")
else
  for req in "${REQUESTED[@]}"; do
    found=false
    for s in "${ALL_SCENARIOS[@]}"; do
      [[ "$s" == "$req" ]] && found=true && break
    done
    if ! $found; then
      echo "ERROR: unknown scenario: ${req}"
      echo "       valid: ${ALL_SCENARIOS[*]}"
      exit 1
    fi
    SCENARIOS+=("$req")
  done
fi

# ── --backfill short-circuit (not yet wired) ────────────────────────────────
# NOTE: The rudder-iac client does NOT expose a sync-trigger/backfill endpoint.
# The RETLConnectionStore interface (api/client/retl/retl.go) only covers CRUD +
# SetConnectionExternalId; there is no Trigger/Run/StartSync method and no HTTP
# path. Until that endpoint exists (or is confirmed via the public API docs) this
# branch cannot be wired. What needs confirming:
#   - HTTP method/path to trigger a manual rETL sync
#   - polling endpoint + terminal state field
#   - request/response shape
# File to watch: api/client/retl/connections.go in the rudder-iac repo.
if [[ "${BACKFILL}" == "true" ]]; then
  echo ""
  echo "ERROR: --backfill is not yet wired."
  echo "  No sync-trigger endpoint exists in the rudder-iac client"
  echo "  (api/client/retl/connections.go). The API path + polling contract"
  echo "  must be confirmed first. See the comment block in run.sh."
  exit 3
fi

# ── Temp dev-override CLI config (cleaned up by the EXIT trap) ───────────────
# No .tfrc suffix: BSD/macOS mktemp requires the X's at the end. TF_CLI_CONFIG_FILE
# accepts any path, so the extension is unnecessary.
OVERRIDE_CFG="$(mktemp "${TMPDIR:-/tmp}/tf-dev-override-XXXXXX")"
CURRENT_SCENARIO=""

# ── Cleanup trap (always runs on EXIT) ──────────────────────────────────────
# Best-effort: destroy any in-flight scenario (set before apply, cleared after a
# clean destroy) so a mid-run failure doesn't leak staging resources. Errors here
# must NOT mask the original exit code.
_cleanup() {
  local original_exit=$?
  if [[ -n "${CURRENT_SCENARIO}" ]]; then
    echo "==> [trap] Destroying in-flight scenario '${CURRENT_SCENARIO}' …"
    terraform -chdir="${SCENARIOS_DIR}/${CURRENT_SCENARIO}" destroy -auto-approve \
      -var-file="${TFVARS_FILE}" || true
  fi
  rm -f "${OVERRIDE_CFG}"
  exit "${original_exit}"
}
trap '_cleanup' EXIT

# ── Build the provider once ─────────────────────────────────────────────────
echo "==> Building provider binary …"
mkdir -p "${BIN_DIR}"
go build -o "${PROVIDER_BINARY}" "${REPO_ROOT}"
echo "    Provider written to: ${PROVIDER_BINARY}"

# ── Write dev-override config once and export for Terraform ──────────────────
cat > "${OVERRIDE_CFG}" <<HCL
provider_installation {
  dev_overrides {
    "rudderstack.com/rudderlabs/rudderstack" = "${BIN_DIR}"
  }
  direct {}
}
HCL
export TF_CLI_CONFIG_FILE="${OVERRIDE_CFG}"
export TF_IN_AUTOMATION=1
# TF_IN_AUTOMATION only adjusts messaging; TF_INPUT=0 disables interactive prompts
# so a missing variable (e.g. absent creds) fails fast instead of hanging.
export TF_INPUT=0
echo "==> TF_CLI_CONFIG_FILE=${TF_CLI_CONFIG_FILE}"

# Unique per-invocation token folded into resource names (via the run_id var) so
# a run never collides with soft-deleted resources whose names the API keeps
# reserved from earlier runs. Shared across all scenarios in this invocation;
# each scenario has a distinct name prefix, so one token is enough.
export TF_VAR_run_id="${TF_VAR_run_id:-$(date +%s)-$$}"
echo "==> run_id=${TF_VAR_run_id}"

# ── Per-scenario runner ─────────────────────────────────────────────────────
# Returns non-zero on any failure; leaves CURRENT_SCENARIO set so the caller (or
# the EXIT trap) destroys the partially-applied scenario.
run_scenario() {
  local name="$1"
  local dir="${SCENARIOS_DIR}/${name}"
  echo ""
  echo "════════════════════════════════════════════════════════════"
  echo "==> SCENARIO: ${name}"
  echo "════════════════════════════════════════════════════════════"
  CURRENT_SCENARIO="${name}"

  # init wires .terraform/modules for the local bigquery_source module. The
  # dev-override keeps this offline (rudderstack is the only provider referenced).
  echo "==> [${name}] terraform init …"
  # NOTE: run_scenario is called as `if ! run_scenario …`, which disables `set -e`
  # for the whole function body — so every terraform call here must check its own
  # exit status explicitly (an unchecked failure would fall through and be
  # misreported as success).
  terraform -chdir="${dir}" init -input=false \
    || { echo "FAIL [${name}]: terraform init failed."; return 1; }

  echo "==> [${name}] terraform apply …"
  terraform -chdir="${dir}" apply -auto-approve -var-file="${TFVARS_FILE}" \
    || { echo "FAIL [${name}]: terraform apply failed."; return 1; }

  echo "==> [${name}] verifying standardized outputs …"
  local out val
  for out in account_id source_id destination_id connection_id; do
    val="$(terraform -chdir="${dir}" output -raw "$out" 2>/dev/null || true)"
    if [[ -z "$val" ]]; then
      echo "FAIL [${name}]: output '$out' is empty — resource may not have been created."
      return 1
    fi
    echo "    ${out} = ${val}"
  done

  # terraform plan -detailed-exitcode: 0=no diff, 1=error, 2=diff (drift).
  echo "==> [${name}] asserting no drift (terraform plan -detailed-exitcode) …"
  set +e
  terraform -chdir="${dir}" plan -detailed-exitcode -var-file="${TFVARS_FILE}"
  local plan_exit=$?
  set -e
  case "${plan_exit}" in
    0) echo "==> [${name}] ASSERT PASSED: plan reports no drift." ;;
    1) echo "FAIL [${name}]: terraform plan returned an error (exit 1)."; return 1 ;;
    2) echo "FAIL [${name}]: terraform plan detected drift (exit 2) — provider Read bug."; return 1 ;;
    *) echo "FAIL [${name}]: terraform plan returned unexpected exit ${plan_exit}."; return 1 ;;
  esac

  echo "==> [${name}] resource IDs:"
  terraform -chdir="${dir}" output

  if [[ "${PAUSE:-false}" == "true" ]]; then
    echo "==> [${name}] PAUSED. Resources are live in staging. Press Enter to destroy."
    echo "    Inspect in another terminal:"
    echo "      export TF_CLI_CONFIG_FILE='${TF_CLI_CONFIG_FILE}'"
    echo "      terraform -chdir='${dir}' show"
    read -r || true   # blocks; || true prevents set -e from exiting on EOF/signal
  fi

  echo "==> [${name}] terraform destroy …"
  if terraform -chdir="${dir}" destroy -auto-approve -var-file="${TFVARS_FILE}"; then
    CURRENT_SCENARIO=""   # destroyed cleanly — trap must not re-destroy
  else
    # Leave CURRENT_SCENARIO set so the loop / EXIT trap retries cleanup, and fail
    # the scenario instead of silently reporting success with leaked resources.
    echo "FAIL [${name}]: terraform destroy failed — leaving CURRENT_SCENARIO set for cleanup."
    return 1
  fi
  echo "==> [${name}] done."
  return 0
}

# ── Loop: keep-going, fail-at-end ───────────────────────────────────────────
FAILED=()
for name in "${SCENARIOS[@]}"; do
  # `if !` guards the call so set -e doesn't abort the loop on a graceful failure.
  if ! run_scenario "${name}"; then
    FAILED+=("${name}")
    # run_scenario left CURRENT_SCENARIO set on a graceful failure; destroy it now
    # and clear so the EXIT trap doesn't double-destroy.
    if [[ -n "${CURRENT_SCENARIO}" ]]; then
      echo "==> [${name}] cleaning up failed scenario …"
      terraform -chdir="${SCENARIOS_DIR}/${name}" destroy -auto-approve \
        -var-file="${TFVARS_FILE}" || true
      CURRENT_SCENARIO=""
    fi
  fi
done

echo ""
if [[ ${#FAILED[@]} -gt 0 ]]; then
  echo "==> SMOKE FAILED for: ${FAILED[*]}"
  exit 1
fi
echo "==> All scenarios passed: ${SCENARIOS[*]}"
