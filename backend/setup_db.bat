@echo off
echo Setting up TerraPaw Database...

REM Set environment variables for default postgres user
set PGUSER=postgres
set PGPASSWORD=postgres

REM Try to create the database (ignoring error if it already exists)
echo Creating database 'terrapaw' if it doesn't exist...
createdb -U postgres -h localhost terrapaw 2>NUL

if %ERRORLEVEL% EQU 0 (
    echo Database 'terrapaw' created successfully.
) else (
    echo Database 'terrapaw' might already exist or connection failed.
)

REM Verify connection
echo Verifying connection...
psql -U postgres -h localhost -c "\l" >NUL 2>&1

if %ERRORLEVEL% EQU 0 (
    echo Connection successful!
) else (
    echo Connection failed. Please check if your PostgreSQL is running and the password for user 'postgres' is 'postgres'.
    echo If your password is different, please update the .env file in the backend folder.
)

pause
