# 修复markdown文件中的数学表达式
# 添加缺失的LaTeX标签

$refactorDir = "."
$markdownFiles = Get-ChildItem -Path $refactorDir -Filter "*.md" -Recurse

Write-Host "🔍 找到 $($markdownFiles.Count) 个markdown文件"

$fixedCount = 0

foreach ($file in $markdownFiles) {
    try {
        $content = Get-Content -Path $file.FullName -Raw -Encoding UTF8
        $originalContent = $content
        
        # 修复行内数学表达式 $...$ 格式
        $content = $content -replace '(?<!```latex\s*\n)\$([^$]+)\$(?!\s*\n```)', "```latex`n`$1`$`n```"
        
        # 修复块级数学表达式 $$...$$ 格式
        $content = $content -replace '(?<!```latex\s*\n)\$\$([^$]+)\$\$(?!\s*\n```)', "```latex`n````1````n```"
        
        if ($content -ne $originalContent) {
            Set-Content -Path $file.FullName -Value $content -Encoding UTF8
            Write-Host "✅ 修复: $($file.FullName)"
            $fixedCount++
        } else {
            Write-Host "⏭️  跳过: $($file.FullName)"
        }
    }
    catch {
        Write-Host "❌ 错误: $($file.FullName) - $_"
    }
}

Write-Host "`n📊 修复完成: $fixedCount/$($markdownFiles.Count) 个文件" 