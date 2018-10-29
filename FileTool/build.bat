set GOARCH=amd64
set GOOS=windows
set CURR=%cd%
cd ../../../../../

set GOPATH=%cd%
cd %CURR%

go build -o fileTool.exe github.com/Blizzardx/simpleTool/FileTool

@IF %ERRORLEVEL% NEQ 0 pause
