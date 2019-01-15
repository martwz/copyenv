package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

type CopyEnv struct{}

func (c *CopyEnv) retrieveAppNameEnv(cliConnection plugin.CliConnection, appName string) ([]string, error) {
	var guid string

	apps, err := cliConnection.GetApps()
	if err != nil {
		msg := fmt.Sprintf("Failed to retrieve environment for \"%s\", is the app name correct?", appName)
		err = errors.New(msg)
	}

	for _, element := range apps {
		if element.Name == appName {
			guid = element.Guid
			break
		}
	}

	if guid == "" {
		msg := fmt.Sprintf("Failed to retrieve environment for \"%s\", is the app name correct?", appName)
		err = errors.New(msg)
	} else {
		url := fmt.Sprintf("/v2/apps/%s/env", guid)
		output, err := cliConnection.CliCommandWithoutTerminalOutput("curl", url)

		if err != nil {
			msg := fmt.Sprintf("Failed to retrieve environment for \"%s\", is the app name correct?", appName)
			err = errors.New(msg)
		}

		return output, err
	}
	return nil, err
}

func (c *CopyEnv) extractCredentialsJSON(envParent string, credKey string, output []string) ([]byte, error) {
	err := errors.New("Missing service credentials for application")
	var envJson []byte

	envKey := strings.Join(output, "")
	if strings.Contains(envKey, credKey) {
		var f interface{}
		err = json.Unmarshal([]byte(envKey), &f)
		if err != nil {
			return nil, err
		}

		envJSON := f.(map[string]interface{})
		envParentJSON := envJSON[envParent].(map[string]interface{})
		envJson, err = json.Marshal(envParentJSON[credKey])
		if err != nil {
			return nil, err
		}
	}

	return envJson, err
}

func (c *CopyEnv) extractAndExportCredentials(envParent string, credKey string, appEnv []string, plain bool) {
	creds, err := c.extractCredentialsJSON(envParent, credKey, appEnv)
	checkErr(err)

	vcapServices := ""
	if !plain {
		vcapServices = fmt.Sprintf("export %s='%s';", credKey, string(creds[:]))
	} else {
		vcapServices = fmt.Sprintf("%s", string(creds[:]))
	}
	fmt.Println(vcapServices)
}

func (copy *CopyEnv) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] != "CLI-MESSAGE-UNINSTALL" {
		appName := args[1]

		if len(args) < 1 {
			checkErr(errors.New("Missing application name"))
		}

		appEnv, err := copy.retrieveAppNameEnv(cliConnection, appName)
		checkErr(err)

		if contains(args, "--all") {
			copy.extractAndExportCredentials("application_env_json", "VCAP_APPLICATION", appEnv, contains(args, "--plain"))
			fmt.Println("")
		}
		copy.extractAndExportCredentials("system_env_json", "VCAP_SERVICES", appEnv, contains(args, "--plain"))
	}
}

func (c *CopyEnv) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "copyenv",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "copyenv",
				HelpText: "Export application VCAP_SERVICES.",
				UsageDetails: plugin.Usage{
					Usage: "copyenv APP_NAME [--all] [--plain] - Retrieve and export remote application VCAP_SERVICES.",
					Options: map[string]string{
						"all":   "Retrieve both VCAP_SERVICES and VCAP_APPLICATION from remote application",
						"plain": "Return plain JSON",
					},
				},
			},
		},
	}
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: ", err)
		os.Exit(1)
	}
}

func main() {
	plugin.Start(new(CopyEnv))
}
