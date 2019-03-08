package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

type config struct {
	InputDir             string
	OutputDir            string
	IgnoreNameList       []string
	TargetFileSuffix     string
	ParameterKeyPrefix   string
	ParameterValuePrefix string
}

var parameterKeyPrefix = "randomKey_"
var parameterValuePrefix = "randomValue_"

func main2() {
	testConfig := &config{
		InputDir:             "E:/porject/zombieDefense/MainProject/assets/script",
		OutputDir:            "E:/porject/koudaitafang/MainProject/assets/script",
		TargetFileSuffix:     "ts",
		ParameterKeyPrefix:   "randomKey_",
		ParameterValuePrefix: "randomValue_",
	}
	testConfig.IgnoreNameList = append(testConfig.IgnoreNameList, "ald-game.js")
	testConfig.IgnoreNameList = append(testConfig.IgnoreNameList, "ald-game-conf.js")
	testConfig.IgnoreNameList = append(testConfig.IgnoreNameList, "wx_mini_game.d.ts")
	testConfig.IgnoreNameList = append(testConfig.IgnoreNameList, "wxAPI.ts")
	testConfig.IgnoreNameList = append(testConfig.IgnoreNameList, "gcfg.ts")
	content, _ := json.Marshal(testConfig)
	ioutil.WriteFile("testConfig.json", content, 0666)
}
func main() {
	if len(os.Args) != 2 {
		fmt.Println(" not found config dir ")
		return
	}
	configContent, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(" can't open config file ", os.Args[1])
		return
	}
	configInfo := &config{}
	err = json.Unmarshal(configContent, configInfo)
	if err != nil {
		fmt.Println(" error on decode config  ", os.Args[1])
		return
	}

	//inputDir := "E:/porject/zombieDefense/MainProject/assets/script"
	//outputDir := "E:/porject/koudaitafang/MainProject/assets/script"
	//var ignoreNameList []string
	//ignoreNameList = append(ignoreNameList, "ald-game.js")
	//ignoreNameList = append(ignoreNameList, "ald-game-conf.js")
	//ignoreNameList = append(ignoreNameList, "wx_mini_game.d.ts")
	//ignoreNameList = append(ignoreNameList, "wxAPI.ts")
	//ignoreNameList = append(ignoreNameList, "gcfg.ts")
	//targetFileSuffix := "ts"
	inputDir := configInfo.InputDir
	outputDir := configInfo.OutputDir
	ignoreNameList := configInfo.IgnoreNameList
	targetFileSuffix := configInfo.TargetFileSuffix
	parameterKeyPrefix = configInfo.ParameterKeyPrefix
	parameterValuePrefix = configInfo.ParameterValuePrefix

	// make sure out put dir is not exist
	if _, err = os.Stat(outputDir); !os.IsNotExist(err) {
		fmt.Println("out put dir ", outputDir, " must be empty")
		return
	}

	// load all file
	allFileList := loadAllFile(inputDir, "")

	totalCount := len(allFileList)
	currentCount := 0

	for _, fileInfo := range allFileList {
		currentCount++
		if fileInfo.isDir {

			continue
		}
		// make sure target dir exist
		makeSureDir(outputDir + fileInfo.subDir + "/")
		// load file
		fileContent, err := ioutil.ReadFile(fileInfo.dir + "/" + fileInfo.name)
		if err != nil {
			fmt.Println("error on load file ", fileInfo.dir, fileInfo.name)
			continue
		}
		// check is target suffix & not in ignore list
		if isNameMatchSuffix(targetFileSuffix, fileInfo.name) &&
			!isNameInIgnoreList(ignoreNameList, fileInfo.name) {
			// do fix

			fileContent = ([]byte)(handleTSContent((string)(fileContent)))
		}
		// do copy
		err = ioutil.WriteFile(outputDir+fileInfo.subDir+"/"+fileInfo.name, fileContent, 0666)
		if err != nil {
			fmt.Println("error on write file ", err, outputDir+fileInfo.subDir+"/"+fileInfo.name)
			return
		}
		// print process
		fmt.Println("process: ", currentCount, totalCount, outputDir+fileInfo.subDir+"/"+fileInfo.name)
	}

	// print done
	fmt.Println("succeed")
}
func makeSureDir(dirPath string) {
	os.MkdirAll(dirPath, 0777)
}

type customFileInfo struct {
	isDir  bool
	name   string
	dir    string
	subDir string
}

func loadAllFile(filePath string, subDir string) []*customFileInfo {
	fileList, err := ioutil.ReadDir(filePath)
	if err != nil {
		return nil
	}
	var fileInfoList []*customFileInfo

	for _, file := range fileList {
		if file.IsDir() {
			customDir := filePath + "/" + file.Name()
			fileInfoList = append(fileInfoList, &customFileInfo{isDir: true, dir: customDir, subDir: subDir})
			childFilePath := loadAllFile(customDir, subDir+"/"+file.Name())
			if nil != childFilePath {
				fileInfoList = append(fileInfoList, childFilePath...)
			}
		} else {
			fileInfoList = append(fileInfoList, &customFileInfo{isDir: false, dir: filePath, name: file.Name(), subDir: subDir})
		}

	}
	return fileInfoList
}
func isNameInIgnoreList(ignoreList []string, name string) bool {
	for i := 0; i < len(ignoreList); i++ {
		if name == ignoreList[i] {
			return true
		}
	}
	return false
}
func isNameMatchSuffix(targetSuffix string, name string) bool {
	if strings.HasSuffix(name, targetSuffix) {
		return true
	}
	return false
}
func main1() {
	inputFile := "autoFixTs/AudioManager.ts"
	outputFile := "autoFixTs/output.ts"

	fileContent, err := ioutil.ReadFile(inputFile)
	if nil != err {
		fmt.Println(err)
		return
	}
	fixedFileContent := handleTSContent((string)(fileContent))
	ioutil.WriteFile(outputFile, ([]byte)(fixedFileContent), 0666)
}
func handleTSContent(content string) string {
	fixedContent := ([]rune)(content)
	result := ""
	status := "left"
	var contentList []string
	tmpChar := ""
	targetLeftIncludeIndex := 0
	targetRightIncludeIndex := 0
	targetResultRightIncludeIndex := 0
	for i := 0; i < len(fixedContent); i++ {
		char := string(fixedContent[i])
		result += char
		if isSpace(char) {
			if !isSpace(tmpChar) {
				contentList = append(contentList, tmpChar)
				tmpChar = ""
			}
		} else if char == "(" || char == ")" || char == "{" || char == "}" || char == "," {
			if !isSpace(tmpChar) {
				contentList = append(contentList, tmpChar)
				tmpChar = ""
			}
			contentList = append(contentList, "(")
		} else {
			tmpChar += char
		}
		if status == "left" {
			if char == "(" {
				status = "right"
				targetLeftIncludeIndex = len(contentList) - 1
				continue
			}
		}
		if status == "right" {
			if char == ")" {
				status = "leftPart"
				targetRightIncludeIndex = len(contentList) - 1
				targetResultRightIncludeIndex = len(result) - 1
				continue
			}
			if !isValidChar(char) && !isValidParameterChar(char) {
				status = "left"
				continue
			}
		}
		if status == "leftPart" {
			if char == "{" {
				status = "left"
				result = doFixResult(result, targetResultRightIncludeIndex, contentList, targetLeftIncludeIndex, targetRightIncludeIndex)
				continue
			}
			if !isValidChar(char) && char != ":" && char != "<" && char != ">" {
				status = "left"
				continue
			}
		}
	}
	return result
}
func doFixResult(result string, resultRightIncludeIndex int, contentList []string, leftIncludeIndex int, rightIncludeIndex int) string {
	tmpIndex := leftIncludeIndex - 1
	if tmpIndex < 0 || tmpIndex >= len(contentList) {
		return result
	}
	targetKeyWord := contentList[tmpIndex]
	if targetKeyWord == "if" {
		return result
	}
	if targetKeyWord == "for" {
		return result
	}
	if targetKeyWord == "switch" {
		return result
	}
	if targetKeyWord == "while" {
		return result
	}
	if targetKeyWord == "catch" {
		return result
	}
	if targetKeyWord == "function" {
		return result
	}
	tmpIndex -= 1
	if tmpIndex >= 0 && tmpIndex < len(contentList) && contentList[tmpIndex] == "get" {
		return result
	}
	if tmpIndex >= 0 && tmpIndex < len(contentList) && contentList[tmpIndex] == "set" {
		return result
	}
	if tmpIndex >= 0 && tmpIndex < len(contentList) && contentList[tmpIndex] == "new" {
		return result
	}
	return addRandomParameterToString(result, resultRightIncludeIndex, rightIncludeIndex-leftIncludeIndex <= 1)
}
func addRandomParameterToString(result string, rightIncludeIndex int, isEmptyParameter bool) string {
	randomParameter := getRandomParameter(isEmptyParameter)

	prefix := result[:rightIncludeIndex]
	suffix := result[rightIncludeIndex:]

	return prefix + randomParameter + suffix
}
func getRandomParameter(isEmptyParameter bool) string {
	result := ""
	valueCount := 3
	if !isEmptyParameter {
		result += ","
	}
	randomCount := int(rand.Int31()%3) + valueCount
	for i := 0; i < randomCount; i++ {
		k := fmt.Sprintf("%s%d", parameterKeyPrefix, i)
		v := fmt.Sprintf("%s%d", parameterValuePrefix, i)
		result += k + "=" + "\"" + v + "\""
		if i != randomCount-1 {
			result += ","
		}
	}
	return result
}
func isNumber(char string) bool {
	min := "0"
	max := "9"

	if char >= min && char <= max {
		return true
	}
	return false
}
func isLetter(char string) bool {
	min := "a"
	max := "z"
	if char >= min && char <= max {
		return true
	}
	min = "A"
	max = "Z"
	if char >= min && char <= max {
		return true
	}
	return false
}
func isSpace(char string) bool {
	if char == "" {
		return true
	}
	if char == " " {
		return true
	}
	if char == "\t" {
		return true
	}
	if char == "\n" {
		return true
	}
	if char == "\r" {
		return true
	}
	return false
}
func isValidChar(char string) bool {
	if isNumber(char) {
		return true
	}
	if isLetter(char) {
		return true
	}
	if isSpace(char) {
		return true
	}
	return false
}
func isValidParameterChar(char string) bool {
	if char == "," {
		return true
	}
	if char == ":" {
		return true
	}
	if char == "=" {
		return true
	}
	if char == "<" {
		return true
	}
	if char == ">" {
		return true
	}
	if char == "." {
		return true
	}
	if char == "_" {
		return true
	}
	if char == "\"" {
		return true
	}
	return false
}
