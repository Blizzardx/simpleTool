package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
)

func tmp() {
	//inputDir := ""
	//outputDir := ""
	//var ignoreNameList []string
	//targetFileSuffix := ""

	// make sure out put dir is not exist
	// load all file
	// make sure target dir exist
	// check is target suffix & not in ignore list
	// do fix
	// do copy
	// print process
	// print done
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
func main() {
	inputFile := "autoFixTs/SDataShare.ts"
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
		} else if char == "(" {
			contentList = append(contentList, "(")
			if !isSpace(tmpChar) {
				contentList = append(contentList, tmpChar)
				tmpChar = ""
			}
		} else if char == ")" {
			if !isSpace(tmpChar) {
				contentList = append(contentList, tmpChar)
				tmpChar = ""
			}
			contentList = append(contentList, ")")
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
			if !isValidChar(char) {
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
	if targetKeyWord == "while" {
		return result
	}
	if targetKeyWord == "catch" {
		return result
	}
	tmpIndex -= 1
	if tmpIndex >= 0 && tmpIndex < len(contentList) && contentList[tmpIndex] == "get" {
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
		k := fmt.Sprintf("%s%d", "randomKey_", i)
		v := fmt.Sprintf("%s%d", "randomValue_", i)
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
	if char == "\"" {
		return true
	}
	return false
}
