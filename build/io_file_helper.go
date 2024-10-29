package build

import (
	"encoding/json"
	"fmt"
	"os"
)

func readIpTable(ipTablePath string) map[uint64]map[uint64]string {
	// Read the contents of ipTable.json
	file, err := os.ReadFile(ipTablePath)
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	// Create a map to store the IP addresses
	var ipMap map[uint64]map[uint64]string
	// Unmarshal the JSON data into the map
	err = json.Unmarshal(file, &ipMap)
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	return ipMap
}

func attachLineToFile(filePath string, line string) error {
	// 以追加模式打开文件，如果文件不存在则创建
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close() // 确保函数结束时关闭文件

	// 写入文件的内容，附加一个换行符
	if _, err := file.WriteString(line + "\n"); err != nil {
		return err
	}

	return nil
}
