package tag

import (
	"fmt"
	"math"
	"os/exec"
	"os"
	"strconv"
	"strings"
)

// currentVersions return map[string]string like {"ios": "1.23.4"} or Error
func currentVersions() (map[string]string, error) {
	if os.Getenv("GITHUB_WORKFLOW") == "Cron" {
		exec.Command("git", "pull", "--tags").Run()
	} 

	version := map[string]string{}

	lines, err := exec.Command("git", "tag").Output()
	if err != nil {
		return version, err
	}

	for _, line := range strings.Split(strings.TrimSuffix(string(lines), "\n"), "\n") {
		if strings.Contains(line, "ios") {
			version["ios"] = strings.Split(line, "-")[1]
		}
		if strings.Contains(line, "android") {
			version["android"] = strings.Split(line, "-")[1]
		}
	}
	return version, nil
}

// UpdateVersionTags has a side effect to git tags.
// return messages []strings or error
func UpdateVersionTags(retrived map[string]string) ([]string, error) {
	messages := []string{}

	current, err := currentVersions()
	if err != nil {
		return messages, err
	}

	iosMsg := doAction(current, retrived, "ios")
	if len(iosMsg) > 0 {
		messages = append(messages, iosMsg)
	}

	androidMsg := doAction(current, retrived, "android")
	if len(androidMsg) > 0 {
		messages = append(messages, androidMsg)
	}

	return messages, nil
}

func doAction(current map[string]string, retrived map[string]string, osType string) string {

	cron := os.Getenv("GITHUB_WORKFLOW") == "Cron"
	fmt.Printf("Cron: %s\n", os.Getenv("GITHUB_WORKFLOW"))

	act := selectAction(current[osType], retrived[osType])
	switch act {
	case insert:
		newTag := fmt.Sprintf("%s-%s", osType, retrived[osType])
		exec.Command("git", "tag", newTag).Run()

		if cron {
			fmt.Printf("newTag:  %s\n", newTag)
			err := exec.Command("git", "push", "origin", newTag).Run()
			if err != nil {
				fmt.Println(err)
			}
		}

		return fmt.Sprintf("%s: %s を登録しました", osType, retrived[osType])
	case update:
		oldTag := fmt.Sprintf("%s-%s", osType, current[osType])
		exec.Command("git", "tag", "-d", oldTag).Run()
		newTag := fmt.Sprintf("%s-%s", osType, retrived[osType])
		exec.Command("git", "tag", newTag).Run()

		if cron {
			exec.Command("git", "push", "origin", fmt.Sprintf(":%s", oldTag)).Run()
			exec.Command("git", "push", "origin", newTag).Run()
		}

		return fmt.Sprintf("%s: %s が公開されました", osType, retrived[osType])
	default:
		return ""
	}
}

type action int

const (
	insert action = iota
	update
	ignore
)

func selectAction(current string, retrived string) action {
	if len(retrived) == 0 {
		return ignore
	}

	if len(current) == 0 { // len(retrieved) > 0
		return insert
	}

	if newerThanRight(retrived, current) {
		return update
	}

	return ignore
}

func newerThanRight(lhs string, rhs string) bool {
	left := strings.Split(lhs, ".")
	right := strings.Split(rhs, ".")

	for i := 0; i < int(math.Max(float64(len(left)), float64(len(right)))); i++ {
		if i > len(left)-1 {
			return false
		}

		if i > len(right)-1 {
			return true
		}

		lVal, _ := strconv.Atoi(left[i])
		rVal, _ := strconv.Atoi(right[i])

		if lVal > rVal {
			return true
		} else if lVal < rVal {
			return false
		}
	}

	return false
}
