# Public catalog

Nico Stranquist — senior software engineer, St. Louis.

I build developer platforms, local-first AI infrastructure, and full-stack products that can prove what they claim.

This page is the inspectable extract of that practice. Each product has a license, a proof tag, and a claim boundary. Private systems and unpublished extracts stay off this surface.

<p align="center">
  <img src="assets/banner-dark.svg#gh-dark-mode-only" alt="Nico Stranquist — Senior software engineer" width="880" />
  <img src="assets/banner-light.svg#gh-light-mode-only" alt="Nico Stranquist — Senior software engineer" width="880" />
</p>

<p align="center">
  <img src="assets/catalog-dark.svg#gh-dark-mode-only" alt="Public catalog of inspectable products" width="880" />
  <img src="assets/catalog-light.svg#gh-light-mode-only" alt="Public catalog of inspectable products" width="880" />
</p>

## Selected work

### [docs-puller](https://github.com/nstranquist/docs-puller)

Mirrors vendor and project docs to Markdown, indexes them with SQLite FTS5, and searches them privately. Quality is a product feature: the public BM25 sample reaches 95.8% Hit@1 and 100% Hit@5, and anyone can rerun it against pinned pages.

`Go` · Apache-2.0 · [v0.5.0](https://github.com/nstranquist/docs-puller/releases/tag/v0.5.0)

Public BM25 sample: **95.8% Hit@1 · 100% Hit@5**

### [Nicos Catalog](https://github.com/nstranquist/nicos-catalog)

Provider plugins, typed entities, full-text and relationship search, and drift checks that fail closed when the catalog and the source system disagree. The public core excludes personal telemetry, valuation, and host-specific policy.

`Go` · Apache-2.0 · [v0.1.1](https://github.com/nstranquist/nicos-catalog/releases/tag/v0.1.1)

### [Openbook](https://github.com/nstranquist/openbook)

Identity, a friend graph, a visibility-scoped feed, comments, notifications, and messages — live-synced through reactive queries. MIT, with a public v0.1.0 release.

`TypeScript` · MIT · [v0.1.0](https://github.com/nstranquist/openbook/releases/tag/v0.1.0)

### [agent-ops](https://github.com/nstranquist/agent-ops)

Trace what an agent did, what it saw, what it was allowed to call, and what was proven. Every answer is a durable local artifact. Nothing phones home, and nothing asks the model to grade itself.

`Go` · MIT · [v0.3.2](https://github.com/nstranquist/agent-ops/releases/tag/v0.3.2)

### [Nicos Hidden Bar](https://github.com/nstranquist/nicos-hidden-menubar)

Visible, hidden, and always-hidden lanes, notch-aware overflow, and a searchable reveal panel. Local-first: no account, no network service. Public source release under GPL-3.0.

`Swift` · GPL-3.0 · [v0.2.0](https://github.com/nstranquist/nicos-hidden-menubar/releases/tag/v0.2.0)

### [JobKit](https://github.com/nstranquist/jobkit)

Profiles, eligibility gates, claim allowlists that fail closed, and human-reviewed apply plans. The public tree is synthetic-fixture only — no personal profile or job-application data.

`Go` · MIT · [v0.9.0](https://github.com/nstranquist/jobkit/releases/tag/v0.9.0)

## How I work

- **Local-first.** The default path works offline. Cloud is optional, not a prerequisite.
- **Fail-closed.** Catalogs drift-check. Claims need evidence. Ambiguous state does not silently pass.
- **Proof, not posture.** Public evaluations, green CI, and tagged releases. No invented adoption, customers, or revenue.

## Toolbox

Go · TypeScript · React · Node.js · Convex · Swift · Python · PostgreSQL · SQLite · AWS · Cloudflare

## Contact

[LinkedIn](https://linkedin.com/in/nstranquist)

<!-- Generated from catalog.yaml by tools/render. Edit the catalog, then `make render`. -->
