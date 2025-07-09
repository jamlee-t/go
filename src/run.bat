:: Copyright 2012 The Go Authors. All rights reserved.
:: Use of this source code is governed by a BSD-style
:: license that can be found in the LICENSE file.

@echo off

if exist ..\bin\go.exe goto ok
echo Must run run.bat from Go src directory after installing cmd/go.
goto fail
:ok

:: Keep environment variables within this script
:: unless invoked with --no-local.
if x%1==x--no-local goto nolocal
if x%2==x--no-local goto nolocal
setlocal
:nolocal

set GOBUILDFAIL=0

set GOENV=off
..\bin\go tool dist env > env.bat
if errorlevel 1 goto fail
call .\env.bat
del env.bat

set GOPATH=c:\nonexist-gopath

@REM if x%1==x--no-rebuild goto norebuild
@REM ..\bin\go tool dist test --rebuild
@REM if errorlevel 1 goto fail
@REM goto end

@REM :norebuild
@REM ..\bin\go tool dist test
@REM if errorlevel 1 goto fail
@REM goto end

:fail
set GOBUILDFAIL=1

:end
