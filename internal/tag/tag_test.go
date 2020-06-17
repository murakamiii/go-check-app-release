package tag

import (
	"log"
	"os/exec"
	"reflect"
	"strings"
	"testing"

)

func getTagLines() []string {
	output, _ := exec.Command("git", "tag").Output()
	return strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
}

func deleteAllTags() {
	cmdstr := "git tag | xargs git tag -d"
	err :=exec.Command("sh", "-c", cmdstr).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func insertTags(tags []string) {
	for _, line := range tags {
		exec.Command("git", "tag", line).Run()
	}
}

// http://kitakitabauer.hatenablog.com/entry/2017/04/04/204701
func makePseudoSet(stringSlice []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, v := range stringSlice {
		set[v] = struct{}{}
	}
	return set
}

func equalsStringSlice(lhs []string, rhs []string) bool {
	return reflect.DeepEqual(makePseudoSet(lhs), makePseudoSet(rhs))
}

func setup() []string {
	lines := getTagLines()
	deleteAllTags()
	return lines
}

func tearDown(lines []string) {
	insertTags(lines)
}

func TestCurrentVersions(t *testing.T) {
	lines := setup()
	defer tearDown(lines)

	cases := []struct {
		tags []string
		expects map[string]string
	}{
		{ []string{}, map[string]string{} },
		{ []string{"ios-1.2.3"}, map[string]string{"ios": "1.2.3"} },
		{ []string{"ios-1.2.3", "android-2.3.4"}, map[string]string{"ios": "1.2.3", "android":"2.3.4"} },
		{ []string{"ios-1.2.3", "android-2.3.4", "release-1.0.0"}, map[string]string{"ios": "1.2.3", "android":"2.3.4"} },
	}

	for _, c := range cases {
		insertTags(c.tags)
		versions, err := currentVersions()
		if err != nil {
			log.Fatal(err)
		}
		if !reflect.DeepEqual(versions, c.expects) {
			t.Errorf("test failed: expects: %#v, actuals: %#v", c.expects, versions)
		}
		deleteAllTags()
	}
}

func TestUpdateVersionTags(t *testing.T) {
	lines := setup()
	defer tearDown(lines)

	// TODO: 重い
}



func TestDoAction(t *testing.T) {
	lines := setup()
	defer tearDown(lines)

	cases := []struct {
		tags []string
		retrived map[string]string
		osType string
		expectsMsg string
		expectsTags []string
	}{
		{ []string{}, map[string]string{}, "ios", "", []string{""} },
		{ []string{"ios-0.0.1"}, map[string]string{}, "ios", "", []string{"ios-0.0.1"} },
		{ []string{"ios-0.0.1", "android-0.1.2"}, map[string]string{}, "ios", "", []string{"ios-0.0.1", "android-0.1.2"} },
		{ []string{"android-0.1.2"}, map[string]string{"ios": "0.0.2"}, "ios", "ios: 0.0.2 を登録しました", []string{"android-0.1.2", "ios-0.0.2"} },
		{ []string{"ios-0.0.1"}, map[string]string{"android": "0.2.0"}, "android", "android: 0.2.0 を登録しました", []string{"android-0.2.0", "ios-0.0.1"} },
		{ []string{"ios-0.0.1", "android-0.1.2"}, map[string]string{"ios": "0.0.2"}, "ios", "ios: 0.0.2 が公開されました", []string{"android-0.1.2", "ios-0.0.2"} },
		{ []string{"ios-0.0.1", "android-0.1.2"}, map[string]string{"android": "0.2.0"}, "android", "android: 0.2.0 が公開されました", []string{"android-0.2.0", "ios-0.0.1"} },
	}

	for _, c := range cases {
		insertTags(c.tags)
		current, err := currentVersions()
		if err != nil {
			log.Fatal(err)
		}

		actualMsg := doAction(current, c.retrived, c.osType)
		actualTags := getTagLines()
		if actualMsg != c.expectsMsg || !equalsStringSlice(actualTags, c.expectsTags) {
			t.Errorf("test failed: \n\tcase: %#v, \n\tactualMsg: %s \n\tactualTags: %#v", c, actualMsg, actualTags)
		}

		deleteAllTags()
	}
}

func TestSelectAction(t *testing.T) {
	cases := []struct {
		current string
		retrived string
		expects action
	}{
		{ "", "", ignore },
		{ "", "0.0.1", insert },
		{ "0.0.1", "", ignore },
		{ "0.0.1", "0.0.1", ignore },
		{ "0.0.1", "0.0.2", update },
		{ "1.0.99", "1.1.0", update },
		{ "1.99.99", "2.0.0", update },
		{ "1.1.1", "1.1.0", ignore },
		{ "1.1.0", "1.00.99", ignore },
		{ "1.0.0", "0.99.99", ignore },
	}

	for _, c := range cases {
		if selectAction(c.current, c.retrived) != c.expects {
			// TODO: https://stackoverflow.com/questions/30177344/how-to-print-the-string-representation-of-an-enum-in-go
			t.Errorf("test failed: current: %s, retrived: %s, expects: %#v, actual: %#v", c.current, c.retrived, c.expects, selectAction(c.current, c.retrived))
		}
	}
}

func TestNewerThanRight(t *testing.T) {
	cases := []struct {
		lhs string
		rhs string
		expects bool
	}{
		{ "0.0.1", "0.0.2", false },
		{ "0.0.3", "0.0.2", true },
		{ "1.1.1", "0.99.99", true },
	}
	for _, c := range cases {
		if newerThanRight(c.lhs, c.rhs) != c.expects {
			t.Errorf("test failed: lhs: %s, rhs: %s, expects: %t", c.lhs, c.rhs, c.expects)
		}
	}
}