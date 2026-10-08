# Model Cooldown

Models/keys pulled out of `config.yaml`'s active rotation because they failed
a real, direct test — not deleted from the environment, just held here until
someone re-checks and moves them back. **When re-adding one to `config.yaml`,
remove its entry here.** Don't leave a working model listed here alongside a
duplicate in `config.yaml` — one source of truth at a time.

## RESOLVED 2026-08-31 (TIC 1) — the "unresolved Go registration bug" was never a Go bug

**Real root cause found, not source-level at all.** The actual running CLIProxyAPI
process is NOT started from this repo. It's `EasyCLIProxyAPI` (a GUI wrapper,
installed separately), whose core binary runs from
`C:\Users\albat\apps\EasyCLIProxyAPI\EasyCLIProxyAPI-v0.2.64-Windows-amd64\cpa-core\`
with its own `config.yaml` — a completely different file from this repo's
`config.yaml`. **Every edit made to this repo's `config.yaml` since roughly
2026-08-25/28 never reached the live service**, because nobody was editing the
file the running process actually reads.

The live file stores its provider list as one flattened single-line JSON value
(`openai-compatibility: [{...}]`, line ~1466) rather than normal multi-line YAML
— a format the EasyCLIProxyAPI GUI itself writes/reads, not something a human
naturally notices while scrolling a mostly-commented 1400-line template. That
line was a stale snapshot missing every provider/key added since ~2026-08-25/28:
`qwen`, `llm7`, `chinaapi`, `aionlabs` (whole blocks), `openrouter-hy3` /
`orcarouter-free` (aliases on existing blocks), plus 6 individual keys added to
existing blocks (zai key #4, bazaarlink key #4, orcarouter keys #3-5, deepseek
keys #2-3). It also still carried two aliases (`nemotron-nano`, `gpt-oss-20b`)
that were correctly removed from this repo's config on 2026-08-30 after real
404s — proof the two files had been diverging in both directions.

**Fix applied:** wrote a script (`sync_config.py`, scratchpad) that parses this
repo's `config.yaml` `openai-compatibility` section with PyYAML and replaces the
live file's single JSON line with it verbatim, byte-for-byte structurally
equivalent. Backed up the live file first. Validated the new line re-parses as
JSON before restarting. Restarted `cli-proxy-api.exe` directly (the
`CLIProxyAPIStartup` scheduled task that used to manage this is `Disabled` — the
EasyCLIProxyAPI GUI, PID confirmed separately running, is the real supervisor
now; manual restart did not orphan or duplicate it, verified via `Get-Process`
before/after).

**Every previously "unknown provider" alias now genuinely works, confirmed via
real completions post-restart, not just `/v1/models` listing:**
`chinaapi-deepseek-flash` (200), `qwen-flash` (200), `aion-2.0` (200),
`orcarouter-free` (200), `openrouter-hy3` (200). `llm7-claude-sonnet-5` now
loads and routes correctly too, but still legitimately 402s (see Category B
below) — that part was real and unrelated to the sync bug.

Two "real connection failures (HTTP 000)" logged below on 2026-08-30 also
turned out to be casualties of this same stale-config state, not genuine dead
endpoints: **`nemotron-ultra` and `requesty-nemotron-ultra` both now return
real 200s** on retest post-fix. Moved to "confirmed live" below.

**Going forward: edit `C:\Users\albat\apps\EasyCLIProxyAPI\EasyCLIProxyAPI-v0.2.64-Windows-amd64\cpa-core\config.yaml`
directly (or re-run the sync script after editing this repo's `config.yaml`) —
this repo's `config.yaml` is documentation/tracking of intent, not the live
source of truth, until a better sync habit exists.** Consider this the standing
correction to every prior session's assumption that editing this repo's file
was sufficient.

---

## Category A — Needs money / real balance-blocked (key valid, would work if funded)

Confirmed via direct curl to each, individually, with delay between requests:

- **`aihubmix-grok-4.6`, `aihubmix-grok-build`** — HTTP 403, "Your account balance is insufficient."
- **`mistral-large`** — HTTP 403, "This model is not available in your subscription tier" (a tier gate, not balance — the other 4-5 Mistral aliases are fine).
- **`glm-coding`, `glm-4.6`** — HTTP 429, "余额不足或无可用资源包" (insufficient balance/no resource package) — matches the pre-existing config comment marking both PAID despite the account's free-tier signup.
- **`cerebras-oss`** — HTTP 402, payment required.
- **`siliconflow-qwen`** — HTTP 402, insufficient balance.
- **`vercel-poolside-free`** — HTTP 403, needs a credit card on file account-wide (blocks every model, not just this one).
- **`pollinations-openai`** — HTTP 402, insufficient balance — real free weekly grant exists but was never claimed (check enter.pollinations.ai/keys dashboard).
- **`deepseek-chat`, `deepseek-reasoner`** — HTTP 402, insufficient balance on all 3 configured keys.
- ~~`orcarouter-auto` — HTTP 402, out of credits~~ **FIXED 2026-08-31**: repointed the alias from upstream `orcarouter/auto` (credit-gated) to `orcarouter/free` (the same working free route as `orcarouter-free`), per Chris's request. Verified live via a real completion. No longer in this category.
- **`llm7-claude-sonnet-5`, `llm7-deepseek-flash`, `llm7-gemini-flash`** — HTTP 402, insufficient balance. Real finding (unchanged from before): the dashboard's advertised "1,000,000 token/day Free plan" doesn't apply to any model this key can actually reach — every exposed model is pro/turbo-tier, real per-token cost. Not a renewing free provider despite the signup page.

## Category B — Cooling down / transient, likely self-resolving, no action needed

- **`mistral` key #1** (`rNDY0...`) — exhausted, documented to reset around 2026-08-31 (today). Worth a spot-check next session; keys #2-6 cover the pool regardless.
- **`zai` / `zai-spare` (`glm-free`, `glm-free-spare`)** — both returned HTTP 000 (connection failure, not a 4xx) on 2026-08-31 02:28 CDT, tested seconds apart against `open.bigmodel.cn` on two DIFFERENT keys/blocks on the same host. Same-host-both-keys-down pattern points to a real, current outage or throttle on Z.AI's endpoint itself, not a config or key problem — `glm-free` has real prior successful calls logged (`usage.db`), so this reads as transient. Re-check next session before assuming anything is broken.

## Category C — Confirmed genuinely live (real completions, not just listed)

`aihubmix-gemini-free`, `mistral-medium`, `mistral-small`, `mistral-codestral`, `mistral-devstral`, `glm-free` (intermittent — see Category B), `groq-oss`, `cf-oss`, `cf-oss-2`, `bazaarlink-auto`, `nararouter-agnes`, `agnes-ai-flash`, `hiapi-deepseek-flash`, `airforce-gpt-oss-20b`, `orca-qwen-free`, `orca-hy3-free`, `orcarouter-free`, `omniroute-auto`, `omniroute-hy3`, `chinaapi-deepseek-flash`, `qwen-flash`, `aion-2.0`, `openrouter-hy3`, `nemotron-ultra`, `requesty-nemotron-ultra`.

## Category D — Working, confirmed reachable, but genuinely never called yet (per `usage.db`, 96 lifetime events)

Not a mystery — these only became reachable tonight (2026-08-31) once the config-sync bug above was fixed, or are untested siblings of an alias that just got confirmed:
- `aion-3.0` (sibling of `aion-2.0`, same key, same account — expect it to work)
- `qwen-plus`, `qwen-coder-plus`, `qwen-coder-480b` (siblings of `qwen-flash` — each has its OWN separate 1M-token quota per Chris's dashboard, genuinely worth exercising independently, not redundant with `qwen-flash`)
- `omniroute-hy3` (sibling of `omniroute-auto`, added 2026-08-31 by TIC 2, confirmed live via direct curl to OmniRoute itself before this fix — just never reachable through CLIProxyAPI until tonight)
- `llm7-deepseek-flash`, `llm7-gemini-flash` (siblings of the confirmed-402 `llm7-claude-sonnet-5` — expect the same balance block, not independently re-tested)

## Manus

**Real API shape confirmed 2026-08-31 (TIC 2) — structurally incompatible with a normal CLIProxyAPI provider block, same category as the AI Horde evaluation.** Manus's real API (`https://api.manus.ai` or `https://api.manus.im`, both resolve the same) is an **async task system** (`POST /v1/tasks` → poll `GET /v1/tasks/{id}`), not a synchronous `/v1/chat/completions` endpoint like every other provider block in this file. It also requires a non-standard `API_KEY` header instead of `Authorization: Bearer`. Not added to `config.yaml`. Would need a real custom adapter (async task creation + polling wrapper) to ever route through CLIProxyAPI, same as AI Horde — deprioritized pending Chris's call on whether that's worth building.

**Re-verified 2026-09-03 (TIC 1): ALL 3 keys live** — `GET https://api.manus.ai/v1/tasks` with the `API_KEY` header returned `HTTP 200` for each. `albatrossflyon1` is NOT exhausted (the 2026-08-30 assumption was wrong). Still blocked from CLIProxyAPI by the async-API-shape issue above, not by key health.

- **`albatrossaionline` Manus key** — `sk-...Trsu`. Live 2026-08-31 and 2026-09-03 (HTTP 200, task history present).
- **`albatrossflyon1` Manus key** — `sk-...-EMY`. Live 2026-09-03 (HTTP 200, task history present).
- **`albatrossflyonGMX` Manus key** — `sk-...wvZt`. Live 2026-09-03 (HTTP 200, empty task list — still unused).
