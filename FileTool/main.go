package main

import (
	"io/ioutil"
	"log"
	"os"
	"io"
	"fmt"
	"github.com/Blizzardx/ConfigProtocol/common"
)
type FileInfo struct{
	Name string
	FullName string
	Dir string
}
func main() {
	if len(os.Args) < 4 {
		log.Fatalln("wrong input args ,need inputdir ,outputdir,optionType(collect,return)")
		return
	}
	inputDir := os.Args[1]
	outputDir := os.Args[2]
	optionType := os.Args[3]

	inputDir = common.FormatePath(inputDir)
	outputDir = common.FormatePath(outputDir)

	if optionType == "collect"{
		collectFileToTargetDir(inputDir,outputDir)

	}else if optionType == "return"{
		returnFileToSourceDir(inputDir,outputDir)
	}
}
func collectFileToTargetDir(inputDir ,outputDir string){
	fileList := getAllFileFromDir(inputDir)
	for _,fileInfo := range fileList{
		sourcePath := fileInfo.FullName
		targetPath := outputDir + "/" + fileInfo.Name
		_,err := CopyFile(targetPath,sourcePath)
		if err != nil{
			log.Fatalln(err)
			return
		}
	}
	fmt.Println("succeed")
}
func returnFileToSourceDir(inputDir ,outputDir string){
	var searchMap = map[string]*FileInfo{}
	targetFileList := getAllFileFromDir(outputDir)
	for _,fileInfo := range targetFileList{
		searchMap[fileInfo.Name] = fileInfo
	}

	fileList := getAllFileFromDir(inputDir)
	for _,fileInfo := range fileList{
		if tmpDir,ok:= searchMap[fileInfo.Name];ok{
			sourcePath := fileInfo.FullName
			targetPath := tmpDir.Dir + "/" + fileInfo.Name
			_,err := CopyFile(sourcePath,targetPath)
			if err != nil{
				log.Fatalln(err)
				return
			}
		}else{
			log.Fatalln("cant' find file by name " + fileInfo.Name + " at target " + outputDir)
			return
		}


	}
	fmt.Println("succeed")
}

func getAllFileFromDir(path string)[]*FileInfo{
	fileList,err := ioutil.ReadDir(path)
	if err != nil{
		log.Fatalln(err)
		return nil
	}
	var resultFileList []*FileInfo = nil
	for _,file := range fileList{
		if file.IsDir(){
			resultFileList = append(resultFileList,getAllFileFromDir(path + "/" + file.Name())...)
		}else{
			resultFileList = append(resultFileList,&FileInfo{Name:file.Name(),FullName:path+"/" + file.Name(),Dir:path})
		}
	}
	return resultFileList
}
func CopyFile(dstName, srcName string) (written int64, err error) {
	_,tmpErr := os.Stat(dstName)
	if tmpErr == nil{
		os.Remove(dstName)
	}
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}