# 批量修正数学表达式格式的PowerShell脚本

# 获取所有markdown文件
$files = Get-ChildItem -Recurse -Filter "*.md"

foreach ($file in $files) {
    $content = Get-Content $file.FullName -Raw -Encoding UTF8
    
    # 修正数学表达式格式
    $modified = $false
    
    # 修正 \text{} 格式
    if ($content -match "\\text\{") {
        $content = $content -replace '\\text\{([^}]*)\}', '\\text{$1}'
        $modified = $true
    }
    
    # 修正未正确包围的数学表达式
    if ($content -match '(?<!\$)[^$]*\\[a-zA-Z]+\{[^}]*\}[^$]*(?!\$)') {
        # 这里需要更复杂的正则表达式处理
        $modified = $true
    }
    
    # 修正表格中的数学表达式
    if ($content -match '&.*\\text\{') {
        $content = $content -replace '& ([^&]*\\text\{[^}]*\}[^&]*) &', '& $\1$ &'
        $modified = $true
    }
    
    # 如果内容被修改，写回文件
    if ($modified) {
        Set-Content $file.FullName $content -Encoding UTF8
        Write-Host "Fixed: $($file.FullName)"
    }
}

Write-Host "Math expression format fixing completed!" 