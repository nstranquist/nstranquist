# Public catalog

Nico Stranquist — senior software engineer, St. Louis.

I build developer platforms, local developer tools, and full-stack products.

These are the public products I will stand behind. Each product has a license and a tagged release you can inspect. Private work and unfinished projects stay off this page.

<p align="center">
  <img src="assets/banner-dark.svg#gh-dark-mode-only" alt="Nico Stranquist — Senior software engineer" width="880" />
  <img src="assets/banner-light.svg#gh-light-mode-only" alt="Nico Stranquist — Senior software engineer" width="880" />
</p>

<p align="center">
  <img src="assets/catalog-dark.svg#gh-dark-mode-only" alt="Public catalog of products" width="880" />
  <img src="assets/catalog-light.svg#gh-light-mode-only" alt="Public catalog of products" width="880" />
</p>

## Selected work

### [docs-puller](https://github.com/nstranquist/docs-puller)

Copies vendor and project docs into Markdown, indexes them with SQLite FTS5, and searches them locally. Quality is measured: the public BM25 sample reaches 95.8% Hit@1 and 100% Hit@5, and anyone can rerun it against pinned pages.

`Go` · Apache-2.0 · [v0.5.0](https://github.com/nstranquist/docs-puller/releases/tag/v0.5.0)

Public BM25 sample: **95.8% Hit@1 · 100% Hit@5**

### [Nicos Catalog](https://github.com/nstranquist/nicos-catalog)

Plugins, typed records, full-text and relationship search, and a drift check that fails when the catalog and the source files disagree. The public core leaves out personal telemetry, valuation, and host-only policy.

`Go` · Apache-2.0 · [v0.1.1](https://github.com/nstranquist/nicos-catalog/releases/tag/v0.1.1)

### [Openbook](https://github.com/nstranquist/openbook)

Accounts, a friend graph, a feed with friends-or-public visibility, comments, notifications, and messages. Updates live through Convex queries. MIT, with a public v0.1.0 release.

`TypeScript` · MIT · [v0.1.0](https://github.com/nstranquist/openbook/releases/tag/v0.1.0)

### [agent-ops](https://github.com/nstranquist/agent-ops)

Trace what an agent did, what it saw, what it was allowed to call, and what was proven. Every answer is saved as a local file (a receipt). Nothing sends data off the machine, and nothing asks the model to score itself.

`Go` · MIT · [v0.3.2](https://github.com/nstranquist/agent-ops/releases/tag/v0.3.2)

### [Nicos Hidden Bar](https://github.com/nstranquist/nicos-hidden-menubar)

Visible, hidden, and always-hidden regions, notch-safe overflow, and a searchable reveal panel. It stays on your Mac: no account, no network service. Public source release under GPL-3.0.

`Swift` · GPL-3.0 · [v0.2.0](https://github.com/nstranquist/nicos-hidden-menubar/releases/tag/v0.2.0)

### [JobKit](https://github.com/nstranquist/jobkit)

Profiles, eligibility rules, claim lists that reject unproven numbers, and apply packages a person reviews before sending. The public tree is synthetic-fixture only (example data, not a real resume or job history).

`Go` · MIT · [v0.9.0](https://github.com/nstranquist/jobkit/releases/tag/v0.9.0)

## How I work

- **Local-first.** The default path works on your computer. A cloud account is not required.
- **Fail-closed.** If a check cannot prove the result is safe, the tool stops. Missing evidence is treated as a failure.
- **Proof, not posture.** I publish evaluations, passing CI, and tagged releases. I do not invent users, customers, or revenue.

## Terms

Short definitions for words used above.

| Term | Meaning |
| --- | --- |
| Local-first | The tool works on your computer. A cloud account is not required. |
| Fail-closed | If a check cannot prove the result is safe, the command stops. Missing evidence is treated as a failure. |
| Drift check | A comparison that fails when the catalog no longer matches the files it describes. |
| BM25 | A standard keyword search ranker. Hit@1 means the right page ranked first. Hit@5 means it ranked in the top five. |
| Synthetic fixture | Made-up example data used in tests and screenshots. Not a real resume or job history. |
| Receipt | A local file that records what happened so you can check it later. |

## Toolbox

Go · TypeScript · React · Node.js · Convex · Swift · Python · PostgreSQL · SQLite · AWS · Cloudflare

## Contact

[LinkedIn](https://linkedin.com/in/nstranquist)

<!-- Generated from catalog.yaml by tools/render. Edit the catalog, then `make render`. -->
