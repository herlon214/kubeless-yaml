package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var file string
var deployedFunctions []*DeployedFunction

// Config is the YAML file structure
type Config struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Runtime string `yaml:"runtime"`
	File    string `yaml:"file"`
	Handler string `yaml:"handler"`
}

// DeployedFunction contains a simple yaml structure for deployed function
type DeployedFunction struct {
	Metadata struct {
		Name      string            `yaml:"name"`
		Namespace string            `yaml:"namespace"`
		Labels    map[string]string `yaml:"labels"`
	} `yaml:"metadata"`
	Spec struct {
		Checksum string `yaml:"checksum"`
	} `yaml:"spec"`
}

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy folders recursively that contains kubeless.yaml configs",
	Long:  `Example: "kubeless-yaml -f example/" will seek for every kubeless.yaml files and deploy the functions.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetOutput(os.Stdout)
		log.SetLevel(log.InfoLevel)
		deployedFunctions = getDeployedFunctions()

		log.Infof("Found %d deployed functions...", len(deployedFunctions))
		if isDirectory(file) {
			deployDirRecursive(file)
		} else {
			deployDir(file)
		}
		log.Info("All done!")
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "File or folder that contains the kubeless.yaml files")
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func deployDirRecursive(dir string) {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		log.Errorf("Could not read dir '%s'", dir)
	}

	log.Infof("-> [%s]", dir)

	for _, file := range files {
		nextDir := fmt.Sprintf("%s/%s", dir, file.Name())

		if file.Name() == "kubeless.yaml" {
			deployDir(dir)
		} else if isDirectory(nextDir) {
			deployDirRecursive(nextDir)
		}
	}

}

func deployDir(dir string) {
	var configs []*Config
	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, "kubeless.yaml"))
	if err != nil {
		log.Errorf("File kubeless.yaml not found in %s", dir)
	}

	err = yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}

	// Deploy all configs inside the kubeless file
	for _, config := range configs {
		var params []string
		willExecute := true

		params = append(params, "function", "deploy", config.Name)
		params = append(params, "--runtime", config.Runtime)
		params = append(params, "--handler", config.Handler)
		params = append(params, "--from-file", fmt.Sprintf("%s/%s", dir, config.File))
		params = append(params, "--runtime", config.Runtime)
		params = append(params, "--label", fmt.Sprintf("kubeless-yaml=%s-%s", config.Name, config.Version))

		// Pre-flight
		var shouldBeDeployed *DeployedFunction
		params = append(params, "--dryrun")
		out, err := exec.Command("kubeless", params...).CombinedOutput()

		err = yaml.Unmarshal(out, &shouldBeDeployed)

		if err != nil {
			log.Error("Failed to parse preflight yaml")
			return
		}

		// Check if there is a deployed function with the same checksum
		deployedFunction := getFunctionByName(deployedFunctions, config.Name)

		if deployedFunction != nil {
			if deployedFunction.Spec.Checksum == shouldBeDeployed.Spec.Checksum {
				log.Infof("Skipping function '%s' due to no changes", config.Name)
				willExecute = false
			} else {
				// Change the deploy command to update command
				log.Infof("Function '%s' will be updated", config.Name)
				params[1] = "update"
			}
		} else {
			log.Infof("Deploying function '%s'", config.Name)
		}

		// Check if the command needs to be executed
		if willExecute {
			// Execute the command without dryrun
			out, err = exec.Command("kubeless", params[:len(params)-1]...).CombinedOutput()

			fmt.Printf("%s\n", string(out))
		}

	}
}

func getDeployedFunctions() []*DeployedFunction {
	log.Info("Getting deployed functions list...")
	var deployedFunctions []*DeployedFunction
	out, err := exec.Command("kubeless", "function", "ls", "-o", "yaml").CombinedOutput()

	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = yaml.Unmarshal(out, &deployedFunctions)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return deployedFunctions
}

func getFunctionByName(functions []*DeployedFunction, name string) *DeployedFunction {
	for _, function := range functions {
		if function.Metadata.Name == name {
			return function
		}
	}

	return nil
}
