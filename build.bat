REM https://github.com/tc-hib/go-winres
go-winres make

REM https://stackoverflow.com/questions/3861634/how-to-reduce-go-compiled-file-size
go build -ldflags="-s -w"