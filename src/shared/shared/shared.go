package shared

import (
	"reflect"
	"shared/parameters"
	"strings"
	"framework/messages"
	"strconv"
	"time"
	"fmt"
	"os"
	"io/ioutil"
	"shared/errors"
	"plugin"
)

type Invocation struct {
	Method  reflect.Value
	InArgs  []reflect.Value
	OutArgs [] reflect.Value
}

// MAPE-K Types
type MonitoredCorrectiveData string   // used in channel Monitor -> Analyser (Corrective)
type MonitoredEvolutiveData [] string // used in channel Monitor -> Analyser (Evolutive)
type MonitoredProactiveData [] string // used in channel Monitor -> Analyser (Proactive)

type AnalysisResult struct {
	Result   interface{}
	Analysis int
}

type AdaptationPlan struct {
	Plan string
}

var ValidActions = map[string]bool{
	parameters.INVP: true,
	parameters.TERP: true,
	parameters.INVR: true,
	parameters.TERR: true}

/*
func IsInternal(action string) bool {
	r := false
	if len(action) >= 2 {
		if action[0:2] == parameters.PREFIX_INTERNAL_ACTION {
			r = true
		}
	} else {
		r = false
	}
	return r
}
*/

func IsInternal(action string) bool {
	if action[0:2] == parameters.PREFIX_INTERNAL_ACTION {
		return true
	}
	return false
}

func IsExternal(action string) bool {
	r := false

	if len(action) >= 2 {
		action := strings.TrimSpace(action)
		if strings.Contains(action, ".") {
			action = action[:strings.Index(action, ".")]
		}

		if action == parameters.INVP || action == parameters.TERP || action == parameters.INVR || action == parameters.TERR {
			r = true
		} else {
			r = false
		}
	} else {
		r = false
	}
	return r
}

func IsAction(action string) bool {
	r := false

	action = strings.TrimSpace(action)
	if len(action) > 2 {
		if strings.Contains(action, parameters.INVP) || strings.Contains(action, parameters.TERP) || strings.Contains(action, parameters.INVR) || strings.Contains(action, parameters.TERR) {
			r = true
		} else {
			r = false
		}
	}
	return r
}

type ParMapActions struct {
	F1 func(*chan messages.SAMessage, *messages.SAMessage)      // External action
	F2 func(any interface{}, name string, args ... interface{}) // Internal action
	P1 *messages.SAMessage
	P2 *chan messages.SAMessage
	P3 interface{}
	P4 string
}

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func IsToElement(action string) bool {
	if action[:2] == parameters.PREFIX_INTERNAL_ACTION || action[:4] == parameters.INVP || action[:4] == parameters.TERR {
		return true
	} else { // TerP and InvR
		return false
	}
}

func LoadParameters(args []string) {
	for i := range args {
		variable := strings.Split(args[i], "=")
		switch strings.TrimSpace(variable[0]) {
		case "SAMPLE_SIZE":
			parameters.SAMPLE_SIZE, _ = strconv.Atoi(variable[1])
		case "REQUEST_TIME":
			temp1, _ := strconv.Atoi(variable[1])
			temp2 := time.Duration(temp1)
			parameters.REQUEST_TIME = temp2
		case "INJECTION_TIME":
			temp1, _ := strconv.Atoi(variable[1])
			temp2 := time.Duration(temp1)
			parameters.INJECTION_TIME = temp2
		case "MONITOR_TIME":
			temp1, _ := strconv.Atoi(variable[1])
			temp2 := time.Duration(temp1)
			parameters.MONITOR_TIME = temp2
		case "STRATEGY":
			parameters.STRATEGY, _ = strconv.Atoi(variable[1])
		case "NAMING_HOST":
			parameters.NAMING_HOST = variable[1]
		case "QUEUEING_HOST":
			parameters.QUEUEING_HOST = variable[1]
		default:
			fmt.Println("Shared:: Parameter '" + variable[0] + "' does not exist")
			os.Exit(0)
		}
	}
}

func ShowExecutionParameters(s bool) {
	if s {
		fmt.Println("******************************************")
		fmt.Println("Sample size                : " + strconv.Itoa(parameters.SAMPLE_SIZE))
		fmt.Println("Direrctory of base code    : " + parameters.BASE_DIR)
		fmt.Println("Directory of plugins       : " + parameters.DIR_PLUGINS)
		fmt.Println("Directory of CSP specs     : " + parameters.DIR_CSP)
		fmt.Println("Directory of Configurations: " + parameters.DIR_CONF)
		fmt.Println("Directory of Go compiler   : " + parameters.DIR_GO)
		fmt.Println("Directory of FDR           : " + parameters.DIR_FDR)
		fmt.Println("------------------------------------------")
		fmt.Println("Naming Host     : " + parameters.NAMING_HOST)
		fmt.Println("Naming Port     : " + strconv.Itoa(parameters.NAMING_PORT))
		fmt.Println("Calculator Port : " + strconv.Itoa(parameters.CALCULATOR_PORT))
		fmt.Println("Fibonacci Port  : " + strconv.Itoa(parameters.FIBONACCI_PORT))
		fmt.Println("Queueing Port   : " + strconv.Itoa(parameters.QUEUEING_PORT))
		fmt.Println("------------------------------------------")
		fmt.Println("Plugin Base Name: " + parameters.PLUGIN_BASE_NAME)
		fmt.Println("Max Graph Size  : " + strconv.Itoa(parameters.GRAPH_SIZE))
		fmt.Println("------------------------------------------")
		fmt.Println("Adaptability  ")
		fmt.Println("Corrective        : " + strconv.FormatBool(parameters.IS_CORRECTIVE))
		fmt.Println("Evolutive         : " + strconv.FormatBool(parameters.IS_EVOLUTIVE))
		fmt.Println("Proactive         : " + strconv.FormatBool(parameters.IS_PROACTIVE))
		fmt.Println("Monitor Time (s)  : " + (parameters.MONITOR_TIME * time.Second).String())
		fmt.Println("Injection Time (s): " + (parameters.INJECTION_TIME * time.Second).String())
		fmt.Println("Request Time (ms) : " + parameters.REQUEST_TIME.String())
		fmt.Println("Strategy (0-NOT DEFINED 1-No change 2-Change once 3-change same plugin 4-alternate plugins): " + strconv.Itoa(parameters.STRATEGY))
		fmt.Println("******************************************")
	}
}

func Log(args ...string) {
	if strings.Contains(args[1], "Proxy") || strings.Contains(args[1], "XXX") {
		fmt.Println(args[0] + ":" + args[1] + ":" + args[2] + ":" + args[3])
	}
}

func IsConnectorInAttachments(atts []string, t string) bool {
	foundConnector := false

	for a := range atts {
		att := strings.Split(atts[a], ",")
		if (strings.TrimSpace(att[1]) == t) {
			foundConnector = true
		}
	}

	return foundConnector
}

func IsComponentInAttachments(atts []string, c string) bool {
	foundComponent := false

	for a := range atts {
		att := strings.Split(atts[a], ",")
		if (strings.TrimSpace(att[0]) == c || strings.TrimSpace(att[2]) == c) {
			foundComponent = true
		}
	}

	return foundComponent
}

func IsInConnectors(conns map[string]string, t string) bool {
	foundConnector := false

	for i := range conns {
		if t == i {
			foundConnector = true
			break
		}
	}
	return foundConnector
}

type QueueingInvocation struct {
	Op   string
	Args []interface{}
}

type QueueingTermination struct {
	R interface{}
}

func Invoke(any interface{}, name string, args ... interface{}) {
	inputs := make([]reflect.Value, len(args))

	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(any).MethodByName(name).Call(inputs)

	inputs = nil
	return
}

func LoadPlugins(confName string) map[string]time.Time {
	listOfPlugins := make(map[string]time.Time)

	pluginsDir := parameters.DIR_PLUGINS + "/" + confName
	OSDir, err := ioutil.ReadDir(pluginsDir)
	if err != nil {
		myError := errors.MyError{Source: "Shared", Message: "Folder '" + pluginsDir + "' not read"}
		myError.ERROR()
	}
	for i := range OSDir {
		fileName := OSDir[i].Name()
		if strings.Contains(fileName, "_plugin") {
			pluginFile := pluginsDir + "/" + fileName
			info, err := os.Stat(pluginFile)
			if err != nil {
				myError := errors.MyError{Source: "Shared", Message: "Plugin '" + pluginFile + "'not read"}
				myError.ERROR()
			}
			listOfPlugins[fileName] = info.ModTime()
		}
	}
	return listOfPlugins
}

func CheckForNewPlugins(listOfOldPlugins map[string]time.Time, listOfNewPlugins map[string]time.Time) [] string {
	var newPlugins [] string

	// check for new plugins
	for key := range listOfNewPlugins {
		val1, _ := listOfNewPlugins[key]
		val2, ok2 := listOfOldPlugins[key]
		if ok2 {
			if val1.After(val2) { // newer version of an old plugin is available
				newPlugins = append(newPlugins, key)
			}
		} else {
			newPlugins = append(newPlugins, key) // a new plugin is available
		}
	}
	return newPlugins
}

func LoadPlugin(confName string, pluginName string, symbolName string) (plugin.Symbol) {

	var lib *plugin.Plugin
	var err error

	pluginFile := parameters.DIR_PLUGINS + "/" + confName + "/" + pluginName
	attempts := 0
	for {
		lib, err = plugin.Open(pluginFile)

		if err != nil {
			if attempts >= 3 {
				fmt.Println(err)
				myError := errors.MyError{Source: "Shared", Message: "Error on trying open plugin " + pluginFile + " " + strconv.Itoa(attempts) + " times"}
				myError.ERROR()
			} else {
				attempts++
				time.Sleep(10 * time.Millisecond)
			}
		} else {
			break
		}
	}

	fx, err := lib.Lookup(symbolName)
	if err != nil {
		myError := errors.MyError{Source: "Shared", Message: "Function " + symbolName + " not found in plugin"}
		myError.ERROR()
	}
	return fx
}
