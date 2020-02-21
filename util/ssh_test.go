package util

import (
	"fmt"
	"testing"
)

func TestSshClient_SendCommands(t *testing.T) {
	version := "UMER18V16.19.31SP03-NIA01"
	apiKey:="AKCp5ejxxUhRuTtpjJkdDeYrjjC6SQ8A8VPrJZzzpVujLSvyeuzrwSTT8uSFd3YtyN4fZuw8r"
	getToolUrl := "https://artsh.zte.com.cn/artifactory/ume-snapshot-generic/TOOLS/get_tool.py"
	cmd1 := fmt.Sprintf("sudo curl --header X-JFrog-Art-Api:%s -o /home/ubuntu/get_tool.py %s",
		apiKey, getToolUrl)
	cmd2 := "cd /home/ubuntu/ && sudo python get_tool.py"
	cmd3 := fmt.Sprintf("cd /home/ubuntu/deploy-tool && sudo python ume_patch_with_json_tool.py -v %s", version)


	conn, err := NewSshClient("10.89.138.241:22", "ubuntu", "cloud")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.CloseClient()

	go func() {
		bs := make([]byte, 0)
		for b := range conn.GetOutput() {
			if b == byte('\n') {
				fmt.Println(string(bs))
				bs = make([]byte, 0)
				continue
			}
			bs = append(bs, b)
		}
	}()

	if err := conn.SendCommands(cmd1, cmd2, cmd3); err != nil {
		panic(err)
	}
}
