@echo off

set ARGC=0
for %%x in (%*) do set /A ARGC+=1

if not %ARGC% == 1 (
	echo "usage: %0 <install | clean>"
	goto over
)

set GOPATH=%cd%
set COMMAND=%1%

if %COMMAND% == install (	
	go install microiot.com/center
) else if %COMMAND% == clean (
	go clean -i -n -x microiot.com/center
)

:over
