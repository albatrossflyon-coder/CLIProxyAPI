# Monthly Provider / Free-Model Refresh — Procedure

**Why:** the free-LLM providers rotate models on a roughly monthly cadence — models get
renamed, pulled to paid, or added. A pinned model ID that worked last month can 404/403
this month with no warning. Run this at the **start or end of every month**.

- **Last full pass:** 2026-08-31 (all providers). **Partial (mistral, aihubmix, +B.ai, +OmniRoute merge):** 2026-09-03.
- **Next full pass due:** 2026-10-01.
- Tracked as a recurring item in `Obsidian Vault\Projects\CLIProxyAPI.md` and at the top of `FREE_MODELS_REGISTRY.md`.

## The two config files (keep in sync — every change goes in BOTH)

| File | Role |
|---|---|
| `C:\Users\albat\apps\EasyCLIProxyAPI\EasyCLIProxyAPI-v0.2.64-Windows-amd64\cpa-core\config.yaml` | **LIVE** — what the running proxy reads. JSON-ish YAML, quoted keys. |
| `C:\Repos\CLIProxyAPI\config.yaml` | repo copy (gitignored, documentation-of-intent). Plain YAML, unquoted keys, inline `# comments`. |

Key adds hot-reload after a save. **New provider blocks and new model aliases need a proxy restart** (Chris's 2-click restart — do NOT `Stop-Process` it).

## Concepts (keys vs aliases)

- **Provider block** = one upstream service + one account (`base-url` + keys + models).
- **Key** (`api-key-entries`) = one credential. Multiple keys → the proxy rotates them for RPM headroom. Does NOT add models.
- **Alias** (`models:`) = one model exposed through the proxy, mapped to a real upstream model name. One key can serve many aliases.
- Client key for testing through the proxy: `clip-...a079` → `http://localhost:8317/v1`.

## Harness scripts (in the session scratchpad; recreate from here if gone)

- `aliastest.py <provider...>` — tests every configured alias for a provider **through the proxy** (the real path; avoids the Cloudflare `error code: 1010` bot-block that hits direct `urllib` calls to groq/airforce/aionlabs).
- `provtest.py <provider...>` / `provtest.py --models <provider> <id...>` — tests **directly** against a provider's raw API + dumps its `/models` catalog. Use for catalog discovery.
- `inventory.py` — dumps the provider/key/alias table from the live config.
- `keyinv.py` — masked key inventory with account labels.

## Per-provider loop (one provider at a time — never batch, attribution gets confused)

1. **`aliastest.py <provider>`** — which configured aliases still return 200 through the proxy.
2. For any FAIL: read the error.
   - `403 ... "not available in your subscription tier"` / `"Deposit required"` / `"balance insufficient"` → **paid now, remove the alias.**
   - `404` / `"model_not_found"` / `"currently unavailable"` → **model ID rotated out.** Pull the provider's `/models` catalog (`provtest.py`), find the successor, repoint.
   - `400 "no_available_channel"` (aihubmix) → the free channel is intermittently down. Verify directly; if it works directly, keep but note the instability.
   - `error code: 1010` → Cloudflare bot-block on the *test*, not a real failure. Re-test through the proxy.
   - `429` in a rapid batch → your own rate-limiting. Re-test slowly (5s apart) before concluding.
3. **Expand:** pull `/models`, look for the free set:
   - Explicit `-free` / `:free` suffix (aihubmix, orcarouter, openrouter).
   - Provider is wholly free/rate-limited (mistral "Experiment", pollinations, airforce, llm7 "turbo" tier).
   - Credit-based with some 0-cost models (B.ai) — test to find which are free.
   - Add an alias for every free chat/coding/reasoning model that returns 200 on a real completion. Skip embed / vision-only / audio / TTS / moderation unless asked. Prefer `-latest` over dated IDs where the provider offers it.
4. **Update BOTH config files.** Alias naming: `<provider>-<short-model-name>`.
5. Re-run `aliastest.py <provider>` after the restart to confirm.

## After the pass

- `python -c "import yaml; yaml.safe_load(open(<file>))"` on **both** files.
- Ask Chris to restart the core.
- `aliastest.py` (no args) — full sweep, everything green or noted.
- Update `FREE_MODELS_REGISTRY.md` (per-provider findings + date) and `MODEL_COOLDOWN.md` (dead/paid models).
- BUILDLOG entry: what rotated, what was added, what's flaky.
- Bump the dates at the top of this file, `FREE_MODELS_REGISTRY.md`, and the vault project file.

## Known provider quirks (2026-09)

| Provider | Note |
|---|---|
| mistral | Free "Experiment" tier = most of the catalog, rate-limited. `mistral-large` is 403 (paid). Pin `-latest`, they churn dated IDs. |
| aihubmix | 55 `-free` models, explicitly "trial only, expect 429s". Fallback pool. Grok models are paid. |
| zai / zai-spare | `glm-*-flash` are genuinely free; `glm-5.3` / `glm-4.6` are PAID per-token (429 code 1113 "余额不足" = insufficient balance). |
| groq | `error code: 1010` on direct calls = Cloudflare block, test via proxy. |
| B.ai | Credit-based. Free: hy3, mimo-v2.5, glm-5.3-flash, qwen3.8-flash. Rest 403 "Deposit required". |
| cloudflare / cloudflare-2 | Two separate accounts (account id in the base-url) — cannot merge. |
| clod / clod-2 / clod-3 | Same `api.clod.io`, 3 one-key blocks. COULD merge into one block + 3 keys (would rename `clod2-*`/`clod3-*` aliases — breaking). Left split. |
