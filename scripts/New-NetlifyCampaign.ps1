#Requires -Version 5.1
<#
.SYNOPSIS
  Creates data/campaigns/<id>.toml and content/<lang>/<slug>/ (landing, thank-you, material, qr).
  Run from the Hugo site root (or pass -RepoRoot).

.EXAMPLE
  cd C:\my-site; .\scripts\New-NetlifyCampaign.ps1 -Slug "workshop-ai"
.EXAMPLE
  .\New-NetlifyCampaign.ps1 -Slug "meetup-sp" -ContentLangDir "en" -Title "Workshop materials"
#>
param(
    [Parameter(Mandatory = $true)]
    [string]$Slug,
    [string]$RepoRoot = "",
    [string]$ContentLangDir = "en",
    [string]$CampaignId = "",
    [string]$FormName = "",
    [string]$Title = "Event landing",
    [string]$Description = "Enter your details to access the material."
)

$ErrorActionPreference = "Stop"

function Write-Utf8NoBom {
    param([string]$LiteralPath, [string]$Value)
    $utf8 = New-Object System.Text.UTF8Encoding $false
    [System.IO.File]::WriteAllText($LiteralPath, $Value, $utf8)
}

if (-not $RepoRoot) {
    $RepoRoot = (Get-Location).Path
}
$RepoRoot = (Resolve-Path $RepoRoot).Path

if (-not $CampaignId) { $CampaignId = $Slug }
if (-not $FormName) { $FormName = "$CampaignId-leads" }

$dataPath = Join-Path $RepoRoot "data/campaigns/$CampaignId.toml"
if (Test-Path $dataPath) {
    throw "Already exists: $dataPath"
}

$tplPath = Join-Path $PSScriptRoot "templates/campaign.toml.tpl"
if (-not (Test-Path $tplPath)) {
    throw "Template not found: $tplPath"
}

$toml = Get-Content -LiteralPath $tplPath -Raw -Encoding UTF8
$toml = $toml.Replace("__SLUG__", $Slug).Replace("__CAMPAIGN_ID__", $CampaignId).Replace("__FORM_NAME__", $FormName).Replace("__LANG__", $ContentLangDir.TrimStart('/'))
Write-Utf8NoBom -LiteralPath $dataPath -Value $toml

$contentDir = Join-Path $RepoRoot "content/$ContentLangDir/$Slug"
New-Item -ItemType Directory -Force -Path $contentDir | Out-Null

$common = @"
ShowBreadCrumbs: false
hideMeta: true
ShowPostNavLinks: false
campaign: $CampaignId
"@

$indexMd = @"
---
title: "$Title"
description: "$Description"
$common
---
"@
Write-Utf8NoBom -LiteralPath (Join-Path $contentDir "_index.md") -Value $indexMd

$thanksMd = @"
---
title: "Thank you"
slug: thank-you
$common
---
"@
Write-Utf8NoBom -LiteralPath (Join-Path $contentDir "thank-you.md") -Value $thanksMd

$materialMd = @"
---
title: "Material"
slug: material
$common
---
"@
Write-Utf8NoBom -LiteralPath (Join-Path $contentDir "material.md") -Value $materialMd

$qrMd = @"
---
title: "Redirecting"
slug: qr
$common
_build:
  list: never
---
"@
Write-Utf8NoBom -LiteralPath (Join-Path $contentDir "qr.md") -Value $qrMd

Write-Host "OK: $dataPath"
Write-Host "OK: $contentDir"
Write-Host "Set file_id in data/campaigns/$CampaignId.toml; run hugo server; add host redirects if needed."
