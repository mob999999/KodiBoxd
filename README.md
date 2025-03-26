[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
# KodiBoxd

A Go-based tool to compare your Letterboxd watchlist with your Kodi movie library. Identify which movies you already own and which ones are missing!

---

## Features

- **Fetch Letterboxd Watchlist**: Scrapes your public Letterboxd watchlist, including movie titles and release years
- **Fetch Kodi Library**: Retrieves your Kodi movie library via JSON-RPC API
- **Comparison**: Highlights movies present in both lists and flags missing ones
- **Configuration Setup**: Guides you through initial setup and saves settings to `config.json`

---

## Installation

### Prerequisites
- [Go](https://go.dev/dl/) (1.20+ recommended)
- [Colly](https://github.com/gocolly/colly) (installed automatically via `go mod tidy`)

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/mob999999/kodiboxd.git
   cd kodiboxd
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the executable:
   ```bash
   go build kodiBoxd
   ```

## Usage

Run the compiled executable:
   ```bash
   Windows: 
   kodiBoxd.exe
   linux:
   kodiBoxd
   ```
---

## Configuration

### Auto-Setup
If no `config.json` exists, the program will prompt you to enter details on first run.

### Manual Setup
Create a `config.json` file in the project root with the following structure:
```json
{
  "LetterBoxdUsername": "your_letterboxd_username",
  "KodiIP": "192.168.1.100",
  "KodiPort": "8080",
  "KodiUsername": "optional_username",
  "KodiPassword": "optional_password"
}
```

---

### Example Output
```
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
```

---

## Troubleshooting

- **Parsing Errors**: Ensure your Letterboxd watchlist is public
- **Connection Issues**: Verify Kodi's JSON-RPC API is enabled
- **Title Mismatches**: Movies are matched using normalized titles

---

## License
This project is licensed under the [MIT License](LICENSE) - see the [LICENSE](LICENSE) file for details.