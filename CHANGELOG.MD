# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.0.0] - 2025-03-26
### Added
- Initial release of Letterboxd web scraper with support for:
  - Watchlist movie extraction
- Kodi JSON-RPC API integration:
  - compare your Kodi Library with your Letterboxd Watchlist and find out which movies you want to watch, but don't yet have
- Go module setup with core dependencies:
  - `github.com/gocolly/colly/v2` for scraping
- MIT License and contribution guidelines

### Security
- Removed hardcoded credentials from early development versions
- Added `.gitignore` rules for environment/config files

---

[v1.0.0]: https://github.com/mob999999/KodiBoxd/releases/tag/1.0.0