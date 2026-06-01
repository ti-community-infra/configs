---
name: tidb-eol-cleanup
description: Perform end-of-maintenance (EOM) cleanup for a TiDB LTS version (X.Y) — remove version labels, drop branch protection rules, and update blocker configs across prow configs.
---

# TiDB Version End-of-Maintenance (EOM) Cleanup

Use this skill when a TiDB LTS version reaches end-of-maintenance (EOM). Replace `X.Y` with the actual version number (e.g., `7.1`). Note: this is EOM, not EOL — the version is no longer under maintenance support but the code is not deleted.

## Files Modified

All paths are relative to `prow/config/`:

| File | Changes |
|------|---------|
| `external_plugins_config.yaml` | Issue triage, label plugin, label blocker |
| `labels.yaml` | Remove version-specific label definitions |
| `labels.md` | Regenerated via script |
| `config.yaml` | Branch protection rules (branch-protection section only) |

## Step 1: `external_plugins_config.yaml`

### 1a. Remove `X.Y` from `ti-community-issue-triage` `maintain_versions`

```yaml
# Before:
maintain_versions: ["X.Y", "7.5", "8.1", "8.5"]
# After:
maintain_versions: ["7.5", "8.1", "8.5"]
```

### 1b. Remove version labels from `ti-community-label` `additional_labels`

Search for and remove these labels from ALL `additional_labels` lists under `ti-community-label`:

- `needs-cherry-pick-release-X.Y`
- `affects-X.Y`
- `may-affects-X.Y`

These are quoted or unquoted depending on the list. Remove each occurrence line, leaving any `# - affects-{{.ver}}` comments in place. Expect ~18 occurrences across ~8 repo groups.

Also remove `type/cherry-pick-for-release-X.Y` from the global `cherryPickerExcludeLabels` and `cherryPickerExcludeLabelsWithoutLgtm` lists.

### 1c. Add `X.Y` to `ti-community-label-blocker`

Add a blocking rule at the top of `block_labels` for the `pingcap`/`tikv`/`pingcap-inc`/`PingCAP-QE` orgs group:

```yaml
      - regex: ^(may-)?affects-X\.Y$
        actions:
          - labeled
          - unlabeled
        trusted_users:
          - ti-chi-bot
          - ti-chi-bot[bot]
        message: "You cannot manually add or delete the affects/may-affects version labels, only I and the trusted members have permission to do so."
```

## Step 2: `labels.yaml`

Remove these label definitions from the default labels section:

- `needs-cherry-pick-release-X.Y` (4 lines: name, color, description, target)
- `affects-X.Y` (4 lines: name, color, description, target)
- `type/cherry-pick-for-release-X.Y` (8 lines: name, color, description, target, prowPlugin, isExternalPlugin, addedBy)

After editing, regenerate `labels.md`:

```bash
cd configs
bash scripts/update-labels.sh
```

## Step 3: `config.yaml` — Branch Protection

Remove all `release-X.Y` and older branches from the **branch-protection** section only (do NOT touch the `tide:` section).

### Branches to remove

All `release-M.N` where `M.N <= X.Y` in TiDB version numbering. Keep:
- `release-1.x`, `release-2.x` for tidb-operator/tispark (different product versioning)
- Anything above `X.Y`

### CRITICAL: YAML Anchor Handling

Many branch entries use YAML anchors. When a branch being removed defines an anchor (`&name`), you MUST relocate that anchor to the next-highest branch that remains.

**Step 3a. Simple refs (no anchor)** — just remove the line:
```yaml
# Remove:
            release-X.Y: *some-anchor
```

**Step 3b. Anchor definitions with inline blocks** — remove the block AND move the anchor:
```yaml
# Before (remove release-X.Y with inline block):
            release-X.Y: &anchor-name
              protect: true
              required_status_checks:
                contexts:
                  - "tide"
                strict: false
              restrictions: *branch-protection_restrictions
            release-X.Z: *anchor-name   # next branch
# After (move anchor to next branch, keep inline block):
            release-X.Z: &anchor-name
              protect: true
              required_status_checks:
                contexts:
                  - "tide"
                strict: false
              restrictions: *branch-protection_restrictions
```

**Step 3c. After all edits, verify no null anchors exist.** Every `&anchor-name` must have `protect: true` (or similar non-null content) immediately below it. A bare `release-X.Y: &anchor-name` with no following block means all `*anchor-name` references get `null` protection.

### Repos to check

Common repos with release branches: tikv/client-c, tikv/copr-test, tikv/pd, tikv/tikv, pingcap/monitoring, pingcap/tidb, pingcap/tiflow, pingcap/tidb-binlog, pingcap/tiflash, pingcap-inc/enterprise-extensions, pingcap-inc/enterprise-plugin, pingcap-inc/tiflash-scripts.

### Exceptions

- `pingcap/docs` and `pingcap/docs-cn`: docs team maintains these independently. Do NOT modify.
- `tidb-operator` and `tispark` branches use different version schemes — do NOT touch.

## Gotchas

- The edit tool's fuzzy matching can reformat multi-line YAML strings from single-line `description: text` to multi-line `description:\n  text`. This is semantically equivalent but makes diffs noisy. Be prepared for this.
- Always run `update-labels.sh` after editing `labels.yaml` to keep `labels.md` in sync.
- After branch-protection edits, verify: no null anchors, all `release-7.2`/`7.3`/`7.4` are present if they existed before, and no release ≤ X.Y branches remain (except docs).
