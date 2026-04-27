# staticLP

Hugo module and companion CLI for **static landing funnels**: **landing** → **thank-you** → **material** (for example an embedded Google Drive PDF), with an optional **QR** hop page.

The default form markup targets [**Netlify Forms**](https://docs.netlify.com/forms/setup/) (`netlify`, `data-netlify`, honeypot). The HTML is generated at **build time** so Netlify can detect fields on deploy. You can reuse the same layouts on **any** static host; only the **submission endpoint** changes if you are not on Netlify (see [Without Netlify](#without-netlify)).

**Requirements**

- **Hugo** ≥ 0.146 (modules / `go.mod` layout).
- **`staticlp` CLI** (optional but recommended): **Go 1.23+** — `go install` only. On Windows, add `%USERPROFILE%\go\bin` (or `%GOPATH%\bin`) to `PATH`.

---

## Install the Hugo module

In your site’s `hugo.toml` or `config.toml`:

```toml
[module]
[[module.imports]]
path = "github.com/uncle-rennus/staticlp"
```

From the **site root** (where your `go.mod` lives):

```bash
hugo mod get github.com/uncle-rennus/staticlp@v1.0.2
hugo mod tidy
```

For local development of the module itself, use a `replace` in the site `go.mod` pointing at your checkout.

---

## Theme integration

Wire these into your theme (or site `layouts/`) so campaign pages render.

1. **CSS + QR meta** — inside `<head>` (e.g. `extend_head.html`):

   ```go-html-template
   {{- partial "staticlp/head_campaign.html" . -}}
   ```

2. **Optional `body` class** from `analytics.body_class` in the campaign TOML:

   ```go-html-template
   <body class="...{{ partial "staticlp/body_extra_class.html" . }}">
   ```

3. **Landing (section)** — when the page has `campaign` in front matter (usually the section `_index.md`):

   ```go-html-template
   {{- partial "staticlp/landing.html" . }}
   ```

4. **Thank-you / material / QR (single)** — for singles under that section:

   ```go-html-template
   {{- partial "staticlp/flow_single.html" . }}
   ```

Styles ship as `assets/css/staticlp.css` in the module; `head_campaign.html` fingerprints them when `campaign` is set.

`flow_single.html` treats slugs **`thank-you`** and **`material`** (and **`qr`**). Legacy slug **`obrigado`** is still accepted for the thank-you step.

---

## Campaign config and URLs

Each funnel is **`data/campaigns/<id>.toml`** plus **`content/<lang>/<slug>/`** (`_index.md`, `thank-you.md`, `material.md`, `qr.md`).

**Critical:** `[paths]` in the TOML must match **exactly** the URLs Hugo emits for that site (including language prefix if you use `defaultContentLanguageInSubdir`).

Examples:

| Site setup | Example `landing` | Example `thank_you` | Example `material` |
|------------|-------------------|----------------------|----------------------|
| Single language, no subdir | `/webinar/` | `/webinar/thank-you/` | `/webinar/material/` |
| Default language in subdir (`pt`) | `/pt/webinar/` | `/pt/webinar/thank-you/` | `/pt/webinar/material/` |

- **`paths.thank_you`** is the form **`action`**. It must be a **real** page URL. Avoid pointing it at a path that only exists as a **301 redirect** (redirects can strip **POST** bodies).
- **`staticlp new`** fills these from the template using your `--lang` and `--slug`.

### Short URLs (`/webinar` vs `/pt/webinar`)

`hugo server` does **not** read `netlify.toml`. If you use host redirects for short links, use the **canonical** paths in `[paths]` and add **Hugo `aliases`** on the content files (see your site’s Hugo docs) so local builds and non-Netlify hosts behave the same.

---

## `staticlp` CLI

Run from the **Hugo site root**:

```bash
go install github.com/uncle-rennus/staticlp/cmd/staticlp@v1.0.2
staticlp                    # TUI: list / new / edit / remove
staticlp new --slug my-campaign --lang en
staticlp list
staticlp edit my-campaign
staticlp rm my-campaign
```

Scaffold output is **UTF-8 without BOM**. Optional: **`scripts/New-NetlifyCampaign.ps1`** for PowerShell.

---

## With Netlify

Use this path when you want **Netlify Forms** (dashboard, notifications, spam filtering, exports).

### Steps

1. **Add the module** and **theme partials** (sections above).
2. **Create a campaign** — `staticlp new --slug my-event --lang en` (or your real default language code). Edit `data/campaigns/<id>.toml`: set `delivery.file_id` (or `REPLACE_WITH_FILE_ID` until ready).
3. **Confirm `[paths]`** match production URLs (with `/en/` or `/pt/` if applicable).
4. **Build** the same way Netlify will — for example `hugo --gc --minify` (your `netlify.toml` `command`).
5. **Deploy** the `public/` output to Netlify. On first deploy, open **Forms** in the Netlify UI; your form name (`form_name` in TOML) should appear once the HTML is live.
6. **Test** a submission, then check **Form notifications** / **spam** settings in the dashboard.

### Netlify-only form attributes

The bundled partial uses `method="POST"`, `netlify`, `data-netlify="true"`, and `netlify-honeypot`. Those are for Netlify’s parser; they are harmless on other hosts but **do nothing** for submission handling elsewhere.

### Optional: short URLs on Netlify

You can add **`[[redirects]]`** in `netlify.toml` from `/event` → `/pt/event/` (301/302) for **GET** marketing links. Do **not** use a short URL as the form **`action`** if that redirect would run on POST.

### Local preview

`netlify dev` can approximate the production host (including some Forms behavior). Plain `hugo server` still builds the funnel; submissions go nowhere unless you point the form elsewhere (custom partial).

---

## Without Netlify

You can deploy the same static site to **GitHub Pages**, **Cloudflare Pages**, **S3**, **any web server**, etc. The **funnel pages**, **redirect scripts**, and **PDF embed** work like any other Hugo output.

**Netlify Forms will not run.** You need a different strategy for the `<form>`:

| Approach | Notes |
|----------|--------|
| **Override the form partial** | Copy `layouts/partials/staticlp/form.html` from this module into **your site** at the same path. Hugo uses your copy. Change `action` to your provider (Formspree, Getform, Basin, a custom API, etc.) and adjust attributes per their docs. Remove `netlify` / `data-netlify` / `netlify-honeypot` if they are not supported. |
| **Serverless backend** | POST to your own function (AWS Lambda, Cloudflare Worker, etc.); keep fields `name` matching what your endpoint expects. |
| **No form** | Replace `landing.html` in your site with a custom partial that only shows CTAs (mailto, external link). |

The **`thank-you`** and **`material`** pages use normal links and `location.replace` in `flow_single.html`; they do not depend on Netlify.

---

## `exampleSite`

```bash
cd exampleSite
hugo server
```

Use `data/campaigns/demo.toml` as a reference. Set a real Google Drive `file_id` if you want the PDF embed to load.

---

## License

MIT — see [LICENSE](LICENSE).
