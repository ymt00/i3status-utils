// Package utils for i3status bar
package utils

import (
	"encoding/json"
	"io"
	"os/exec"

	// "reflect"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	wpaCliBinary = "wpa_cli"
)

func contains(s []interface{}, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsRegex(s []interface{}, e string) bool {
	re := regexp.MustCompile(e)

	for _, a := range s {
		if re.MatchString(a.(string)) {
			return true
		}
	}
	return false
}

// func roundFloat(val float64, precision uint) float64 {
// 	ratio := math.Pow(10, float64(precision))
// 	return math.Round(val*ratio) / ratio
// }

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
		"--no-overlap",
		"--list", "10",
		"--prefix", "",
		"--prompt", " Menu",
		// "--no-spacing",
		"--line-height", "26",
		"--tb", "#6272a4",
		"--tf", "#f8f8f2",
		"--fb", "#282a36",
		"--ff", "#f8f8f2",
		"--nb", "#282a36",
		"--nf", "#6272a4",
		"--hb", "#44475a",
		"--hf", "#50fa7b",
		"--sb", "#44475a",
		"--sf", "#50fa7b",
		"--scb", "#282a36",
		"--scf", "#ff79c6",
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

// IdleInhibitors struct
type IdleInhibitors struct {
	User string `json:"user"`
}

// Tree struct
type Tree struct {
	Nodes          []interface{}  `json:"nodes"`
	FloatingNodes  []interface{}  `json:"floating_nodes"`
	InhibitIdle    string         `json:"inhibit_idle"`
	IdleInhibitors IdleInhibitors `json:"idle_inhibitors"`
}

// SwayMsgTree returns the sway message tree
func SwayMsgTree() []interface{} {
	cmd := exec.Command("swaymsg", "-r", "-t", "get_tree")

	out, _ := cmd.CombinedOutput()

	// var tree map[string]interface{}
	var tree Tree

	json.Unmarshal(out, &tree)

	// return tree
	// return tree["nodes"].([]interface{})
	return tree.Nodes
}

// SwayMsgTreeTrue return the sway message tree
func SwayMsgTreeTrue() Tree {
	cmd := exec.Command("swaymsg", "-r", "-t", "get_tree")

	out, _ := cmd.CombinedOutput()

	// var tree map[string]interface{}
	var tree Tree

	json.Unmarshal(out, &tree)

	// return tree
	// return tree["nodes"].([]interface{})
	return tree
}

// Output struct
type Output struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

// SwayMsgOutputs returns the sway message outputs
func SwayMsgOutputs() []Output {
	cmd := exec.Command("swaymsg", "-r", "-t", "get_outputs")

	out, _ := cmd.CombinedOutput()

	var outputs []Output

	json.Unmarshal(out, &outputs)

	return outputs
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

// GetNextWorkspace return the next available workspace not in use
func GetNextWorkspace() string {
	workspaces := SwayMsgWorkspaces()

	// fmt.Printf("%s\n", workspaces)
	numbers := make([]int, 0)

	for _, workspace := range workspaces {
		i := workspace.Num
		numbers = append(numbers, i)
	}
	sort.Ints(numbers)

	next := 1
	for _, v := range numbers {
		if next < v {
			return strconv.Itoa(next)
		}
		next++
	}

	return strconv.Itoa(len(numbers) + 1)
}

// // SwayMsgInputs returns the sway message inputs
// func SwayMsgInputs() []interface{} {
// 	cmd := exec.Command("swaymsg", "-r", "-t", "get_inputs")
// 	out, _ := cmd.CombinedOutput()
// 	var inputs []interface{}
// 	json.Unmarshal(out, &inputs)
// 	return inputs
// }

// IsAppRunning check if an application 'app' is currently running
func IsAppRunning(app string) bool {
	nodes := SwayMsgTree()

	var apps []interface{}
	appIds(nodes, &apps)

	// if contains(apps, app) {
	// 	return true
	// }

	// return false

	return contains(apps, app)
}

// IsAppRunningRegex check if an application 'app' is currently running
func IsAppRunningRegex(app string) bool {
	nodes := SwayMsgTree()

	var apps []interface{}
	appIds(nodes, &apps)

	// if containsRegex(apps, app) {
	// 	return true
	// }

	// return false

	return containsRegex(apps, app)
}

// LaunchTerminalScript launches a script in footclient terminal
func LaunchTerminalScript(app string, args ...string) {
	arguments := []string{"--app-id=\"" + app + "\""}
	arguments = append(arguments, args...)
	exec.Command("footclient", arguments...).Start()
}

// GetUsdUsd24hChange get the usd and usd_24h_change from json results from Coingecko API
func GetUsdUsd24hChange(coin interface{}) (float64, float64) {
	c := coin.(map[string]interface{})

	return c["usd"].(float64), c["usd_24h_change"].(float64)
}

// ParsePrice parses price
func ParsePrice(p float64) string {
	price := strconv.FormatFloat(p, 'f', 6, 64)
	price = strings.TrimRight(price, "0")
	return strings.TrimRight(price, ".")
}

// send command to wpa_cli
func WpaCliCommand(command string) string {
	cmd := exec.Command("sudo", wpaCliBinary)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, command)

	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	return string(out)
}
