package monit

import (
	"bytes"
	"gitlab.com/pmoscodegrp/common/heartbeat"
	mqtt2 "gitlab.com/pmoscodegrp/common/mqtt"
	"monitor/monit/storage"
	"os/exec"
	"path"
	"strconv"
	"time"
)

var manager = storage.NewManager()

func OnMessageReceivedHeartbeat(message mqtt2.Message) {
	processName := path.Base(message.Topic)
	manager.IncrementHeartbeat(processName)
}

func RunMonitor() {
	monitorTimer := heartbeat.New(10*time.Second, checkInstances)
	monitorTimer.Run()
}

func checkInstances() {
	manager.IncrementCheck()
	//manager.Debug()

	deadProcesses := manager.GetExceededThresholdNames()

	for _, processName := range deadProcesses {
		processId := getProcessIdOf(processName)
		if processId > -1 {
			killed := killProcessWith(processId)
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
		return -1
	}

	processId, err := strconv.Atoi(out.String())

	if err != nil {
		processId = -1
	}

	return processId
}

func killProcessWith(id int) bool {
	cmd := exec.Command("kill", "-9", string(rune(id)))
	err := cmd.Run()

	return err == nil
}
