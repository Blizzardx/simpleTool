set GOARCH=amd64
set GOOS=windows
set CURR=%cd%
cd ../../../../../

set GOPATH=%cd%
cd %CURR%

go build -o autoFixTs.exe github.com/Blizzardx/simpleTool/autoFixTs

@IF %ERRORLEVEL% NEQ 0 pause
