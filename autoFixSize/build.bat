set GOARCH=amd64
set GOOS=windows
set CURR=%cd%
cd ../../../../../

set GOPATH=%cd%
cd %CURR%

go build -o autoFixSize.exe github.com/Blizzardx/simpleTool/autoFixSize

@IF %ERRORLEVEL% NEQ 0 pause
