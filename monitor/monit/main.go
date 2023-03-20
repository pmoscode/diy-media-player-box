package monit

import (
	"bytes"
	"fmt"
	"gitlab.com/pmoscodegrp/common/heartbeat"
	mqtt2 "gitlab.com/pmoscodegrp/common/mqtt"
	"log"
	"monitor/monit/storage"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

var manager = storage.NewManager()

func OnMessageReceivedHeartbeat(message mqtt2.Message) {
	processName := path.Base(message.Topic)
	manager.IncrementHeartbeat(processName)
}

func RunMonitor(processNames string) {
	manager.Init(processNames)

	monitorTimer := heartbeat.New(10*time.Second, checkInstances)
	monitorTimer.Run()
}

func checkInstances() {
	manager.IncrementCheck()
	//manager.Debug()

	deadProcesses := manager.GetExceededThresholdNames()
	//log.Println("Dead processes: ", deadProcesses)

	for _, processName := range deadProcesses {
		processId := getProcessIdOf(processName)
		//log.Println("Process id '", processId, "' of '", processName, "'")
		if processId > -1 {
			killed := killProcessWith(processId)
			//log.Println("Process killed?: ", killed)
			if killed {
				manager.ResetAllCounterFor(processName)
			}
		}
	}
}

func getProcessIdOf(processName string) int {
	cmd := exec.Command("pidof", processName)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		log.Println(err)
		return -1
	}

	processId, err := strconv.Atoi(strings.Trim(out.String(), "\n"))

	if err != nil {
		log.Println(err)
		processId = -1
	}

	return processId
}

func killProcessWith(id int) bool {
	parameterId := fmt.Sprintf("%d", id)
	cmd := exec.Command("kill", "-9", parameterId)
	err := cmd.Run()

	return err == nil
}
