package tag

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// Message contains os("ios", "android"), version, template.
type Message struct {
	OS 		string
	Version string
	tmpl	*template.Template
}

func (m Message) message() (string, error) {
	var message bytes.Buffer
	if err := m.tmpl.Execute(&message, m); err != nil {
		return "", err
	}

	return message.String(), nil
}

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
func UpdateVersionTags(retrived map[string]string, registerdMessage string, updateMessage string) ([]string, error) {
	messages := []string{}

	current, err := currentVersions()
	if err != nil {
		return messages, err
	}

	iosMsg := doAction(current, retrived, "ios", registerdMessage, updateMessage)
	if len(iosMsg) > 0 {
		messages = append(messages, iosMsg)
	}

	androidMsg := doAction(current, retrived, "android", registerdMessage, updateMessage)
	if len(androidMsg) > 0 {
		messages = append(messages, androidMsg)
	}

	return messages, nil
}

func doAction(current map[string]string, retrived map[string]string, osType string, registerdMessage string, updateMessage string) string {
	cron := os.Getenv("GITHUB_WORKFLOW") == "Cron"

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

		tmpl, err := template.New("insert").Parse(registerdMessage)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		message, err := Message { osType, retrived[osType], tmpl}.message()
		if err != nil {
			fmt.Println(err)
			return ""
		}

		return message
	case update:
		oldTag := fmt.Sprintf("%s-%s", osType, current[osType])
		exec.Command("git", "tag", "-d", oldTag).Run()
		newTag := fmt.Sprintf("%s-%s", osType, retrived[osType])
		exec.Command("git", "tag", newTag).Run()

		if cron {
			exec.Command("git", "push", "origin", fmt.Sprintf(":%s", oldTag)).Run()
			exec.Command("git", "push", "origin", newTag).Run()
		}

		tmpl, err := template.New("insert").Parse(updateMessage)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		message, err := Message { osType, retrived[osType], tmpl}.message()
		if err != nil {
			fmt.Println(err)
			return ""
		}

		return message
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

	if retrived != current {
		return update
	}

	return ignore
}
