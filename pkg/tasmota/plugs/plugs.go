package plugs

import (
	"fmt"
	"regexp"
)

type Status struct {
	ID          string
	PowerOn     bool
	ActivePower int
}

func ParseStatus(data string) (*Status, error) {
	status := &Status{}

	re := regexp.MustCompile(`(?i)<tr><td[^>]*>(ON|OFF)</td></tr>`)
	matches := re.FindStringSubmatch(data)
	const expectedMatches = 2
	if len(matches) == expectedMatches {
		status.PowerOn = matches[1] == "ON"
	} else {
		return nil, fmt.Errorf(
			"unexpected number of matches: expected %d, got %d",
			expectedMatches,
			len(matches),
		)
	}

	rePower := regexp.MustCompile(`(?i){e}{s}Active Power{m}</td><td style='text-align:left'>(.+)</td>`)
	matchesPower := rePower.FindStringSubmatch(data)
	if len(matchesPower) != expectedMatches {
		return nil, fmt.Errorf(
			"unexpected number of matches for active power: expected %d, got %d",
			expectedMatches,
			len(matchesPower),
		)
	}
	_, err := fmt.Sscanf(matchesPower[1], "%d", &status.ActivePower)
	if err != nil {
		return nil, fmt.Errorf("error parsing active power: %v", err)
	}

	return status, nil
}
