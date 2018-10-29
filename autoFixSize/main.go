package main

import (
	"io/ioutil"
	"fmt"
	"strings"
	"sort"
	"os"
	"strconv"
)
type FileDetail struct{
	FileName string
	FilePath string
	FileSize int64
}
type FileDetailPool []*FileDetail

// Len is the number of elements in the collection.
func (self FileDetailPool)Len() int{
	return len(self)
}
// Less reports whether the element with
// index i should sort before the element with index j.
func (self FileDetailPool)Less(i, j int) bool{
	return self[i].FileSize < self[j].FileSize
}
// Swap swaps the elements with indexes i and j.
func (self FileDetailPool) Swap(i, j int){
	self[i], self[j] = self[j], self[i]
}

func main()  {
	if len(os.Args) < 4{
		fmt.Println("err input arg")
		return
	}

	sourceDir := strings.Replace(os.Args[1],"\\","/",-1)
	canFixDir := strings.Replace(os.Args[2],"\\","/",-1)
	targetSize,err := strconv.ParseInt(os.Args[3],10,64)
	if err != nil{
		fmt.Println("error on parser dir size ",os.Args[3])
	}

	ignoreFileList := map[string]int{}
	if len(os.Args) > 4{
		loadIgnoreInfoAtPath(os.Args[4],ignoreFileList)
	}
	fmt.Println("begin cut dir ",sourceDir,canFixDir,targetSize)

	filePool := loadAllFile(sourceDir)

	canFixFilePool,leftSize:= getFileFromTargetDir(filePool,canFixDir,targetSize,ignoreFileList)
	if leftSize <= 0{
		fmt.Println("error file size ",targetSize)
		return
	}

	fmt.Println("total size " ,targetSize, " can fix size " , leftSize,"begin cut size")

	cutFileInSize(leftSize,canFixFilePool)

	fmt.Println("done")
}
func loadIgnoreInfoAtPath(targetPath string, ignoreFileList map[string]int) bool {
	filePool := loadAllFile(targetPath)
	for _,file := range filePool{
		ignoreFileList[file.FileName] = 1
	}
	return true
}
func cutFileInSize(leftSize int64,canFixFilePool FileDetailPool){
	sort.Sort(canFixFilePool)
	var currentSize int64 = 0
	var startIndex = 0
	for index,file := range canFixFilePool{
		currentSize += file.FileSize
		startIndex = index
		if currentSize > leftSize{
			break
		}
	}
	if currentSize <= leftSize{
		return
	}
	// do delete
	size := len(canFixFilePool)
	for i:=startIndex;i<size;i++{
		// do delete this file
		err := os.Remove(canFixFilePool[i].FilePath)
		if err != nil{
			fmt.Println("error " , err)
		}
		fmt.Println("delete file " + canFixFilePool[i].FilePath ,"size ",canFixFilePool[i].FileSize)
	}
}
func getFileFromTargetDir(filePool []*FileDetail, targetDirPath string,totalSize int64,ignoreFileList map[string]int) ([]*FileDetail,int64) {
	var canFixFilePool []*FileDetail
	var tmpSize int64 = 0
	for _,file := range filePool{
		if strings.HasPrefix(file.FilePath,targetDirPath) && !isFileInIgnoreList(file,ignoreFileList){

			canFixFilePool = append(canFixFilePool,file)
		}else{
			tmpSize += file.FileSize
		}
	}
	return canFixFilePool,totalSize-tmpSize
}
func isFileInIgnoreList(fileDetail *FileDetail, ignoreList map[string]int) bool {
	for ignoreFile := range ignoreList{
		if strings.Index(fileDetail.FilePath,ignoreFile) != -1{
			return true
		}
	}
	return false
}
func loadAllFile(dirPath string)[]*FileDetail{
	files,err := ioutil.ReadDir(dirPath)
	if err != nil{
		fmt.Println(err)
		return nil
	}
	var filePool []*FileDetail
	for _,file := range files{
		if file.IsDir(){
			filePool = append(filePool,loadAllFile(dirPath + "/" + file.Name())...)
		}else{
			filePool = append(filePool,&FileDetail{FileName : file.Name(),FilePath:dirPath + "/" + file.Name(),FileSize:file.Size()})
		}
	}
	return filePool
}