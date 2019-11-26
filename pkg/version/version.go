package version

import (
	"fmt"
	"log"
)

type BumpType int

const (
	MAJOR BumpType = iota + 1
	MINOR
	PATCH
)

func BumpVersion(prevVersion string, kind BumpType) (string, error) {
	var major, minor, patch int
	n, err := fmt.Sscanf(prevVersion, "v%d.%d.%d", &major, &minor, &patch)
	if err != nil {
		log.Fatalf("parsed (%v), err: %v", n, err)
	}

	switch kind {
	case MAJOR:
		major++
	case MINOR:
		minor++
	case PATCH:
		patch++
	default:
		return "", fmt.Errorf("unknown bump type")
	}

	return fmt.Sprintf("v%d.%d.%d", major, minor, patch), nil
}
