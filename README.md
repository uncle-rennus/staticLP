# staticLP

Hugo module and companion CLI for **static landing funnels** that use [**Netlify Forms**](https://docs.netlify.com/forms/setup/). The HTML form is emitted at **build time** (required for Netlify to detect it). Typical flow: **landing** → **thank-you** (slug `thank-you`; legacy `obrigado` is still supported in `flow_single.html`) → **material** (e.g. embedded Google Drive PDF), with an optional **QR** hop page.

**Requirements**

- **Hugo** ≥ 0.146 (modules / `go.mod` layout).
- **`staticlp` CLI**: **Go 1.23+** only — install with `go install` (no separate binary releases). On Windows, ensure `%GOPATH%\bin` (or `%USERPROFILE%\go\bin`) is on your `PATH`.

## Why Netlify deploys

Netlify Forms only processes forms present in the **deployed static HTML**. This module is aimed at sites you **build with Hugo and deploy on Netlify**, so submissions, spam filtering, and notifications work without a custom backend. Use another host only if you replicate the same “static HTML form in the build output” contract.

## Install the Hugo module

In your site’s `hugo.toml` (or `config.toml`):

```toml
[module]
[[module.imports]]
path = "github.com/uncle-rennus/staticlp"
```

Then:

```bash
hugo mod get github.com/uncle-rennus/staticlp@v1.0.1
hugo mod tidy
```

If you vendor the module locally during development, use a `replace` in the site’s `go.mod` pointing at your checkout.

## Theme integration

1. **CSS + QR meta** — in `<head>` (e.g. theme `extend_head` or `head.html`):

   ```go-html-template
   {{- partial "staticlp/head_campaign.html" . -}}
   ```

2. **Optional body class** from campaign analytics:

   ```go-html-template
   <body class="...{{ partial "staticlp/body_extra_class.html" . }}">
   ```

3. **Landing (section)** — for section templates when `campaign` is set:

   ```go-html-template
   {{- partial "staticlp/landing.html" . }}
   ```

4. **QR / thank-you / material (single)** — for single pages with `campaign`:

   ```go-html-template
   {{- partial "staticlp/flow_single.html" . }}
   ```

Styles live in `assets/css/staticlp.css` inside the module; the partial fingerprints them when `campaign` is present.

## `staticlp` CLI

From the **Hugo site root** (directory containing `hugo.toml` / `config.*`):

```bash
go install github.com/uncle-rennus/staticlp/cmd/staticlp@v1.0.1
staticlp              # TUI: list / new / edit / remove
staticlp new --slug my-campaign --lang en
staticlp list
staticlp edit my-campaign
staticlp rm my-campaign
```

The CLI scaffolds `data/campaigns/<id>.toml` and `content/<lang>/<slug>/` (`_index.md`, `thank-you.md`, `material.md`, `qr.md`) as **UTF-8 without BOM**.

Optional: **`scripts/New-NetlifyCampaign.ps1`** mirrors the same scaffold for PowerShell users.

## `exampleSite`

```bash
cd exampleSite
hugo server
```

Use `data/campaigns/demo.toml` as a reference; set a real `file_id` for the Google Drive PDF if you want the embed to load.

## License

MIT — see [LICENSE](LICENSE).
