package main

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var regExHyphen = regexp.MustCompile("([a-z])([A-Z])")
const nodeCmd = "docker-machine"

func main() {
	setLogLevel()

	configMap := make(map[interface{}]interface{})
	pwd, _ := os.Getwd()
	yamlFile, err := ioutil.ReadFile(pwd + "/example/aws-cn.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &configMap)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	// set cloud credentials
	setCloudCredentials(&configMap)

	// build create machine command
	cmdArgs := buildCreateCommand(&cloudConfig{
		"amazonec2",
		"us",
		"dev",
	}, configMap)

	// set machine hostname
	cmdArgs	= append(cmdArgs, "test-docker-machine")

	cmd, err := buildCommand(cmdArgs)

	if err != nil {
		return
	}

	//logrus.Infof("Provisioning node %s", obj.Spec.RequestedHostname)

	stdoutReader, stderrReader, err := startReturnOutput(cmd)
	if err != nil {
		return
	}
	defer stdoutReader.Close()
	defer stderrReader.Close()
	defer cmd.Wait()

	scanner := bufio.NewScanner(stdoutReader)
	for scanner.Scan() {
		msg := scanner.Text()

		fmt.Println(msg)

	}
	scanner = bufio.NewScanner(stderrReader)
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Println(msg)
	}

	if err := cmd.Wait(); err != nil {
		return
	}
}

type cloudConfig struct {
	cloudDriver string
	region string
	env string
}

func buildCreateCommand(cloud *cloudConfig, configMap map[interface{}]interface{}) []string {
	sDriver := strings.ToLower(cloud.cloudDriver)
	cmd := []string{"create", "-d", sDriver}
	//
	//cmd = append(cmd, buildEngineOpts("--engine-install-url", []string{node.Spec.EngineInstallUrl})...)
	//cmd = append(cmd, buildEngineOpts("--engine-opt", mapToSlice(node.Status.NodeTemplateSpec.EngineOpt))...)
	//cmd = append(cmd, buildEngineOpts("--engine-env", mapToSlice(node.Status.NodeTemplateSpec.EngineEnv))...)
	//cmd = append(cmd, buildEngineOpts("--engine-insecure-registry", node.Status.NodeTemplateSpec.EngineInsecureRegistry)...)
	//cmd = append(cmd, buildEngineOpts("--engine-label", mapToSlice(node.Status.NodeTemplateSpec.EngineLabel))...)
	//cmd = append(cmd, buildEngineOpts("--engine-registry-mirror", node.Status.NodeTemplateSpec.EngineRegistryMirror)...)
	//cmd = append(cmd, buildEngineOpts("--engine-storage-driver", []string{node.Status.NodeTemplateSpec.EngineStorageDriver})...)

	for k, v := range configMap {
		dmField := "--" + sDriver + "-" + strings.ToLower(regExHyphen.ReplaceAllString(k.(string), "${1}-${2}"))
		if v == nil {
			continue
		}

		switch v.(type) {
		case float64:
			cmd = append(cmd, dmField, fmt.Sprintf("%v", v))
		case string:
			if v.(string) != "" {
				cmd = append(cmd, dmField, v.(string))
			}
		case bool:
			if v.(bool) {
				cmd = append(cmd, dmField)
			}
		case []interface{}:
			for _, s := range v.([]interface{}) {
				if _, ok := s.(string); ok {
					cmd = append(cmd, dmField, s.(string))
				}
			}
		}
	}
	logrus.Debugf("create cmd %v", cmd)
	//cmd = append(cmd, node.Spec.RequestedHostname)
	return cmd
}

func buildCommand(cmdArgs []string) (*exec.Cmd, error) {
	command := exec.Command(nodeCmd, cmdArgs...)
	//command.SysProcAttr = &syscall.SysProcAttr{}
	//	//command.Env = []string{
	//	//}
	return command, nil
}

