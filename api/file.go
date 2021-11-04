package api

import (
	"text/template"
)

type file struct {
	Path          string
	Source        string
	Name          string
	Size          uint64
	DirName       string
	Exists        bool
	CreationTime  string
	LastWriteTime string
}

type createOrUpdateFileArgs struct {
	Source string
	Path   string
	//FileJson    string
}

var createOrUpdateFileTemplate = template.Must(template.New("CreateOrUpdateFile").Parse(`
$ErrorActionPreference = 'Stop'

$source='{{.Source}}'
$path='{{.Path}}'

function Get-FileFromUri {
    param(
        [Parameter(Mandatory = $true, Position = 0, ValueFromPipeline = $true, ValueFromPipelineByPropertyName = $true)]
        [string]
        [Alias('Uri')]
        $Url,
        [Parameter(Mandatory = $false, Position = 1)]
        [string]
        [Alias('Folder')]
        $FolderPath
    )
    process {
        $req = [System.Net.HttpWebRequest]::Create($Url)
        $req.Method = "HEAD"
        $response = $req.GetResponse()
        $fUri = $response.ResponseUri
        $filename = [System.IO.Path]::GetFileName($fUri.LocalPath);
        $response.Close()

        $destination = (Get-Item -Path ".\" -Verbose).FullName
        if ($FolderPath) { $destination = $FolderPath }
        if ($destination.EndsWith('\')) {
            $destination += $filename
        }
        else {
            $destination += '\' + $filename
        }
        $webclient = New-Object System.Net.webclient
        $webclient.downloadfile($fUri.AbsoluteUri, $destination)
    }
}

function Test-Uri {
    param(
        [Parameter(Mandatory = $true, Position = 0, ValueFromPipeline = $true, ValueFromPipelineByPropertyName = $true)]
        [string]
        [Alias('Uri')]
        $Url
    )
    process {
        $testUri = $Url -as [System.URI]
        $null -ne $testUri.AbsoluteURI -and $testUri.Scheme -match '[http|https]' -and ($testUri.ToString().ToLower().StartsWith("http://") -or $testUri.ToString().ToLower().StartsWith("https://"))
    }
}

if (!(Test-Path -Path $path)) {
    $pathDirectory = [System.IO.Path]::GetDirectoryName($path)
	$pathFilename = [System.IO.Path]::GetFileName($path)

    if (!(Test-Path $pathDirectory)) {
        New-Item -ItemType Directory -Force -Path $pathDirectory
    }

	Push-Location $pathDirectory
	
	if (Test-Uri -Url $source) {
		Get-FileFromUri -Url $source -FolderPath $pathDirectory
	}
	else {
		Copy-Item $source "$pathDirectory\$pathFilename" -Force
	}

	Pop-Location
}
`))

func (c *HypervClient) CreateOrUpdateFile(path string, source string) (err error) {
	err = c.runFireAndForgetScript(createOrUpdateFileTemplate, createOrUpdateFileArgs{
		Source: source,
		Path:   path,
		//FileJson:    string(fileJson),
	})

	return err
}

type getFileArgs struct {
	Path string
}

var getFileTemplate = template.Must(template.New("GetFile").Parse(`
$ErrorActionPreference = 'Stop'
$path='{{.Path}}'

$fileObject = $null
if (Test-Path $path) {
	$fileObject = Get-ChildItem -path $path | %{ @{
		Name=$_.Name;
		Size=$_.Length;
		DirName=$_.DirectoryName;
		Exists=$_.Exists;
		CreationTime=$_.CreationTime;
		LastWriteTime=$_.LastWriteTime;
	}}
}

if ($fileObject){
	$file = ConvertTo-Json -InputObject $fileObject
	$file
} else {
	'{"Exists": false}'
}
`))

func (c *HypervClient) GetFile(path string) (result file, err error) {
	err = c.runScriptWithResult(getFileTemplate, getFileArgs{
		Path: path,
	}, &result)

	return result, err
}

type deleteFileArgs struct {
	Path string
}

var deleteFileTemplate = template.Must(template.New("DeleteFile").Parse(`
$ErrorActionPreference = 'Stop'

$targetDirectory = (split-path '{{.Path}}' -Parent)
$targetName = (split-path '{{.Path}}' -Leaf)
$targetName = $targetName.Substring(0,$targetName.LastIndexOf('.')).split('\')[-1]

Get-ChildItem -Path $targetDirectory |?{$_.BaseName.StartsWith($targetName)} | %{
	Remove-Item $_.FullName -Force
}
`))

func (c *HypervClient) DeleteFile(path string) (err error) {
	err = c.runFireAndForgetScript(deleteFileTemplate, deleteFileArgs{
		Path: path,
	})

	return err
}
