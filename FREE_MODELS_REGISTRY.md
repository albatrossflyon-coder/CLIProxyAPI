# Free Models Registry

Living record of what's actually genuinely free right now on each provider that
was audited, vs. what's blocked and why. Providers change their free lineups
over time — this file exists so a future session can re-check and update
instead of re-deriving everything from scratch. **When re-checking a provider,
update this file's date and findings, don't just trust the old entry.**

> ## 🔁 RE-CHECK CADENCE — monthly
> Chris's call (2026-09-03): the providers rotate models on a roughly monthly
> basis (Mistral dropped 3 pinned versions and pulled `large` to paid between
> 2026-08-17 and 2026-09-03; new `magistral`/`ministral` families appeared).
> **At the start OR end of each month, run a full provider-by-provider
> re-verification pass** (catalog pull + a real completion per configured model,
> for every provider in `config.yaml`), and update this file + `MODEL_COOLDOWN.md`.
> - **Last full pass:** 2026-08-31 (all providers). **Partial 2026-09-03:** full 28-provider alias health sweep + deep refresh of mistral / aihubmix / zai + new provider B.ai + OmniRoute key merge. NOT re-done: full free-model expansion for the other ~24, and the per-key pass.
> - **Next full pass due:** 2026-10-01. Follow the checklist in `PROVIDER_REFRESH_PROCEDURE.md`.
> - Also tracked as a recurring item in `Obsidian Vault\Projects\CLIProxyAPI.md`.

---

## xKiro (xkiro.com) — new provider, wired 2026-09-10

- **Accounts:** 3 keys, all verified live 2026-09-10, separate accounts = separate free quota -- `XkiroAlbatrossflyonGmx` (ends `0cbf`), `XkiroAlbatrossflyon1` (ends `536d`), `XkiroAlbatrossaionline` (ends `81cd`).
- **Base URL:** `https://xkiro.com/v1` (OpenAI-compatible). **Models list:** `GET https://xkiro.com/v1/models` (auth required) — 112 models, `access_tier` field (`free` / `premium` / `paid`) + `pricing.input`/`pricing.output` ($0 = free). ~40 tagged free.
- **Verified free through a real completion 2026-09-10** (returned "OK", usage cost 0), wired in `config.yaml`:
  - `qwen/qwen3.5-flash:free` → `xkiro-qwen-flash` (1M ctx, tools, reasoning)
  - `qwen/qwen3.8-max:free` → `xkiro-qwen-max` (1M ctx)
  - `qwen/qwen3-coder-plus:free` → `xkiro-qwen-coder` (1M ctx)
  - `qwen/qwen3.5-397b-a17b:free` → `xkiro-qwen-397b` (262k ctx)
  - `deepseek/deepseek-v4-flash` → `xkiro-deepseek-flash` (1M ctx)
- **Tagged free but NOT actually free (2026-09-10):** `openai/gpt-5.3-codex-spark` → 402 "This model is available to paying customers only." Do not wire without a wallet balance.
- **Flaky 2026-09-10 (transient 500 "server error"):** `minimax/minimax-m3:free`, `minimax/minimax-m2.5:free` — retry at next pass, wire if stable.
- **Not yet on the VPS instance** — same block goes there during the herdr-VPS Phase 6 session (no IP whitelist, so the key is reusable as-is).
- **Next re-test:** 2026-10-01 full pass — re-pull `/v1/models`, re-run a completion per wired alias, expand the free set if the qwen/minimax free lineup holds.

---

## 2026-09-03 partial pass — findings

**aihubmix:** free tier = every model with a `-free` suffix (55 of them; docs: "limited, trial only, expect 429s"). Config: `aihubmix-gemini-free` (gemini-3.8-flash-free), `-gpt-5.5`, `-gpt-4.1-mini`, `-glm-coding` (coding-glm-5.3-free), `-qwen` (qwen3.6-plus-preview-free), `-dots-3` — all verified through the proxy. Also free but flaky through the proxy (400 "no_available_channel", worked direct): `gpt-oss-20b-free`, `nemotron-3-ultra-550b-a55b-free`, `hy3-free`, `gemma-4-31b-it-free` — re-add next pass if stable. `grok-*` are PAID (403 balance).

**B.ai** (`bai`, `https://api.b.ai/v1`) — NEW. Credit-based relay, OpenAI + Anthropic compatible. Free (0-credit, verified): `glm-5.3-flash`, `qwen3.8-flash`, `hy3`, `mimo-v2.5`. Everything else (all Claude, GPT-5.5+, kimi, deepseek, gemini, minimax) → 403 "Access restricted. Deposit required to unlock premium models" or "credit insufficient balance". 3 keys (albatrossflyon1 / GMX / albatrossai).

**zai** — `glm-5.3` (`glm-coding` alias) and `glm-4.6` (`glm-4.6` alias) removed: both PAID per-token, all 4 keys now 429 code 1113 "余额不足或无可用资源包" (no balance / no resource pack). `glm-4.7-flash` (`glm-free`) still genuinely free. zai-spare `glm-4.5-flash` (`glm-free-spare`) still free.

**OmniRoute merge** — the retired OmniRoute Fly deploy held 9 provider connections; the only ones not already in CLIProxyAPI and still usable: mistral `GzYjgcs15…` (JuneBugMistralGmx acct) and groq `gsk_cd3shnw8…`. Merged both. An OpenAI `sk-proj-…` key there is dead (429 "no credits"). No Manus keys in OmniRoute.

**Alias sweep — all 28 providers reachable through the proxy.** Flaky/noted, not removed: `llm7-llama3.1-8b` (400 "currently unavailable"), `orca-deepseek-free` (503 auth_unavailable — the `deepseek/deepseek-v4-pro-free` model ID likely rotated on orcarouter).

**Not verified this pass:** individual key health for the ~65 pre-existing keys (the proxy sweep only proves each provider is reachable via ≥1 key). Last per-key check: 2026-08-31.

Built 2026-08-31 (TIC 1) during a full "get every key sorted" pass, prompted by
Chris after several providers configured with paid-tier models kept 402/403ing
through CLIProxyAPI while apparently having real free tiers on their own docs.

---

## Mistral — free "Experiment" tier = rate-limited access to MOST of the catalog (re-verified 2026-09-03)

**Real offering:** Mistral's free plan is rate-limited access to nearly the whole
catalog — there are no per-model "free" tags (unlike Z.AI). So "which are free" =
almost all of them, subject to rate limits. Confirmed via OmniRoute's own provider
metadata (`hasFree=true`, "rate-limited access to all models") and, 2026-09-03, by
a real chat completion on each model ID below.

**The one confirmed exclusion:** `mistral-large` (any form, incl. `-latest`) now
returns **HTTP 403 `tier_not_allowed` code 1910** — "not available in your
subscription tier." Mistral pulled the `large` tier off the free plan sometime
between 2026-08-17 and 2026-09-03. Removed from config.

**Catalog churn (2026-08-17 → 2026-09-03):** the dated versions we had pinned
rotated out — `mistral-medium-2508` GONE, `devstral-2512` GONE (no replacement;
`mistral-code-latest` is the nearest successor). `mistral-small-2603` and
`codestral-2508` survived. Lesson: **pin `-latest`, not dated versions**, for
Mistral — they churn fast.

**Config now (6 models, all verified with a real completion 2026-09-03):**
| alias | model ID | notes |
|---|---|---|
| `mistral-medium` | `mistral-medium-latest` | resolves to mistral-medium-3.5 |
| `mistral-small` | `mistral-small-latest` | |
| `mistral-codestral` | `codestral-latest` | coding, FIM |
| `mistral-devstral` | `mistral-code-latest` | devstral's successor (alias name kept for back-compat) |
| `mistral-magistral` | `magistral-medium-latest` | reasoning model — NEW since last pass |
| `mistral-ministral` | `ministral-8b-latest` | small/fast — NEW |

**6 keys on file** (see `config.yaml` mistral block). All in the same free pool.

**Re-check procedure:** `curl -s https://api.mistral.ai/v1/models -H "Authorization: Bearer <key>"`
for the catalog, then POST a 5-token completion to `/v1/chat/completions` per
candidate — a model in the list can still 403 (`large`) or be tier-gated.

---

## TokenRouter — 1 genuinely free model right now, everything else paid

**Real offering:** `z-ai/glm-5.3-free` is the ONLY model on the whole 126-model catalog priced at $0.0000/$0.0000 (input/output). Confirmed via the public `/models` catalog page (no login needed) and a real completion (HTTP 200, real content back). The site's own homepage banner ("LIMITED OFFER — GLM-5.3 Model Now Free to Use!") only ever promotes this one model, singular — no other free-tier promotion found anywhere on the site.

**This is a promotional, time-limited offer, not a standing free tier** — the site literally labels it "LIMITED OFFER"/"LIMITED TIME". Chris confirmed 2026-09-01: **free until September 30, 2026.** ⚠️ **Re-check around that date** whether TokenRouter is offering another free model to rotate to — don't assume it just goes away with nothing to replace it. Note it's a genuinely separate model ID from `z-ai/glm-5.3` ($1.40/$4.40 per 1M, paid) and `z-ai/glm-5.3-flash` (50%-off promo, still paid at $0.075/$0.25) — using the wrong ID will bill you.

**Not checked exhaustively:** the full 126-model catalog is paginated client-side (JS-rendered, URL query params don't page it), so only the first 24 (sorted "Newest") were confirmed. No other `$0.0000` entries appeared there, and the single-model marketing banner is strong evidence this is the only free one, but a full pass across all 6 pages hasn't been done.

**Keys on file, all wired into `config.yaml` (+ live `EasyCLIProxyAPI` copy), model alias `tokenrouter-glm-5.3-free`:**
- Key #1: `sk-...kKon` (label `TokenRouterTIC1@2`), verified live via a real completion.
- Key #2: `sk-...YgBt`, added 2026-09-01, verified live via a real completion.
- Key #3: `sk-...mfIu` (label `TokenRouterAlbatrossAI`), added 2026-09-01, verified live via a real completion.
- Key #4: `sk-...lY0J` (label `ClaudeCodeArmyTokenRouter`), added 2026-09-01 (separate session, same day), verified live via a real completion, both through the raw provider endpoint and end-to-end through the local proxy (`tokenrouter-glm-5.3-free` alias, port 8317).

- Key #5: `sk-...cSKk` (label `AlbatrossflyonGMXTokenrouter`), added 2026-09-01, verified live via a real completion.

**Standing rule, per Chris 2026-09-01: every TokenRouter key added — this one and any future ones — points at `z-ai/glm-5.3-free` only, never another model on this provider**, unless a future re-check finds a real reason to pay for one. The `TIC1@2` naming on this key suggests more TokenRouter accounts are likely coming (same pattern as the xAI/Mistral/ChinaAPI multi-account keys already in this file) — apply the same model to each one as it's added, don't re-derive per key.

---

## The Grid (thegrid.ai) — $25 signup credit, spot-market pricing, not a standing free tier

**Real offering:** a spot market for AI inference — OpenAI-compatible Chat Completions API at `https://api.thegrid.ai/v1`, pricing is live market-clearing rate per token, not fixed. New accounts get a **$25 signup credit**, confirmed via the site itself. This is credit-based, not a permanently-free model like TokenRouter's `glm-5.3-free` above — it runs out.

**Model names are "instruments", not pinned models** — `text-prime`/`text-max`/`text-standard`, `code-prime`/`code-max`/`code-standard`, `agent-prime`/`agent-max`/`agent-standard`, plus lab-specific "latest" routes (`claude-opus-latest`, `gpt-sol-latest`, `kimi-latest`, `glm-latest`, `deepseek-pro-latest`, `gemini-pro-latest`, `bytedance-pro-latest`). Each instrument routes to whichever backing model is cheapest/available at request time — confirmed via a real completion on `text-prime`, which routed to `minimax/minimax-m3` under the hood.

**Wired into `config.yaml` (+ live `EasyCLIProxyAPI` copy)** as provider `thegrid`, 3 instrument aliases: `grid-text-prime`, `grid-code-prime`, `grid-agent-prime`.
- Key: `STTrJS2TkK4scYOIV1ue9xMm2WZa6ga/8d4ByOyvr1w` (label `thegridAlbatrossGmx`), added 2026-09-01, verified live via a real completion.

**Not tracked:** remaining credit balance — no endpoint checked yet for this. Since it's a spend-down credit (not a monthly-reset free tier like TokenRouter's), re-check balance periodically rather than assuming it's still available indefinitely.

---

## Apinex (api.apinex.bond): 8 genuinely-free models, `free/` prefix, wired 2026-09-08

**Real offering:** an OpenAI-compatible LLM aggregator (`https://api.apinex.bond/v1`,
also `https://apinex.bond/v1`). Auth `Authorization: Bearer <key>`. `/v1/models` and
`/v1/chat/completions` both work. The public models page (`apinex.bond/models`) is a
JS SPA and does not render via curl, so the live source of truth for the free set is
the `/v1/models` catalog: every model whose id carries the `free/` prefix returns
`cost_tokens: 0` on a real completion. Paid tier (from $0.07/1M, weight-based billing):
`claude/opus-5`, `claude/sonnet-5`, `gemini/3.8-flash`, `deepseek/v4-pro`, `kimi/k3`,
`gpt/5.6-{sol,terra,luna}`, `gpt/6-astra`, `glm/5.3`, `grok/4.6`, etc.

**Live free set as of 2026-09-08 (`GET https://api.apinex.bond/v1/models`, 8 of 22 ids
carry the `free/` prefix):**

| model id | alias in config |
|---|---|
| `free/gemini-3.8-flash` | `apinex-gemini-3.8-flash` |
| `free/gemini-3.1-pro` | `apinex-gemini-3.1-pro` |
| `free/deepseek-v4-pro-0813` | `apinex-deepseek-v4-pro` |
| `free/deepseek-v4-flash-0731` | `apinex-deepseek-v4-flash` |
| `free/glm-5.3-flash` | `apinex-glm-5.3-flash` |
| `free/gpt-5.6-luna` | `apinex-gpt-5.6-luna` |
| `free/qwen-3.8-max` | `apinex-qwen-3.8-max` |
| `free/muse-spark-1.3` | `apinex-muse-spark-1.3` |

**3 keys wired as a rotation pool** (separate accounts / signups, so likely separate
free-quota pools), all smoke-tested live 2026-09-08 against `/v1/chat/completions`
(a ~5-token prompt, `free/gemini-3.8-flash` and `free/glm-5.3-flash`):

| key | account label | smoke test |
|---|---|---|
| `sk-...497a` | `ApinexAlbatrossflyon1` | HTTP 200, real "pong", `cost_tokens: 0` |
| `sk-...fbb5` | (new, from Chris 2026-09-08) | HTTP 200, real "pong", `cost_tokens: 0` |
| `sk-...09ef` | `ApinexClaudeCodeArmy` | HTTP 200, real "pong", `cost_tokens: 0` |

Wired into BOTH `config.yaml` (repo) and the live `EasyCLIProxyAPI` copy, provider block
`apinex`. NOT restarted (Chris's 2-click restart) as of this write, so the new provider
block and aliases are not live in the running proxy yet.

**Rate limits:** NOT published anywhere on the site or in the API response. Treat this
as a fallback rotation route, not a primary. **Re-check procedure:** `curl -s
https://api.apinex.bond/v1/models -H "Authorization: Bearer <key>"`, filter ids for the
`free/` prefix, then POST a 5-token completion per id to confirm `cost_tokens: 0`.

---

## UnoRouter (api.unorouter.com): 3 free aliases wired 2026-09-09

**Real offering:** OpenAI-compatible gateway (`https://api.unorouter.com/v1`, `/v1/models`
+ `/v1/chat/completions`, `Authorization: Bearer <key>`). ~213 models, ~115 carry a
`:free` suffix ($0). Free tier is ~1 req/min/model/user; a cap returns `429` + a
`Retry-After` header. Pay-as-you-go from $1, credit never expires. Open source / self-hostable.
Terms checked 2026-09-09: no training on user content; private personal gateway use is fine
(the only relevant prohibition is reselling account access to third parties, which this is not).

**3 aliases wired** (kept small on purpose - the 1 req/min cap makes this a fallback route,
and `glm-5.3-flash` already has 3 other CLIProxyAPI routes: `tokenrouter-glm-5.3-free`,
`bai-glm-flash`, `apinex-glm-5.3-flash`):

| model id | alias in config |
|---|---|
| `glm-5.3-flash:free` | `unorouter-glm-flash` |
| `gemini-3.1-flash-lite:free` | `unorouter-gemini-flash-lite` |
| `gpt-oss-120b:free` | `unorouter-gpt-oss-120b` |

Other `:free` ids of note if the pool ever needs widening: `glm-5.3:free`,
`glm-5.3-thinking:free`, `gpt-oss-20b:free`, `llama-4-maverick-17b-128e-instruct:free`,
plus DeepSeek / Qwen / Kimi variants.

**3 keys wired as a rotation pool** (separate accounts, separate free quota):

| key | account label | smoke test |
|---|---|---|
| `sk-...ClrE` | `UnoRouterAlbatrossflyon1` | 2026-09-09: `/v1/models` 215, `glm-5.3-flash:free` -> "PONG" |
| `sk-...GL9V` | `UnoRouterAlbatrossflyonGmx` | 2026-09-09: `gemini-3.1-flash-lite:free` -> "PONG" HTTP 200 |
| `sk-...QuWN` | `claudecodearmyUnorouter` | 2026-09-09: `/v1/models` 213, `glm-5.3-flash:free` -> "PONG" |

Wired into BOTH `config.yaml` (repo) and the live `EasyCLIProxyAPI` copy, provider block
`unorouter` (last entry under `openai-compatibility`, now 30 providers). Both re-validated
with `yaml.safe_load`. **Local restart done 2026-09-09** - all 3 `unorouter-*` aliases
confirmed in `localhost:8317/v1/models` (116 aliases total) and routing verified: a real
`/v1/chat/completions` reaches UnoRouter and comes back with its own free-tier throttle
(`get_channel_failed` / "all providers busy, rate limit") after the smoke-test calls
earlier the same day - expected for a ~1 req/min/model fallback route, clears on its own.
**NOT yet on the VPS instance** - the VPS `config.yaml` needs the same block + a redeploy
of the Coolify service (`albatross-vps` / `dwpx7ylvitd0zmmeg6gtwygk`).

**Rate limits:** ~1 req/min/model/user, `429` + `Retry-After`. Fallback route only.
**Re-check:** `curl -s https://api.unorouter.com/v1/models -H "Authorization: Bearer <key>"`,
count `:free` ids, POST a 5-token completion on `glm-5.3-flash:free`.

---

## Cerebras — REMOVED from config.yaml (no free tier at all right now)

**Real offering:** a ONE-TIME $5 signup trial credit, "access to all models" —
not a renewing free tier. Once spent, it's gone until a fresh account.

**Tested 2026-08-31:** both models on this account (`gpt-oss-120b`,
`gemma-4-31b` — confirmed via `GET /v1/models`) returned the identical
`payment_required_error` directly against `api.cerebras.ai`. Account-wide
block, not per-model — the trial credit is exhausted on this specific key.

**To bring back:** would need a genuinely fresh Cerebras signup (new email) for
a new $5 trial, or real payment. Key preserved below for that path if ever
taken.
- `csk-...tepm`

**Re-check procedure:** `curl https://api.cerebras.ai/v1/models -H "Authorization: Bearer <key>"` then test each returned model id directly with a chat completion.

---

## Free-LLM.com / github.com/nejib1/Free-LLM — external reference directory (checked 2026-09-01)

A community-maintained directory of 120+ free LLM models across 41 providers, MIT-licensed, README confirmed legitimate (real rate limits/credit amounts per provider, code snippets in `code-examples/`). Cross-checked against this file's own providers — heavy overlap (Groq, Cohere, OpenRouter, Cloudflare, LLM7, Mistral, Pollinations, DeepSeek, Qwen, SiliconFlow all match what's already here). Its Cerebras entry ($5 one-time trial) matches this file's own Cerebras section above exactly — not a discrepancy, just confirms it.

**Providers listed there NOT currently in `config.yaml`, worth a look sometime:**
- **Ollama Cloud** (`ollama.com/cloud`) — no card, light usage tier, GPT-OSS 120B/20B, Qwen3.5, DeepSeek V4 Flash
- **Hetzner Inference API** (`experiments.hetzner.com/inference`) — no card, free during experimental phase, Qwen3.6 35B A3B
- **Inference.net** — no card, fair-use rate limit, DeepSeek-R1, Llama 3.1 8B/70B
- **SambaNova Cloud** — registration, $5 trial credit, 3-month expiry
- **Fireworks AI** — registration, $1 trial credit

Not evaluated further this session — flagged as candidates for a future key-collection pass, not verified live.

---

## SiliconFlow — REMOVED from config.yaml (0 of 77 models free)

**Real offering:** a one-time $1 signup credit — not a renewing free tier,
no per-model $0 pricing found anywhere in their catalog.

**Tested 2026-08-31:** every one of the 77 models on the account (`GET
/v1/models`) was individually tested with a real chat completion. All 61 chat
models returned the identical `{"code":30001,"message":"Sorry, your account
balance is insufficient"}`. The other 16 (image/video/audio/embedding/reranker
models) correctly 400'd as "model does not exist" against the chat endpoint —
expected, they're not chat models. **Zero genuinely free models found.**

**Keys preserved for reference (2 accounts):**
- `sk-...ezpp`
- `sk-...xhpo` ("Albatrossflyon1" account)

**Re-check procedure:** `curl https://api.siliconflow.com/v1/models -H "Authorization: Bearer <key>"`, then loop every returned id through a real chat completion (script used this session, not saved — rebuild similarly if re-checking).

---

## DeepSeek — REMOVED from config.yaml (no free tier, ever, per official docs)

**Real offering:** none. `api-docs.deepseek.com/quick_start/pricing` confirms
pure pay-per-token from day one, no trial credit mentioned anywhere.

**Tested 2026-08-31:** our config's model IDs (`deepseek-chat`,
`deepseek-reasoner`) are stale/renamed — the real current catalog (`GET
/v1/models`) is `deepseek-v4-flash`, `deepseek-v4-pro`,
`deepseek-v4-flash-vision-exp`. Tested the real current model directly on all
3 configured keys — all 3 "Insufficient Balance". Not a naming issue, a real
balance requirement with no free option.

**Already have indirect free DeepSeek access** via `orca-deepseek-free` /
`orcarouter-free` (OrcaRouter's free proxy) and `chinaapi-deepseek-flash` —
no functional loss from removing this provider block.

**Keys preserved for reference (3 accounts):**
- `sk-...777b` ("albatrossflyon1")
- `sk-...510a` ("albatrossaionline")
- `sk-...a56e` ("GMX.com")

**Re-check procedure:** `curl https://api.deepseek.com/v1/models -H "Authorization: Bearer <key>"`, test the real current model IDs it returns, not the old `deepseek-chat`/`deepseek-reasoner` names.

---

## LLM7 — FIXED 2026-08-31, 5 genuinely free models now live in config.yaml

**Real offering:** a real, currently-usable free tier — but only on models the
catalog tags `tier: "turbo"` AND `usage_based_only: false`. `tier: "pro"`
models are real per-token paid regardless of the dashboard's advertised "1M
free tokens/day" (that number applies to turbo-tier usage, not the whole
catalog).

**Live catalog as of 2026-08-31** (`GET https://api.llm7.io/v1/models`, 44
models total, 6 tagged "turbo"): individually tested all 6 turbo models with a
real chat completion directly against `api.llm7.io`.

| Model ID | Result | Now aliased in config as |
|---|---|---|
| `codestral-latest` | ✅ works, real completion | `llm7-codestral` |
| `gpt-oss` (resolves to gpt-oss-20b) | ✅ works, real completion | `llm7-gpt-oss` |
| `minimax-m2.7` | ✅ works, real completion | `llm7-minimax` |
| `mistral-Nemo-Instruct-2407` | ✅ works (tiny tracked cost, still succeeded) | `llm7-mistral-nemo` |
| `meta-Llama-3.1-8B-Instruct-Turbo` | ✅ works (tiny tracked cost, still succeeded) | `llm7-llama3.1-8b` |
| `gemma4:31b` | ❌ "Insufficient balance" | not added — `usage_based_only:true` in the catalog, the real tell (tier label alone isn't reliable) |

Old config (removed): `claude-sonnet-5` → `llm7-claude-sonnet-5`, `deepseek-v4-flash` → `llm7-deepseek-flash`, `gemini-3-flash` → `llm7-gemini-flash` — all 3 were `tier: "pro"`, real paid, correctly 402'd. Not restored.

**Re-check procedure (their catalog rotates — do this again periodically):**
```
curl https://api.llm7.io/v1/models
```
Filter for `"tier": "turbo"` AND `"usage_based_only": false` — those are the real free candidates. Test each individually with a chat completion before trusting the flag (as `gemma4:31b` proved, the tier label alone isn't reliable).

---

## Pollinations — RESOLVED 2026-08-31, working now

Chris claimed the grant on the dashboard (`enter.pollinations.ai`) mid-session
— exact mechanism inside the dashboard not confirmed (it's a JS SPA, not
readable via fetch), but balance went from 0.0000 to enough to serve requests.
**Re-tested and confirmed live**, both directly against `gen.pollinations.ai`
and through CLIProxyAPI's `pollinations-openai` alias — real completions,
routes to `gpt-5.4-nano-2026-03-17` under the `openai` model name (a real,
capable model, not a placeholder). No config change needed, it was always
correctly configured — just had zero balance until claimed.

<details><summary>Original blocked state (for reference)</summary>


**Real offering:** a genuine renewing free grant (1.5 pollen/week base +
tier-based daily grants for registered developers + one-time rewards like
starring their GitHub repo) — but it has to be claimed on the dashboard, it
doesn't apply automatically to a signed-up key.

**Tested 2026-08-31:** account balance is 0.0000 pollen, confirmed via a real
402. Also confirmed their genuinely-free ANONYMOUS tier (no signup, 1
req/15s) is a *different, non-OpenAI-compatible* endpoint — doesn't work
through CLIProxyAPI's `/v1/chat/completions` route regardless (tested
directly with no Authorization header: `401 UNAUTHORIZED`, "a valid API key
is required"). Same incompatibility class as the earlier OVHcloud finding.

Original action item (done): visit `https://enter.pollinations.ai`, claim the
grant, then re-test `pollinations-openai`. Confirmed resolved above.

</details>

---

## Vercel AI Gateway — REMOVED from config.yaml 2026-08-31 (Chris's call)

**Real offering:** Vercel's own docs say every team gets free-tier AI Gateway
credits automatically. In practice, this specific team account still requires
a card before anything (even free-tier models) will serve a request — the
error message itself says so explicitly:

> "AI Gateway requires a valid credit card on file to service requests.
> Please visit https://vercel.com/d?to=%2F[team]%2F~%2Fai%3Fmodal%3Dadd-credit-card
> to add a card and unlock your free credits."

**Confirmed 2026-08-31:** re-tested directly, same error, verbatim. Not a
model-selection problem — `vercel-poolside-free` is already the correct
$0-priced model per Vercel's own model catalog; the gate is account-level.

**Not fixed** — no config/model workaround exists. If a card ever gets added
later, the key below is still valid and `poolside/laguna-s-2.1-free` is
confirmed the correct $0-priced model to re-add.
- `vck_...xcAh`

---

## CLōD — 3 accounts wired 2026-08-31, 7 free models each, 300 free req/day combined

**Real offering:** CLōD (`api.clod.io`, OpenAI-compatible) markets "7 free
models" on its pricing page — free tier is a **100-requests/day quota per
account**, not $0-per-token pricing (their own price calculator still lists
non-zero $/M for every model, including the ones marketed as free).

**The 7 free-tier models, confirmed via `/v1/models` + live completions on
2026-08-31:** `Trinity Mini`, `Gemma 4 31B IT`, `Llama 3.1 8B`, `Meta Llama
3.3 70B Instruct`, `GPT OSS 120B`, `GPT OSS 20B`, `Qwen 3.5 9B`.

Three separate CLōD team accounts now wired into `config.yaml` as `clod` /
`clod-2` / `clod-3` (each account's 7 models aliased `clod{,2,3}-<model>`),
giving 300 combined free requests/day:
- `clod` — "AlbatrossflyonGmxCloD", verified live via a real completion on GPT OSS 20B.
- `clod-2` — "AlbatrossAiOnlineCLoD", separate account, verified live via Qwen 3.5 9B.
- `clod-3` — "CLoDAlbatrossflyon1", a third distinct account (confirmed via a
  different teamId/userId in the JWT payload) added 2026-08-31 after Chris
  couldn't locate an earlier key for this account anywhere else. Verified
  live via a real completion on Llama 3.1 8B, then re-verified end-to-end
  through CLIProxyAPI's own `clod3-llama-3.1-8b` alias (HTTP 200).

**Not verified:** the per-request token cap under the free tier — not worth
burning free requests to find out, same call as other unknowns in this file.

---

## Net effect of this pass

- Removed 4 dead-weight providers (Cerebras, SiliconFlow, DeepSeek, Vercel AI
  Gateway) — 9 keys retired from active rotation, zero functional loss (all
  had real free-provider alternatives elsewhere already, e.g. orcarouter for
  DeepSeek). Vercel's block is real card-required, not a config bug, but
  Chris's call was to drop it rather than wait on adding a card.
- Fixed LLM7 — went from 0-of-3 working models to 5-of-6 real free models
  confirmed live.
- Fixed Pollinations — Chris claimed the grant mid-session, confirmed working live.
- Provider count: 26 → 22 in `config.yaml`'s `openai-compatibility` list.
