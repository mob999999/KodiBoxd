package main

import "strings"

// extractYearFromTitle checks if the movie title ends with a year in parentheses
// and, if so, returns the title without the year and the year itself.
func extractYearFromTitle(title string) (string, string) {
	if strings.HasSuffix(title, ")") {
		openParenIndex := strings.LastIndex(title, "(")
		if openParenIndex != -1 {
			year := title[openParenIndex+1 : len(title)-1]
			// Verify that the extracted part is a 4-digit number.
			if len(year) == 4 && strings.Count(year, "0")+strings.Count(year, "1")+
				strings.Count(year, "2")+strings.Count(year, "3")+strings.Count(year, "4")+
				strings.Count(year, "5")+strings.Count(year, "6")+strings.Count(year, "7")+
				strings.Count(year, "8")+strings.Count(year, "9") == 4 {
				return title[:openParenIndex-1], year
			}
		}
	}
	return title, ""
}

// normalizeTitle converts a movie title to a simplified form for comparison.
func normalizeTitle(title string) string {
	// Remove year if present.
	titleOnly, _ := extractYearFromTitle(title)

	// Convert to lowercase.
	titleOnly = strings.ToLower(titleOnly)

	// Remove common prefixes.
	prefixes := []string{"the ", "a ", "an "}
	for _, prefix := range prefixes {
		if strings.HasPrefix(titleOnly, prefix) {
			titleOnly = titleOnly[len(prefix):]
		}
	}

	// Remove special characters.
	specialChars := []string{":", "-", ",", ".", "'", "\"", "!", "?", "(", ")", "[", "]"}
	for _, char := range specialChars {
		titleOnly = strings.ReplaceAll(titleOnly, char, "")
	}

	// Replace multiple spaces with a single space.
	for strings.Contains(titleOnly, "  ") {
		titleOnly = strings.ReplaceAll(titleOnly, "  ", " ")
	}

	return strings.TrimSpace(titleOnly)
}
