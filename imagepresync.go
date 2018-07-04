package ips

import (
	"os/exec"
	log "github.com/sirupsen/logrus"
	//"fmt"
	"strings"
	//"bytes"
)

//TargetRegistryURL should in the format of "ip:port"
//for example "0.0.0.0:5000"
var TargetRegistryURL = "0.0.0.0:5000"

//the format after "docker push" is registryURL/repository:tag

func PreSync() {
	//this set will force the command history stores in buffer written to .bash_history immediately
	cmdPrompt := "echo PROMPT_COMMAND='history -a'"
	execCmd(cmdPrompt)

	for {
		cmdHistory := exec.Command("tail","-n", "5", "/root/.bash_history")
		out1, err := cmdHistory.Output()
		if err != nil {
			log.Fatalf("cmdHistory.Output failed: %v", err)
		}
		stringCmdHistory := string(out1)
		//fmt.Println(stringCmdHistory)
		//docker push 10.33.12.35:50000/nignx:
		if strings.Contains(stringCmdHistory, "docker push") {
			//var stderr bytes.Buffer
			imageName, fullImageName := getFullImageName(stringCmdHistory)

			cmdDockerPush := "docker push " + fullImageName
			execCmd(cmdDockerPush)


			//cmd.Stderr = &stderr
			if err != nil {
				//fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
				log.Fatalf("cmt.OutPut failed: %v",  err)
			}
			log.Infof("Successfully synchronized image %v to registry %v.", imageName, TargetRegistryURL)

			//fmt.Println(string(out))
			clearCommand := "echo > /root/.bash_history"
			execCmd(clearCommand)
		}
	}
}

func getFullImageName(stringCmdHistory string) (string, string) {
	//dockerPushIndex is the index of 'd'
	dockerPushIndex := strings.Index(stringCmdHistory, "docker push")
	//fmt.Println(dockerPushIndex)
	//enterIndex is the firt indx of "\n" after "docker push"
	enterIndex := strings.Index(stringCmdHistory[dockerPushIndex:],"\n")
	//fmt.Println(enterIndex2)
	//len("docker push ")=12
	//between dockerPushIndex+12 and dockerPushIndex+enterIndex is the name of the pushed image
	imageName := stringCmdHistory[dockerPushIndex+12 : dockerPushIndex+enterIndex]
	//fmt.Println(stringCmdHistory[dockerPushIndex+12 : dockerPushIndex+enterIndex])
	//fullImageName is TargetRegistryURL + imageName
	fullImageName := TargetRegistryURL + "/" + imageName
	//fmt.Println(fullImageName)
	return imageName, fullImageName
}

func execCmd(command string) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Start()
}
