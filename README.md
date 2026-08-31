<div align="center">

# 🗺️ Google Maps Scraper — Business Lead Extractor

### Extract Google Maps leads — names, phones, emails, websites, ratings & reviews — ranked **nearest-first** with distance in m/km. Web UI · REST API · CLI. Written in Go.

[![Go](https://img.shields.io/badge/Go-1.26.5+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker&logoColor=white)](#option-b--docker)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](#-contributing)

**Free · Open-source · Self-hosted Google Maps scraper for lead generation, sales prospecting, local SEO and data enrichment.**

⭐ *If this project helps you, please star it — it helps others discover it.*

</div>

---

![Home page](docs/screenshots/home.png)

## 🔎 What is this?

**Google Maps Scraper** is a fast, open-source tool that turns any Google Maps search into a clean list of business leads. Search *“dentists in Berlin”* or *“cafes near me”* and get **36 data points per business** — phone, website, email, rating, review count, opening hours, GPS coordinates and more — exported to **CSV or JSON**.

Its standout feature: give it a **center point** and every result is ranked **nearest-first**, with the **distance shown in meters/kilometers** and plotted on an interactive map.

> ⚠️ **Use responsibly.** Scraping Google Maps may conflict with Google's Terms of Service. Use this tool for legitimate research and in accordance with applicable laws and terms.

---

## 📊 Example output

A real scrape of **“cafe in athens”** (center point set → sorted by distance):

| # | Business | Distance | Rating | Phone | Website |
|---|----------|----------|--------|-------|---------|
| 1 | Benedict Athens | 285 m | ⭐ 4.8 | +30 210 8083075 | benedict.gr |
| 2 | Taf Coffee | 294 m | ⭐ 4.6 | +30 210 3800014 | cafetaf.gr |
| 3 | Cupaki Street Espresso | 295 m | ⭐ 4.8 | +30 210 3810123 | cupaki.com |
| 4 | Coffee Lab Athens | 307 m | ⭐ 4.9 | +30 210 3318644 | coffeelab.gr |
| 5 | MOTIV | 326 m | ⭐ 4.7 | +30 210 3810899 | motiv.gr |

*…20 unique businesses returned in ~2.4s. Duplicates are removed automatically.*

---

## ✨ Features

| | |
|---|---|
| 📍 **Nearest-first + distance** | Set a center point; results are ranked by proximity with a `m`/`km` distance badge. |
| 🏷️ **36 data points** | Title, category, address, phone, website, rating, reviews, hours, price range, coordinates & more. |
| 🧹 **Automatic de-duplication** | Repeated listings are collapsed so your lead list stays unique. |
| ✉️ **Email extraction** | Optionally visit each business site and pull contact emails. |
| 🗺️ **Interactive map** | Browse places on a live map, filter/search instantly, click to zoom. |
| 🛡️ **Proxy rotation** | SOCKS5 / HTTP / HTTPS proxies to protect your IP and dodge rate limits. |
| ⚡ **REST API in Go** | FastAPI-style JSON API, deployable standalone via `API_ONLY=1`. |
| 📤 **Flexible output** | CSV & JSON; optional PostgreSQL / S3 / LeadsDB backends. |

---

## 🚀 Quick Start

This fork includes the browser lifecycle fixes tested on Windows. See the
[Windows setup guide (Português)](docs/windows-pt-BR.md) and
[dependency provenance and changes](third_party/README.md).

### Option A — Build from source (Go 1.26.5 minimum; Go 1.27 recommended)

```bash
git clone https://github.com/KaiD3v/tools-google-maps.git
cd tools-google-maps
go build -o bin/gms .

# Launch the Web App (UI + API) on http://localhost:8080
./bin/gms -web
```

Then open **http://localhost:8080**:

| Route | Page |
|---|---|
| `/` | Landing page (features, how-it-works, API) |
| `/app` | **Manage** dashboard — create jobs, watch status, view results on a map |
| `/api/docs` | Interactive REST API reference |

### Option B — Docker

```bash
docker build -t gmaps-scraper .
docker run -p 8080:8080 -v "$PWD/out:/app/webdata" gmaps-scraper -web
```

---

## 📍 Nearest-first & distance

Set a **center point** (latitude/longitude) under **Location & Nearest** in the create-job form. The scraper biases the search to that area, then computes the great-circle (haversine) distance from that point to every place, sorts them closest-first, and labels each with `450 m` / `1.20 km` in the map view.

---

## 🔌 REST API — deploy without the UI

The whole scraper is a JSON REST API. Run it **headless** (no HTML, API only):

```bash
API_ONLY=1 ./bin/gms -web -addr :8080
```

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/jobs` | List jobs |
| `POST` | `/api/v1/jobs` | Create a scraping job |
| `GET` | `/api/v1/jobs/{id}` | Get job status + results |
| `DELETE` | `/api/v1/jobs/{id}` | Delete a job |
| `GET` | `/api/v1/jobs/{id}/download` | Download results as CSV |

```bash
curl -X POST http://localhost:8080/api/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{"name":"dentists-berlin","keywords":["dentist in Berlin"],"lang":"en","depth":5,"max_time":"10m"}'
```

---

## ⌨️ CLI

```bash
# queries.txt: one search per line
./bin/gms -input queries.txt -results out.csv -c 1 -depth 10

# JSON + emails
./bin/gms -input queries.txt -results out.json -json -email

# Fast mode (needs a geo center)
./bin/gms -input queries.txt -results out.csv -fast-mode -geo "37.9838,23.7275" -radius 3000
```

| Flag | Description | Default |
|---|---|---|
| `-web` | Run the Web App / API server | off |
| `-addr` | Listen address | `:8080` |
| `-input` / `-results` | Input queries / output file | — |
| `-json` | Write JSON instead of CSV | off |
| `-c` | Concurrency | half CPU cores |
| `-depth` | Max scroll depth | `10` |
| `-geo` / `-radius` / `-zoom` | Center, radius (m), zoom | — / 10000 / 15 |
| `-email` / `-extra-reviews` | Extract emails / more reviews | off |
| `-fast-mode` | Faster, reduced data (needs `-geo`) | off |
| `-proxies` / `-proxies-file` | Proxy URL(s) | — |

---

## 🛡️ Avoiding rate limits / IP blocks

1. **Use proxies** (most important) — residential/mobile SOCKS5 or HTTP:
   `-proxies "socks5://user:pass@host:port"` (or paste them in the Web UI).
2. **Keep concurrency low** (`-c 1`).
3. **Limit depth** (smaller `-depth`).
4. **Prefer the default browser mode** over `-fast-mode` for stealth.

---

## 💡 Use cases

- **Lead generation** — build targeted B2B lead lists by city & niche.
- **Sales prospecting** — phones, emails and websites ready for outreach.
- **Local SEO & market research** — analyze competitors and coverage in an area.
- **Data enrichment** — augment your CRM with ratings, hours and coordinates.
- **Territory planning** — nearest-first ranking for field sales routes.

---

## ❓ FAQ

**Is this Google Maps scraper free?** Yes — open-source under the MIT license, self-hosted, no API key required.

**What data can it extract?** 36 fields including business name, category, address, phone, website, email, rating, review count, opening hours, price range and GPS coordinates.

**Can I export to Excel?** Yes — export CSV and open it in Excel/Google Sheets, or use JSON.

**Does it work without the UI?** Yes — set `API_ONLY=1` to run a headless REST API.

**How do I avoid getting blocked?** Use proxies, keep concurrency low, and limit depth (see above).

---

## 🧱 Tech stack

Go · [scrapemate](https://github.com/gosom/scrapemate) + Playwright (Chromium) · `net/http` + `chi` · htmx UI · Leaflet + OpenStreetMap · SQLite / PostgreSQL · Docker · AWS Lambda.

---

## 🤝 Contributing

Contributions, issues and feature requests are welcome. Open an issue or a PR. If you find the project useful, a ⭐ **star** goes a long way.

## 📄 License

MIT — see [LICENSE](LICENSE). Based on the open-source
[google-maps-scraper](https://github.com/gosom/google-maps-scraper); the MIT license and attribution are retained.

---

<div align="center">

**Keywords:** google maps scraper · google maps data extractor · business lead scraper · lead generation tool · email extractor · local business scraper · maps scraping API · Go web scraper · nearest-first search · CSV/JSON export

</div>
