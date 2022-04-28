<#

.SYNOPSIS
Checks whether the Build Number env var exists and sets default if not

#>

[CmdletBinding()]
param (

    [string]
    # Set the default build number
    $default = "100.98.99",

    [string]
    [Parameter(Position=0)]
    [AllowNull()]
    [AllowEmptyString()]
    # Allow the build number to be passed in
    $number = ""
)

dir env:

if (![String]::IsNullOrEmpty($number)) {
    $result = $number
} else {

    if ([String]::IsNullOrEmpty($env:BUILDNUMBER)) {
        $result = $default
    } else {
        $result = $env:BUILDNUMBER
    }
}

# output the build number
Write-Output $result.Trim()
