
@echo OFF
SETLOCAL ENABLEEXTENSIONS
SET "script_name=%~n0"
SET "script_path=%~0"
SET "script_dir=%~dp0"
rem # to avoid invalid directory name message calling %script_dir%\config.bat
cd %script_dir%
call config.bat
cd ..
set project_dir=%cd%

set module_name=%REPO_HOST%/%REPO_OWNER%/%REPO_NAME%
set bin_dir=%project_dir%\bin

echo script_name   %script_name%
echo script_path   %script_path%
echo script_dir    %script_dir%
echo project_dir   %project_dir%
echo module_name   %module_name%
echo bin_dir       %bin_dir%

IF DEFINED SDLC_GO_VENDOR (
    echo Using Go vendor
    set GOPROXY=off
)

cd %project_dir%

call go test -v -cover ./...
