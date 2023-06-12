// Package utils for i3status bar
package utils

import (
	"encoding/json"
	"io"
	"os/exec"
	"regexp"
	"strings"
)

func appIds(nodes []interface{}, apps *[]interface{}) {
	for i := range nodes {
		n := nodes[i].(map[string]interface{})

		if n["app_id"] != nil {
			*apps = append(*apps, n["app_id"])
		}
		if n["nodes"] != nil {
			appIds(n["nodes"].([]interface{}), apps)
		}
		if n["floating_nodes"] != nil {
			appIds(n["floating_nodes"].([]interface{}), apps)
		}
	}
}

// Bemenu show the bemenu listing 'items'
func Bemenu(items []string, args ...string) string {
	defaultArgs := []string{
		"--ignorecase",
		"--wrap",
		"--fork",
		"--no-exec",
		"--scrollbar", "autohide",
		"--bottom",
		"--grab",
		"--list", "10",
		"--prefix", "",
		"--prompt", " Menu",
	}

	defaultArgs = append(defaultArgs, args...)

	cmd := exec.Command("bemenu", defaultArgs...)

	stdin, _ := cmd.StdinPipe()

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, strings.Join(items, "\n"))
	}()

	out, _ := cmd.CombinedOutput()

	return strings.TrimRight(string(out[:]), "\n")
}

// ScratchpadShow shows/hides 'app' from the scratchpad
func ScratchpadShow(app string) {
	exec.Command("sway", "[app_id="+app+"]", "scratchpad", "show").Start()
}

// // IdleInhibitors struct
// type IdleInhibitors struct {
// 	User string `json:"user"`
// }

// Tree struct
type Tree struct {
	Nodes []interface{} `json:"nodes"`
	// FloatingNodes  []interface{}  `json:"floating_nodes"`
	// InhibitIdle    string         `json:"inhibit_idle"`
	// IdleInhibitors IdleInhibitors `json:"idle_inhibitors"`
}

// SwayMsgTree returns the sway message tree
func SwayMsgTree() []interface{} {
	cmd := exec.Command("swaymsg", "-r", "-t", "get_tree")

	out, _ := cmd.CombinedOutput()

	var tree Tree

	json.Unmarshal(out, &tree)

	return tree.Nodes
}

// Workspace struct
type Workspace struct {
	Name    string `json:"name"`
	Num     int    `json:"num"`
	Focused bool   `json:"focused"`
	Focus   []int  `json:"focus"`
}

// SwayMsgWorkspaces returns the sway message workspaces
func SwayMsgWorkspaces() []Workspace {
	cmd := exec.Command("swaymsg", "-r", "-t", "get_workspaces")

	out, _ := cmd.CombinedOutput()

	var workspaces []Workspace

	json.Unmarshal(out, &workspaces)

	return workspaces
}

// IsAppRunning check if an application 'app' (pattern) is currently running
func IsAppRunning(app string) bool {
	nodes := SwayMsgTree()

	var apps []interface{}
	appIds(nodes, &apps)

	for _, a := range apps {
		matched, _ := regexp.MatchString(app, a.(string))
		if matched {
			return true
		}
	}

	return false
}
