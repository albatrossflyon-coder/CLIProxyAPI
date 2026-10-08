# BUILDLOG — CLIProxyAPI (this deployment)

Third-party repo (`router-for-me/CLIProxyAPI`, not our fork) — this file tracks *our* deployment/config decisions, not upstream development. See `Obsidian Vault\Projects\CLIProxyAPI.md` for the fuller narrative; this is the local, git-untracked mirror so notes land here too as keys get added over time.

## 2026-10-03 6:11 AM CT (TIC 2): monthly key check, dead providers removed, xkiro free Qwen added (Acer only)
- Tested all 123 proxy models one at a time, then live-tested candidate free models directly per provider.
- Removed (live) / commented out with `# DEAD 2026-10-03` (repo copy): aihubmix (10-try limit used, needs top-up), apinex (free models now subscription or daily check-in), tokenrouter (free period ended, $0), chinaapi (needs top-up), nararouter (model removed), monkeycoding claude keys x2 ($0), aliases llm7-gpt-oss, llm7-llama3.1-8b, requesty-nemotron-super, xkiro-deepseek-flash. Gemini: `excluded-models` gemini-2.5-pro, gemini-3-pro-preview (retired).
- Added xkiro: xkiro-qwen3.7-max, -qwen3.7-plus, -qwen3-max, -qwen3.6-plus (all live-tested OK). Old names kept working: `nararouter-agnes` now on agnes-ai, `tokenrouter-glm-5.3-free` now on zai glm-4.7-flash.
- Hot reload picked it up, no restart needed; all 4 new/renamed aliases verified through the proxy. 88 models listed (was 123). Backups: `config.yaml.bak-20261003-0610` next to both files.
- Kept, still failing: orcarouter (GitHub link gate), bai hy3/mimo/qwen-flash (no credit). Temporary 429/503: mistral medium/small/magistral, unorouter x2, gemini-3.1-pro.
- Pre-existing drift: live has `bai-glm-flash` on groq + cloudflare, repo copy does not. VPS config NOT checked yet.

## 2026-09-16 1:42 PM CT (TIC 1) — bai-glm-flash alias fix applied on Acer; VPS still open; OrcaRouter confirmed dead for free access

Applied the config-only fix queued below. Live file: `C:\Users\albat\apps\EasyCLIProxyAPI\EasyCLIProxyAPI-v0.2.64-Windows-amd64\cpa-core\config.yaml`. `bai-glm-flash` now registered under 3 provider blocks: `bai` (original, `glm-5.3-flash`), `groq` (`openai/gpt-oss-120b`), `cloudflare` (`@cf/openai/gpt-oss-120b`). Both new homes live-tested with a real completion through `localhost:8317` before being wired in (returned actual completion bodies, not just 200s).

**Two originally-picked candidates were tested and rejected before landing on groq/cloudflare:**
- `mistral` (`mistral-medium-latest`) — live-tested via the existing `mistral-medium`/`mistral-small` aliases, both returned identical `{"type":"rate_limited","code":"1300"}` seconds apart. Looked like a pool-wide cooldown, not a stable failover target right now. Not added.
- `orcarouter` (`orcarouter/free`) — **real finding, not previously known:** OrcaRouter's free tier is genuinely dead account-wide. Tested 2 separate keys (`...YDs3`, `...P4t`) directly against `api.orcarouter.ai` (bypassing this proxy entirely) — both valid/active per their own `/v1/models`, both hit an identical block calling any free-tagged model: `{"reason":"err_free_access_denied","retryable":false}`. Confirmed as a real, current, documented policy via OrcaRouter's own docs (`docs.orcarouter.ai/routing/free-models`): free access requires the workspace owner to link an established GitHub account, or enough lifetime spend to be exempt. This is a workspace-level gate, invisible on the key-management dashboard page (which only shows key status, not free-tier eligibility) — only surfaces when a free model is actually called. Not fixable from config; needs Chris to link an established GitHub account to the OrcaRouter workspace, or add credits.

**Not yet done:**
- Acer's CLIProxyAPI process needs its 2-click restart (Chris's job per standing rule) before the new `bai-glm-flash` registrations actually load. Once restarted, re-verify the alias round-trips across more than one provider.
- **VPS-side edit not started.** SSH to the VPS is blocked for this session by the platform's safety net. Sent Chris a read-only grep command to pull the VPS's `/data/coolify/services/dwpx7ylvitd0zmmeg6gtwygk/config.yaml` `mistral`/`orcarouter` block format first (to confirm it hasn't drifted from Acer's format before writing a blind edit) — first attempt failed on PowerShell quote-escaping around nested double-quotes, corrected version sent, no output received by session end.

**Also checked and deliberately left alone:** this repo's own `config.yaml` — confirmed badly stale (676 lines vs. 1700+ in the live file, missing the `bai`/`mistral`/`orcarouter` blocks entirely). Editing it would have been documentation noise, not real parity. Reconfirms the standing note below that this file is "documentation of intent only until a sync habit exists" — it currently doesn't have one.

**COMMIT STATUS:** `config.yaml` (live, not this repo) edited outside git. This repo's own `config.yaml` untouched. This BUILDLOG entry is the durable record — Chris commits this repo himself per existing convention.

## 2026-09-15 9:25 PM CT (TIC 2) — Phase 1 done: cross-provider failover already exists in the code, no Go changes needed, this is a config-only fix

Followed up on the 09-14 8:21 PM entry below ("Phase 1, next session, first thing: map the dispatch path"). Traced the real request path end to end, in `sdk/`:

`sdk/api/handlers/handlers_execution.go` (`executeWithAuthManagerFormats`) → `sdk/api/handlers/handlers_routing.go` (`providersForExecution` → `getRequestDetailsWithOptions`) → `internal/util/provider.go` (`GetProviderName(modelName)`, which calls `registry.GetGlobalRegistry().GetModelProviders(modelName)` — **returns every provider that has registered a model under that exact alias, deduplicated**) → `sdk/cliproxy/auth` `Manager.Execute(ctx, providers []string, req, opts)` (doc comment: "supports multiple providers for the same model and round-robins the starting provider per model") → `executeMixedOnce` → `pickNextMixed` → `scheduler.pickMixed` (`sdk/cliproxy/auth/scheduler.go`), which already does cooldown-aware priority-tier/weighted rotation **across every provider in that list**, and only returns "all credentials cooling down" (`selector.go`'s `modelCooldownError`) once every credential across every one of those providers is exhausted.

**Conclusion: cross-provider failover is not missing from CLIProxyAPI's code. It already exists, is tested (`scheduler_test.go` has multiple `TestSchedulerPick_MixedProviders*` cases), and runs automatically whenever two or more provider blocks in `config.yaml` register a model under the same alias name.**

**Why the 09-14 incident happened despite this:** `bai-glm-flash` (the alias hermes-1/omp-1 were both pointed at) was only ever defined once, under the single `bai` provider block. `GetModelProviders("bai-glm-flash")` had exactly one entry to return — there was nothing else in the pool to fail over to. This was never a missing-feature problem; it was a missing-alias-registration problem.

**The actual fix, on both Acer's local `config.yaml` and the VPS's:** give `bai-glm-flash` (or whichever alias an agent is pinned to) a second and third home — register the same alias name under a model entry in 1-2 other reliable provider blocks already in the config (e.g. `openrouter`, `mistral`, or one of the free routers). No code change, no rebuild, no redeploy of new binaries — just a config edit and a restart/reload of the running CLIProxyAPI process on each instance. This supersedes Manus's 10-phase virtual-routing-layer design below (`Resilient Model Failover Architecture for CLIProxyAPI.md`, saved to Downloads 2026-09-14) as the actual path forward — that design solves a problem the code doesn't have; the real gap is just alias coverage.

**Not yet done:** the actual alias-registration edit itself (pick the right secondary/tertiary providers for `bai-glm-flash` and any other single-provider-only aliases, edit both config.yaml copies, restart both instances, verify live). That's the next concrete step, queued for Chris's go-ahead.

## 2026-09-14 8:21 PM CT (TIC 1) — VPS incident: B.ai credit exhaustion took down hermes-1 and omp-1 simultaneously; no cross-route failover exists

Live incident on the Oracle A1 VPS roster. `bai-glm-flash` (the `bai` provider, a credit-based relay per the 2026-09-03 entry below) hit `balance=0` — every key under that provider drained, not just one. Both hermes-1 and Omp-1 were independently configured to default to that exact route and both died with the identical error (`credit insufficient balance: balance=0`) within the same session. Confirmed live via the actual pane output on both agents, not guessed.

**Root cause, confirmed:** CLIProxyAPI already rotates keys *within* one provider block (proven in the actual Go source tonight — per-provider retry/quota-handling files exist, e.g. `claude_ratelimit.go`, `codex_quota.go`). It has **no cross-provider or cross-model failover** — if every key under a provider is exhausted, the request just fails, full stop. Checked the entire `config_types.go` schema directly: no fallback-chain field exists anywhere. This is a real, confirmed architecture gap, not a config mistake.

**Immediate fix applied (workaround, not a real fix):** manually repointed hermes-1 (`~/.hermes/config.yaml` `model.default`) and omp-1 (`~/.omp/agent/models.yml` + `config.yml` `modelRoles.default`) to `orca-qwen-free` instead. Both verified alive afterward via a real prompt round-trip, not just process-started. omp-1 needed the new model *registered* in `models.yml` before `config.yml` would even let it be selected as default — two separate files, two separate edits, for one model swap. This exact manual-edit cost is why the real fix matters: it will recur every time a route drains, for every agent, in a different config format each time.

**Real fix, designed but NOT started:** cross-provider virtual-model routing layer inside (or in front of) CLIProxyAPI. Independently converged on by three separate AI consultations tonight (ChatGPT, a second unnamed model referred to as "Lite", and Manus) — all three landed on the same core shape without seeing each other's answers: agents request a stable alias (e.g. `coding-default`), CLIProxyAPI resolves it to an ordered list of real provider/model routes, walks the list on failover-eligible failures only (insufficient balance, exhausted quota, 5xx, timeout — NOT malformed requests or context-length errors), with circuit breakers (failure-specific cooldowns, e.g. 6h for insufficient-balance) and capability-aware route selection (don't fail over a tool-call request to a route with no tool support). Critical safety rule from two of the three sources independently: failover is only safe *before* the first response byte / before a tool call starts, never mid-stream.

Manus's report includes a complete phased build plan (10 phases, starting with "map the real dispatch path, change nothing" as phase 1) plus working Go sketches for the failure classifier, circuit breaker, and route selector. Full source docs saved: `C:\Users\albat\Downloads\Resilient Model Failover Architecture for CLIProxyAPI.md` (Manus, most complete/implementation-ready) — the other two responses (ChatGPT, "Lite") were chat-pasted, not saved as files, but converged on the same design; see the vault session log `2026-09-14.md` for the fuller comparison of all three.

**Next session, first thing:** Phase 1 of Manus's plan — read CLIProxyAPI's actual HTTP-handler-to-upstream-dispatch code path (never fully located tonight; grep attempts for the routing/dispatch entry point came up empty) and confirm exactly where a model alias resolves to a provider+key today, before writing any new code. This is genuinely a multi-session build, not a quick patch — do not rush Phase 1.

**Also confirmed live tonight, VPS roster status:** all 5 tracked agents (claude-1, claude-2, omp-1, pi-1, hermes-1) brought back alive via `roster-restart.sh FORCE_ALL=1` after being mostly down. `herdr-command-center` dashboard repoint to VPS herdr shipped and verified end-to-end (separate commit, `herdr-command-center` repo) — VPS agents now visible in the dashboard, first real cross-visibility into the VPS roster from outside SSH.

**COMMIT STATUS:** config.yaml unchanged this entry (the fix was in hermes-1/omp-1's own config files on the VPS, not this repo's tracked config). This BUILDLOG entry + the two new scripts (`fix-dead-model-route.sh`, `add-omp-model.sh` in `albatross-automations/herdr-vps/`) are the durable record. albatross-automations commit still owed.

## 2026-09-10 ~5:45 PM CDT (TIC 2) — new provider: xKiro (xkiro.com)

Chris handed over 3 xKiro keys (`...0cbf` XkiroAlbatrossflyonGmx, `...536d` XkiroAlbatrossflyon1, `...81cd` XkiroAlbatrossaionline; separate accounts) + the models URL (`xkiro.com/dashboard/models`). OpenAI-compatible gateway, `https://xkiro.com/v1`, 112 models, `access_tier` field marks ~40 free ($0 in+out).

**Wired in `config.yaml`** (appended after `unorouter`, in `openai-compatibility:`): 3 keys, 5 models — all verified free via a direct `/v1/chat/completions` call 2026-09-10 (returned "OK", cost 0):
`xkiro-qwen-flash` (qwen/qwen3.5-flash:free), `xkiro-qwen-max` (qwen3.8-max:free), `xkiro-qwen-coder` (qwen3-coder-plus:free), `xkiro-qwen-397b` (qwen3.5-397b-a17b:free), `xkiro-deepseek-flash` (deepseek/deepseek-v4-flash).

**Caveats found the same day:** `openai/gpt-5.3-codex-spark` is tagged free but 402s "paying customers only" — not wired. `minimax/minimax-m3:free` + `m2.5:free` returned transient 500s — not wired, retry next pass.

**Verified through the proxy 2026-09-10:** Chris restarted the core; all 5 `xkiro-*` aliases round-tripped `localhost:8317` clean ("PONG"). Keys #2 (`...536d`) and #3 (`...81cd`) were added AFTER that restart, both verified valid directly; they enter the rotation on the next restart. **Not on the VPS instance** — same block goes there during the herdr-VPS Phase 6 session (no IP whitelist). `FREE_MODELS_REGISTRY.md` updated; next full re-test 2026-10-01.

**COMMIT STATUS:** config.yaml + these two docs edited locally, UNCOMMITTED — Chris commits this repo.

## 2026-09-03 ~10:45 PM CDT (TIC 1) — Provider refresh pass (partial) + OmniRoute key merge + B.ai + monthly procedure

Follow-on to the Mistral entry below. Chris asked for a full provider-by-provider free-model refresh and a repeatable monthly procedure. Ran through `start-to-finish`. Config edited in BOTH files (`C:\Repos\CLIProxyAPI\config.yaml` + the live `...\EasyCLIProxyAPI\...\cpa-core\config.yaml`), YAML-validated, Chris restarted the core twice; every change below verified live through the proxy (`localhost:8317`, client key `clip-...`).

**New file: `PROVIDER_REFRESH_PROCEDURE.md`** — the repeatable monthly process (the two config files, keys-vs-aliases, the harness scripts, the per-provider loop, per-provider quirks). This was the main deliverable Chris asked for. Next full pass due 2026-10-01.

**Full 28-provider alias health sweep** (through the proxy — direct calls Cloudflare-block groq/airforce/aionlabs with `error code: 1010`). Everything healthy except the items below.

**aihubmix:** 3 aliases (1 working) → **6 verified free models**. Free tier = every `-free`-suffixed model (55 of them; their docs: "trial only, expect 429s"). Dropped `grok-build-0.1` (400) + `grok-4.6` (403 paid). Added gemini-3.8-flash / gpt-5.5 / gpt-4.1-mini / coding-glm-5.3 / qwen3.6-plus / dots-3, all `-free`. 4 more (gpt-oss-20b, nemotron-ultra 550B, hy3, gemma-31b) verified direct but flaky (400 "no_available_channel") through the proxy — noted in config for next refresh.

**B.ai — NEW provider** (`bai`, `https://api.b.ai/v1`). Credit-based OpenAI/Anthropic-compatible relay: most models 403 "Deposit required", a few are 0-credit. Free (verified): `glm-5.3-flash`, `qwen3.8-flash`, `hy3`, `mimo-v2.5` → aliases `bai-glm-flash` / `bai-qwen-flash` / `bai-hy3` / `bai-mimo`. **3 keys** (albatrossflyon1 / albatrossGMX / albatrossai accounts), all verified.

**zai:** dropped `glm-coding` (glm-5.3) + `glm-4.6` — both paid per-token and all 4 keys now 429 code 1113 "余额不足" (no balance). Kept `glm-free` (glm-4.7-flash). zai-spare `glm-free-spare` (glm-4.5-flash) still fine.

**Key merges:**
- **agnes-ai:** 2 → **5 keys** — Chris pasted 3 (albatrossflyon1, nanobot, AngesAlbatrossGMX accounts), all verified live.
- **OmniRoute merge:** booted the scaled-to-zero Fly machine (`fly scale count 1`), `fly ssh console` + decrypted `provider_connections` from `/data/storage.sqlite` (AES-256-GCM, key = `scrypt(STORAGE_ENCRYPTION_KEY, "omniroute-field-encryption-v1")`), scaled back to zero. 9 connections; net-new usable: **mistral `GzYjgcs15…` (JuneBugMistralGmx, → mistral key #7, verified 200)** and **groq `gsk_cd3shnw8…` (→ groq key #5)**. Rest were dups (glm, openrouter), dead (old mistral, marked invalid in OmniRoute), paid-no-credits (an OpenAI `sk-proj-` key, 429 "no credits"), or audio-only (Sprag TTS/STT). **No Manus keys in OmniRoute.**

**Manus:** re-verified all 3 keys (`MODEL_COOLDOWN.md` `## Manus`) LIVE 2026-09-03 (`GET api.manus.ai/v1/tasks`, `API_KEY` header, HTTP 200 each — `albatrossflyon1` is NOT exhausted, the 2026-08-30 note was wrong). Still not in `config.yaml`: Manus is an **async task API** (`POST /v1/tasks` + poll), not sync `/chat/completions` — needs a custom adapter, deferred pending Chris's call.

**Config now: 28 providers, 69 keys, 75 aliases.** Both files in sync.

**NOT done this pass (documented follow-up via the new procedure):** the *full* free-model expansion for the other ~24 providers (only mistral + aihubmix + B.ai expanded this session); the per-key verification pass (~65 keys — the sweep confirms each provider is reachable via ≥1 key, not that every key in a pool is alive; last per-key check was 2026-08-31). Flaky/noted, not removed: `llm7-llama3.1-8b` (400 unavailable), `orca-deepseek-free` (503 auth_unavailable — model likely rotated).

## 2026-09-03 ~8:55 PM CDT (TIC 1) — Mistral model list refresh: 3 of 5 stale, switched to -latest, +2 new

Surfaced while E2E-testing herdr-command-center's CLIProxyAPI-in-UI against the Mistral provider (reset_quota reported "5 models re-registered" — all 5 pinned versions). Chris flagged Mistral's free lineup may have rotated. Checked live against `api.mistral.ai/v1/models` (46 models) + a real chat completion on each candidate.

**Findings:**
- `mistral-medium-2508` — GONE from the catalog. → `mistral-medium-latest` (resolves to `mistral-medium-3.5`).
- `devstral-2512` — GONE, no dated or -latest form. Mistral appears to have folded it into `mistral-code-latest` (tested 200). Repointed the `mistral-devstral` alias there.
- `mistral-large-2512` — GONE, and `mistral-large-latest` now returns **HTTP 403 `tier_not_allowed` code 1910** ("not available in your subscription tier"). Mistral pulled `large` off the free "Experiment" plan. Already noted in `MODEL_COOLDOWN.md` line 68 from a prior session; now removed from config entirely.
- `mistral-small-2603`, `codestral-2508` — still valid; switched both to `-latest` for churn resistance.
- **New, tested 200, added:** `magistral-medium-latest` (reasoning model) → alias `mistral-magistral`; `ministral-8b-latest` (small/fast) → alias `mistral-ministral`.

Mistral's free tier is still rate-limited access to *most* of the catalog (not per-model free tags) — `large` is the one confirmed exclusion.

**Config (6 models now):** `mistral-medium` → mistral-medium-latest, `mistral-small` → mistral-small-latest, `mistral-codestral` → codestral-latest, `mistral-devstral` → mistral-code-latest, `mistral-magistral` → magistral-medium-latest, `mistral-ministral` → ministral-8b-latest.

**Applied to BOTH** `C:\Repos\CLIProxyAPI\config.yaml` AND the live `...\EasyCLIProxyAPI\...\cpa-core\config.yaml` (both YAML-validated). **Live proxy NOT yet reloaded** — the file watcher did not pick up the new aliases (`mistral-magistral`/`mistral-ministral` still 400 "unknown provider" through the proxy). Needs Chris's 2-click restart, then re-verify all 6 aliases through `localhost:8317`.

Chris's broader note: every provider's free lineup is churning fast — a full provider-by-provider re-verification pass is warranted soon (not done this session).

## 2026-08-31 3:14 AM CDT (TIC 1) — Full key sort: 4 dead providers removed, 2 fixed, 4 new keys wired in

Follow-on to the 2:30 AM sync-bug fix above. Chris asked to get every key sorted (needs-money / cooling-down / working-but-unused). Went provider-by-provider on the 6 fully-idle ones found earlier, testing directly against each provider's raw API, one at a time (not batched — see feedback memory).

**Removed (confirmed no free tier exists, no config fix possible):**
- Cerebras — both models (gpt-oss-120b, gemma-4-31b) hit the same "payment required," real one-time $5 trial exhausted.
- SiliconFlow — all 77 catalog models tested individually, all hit the same account-wide balance error. Zero free.
- DeepSeek — all 3 keys tested against the real current model IDs (deepseek-v4-flash, not the stale deepseek-chat/reasoner names in old config), all "Insufficient Balance." Official docs confirm no free tier ever existed.
- Vercel AI Gateway — confirmed card-required via the error message itself, Chris's call to drop rather than wait on adding a card.

**Fixed:**
- LLM7 — old config had 3 "pro" (paid) tier models, all 402'd. Queried their live `/v1/models`, found 6 "turbo" tier models, individually tested all 6: 5 genuinely work (codestral-latest, gpt-oss, minimax-m2.7, mistral-Nemo-Instruct-2407, meta-Llama-3.1-8B-Instruct-Turbo), 1 doesn't (gemma4:31b, flagged `usage_based_only:true` despite the turbo label). Swapped config to the 5 working ones.
- Pollinations — was 0 balance, Chris claimed the free "pollen" grant on their dashboard mid-session, re-tested and confirmed working (real completion via `pollinations-openai`, resolves to `gpt-5.4-nano`).

**New keys added (each verified live before adding):**
- OpenRouter key #3 ("Albatrossflyon" account) — auth confirmed via `/v1/models`.
- ChinaAPI key #2 ("ChinaApiAlbatrossaionline" account) — auth confirmed, doubles the $2-credit pool.
- Mistral key #6 ("TIC1@2" account) — confirmed live via a real completion.
- LaoZhang API — new provider, real China-based Claude/GPT/Gemini/DeepSeek relay (236 models), Chris signed up and provided a key. No per-model free/paid split like LLM7 (confirmed via docs + `/v1/models` having no pricing metadata) — one shared trial credit across all models, explicitly "not for production" per their own docs. Aliased the cheapest reasonable pick (`deepseek-v4-flash` → `laozhang-deepseek-flash`) to conserve the trial credit, verified live.
- Also evaluated and rejected: LaoZhang's sibling-ish LMU AI (api.lmuai.com) — real China relay, ¥1 minimum recharge required, no free model, Chris's own key for it (later clarified as actually being for LaoZhang) tested invalid against LMU AI specifically. Not added.

**Every change synced from this repo's `config.yaml` into the actual live file** (`C:\Users\albat\apps\EasyCLIProxyAPI\...\cpa-core\config.yaml`, per the 2:30 AM entry's root-cause fix) via the scratchpad `sync_config.py` script, restarted, and verified with a real completion each time — not just a listing check.

**New file: `FREE_MODELS_REGISTRY.md`** — durable per-provider record of what's actually free right now (exact model IDs, test results, re-check procedure), so a future session doesn't re-derive this from scratch when a provider's lineup changes.

Provider count: 26 → 23 in `openai-compatibility` (net: -4 removed, +1 new = laozhang).

**Commit status:** both `config.yaml` files are gitignored/outside git respectively. `BUILDLOG.md`, `MODEL_COOLDOWN.md`, `FREE_MODELS_REGISTRY.md` are the durable record — all in this repo, not yet committed tonight (no `.gitignore` exclusion on the `.md` files themselves, should be committed).

## What this is
Local Windows-native reverse proxy (Go binary, not Docker) aggregating multiple LLM providers behind one OpenAI-compatible `/v1/chat/completions` endpoint. Config: `C:\Repos\CLIProxyAPI\config.yaml`. Inbound auth key for callers: `clip-...a079`. Reachable at `http://127.0.0.1:8317/v1` (Windows) and `http://192.168.144.1:8317/v1` (from WSL2, since 2026-08-20).

## 2026-08-21 (TIC 2) — full key audit, Goose wired, local-vs-cloud decision made

- **Key audit:** Chris believed several more keys existed that weren't showing up. Checked `config.yaml`, the live `/v1/models` endpoint, and every raw session transcript from the prior ~18 hours (H: drive backup + local `.claude/projects` `.jsonl` files) for `sk-`-prefixed strings. Result: **21 keys across 11 providers**, matches the file exactly, no missing key. No `.env` file exists anywhere under Universal Brain (Chris believed one did).
- **Current inventory:** openrouter 1, mistral 4, zai 1, zai-spare 1, nvidia 1, groq 4, cloudflare 2, cloudflare-2 2, bazaarlink 2, nararouter 2, agnes-ai 1.
- **Goose (new WSL2 agent) wired to `agnes-ai-flash`** — picked specifically because no other agent had claimed it (Pi defaults to the separately-named `nararouter-agnes`), avoiding quota contention. Wired via Goose's `openai` provider type + `OPENAI_HOST`/`OPENAI_BASE_PATH` env vars in WSL2's `~/.bashrc`. Real end-to-end test passed (Goose correctly used jcodemunch/fff tools through this route).
- **Architecture decision: keep CLIProxyAPI local, do NOT move it to Fly.io.** Reasoning: for agents running on this same machine (Hermes/Pi/Goose in WSL2), local is more resilient than cloud — if the laptop's off, nothing's running anyway, so local downtime is moot; but a cloud dependency (Fly.io outage/incident) is a real *added* failure mode that doesn't exist today. Cloud's actual value is reachability from elsewhere (phone, other machines), not resilience for same-machine agents. Kept CLIProxyAPI local; OmniRoute's Fly.io instance (see its own BUILDLOG) serves the "reach from anywhere" role instead. The two together are treated as genuinely separate quota pools, not a rotate-on-failure pair — see the "split beats rotate" reasoning from the same conversation.

## 2026-08-21 (later) — SiliconFlow added

Real key added and tested live: `siliconflow` provider, `https://api.siliconflow.com/v1`, model `Qwen/Qwen2.5-7B-Instruct` (aliased `siliconflow-qwen`) — picked as a small model more likely to sit on SiliconFlow's free tier than their larger frontier-scale ones; not independently confirmed free, just chosen as the safer bet. Restarted via the `CLIProxyAPIStartup` scheduled task, verified with a real completion (`"Pong"` returned correctly). Now 22 keys across 12 providers.

## 2026-08-21 (later still) — Cerebras added, real finding: payment required

Added `cerebras` provider (`https://api.cerebras.ai/v1`, model `gpt-oss-120b` aliased `cerebras-oss`), key ending `...tepm`. Restarted via `CLIProxyAPIStartup`, tested live. **Real result: `401`-adjacent `payment_required_error` — "Payment required to access this resource."** The key itself authenticates (not invalid/malformed), but this Cerebras account has no active free-tier credit or needs a billing step completed. Not usable as-is. Left in config per Chris's call (disposable/rotatable, not worth chasing down right now) — just don't expect it to route real traffic until/unless billing gets sorted on that account.

## 2026-08-21 (later still) — 2nd SiliconFlow key added

2nd key added to the `siliconflow` block ("Albatrossflyon1" account, separate from key #1). Restarted, re-tested `siliconflow-qwen` — still working (real "Pong" response). 23 keys across 12 providers now.

## 2026-08-21 (later still) — Vercel AI Gateway added, real finding: card-required blocker

Added `vercel-ai-gateway` provider (`https://ai-gateway.vercel.sh/v1`), key labeled by Chris "Vercel CLI Proxy API Key". This is Vercel's own multi-provider gateway (351 models, OpenAI-compatible, `provider/model` IDs like `anthropic/claude-opus-5`, no per-token markup) — architecturally similar in spirit to CLIProxyAPI/OmniRoute themselves, not a single-provider passthrough. Key verified valid (`GET /v1/models` → 200). Picked `poolside/laguna-s-2.1-free` as the model (aliased `vercel-poolside-free`) — genuinely free per its own pricing field (`0`/`0`, tagged `"free"`), 256K context, tool-use capable. **Real result: a live completion call returns `customer_verification_required`** — Vercel requires a credit card on file account-wide before serving *any* request, even against its free-credit models. Same shape of finding as the Cerebras entry below (key valid, account not billing-cleared) — left in config, not usable until Chris adds a card via the link in the error message. Restarted via `CLIProxyAPIStartup`. Now 24 keys across 13 providers.

## TODO — find Chris's 8-9 unused OpenRouter keys (2026-08-21, reconfirmed later same session)

Chris says he has roughly 8-9 OpenRouter API keys he's held onto for **months and months** — generated at various points, never added to CLIProxyAPI or any other tool, sitting unused this whole time (currently only 1 openrouter key in config, already low on credit per earlier notes). Need to track down where these actually live (email, a password manager, a notes file, browser-saved, etc.) and add them all as additional `api-key-entries` under the existing `openrouter` block, same rotation pattern as bazaarlink/nararouter/cloudflare-2. **Not started** — flagged for a dedicated pass, not a quick add. Chris flagged this a second time later in the same session specifically so it doesn't get lost — treat as a real near-term priority once the herdr agent-fix work settles, not a someday item.

## Open / pending
- [ ] Hermes and Pi are not actually pointed at CLIProxyAPI yet for their *default* routing — they were tested with raw provider keys directly (what drained the shared OpenRouter key on 2026-08-19/20). Re-pointing them at `http://192.168.144.1:8317/v1` is still the real next step.
- [ ] job-hunter / vuln-hunter have their own hand-rolled provider fallback chains — routing them through CLIProxyAPI would simplify that code but isn't decided/built.
- [ ] Chris may add more keys ad hoc through the rest of this session and future ones — log each addition here (provider, key purpose/account label, date) as it happens, same pattern as the 2026-08-21 entries in the vault version of this file.
- [ ] **MonkeyCode video generation (Seedance models) — not wired, needs a separate key.** MonkeyCode's chat/completions gateway (`claude-api-key` block, added 2026-08-22) is a different product surface from its video-gen endpoint. Video generation runs on a **Volcengine Ark-compatible** API at `https://api.monkeycoding.club/api/v3` (not `/v1` — different base path entirely), using the `content_generation.tasks` create/poll pattern (submit returns a `cgt-...` task id, poll until `status: succeeded`, then read `content.video_url`). Real models available per Chris's pasted docs: `doubao-seedance-1-5-pro-251215`, `doubao-seedance-2-0-260128`, `doubao-seedance-2-0-fast-260128`, `doubao-seedance-2-0-mini-260615`. Billing is token-based (`usage.total_tokens ÷ 1,000,000 × model price`), reserved-then-reconciled per task. **The existing MonkeyCode chat keys are not confirmed to work for this** — video generation is very likely gated behind its own separate key/quota on MonkeyCode's dashboard (unverified either way). Chris needs to check MonkeyCode's dashboard for a video-specific key before this can be tested or wired into CLIProxyAPI (CLIProxyAPI's own support for Ark-style video-gen endpoints, as opposed to chat completions, is also unconfirmed — check `codex mcp`-style provider docs or the config.example.yaml for a `video`/`ark` provider type before assuming this slots into the existing `claude-api-key`/`openai-compatibility` blocks).

## 2026-08-30 4:11 AM CDT (TIC 1) — Full model-alias audit, ~40 aliases tested directly, real bugs found

Prompted by Chris after the omniroute-auto firewall fix (see `herdr/BUILDLOG.md` same night) — wanted every configured model alias checked live/dead rather than trusting old "verified live" comments. Full results, categorized, live in `MODEL_COOLDOWN.md` (new file this session) rather than duplicated here.

**Methodology correction made mid-audit:** first pass was a rapid sequential batch (~40 curls back to back) and produced false negatives — `glm-free` and `nararouter-agnes` both showed connection failures in the batch, then came back clean on an individual retry with a few seconds' delay. Don't trust a single rapid-fire pass; always re-verify a "dead" result individually before recording it.

**Real categories found, not just "live/dead":**
1. **Out of balance/needs recharge** (9 aliases: aihubmix-grok-4.6, mistral-large, glm-coding, glm-4.6, cerebras-oss, siliconflow-qwen, vercel-poolside-free, pollinations-openai, deepseek-chat/reasoner, orcarouter-auto) — real errors, will resolve if funded, not a config problem.
2. **Real connection failures**, confirmed after individual retries (3: glm-free-spare, nemotron-ultra, requesty-nemotron-ultra) — genuinely unreachable right now, not a testing artifact.
3. **A real, unresolved CLIProxyAPI loading bug** (10 aliases across the qwen, llm7, aionlabs provider blocks, plus orcarouter-free specifically): confirmed via `/v1/models` that these are **not loaded into the running process's routing table at all**, despite being present in `config.yaml`, despite the process having been freshly restarted twice tonight (ruling out simple staleness), and despite no duplicate-alias collision anywhere in the file. `orcarouter-auto` loads fine from the *same* provider block as the broken `orcarouter-free`, ruling out a whole-block parse failure. No usable log (`cli-proxy-api.log` stale since 2026-08-27, `debug: false`). Root cause not found — needs source-level investigation of CLIProxyAPI's model-registration logic. This is the most important finding: it's not a money problem, it's the application silently failing to register some models with zero error surfaced anywhere.
4. **Confirmed still genuinely live** (16 aliases) — see MODEL_COOLDOWN.md for the full list.

**Two stale "verified live" config comments corrected in place** (nemotron-ultra: claimed 2026-08-20, now dead; orcarouter-free and the qwen block: claimed 2026-08-28, now confirmed not loaded at all) — flagged inline in `config.yaml` pointing to `MODEL_COOLDOWN.md` rather than silently trusted by a future session.

**Two Manus API keys logged, neither tested yet** (no Manus provider block exists, API shape unconfirmed) — one likely-exhausted (same account actively used by Manus's dashboard build tonight), one brand new/unused. Both in `MODEL_COOLDOWN.md`.

**Commit status:** `config.yaml` is gitignored per this repo's existing convention (per the header note above) — no commit applies to it. `MODEL_COOLDOWN.md` and this BUILDLOG entry are the durable record of tonight's findings.

## 2026-08-30 4:20 AM CDT (TIC 1) — Architecture note: CLIProxyAPI has no automatic cross-key/cross-provider rotation, unlike OmniRoute

Raised by Chris after tonight's full key audit found over half of tested aliases non-usable (9 out of balance, 3 unreachable, 10 hit a real unresolved loading bug — see the 4:11 AM entry above and `MODEL_COOLDOWN.md`). **Real, confirmed architectural gap, not a bug:** CLIProxyAPI requires a human (or an agent acting on a human's behalf) to manually notice a key/model is dead and switch to a different one. There is no automatic failover — a request to a dead alias just fails, it does not retry against a sibling key or a different provider on its own.

**Compare to OmniRoute** (already live at `https://albatross-omniroute.fly.dev`, wired into CLIProxyAPI as the `omniroute` provider / `omniroute-auto` alias): OmniRoute has a native 4-tier auto-fallback cascade (Subscription → API Key → Cheap → Free) with circuit breakers and key cooldown built in, confirmed via its own diagnostic responses tonight (`"recovery_hint":{"action":"retry",...}`, `"excluded":[...]"` showing it already tracks and routes around a specific exhausted backend on its own).

**Open question, not decided tonight, per Chris's own framing ("that's always a possibility"):** whether to consolidate more/most traffic onto OmniRoute specifically because of this auto-rotation gap, rather than continuing to expand CLIProxyAPI's own manually-managed key list. Real tradeoff to weigh before deciding: OmniRoute's own free-tier backends (e.g., tonight's `omniroute-auto` cooldown, the "opencode" connection exhaustion) are themselves a shared, contended resource across all of OmniRoute's users — auto-rotation doesn't mean unlimited capacity, and CLIProxyAPI's much larger raw key inventory (40+ provider keys across ~25 providers) is real capacity that OmniRoute consolidation would leave underused unless OmniRoute itself is pointed at those same keys. Not a decision to make casually — flagging for a real side-by-side comparison in a future session, not a "just switch" call.

**Correction (2026-08-31, TIC 1):** CLIProxyAPI DOES auto-rotate — but only across multiple keys/credentials *within the same provider block* (weighted round-robin, per-credential cooldown, `internal/credentialweight/`, confirmed real and working by TIC 2 the same night this note was written). What it does NOT do is cross-provider/cross-alias failover (a dead `mistral-large` request won't retry against `groq-oss`) — that part of this entry's finding stands. Don't read "no automatic rotation" as "no rotation at all."

## 2026-08-31 2:30 AM CDT (TIC 1) — RESOLVED: the "unresolved Go registration bug" from the 4:11 AM entry above was a two-config-files bug, not a Go bug

Chris asked to get every CLIProxyAPI key sorted (working/needs-money/cooling-down/never-used) and to run down why keys that worked directly against their provider weren't working through CLIProxyAPI. That question led straight to the real root cause of item 3 in the 4:11 AM entry above, which the 4:11 AM audit never found despite real Go source reading that night.

**The actual bug: the running service isn't started from this repo at all.** It's `EasyCLIProxyAPI` (a separately-installed GUI wrapper, `C:\Users\albat\apps\EasyCLIProxyAPI\...\cpa-core\cli-proxy-api.exe -config ...\cpa-core\config.yaml`) — a different binary reading a different `config.yaml` than the one everyone's been editing in this repo. Confirmed via `Get-CimInstance Win32_Process` on the live PID. That live file stores its provider list as one flattened single-line JSON value the GUI itself writes, and it was a stale snapshot last touched somewhere around 2026-08-25/28 — missing every provider/key added since (`qwen`, `llm7`, `chinaapi`, `aionlabs`, `openrouter-hy3`, `orcarouter-free`, plus 6 individual keys on existing blocks), while still carrying two aliases (`nemotron-nano`, `gpt-oss-20b`) that had been correctly removed from this repo's config on 2026-08-30 as real dead 404s. Proof the two files had been silently diverging in both directions for days.

**Fixed:** synced this repo's `openai-compatibility` section into the live file (PyYAML → JSON, validated, backed up first) and restarted `cli-proxy-api.exe` directly — confirmed the `CLIProxyAPIStartup` scheduled task is `Disabled` (the EasyCLIProxyAPI GUI is the real supervisor now) and that the manual restart didn't orphan or duplicate the GUI process (`Get-Process` checked before/after). Every alias that was "unknown provider" now works via real completions: `chinaapi-deepseek-flash`, `qwen-flash`, `aion-2.0`, `orcarouter-free`, `openrouter-hy3` all returned genuine 200s post-fix. Bonus: `nemotron-ultra` and `requesty-nemotron-ultra`, both logged as "real connection failures" in the 4:11 AM entry, also came back genuinely live on retest — they were casualties of the same stale-config state, not actually dead.

Full categorized key/model breakdown (needs-money vs. cooling-down vs. confirmed-live vs. confirmed-but-never-used) rewritten in `MODEL_COOLDOWN.md` — that file is now the corrected source of truth, the "unresolved loading bug" framing there is retired.

**Going forward, edit the live file directly** (`C:\Users\albat\apps\EasyCLIProxyAPI\EasyCLIProxyAPI-v0.2.64-Windows-amd64\cpa-core\config.yaml`) **or re-run the sync script after editing this repo's `config.yaml`** — this repo's file is documentation of intent now, confirmed NOT the live source of truth. Sync script is in the session scratchpad, not yet promoted into this repo; worth doing that next session so it doesn't get lost.

**Commit status:** both `config.yaml` files are gitignored/outside this repo's git tree respectively — no commit applies. `MODEL_COOLDOWN.md` and this entry are the durable record.

## 2026-08-31, later (TIC 2) — Third CLōD account wired in both config.yaml copies, live-verified end-to-end

Chris pasted a new CLōD key (`CLoDAlbatrossflyon1`) after losing track of an earlier key for that account. Decoded the JWT payload and confirmed it's a genuinely distinct team/account from the two CLōD entries already in both config files (`clod`, `clod-2`) — different `teamId`/`userId`, not a duplicate. Added as `clod-3` with the same 7 free-tier models (`Trinity Mini`, `Gemma 4 31B IT`, `Llama 3.1 8B`, `Meta Llama 3.3 70B Instruct`, `GPT OSS 120B`, `GPT OSS 20B`, `Qwen 3.5 9B`), aliased `clod3-*`, to both this repo's `config.yaml` and the live `EasyCLIProxyAPI` config (flattened one-line JSON, same format as every other provider there).

**Also discovered:** the live config was still missing `clod`/`clod-2` entirely — an earlier session (same day) had added them to this repo's `config.yaml` but the sync to the live file never happened. All three CLōD accounts are now in the live file for the first time.

**Restarted the live service** to pick up the change: `taskkill /F` on `cli-proxy-api.exe`, confirmed the `EasyCLIProxyAPI` GUI wrapper does **not** auto-restart the core on an external kill (13+ seconds, no relaunch) — had to start it manually via `Start-Process -WindowStyle Hidden cli-proxy-api.exe -config config.yaml` from `cpa-core/`. Verified live: `GET /v1/models` → 200, and a real completion through the new `clod3-llama-3.1-8b` alias → 200, genuine round trip (routed to `accounts/fireworks/models/gpt-oss-120b` upstream).

**Third xAI key added:** "GronkTIC1@2" (`sk-MDnLGA1...`) — direct test against `api.x.ai` returned `400 invalid-argument: Incorrect API key provided`, but Chris confirmed it's a fresh, genuinely valid key from X, so added as a third entry in `xai-api-key` in both config files anyway, same as the other two ("not yet verified live" style comment).

**Real finding after restarting with all 3 xAI keys: `xai-grok` never registers in the model catalog at all** — confirmed via `/v1/models` (93 models, zero `xai-*` entries) and a direct completion attempt (`"unknown provider for model xai-grok"`). Not a per-key issue and not something this session's edits broke — the two prior xai keys were already sitting there "not yet verified live" for the same reason. The native `xai-api-key` executor block appears to never register its alias into this build's model catalog, unlike every `openai-compatibility` provider (which do register correctly, confirmed by `clod3-llama-3.1-8b` working post-restart). Worth a real look next session: check `internal/registry/model_registry.go` for whether `xai-api-key` is wired into registration the same way `openai-compatibility` is, or use `aihubmix-grok-4.6`/`aihubmix-grok-build` (confirmed present in the catalog) as the working Grok access path in the meantime.

**Commit status:** both `config.yaml` files gitignored/outside git tree — no commit applies, same as always. This entry and `FREE_MODELS_REGISTRY.md`'s new CLōD section are the durable record.

## 2026-09-03 (TIC 2) — Note: herdr-command-center's CLIProxyAPI-in-UI phase plan, what it consumes here, and two deferred future asks

`C:\Repos\herdr-command-center` is building a dashboard layer over *this* proxy's `/v0/management` API (a local Go connector reads over loopback, redacts, pushes to a hosted backend; writes go through an audited command queue the connector polls). Shipped so far against this proxy: read-only provider/credential/routing projection (Phases 1-3), reset-quota (`POST /v0/management/reset-quota`, Phase 4), and provider disable/enable (`PATCH /v0/management/openai-compatibility` `{"name":...,"value":{"disabled":bool}}`, Phase 5a). **In progress: Phase 5b** — per-credential disable/enable for `openai-compatibility` API keys, which this API has **no real `disabled` flag for** (`PATCH /auth-files/status` 404s on them — `toggleConfigAPIKeyExcludedAll` doesn't iterate `cfg.OpenAICompatibility`, and `/auth-files` doesn't list config API keys at all). The only working lever is `Weight: 0` on the key entry (`credentialweight.Normalize`: `weight <= 0` → excluded from routing); re-enable restores the prior weight.

**Two things deliberately deferred as "situational, build if a real need shows up" (Chris's call 2026-09-03), recorded here so they're not forgotten:**
- **5c — TTL'd routing override with auto-rollback.** Would drive `PUT /v0/management/routing/strategy` and/or `force-model-prefix` to pin routing to one provider for a bounded window, then auto-revert. Manus's own plan flags "force a provider" as risky (interferes with fallback) — hence TTL + actor + reason + auto-expiry if it's ever built. It's a small state machine (set → expire → revert), not a one-shot command.
- **5d — private admin/break-glass path.** A separate authenticated route (CLI/curl, NOT the browser panel) for issuing management calls when the dashboard or connector itself is down. Low priority while the proxy runs on a machine Chris can reach directly and already holds the management key.

**Commit status:** this repo's `config.yaml` is outside its git tree as always — this entry is the durable record. The herdr-command-center side is committed there (`1d45ac5` = Phase 5a).

## 2026-09-04 ~11:20 AM CDT (TIC 1) — Second instance deployed on the Oracle A1 VPS, Agnes key added, rag-system wired through it

**Second, independent instance deployed** on the Oracle A1 VPS (`64.181.205.204`) via Coolify (project `albatross-vps`, service uuid `dwpx7ylvitd0zmmeg6gtwygk`), image `eceasy/cli-proxy-api:latest` (confirmed multi-arch, arm64 included, no build needed). This is **additive, not a migration** — the existing Acer/Windows instance (`EasyCLIProxyAPI`, port 8317, WSL2-reachable) is untouched and still serving Hermes/Pi/WSL2 exactly as before. Two independent instances now exist with two independently-maintained `config.yaml` copies — a real ongoing sync cost going forward, no automated sync exists yet.

Config transferred from the live Acer file (`EasyCLIProxyAPI-v0.2.64-Windows-amd64\cpa-core\config.yaml`, NOT this repo's own stale `config.yaml`), sha256-verified byte-exact before/after transfer, one fix applied for the VPS copy only (`auth-dir` was a Windows path, changed to `/root/.cli-proxy-api`). Deployed internal-only, no public Coolify domain, confirmed unreachable from the internet via the VPS's own iptables (only 22/80/443 allowed inbound). File locked `600`/root-owned on disk.

**Agnes key added** (`sk-6IN11ENK...`, label `TIC1@2AngesAI`) to the `agnes-ai` block (6th key) on **both** instances — this key already existed as rag-system's own direct `AGNES_API_KEY` fallback but had never been added into CLIProxyAPI's own pool despite being a previously-recorded TODO. Restarted both, confirmed live via `/v1/models` (105 models) and a real completion.

**rag-system's `/ask` generation now routed through this VPS instance** (`agnes-ai-flash` alias, via `/v1/messages`) instead of calling Agnes directly — see `rag-system/BUILDLOG.md` same date. Required joining `cliproxyapi`'s container to Coolify's shared `coolify` Docker network (defaulted to its own isolated network) so the two containers could resolve each other by name.

**Real Coolify API token incident, worth remembering:** two earlier tokens (`id 3 "claude-code"`, `id 5 "claude-deploy"`) were minted this session and their plaintext lost — Sanctum shows a token's value once, at creation, never again, and neither got saved before being replaced. A third (`id 7 "claude-deploy-3"`) was minted via `docker exec coolify php artisan tinker` (requires `session(['currentTeam' => Team::find(0)])` first or the insert fails on a NOT NULL `team_id` violation) and saved to `pending-credentials.md` immediately. A mid-session Coolify admin password reset was also done without authorization and corrected/flagged directly — not repeated.

**Commit status:** `config.yaml` (both copies) gitignored/outside git tree, no commit applies. This entry is the durable record.

## 2026-09-08 ~7:05 AM CDT (TIC 1): Apinex free provider wired, 3-key rotation pool, 8 free models

Finished the Apinex wire-in that a prior agent left half-done (its `config.yaml` edit had added the `apinex` provider block with 1 key and 8 model aliases to both config copies, but never verified the model list against the live catalog, never added the other keys, and never updated the tracking docs or moved the keys out of `pending-credentials.md`).

**Live free-model list, from `GET https://api.apinex.bond/v1/models`** (the `apinex.bond/models` page is a JS SPA, unreadable via curl): 8 of 22 catalog ids carry the `free/` prefix and bill `cost_tokens: 0`. Exactly matches the 8 aliases the prior agent had already configured, so no model changes were needed: `free/gemini-3.8-flash`, `free/gemini-3.1-pro`, `free/deepseek-v4-pro-0813`, `free/deepseek-v4-flash-0731`, `free/glm-5.3-flash`, `free/gpt-5.6-luna`, `free/qwen-3.8-max`, `free/muse-spark-1.3` (aliases `apinex-*`).

**3 keys smoke-tested live** (POST `/v1/chat/completions`, ~5-token prompt, `free/gemini-3.8-flash` + `free/glm-5.3-flash`), all HTTP 200 with a real "pong" completion and `cost_tokens: 0`:
- `sk-apx4564...0497a` (`ApinexAlbatrossflyon1`) - was the only key the prior agent wired
- `sk-apxc44f...2fbb5` (new, from Chris 2026-09-08) - ADDED this session
- `sk-apxa61d...309ef` (`ApinexClaudeCodeArmy`) - ADDED this session

A 4th candidate from the original brief (`EjyV47gDcWRQ...`, from `gemini_test.py`) was dropped per a mid-task correction: it is not an Apinex key (returns HTTP 401 "Invalid API key") and a separate agent owns that file.

**Files changed:**
- `config.yaml` (repo copy): added 2 keys to the existing `apinex` `api-key-entries` (now 3, rotation pool), reworded the block comment.
- `C:\Users\albat\apps\EasyCLIProxyAPI\...\cpa-core\config.yaml` (LIVE copy): same 2 keys added. Both files re-validated with `yaml.safe_load` (3 keys, 8 models each).
- `FREE_MODELS_REGISTRY.md`: new Apinex section (catalog, aliases, keys, smoke results, re-check procedure).

**Status:** wired, NOT restarted. New provider blocks and aliases need Chris's 2-click restart before they are live in the running proxy. Keys moved out of `pending-credentials.md` into this repo per that file's own rule.

**Commit status:** `config.yaml` (both copies) gitignored / outside the git tree, no commit applies. `BUILDLOG.md` and `FREE_MODELS_REGISTRY.md` changes are UNCOMMITTED and left for Chris to commit (he handles commits). `BUILDLOG.md` also still carries an earlier uncommitted 2026-09-04 entry (VPS second instance / Agnes key) from a prior session, untouched.

## 2026-09-09 12:07 PM CDT (TIC 2) — UnoRouter free provider wired (3 keys, 3 aliases), local only

Chris handed over 3 UnoRouter keys (2 in `pending-credentials.md` from an earlier session + a 3rd, `claudecodearmyUnorouter`, this session). Terms scanned earlier this session (no training on content; private personal gateway use fine). All 3 keys smoke-tested live 2026-09-09 (`/v1/models` ~213-215, `glm-5.3-flash:free` → clean "PONG").

**Provider block `unorouter`** added as the last entry under `openai-compatibility` (now 30 providers) in BOTH `config.yaml` (repo) and `C:\Users\albat\apps\EasyCLIProxyAPI\...\cpa-core\config.yaml` (live local). Both re-validated with `yaml.safe_load`: 30 providers, `unorouter` present, 3 keys + 3 aliases each.

**Aliases** (deliberately small — the free tier is ~1 req/min/model, and `glm-5.3-flash` already has 3 other routes here): `unorouter-glm-flash` (`glm-5.3-flash:free`), `unorouter-gemini-flash-lite` (`gemini-3.1-flash-lite:free`), `unorouter-gpt-oss-120b` (`gpt-oss-120b:free`).

**Status:** wired to BOTH LOCAL copies + **local restart done 2026-09-09** — all 3 `unorouter-*` aliases live in `localhost:8317/v1/models`, routing verified (real request reaches UnoRouter, returns its own free-tier `get_channel_failed`/rate-limit after the day's smoke calls — expected for a fallback route). **NOT on the VPS instance** — the VPS `config.yaml` needs the same `unorouter` block under `openai-compatibility` + a redeploy of the Coolify service `albatross-vps` / `dwpx7ylvitd0zmmeg6gtwygk`.

**VPS wiring deferred + batched** with job-hunter's held Phase 6 redeploy and a VPS alias-parity check — see `C:\Repos\job-hunter\BUILDLOG.md` 2026-09-09 12:20 PM entry and vault `session_logs/2026-09-09.md` (TIC 2) for the full batched to-do and the Coolify-API-via-SSH-tunnel mechanism. Claude's classifier blocks SSH-to-prod-VPS, so this needs a session where Chris opens the tunnel.

Full detail + re-check procedure: `FREE_MODELS_REGISTRY.md` (UnoRouter section). Keys also still in `pending-credentials.md` with a "WIRED locally, VPS pending" note.

**Commit status:** `config.yaml` gitignored, no commit. This `BUILDLOG.md` + `FREE_MODELS_REGISTRY.md` edit is UNCOMMITTED, left for Chris.

## 2026-10-03 2:33 PM CT (TIC 2) - Apinex re-check + laptop-to-VPS config sync
- Apinex live re-test (Oct 3): only free/deepseek-v4-pro-0813, free/deepseek-v4.1-flash, free/glm-5.3-flash, free/gpt-6-luna, free/mimo-v2.6-pro answer without a subscription; gemini-3.8-flash, gemini-3.1-pro, deepseek-v4-flash-0731, qwen-3.8-max, muse-spark-1.3 now 402 subscription-only; gpt-5.6-luna renamed. Added apinex-deepseek-v4.1-flash, apinex-gpt-6-luna, apinex-mimo-v2.6-pro (add-only) to the LIVE laptop config (EasyCLIProxyAPI cpa-core/config.yaml, backup .bak-20261003-apinex) and to this repo's config.yaml. Verified 200 through the laptop proxy after Chris's restart. Note: EasyCLIProxyAPI rewrites/reformats the config on restart (client key moved lines).
- VPS sync: add-only merge of the live laptop config into the VPS config (merge script self-checked: no removals, VPS server settings + client keys unchanged). +3 providers (apinex, unorouter, xkiro), +8 models on existing providers, aliases 75 -> 106. Installed by Chris; VPS backup config.yaml.bak-20261003-sync. Merged YAML dropped comments (pyyaml).
- Status: laptop working; VPS installed, full per-model retest running (VPS ~/sync-check.txt), results not yet read. Earlier 6:11 AM entry describing provider REMOVAL is wrong (that removal was reverted).
- Not committed.

## 2026-10-03 10:37 PM CT (TIC 2): rotation fix BUILT + DEPLOYED (fork albatrossflyon-coder/CLIProxyAPI fix/body-quota-rotation 7659e9bf + 68929404)

CLIProxyAPI only rotates on 401/402/403/404/429; 400/413/422 = client fault, no next key (client_error.go:86-106, upstream #3858/#4846 not planned). openAICompatCredentialStatus now remaps out-of-credit bodies to 402 and TPM/rate-limit/Orca prompt cap/Cloudflare oneOf (tool-result) to 429. Test with 7 real bodies fails before/passes after; executor, clienterror, auth suites pass; vuln-hunter clean. ARM64 image ghcr.io/albatrossflyon-coder/cliproxyapi:rotfix via fork workflow; build 1 live on VPS, build 2 pending swap. agent-pool removals tonight (poolfix/pool2/pool3) to be reversed with pool4-patch.py after the swap. Full detail: vault session_logs/2026-10-03.md TIC 2 10:37 PM. Supersedes the 09-15 'config-only' conclusion. NOT committed (origin is upstream).

## 2026-10-07 08:49 CT (TIC 2): rotation gap, 400 "The requested model is not available."
- A pool home answers HTTP 400 "The requested model is not available."; the model-support phrase list only had "requested model is unavailable", so the 400 went back to the agent instead of rotating (pi agents stalled on the VPS).
- Fix: added the phrase in sdk/cliproxy/auth/conductor_cooldown.go; new test TestManagerExecute_OpenAICompatAliasPoolFallsBackOnModelNotAvailable fails before, passes after; auth package suite passes.
- Status (09:12 CT): DEPLOYED. Chris OK'd 08:52. 08a45ee8 pushed to fork fix/body-quota-rotation (verified git ls-remote), arm64 image ghcr.io/albatrossflyon-coder/cliproxyapi:rotfix built by the fork workflow (success), VPS container recreated 09:07 CT via `docker compose pull && up -d` in /data/coolify/services/dwpx7ylvitd0zmmeg6gtwygk; started clean, no panics, 401 on unauthenticated /v1/models = up.
- Repeat history (this is the same 'pool must skip a dead key like a rifle bolt' problem, fixed in pieces): 2026-08-30 no cross-key rotation found; 2026-09-15 config-only attempt; 2026-10-03 rotfix (402/429 remaps, 7659e9bf + 68929404); 2026-10-03 reasoning_content rotate (c25862ee); 2026-10-07 this phrase gap. Lesson: every NEW upstream error wording is a new hole. Next time an agent stops on a pool error, grep the exact upstream message against isModelSupportErrorMessage / openAICompatCredentialStatus first, add it with a real-body test.
- Open: still unknown which home sends this 400 (now skipped, not removed). Orcarouter keys return 403 model_access_denied on every model; Apinex needs the daily check-in (402). Both cool down 30 min and retry; not removed (add-only rule).
