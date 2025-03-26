# KodiBoxd

A Go-based tool to compare your Letterboxd watchlist with your Kodi movie library. Identify which movies you already own and which ones are missing!

---

## Features

- **Configuration Setup**: Guides you through initial setup and saves settings to `config.json`.
- **Fetch Letterboxd Watchlist**: Scrapes your public Letterboxd watchlist, including movie titles and release years.
- **Fetch Kodi Library**: Retrieves your Kodi movie library via JSON-RPC API.
- **Comparison**: Highlights movies present in both lists and flags missing ones.

---

## Installation

### Prerequisites
- [Go](https://go.dev/dl/) (1.20+ recommended)
- [Colly](https://github.com/gocolly/colly) (installed automatically via `go mod tidy`)

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/kodiboxd.git
   cd kodiboxd
2. Install dependencies:
   ```bash
    go mod tidy

3. Build the executable:
   ```bash
    go build

 4. Configuration
    Auto-Setup:
    If no config.json exists, the program will prompt you to enter details on first run.
    Manual Setup:
    Create a config.json file in the project root with the following structure:
    ```json

    {
      "LetterBoxdUsername": "your_letterboxd_username",
      "KodiIP": "111.111.1.111",
      "KodiPort": "8080",
      "KodiUsername": "optional_username",
      "KodiPassword": "optional_password"
    }

### Usage

   Run the compiled executable:
   ```bash
   KodiBoxd.exe


### Output Example:
```bash
Your Letterboxd Watchlist:
- Incendies (2010)
- The Handmaiden (2016)
- ...

Your Kodi Library Movies:
- Inception (2010)
- ...

Comparison Results:
✔ Found in Kodi Library: Incendies (2010)
❌ Movies NOT in Kodi Library:
- The Handmaiden (2016)

Troubleshooting

    Parsing Errors: Ensure your Letterboxd watchlist is public. Titles with non-standard formatting (e.g., colons, commas) may log warnings but are skipped gracefully.

    Connection Issues: Verify Kodi’s JSON-RPC API is enabled and credentials are correct.

    Title Mismatches: Movies are matched using normalized titles (case-insensitive, trimmed). Minor discrepancies (e.g., "The Movie" vs. "Movie, The") may cause false negatives.
