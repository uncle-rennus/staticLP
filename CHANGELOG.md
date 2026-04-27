# Changelog

## v1.0.0

- **Module rename** to `github.com/uncle-rennus/staticlp` (import path for `go get` / `hugo mod get`).
- **Breaking:** Partials moved from `layouts/partials/netlify_campaign/` to `layouts/partials/staticlp/`.
- **Breaking:** Styles renamed to `assets/css/staticlp.css`; CSS classes use the `staticlp-*` prefix (replacing `ibrahort-*`).
- README and default scaffold copy default to English; CLI / TUI defaults updated accordingly.
- **`staticlp` CLI** unchanged in spirit: TUI + `new`, `list`, `edit`, `rm` (Go 1.23+).

## v0.2.0

- **staticLP** Go CLI (`cmd/staticlp`): run `staticlp` for a TUI, or `staticlp new|list|edit|rm` for scripts. Scaffolds `data/campaigns/*.toml` and `content/<lang>/<slug>/` with UTF-8 (no BOM). The PowerShell script remains optional.

## v0.1.0

- Initial Hugo module: Netlify Forms funnel (landing, thank-you redirect, Google Drive PDF embed, optional QR hop).
- Partials under `layouts/partials/netlify_campaign/`.
- Styles `assets/css/netlify_campaign.css`.
- Scaffold script `scripts/New-NetlifyCampaign.ps1`.
- `exampleSite` for local verification.
