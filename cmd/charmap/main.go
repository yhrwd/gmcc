package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gmcc/internal/mcclient"
)

func main() {
	analyzer, err := mcclient.InitializeCharacterMap(".charmap")
	if err != nil {
		fmt.Fprintf(os.Stderr, "初始化字符映射失败: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		listMappings(analyzer)
	case "add":
		if len(os.Args) < 5 {
			fmt.Println("用法: charmap add <unicode> <描述> <替换字符>")
			fmt.Println("示例: charmap add \\uE000 \"方块图标\" \"█\"")
			os.Exit(1)
		}
		unicodeStr := os.Args[2]
		description := os.Args[3]
		replaceWith := os.Args[4]
		addMapping(analyzer, unicodeStr, description, replaceWith)
	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("用法: charmap remove <unicode>")
			fmt.Println("示例: charmap remove \\uE000")
			os.Exit(1)
		}
		unicodeStr := os.Args[2]
		removeMapping(analyzer, unicodeStr)
	case "enable":
		enableReplace(analyzer, true)
	case "disable":
		enableReplace(analyzer, false)
	case "show-unicode":
		if len(os.Args) < 3 {
			showUnicodeInfo(analyzer, true)
		} else {
			enable := os.Args[2] == "true" || os.Args[2] == "on" || os.Args[2] == "1"
			showUnicodeInfo(analyzer, enable)
		}
	case "template":
		outputPath := ""
		if len(os.Args) >= 3 {
			outputPath = os.Args[2]
		}
		generateTemplate(analyzer, outputPath)
	case "interactive":
		interactiveMode(analyzer)
	case "analyze":
		if len(os.Args) < 3 {
			fmt.Println("用法: charmap analyze <文本>")
			os.Exit(1)
		}
		text := strings.Join(os.Args[2:], " ")
		analyzeText(analyzer, text)
	default:
		fmt.Printf("未知命令: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("字符映射管理工具")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  charmap list                              列出所有映射")
	fmt.Println("  charmap add <unicode> <描述> <替换字符>   添加映射")
	fmt.Println("  charmap remove <unicode>                  删除映射")
	fmt.Println("  charmap enable                            启用字符替换")
	fmt.Println("  charmap disable                           禁用字符替换")
	fmt.Println("  charmap show-unicode [on|off]             显示/隐藏Unicode信息")
	fmt.Println("  charmap template [输出文件]               生成映射模板")
	fmt.Println("  charmap interactive                       交互模式")
	fmt.Println("  charmap analyze <文本>                    分析文本中的特殊字符")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  charmap add \\uE000 \"左上角\" \"┌\"")
	fmt.Println("  charmap remove \\uE000")
	fmt.Println("  charmap analyze \"测试文本内容\"")
}

func listMappings(analyzer *mcclient.CharacterAnalyzer) {
	config := analyzer.GetConfig()

	fmt.Printf("字符替换状态: %v\n", config.EnableReplace)
	fmt.Printf("显示Unicode信息: %v\n\n", config.ShowUnicodeInfo)

	if len(config.Mappings) == 0 {
		fmt.Println("暂无字符映射")
		return
	}

	fmt.Println("当前映射:")
	for unicode, mapping := range config.Mappings {
		char := parseUnicodeEscape(unicode)
		fmt.Printf("  %s (%c) -> %q [%s]\n", unicode, char, mapping.ReplaceWith, mapping.Description)
	}
}

func addMapping(analyzer *mcclient.CharacterAnalyzer, unicodeStr, description, replaceWith string) {
	if err := analyzer.AddMapping(unicodeStr, description, replaceWith); err != nil {
		fmt.Fprintf(os.Stderr, "添加映射失败: %v\n", err)
		os.Exit(1)
	}

	char := parseUnicodeEscape(unicodeStr)
	fmt.Printf("已添加映射: %s (%c) -> %q [%s]\n", unicodeStr, char, replaceWith, description)
}

func removeMapping(analyzer *mcclient.CharacterAnalyzer, unicodeStr string) {
	if err := analyzer.RemoveMapping(unicodeStr); err != nil {
		fmt.Fprintf(os.Stderr, "删除映射失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("已删除映射: %s\n", unicodeStr)
}

func enableReplace(analyzer *mcclient.CharacterAnalyzer, enable bool) {
	if err := analyzer.SetEnableReplace(enable); err != nil {
		fmt.Fprintf(os.Stderr, "设置失败: %v\n", err)
		os.Exit(1)
	}

	status := "禁用"
	if enable {
		status = "启用"
	}
	fmt.Printf("已%s字符替换功能\n", status)
}

func showUnicodeInfo(analyzer *mcclient.CharacterAnalyzer, enable bool) {
	if err := analyzer.SetShowUnicodeInfo(enable); err != nil {
		fmt.Fprintf(os.Stderr, "设置失败: %v\n", err)
		os.Exit(1)
	}

	status := "隐藏"
	if enable {
		status = "显示"
	}
	fmt.Printf("已设置为%s Unicode信息\n", status)
}

func generateTemplate(analyzer *mcclient.CharacterAnalyzer, outputPath string) {
	if outputPath == "" {
		outputPath = "charmap_template.json"
	}

	if err := analyzer.GenerateMappingTemplate(outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "生成模板失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("已生成映射模板: %s\n", outputPath)
}

func analyzeText(analyzer *mcclient.CharacterAnalyzer, text string) {
	fmt.Printf("原文: %q\n", text)

	replaced := analyzer.ReplaceText(text)
	fmt.Printf("替换: %q\n", replaced)

	unicodeInfo := analyzer.AnalyzeText(text)
	if unicodeInfo != "" {
		fmt.Println(unicodeInfo)
	} else {
		fmt.Println("未检测到特殊Unicode字符")
	}
}

func interactiveMode(analyzer *mcclient.CharacterAnalyzer) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("字符映射交互模式")
	fmt.Println("输入 'help' 查看命令，'exit' 退出")
	fmt.Println()

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "读取输入失败: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]

		switch command {
		case "exit", "quit":
			fmt.Println("再见！")
			return
		case "help":
			fmt.Println("可用命令:")
			fmt.Println("  list                              列出所有映射")
			fmt.Println("  add <unicode> <描述> <替换字符>   添加映射")
			fmt.Println("  remove <unicode>                  删除映射")
			fmt.Println("  enable                            启用替换")
			fmt.Println("  disable                           禁用替换")
			fmt.Println("  analyze <文本>                    分析文本")
			fmt.Println("  show-unicode [on|off]             显示/隐藏Unicode信息")
			fmt.Println("  help                              显示帮助")
			fmt.Println("  exit                              退出")
		case "list":
			listMappings(analyzer)
		case "add":
			if len(parts) < 4 {
				fmt.Println("用法: add <unicode> <描述> <替换字符>")
				continue
			}
			addMapping(analyzer, parts[1], parts[2], parts[3])
		case "remove":
			if len(parts) < 2 {
				fmt.Println("用法: remove <unicode>")
				continue
			}
			removeMapping(analyzer, parts[1])
		case "enable":
			enableReplace(analyzer, true)
		case "disable":
			enableReplace(analyzer, false)
		case "analyze":
			if len(parts) < 2 {
				fmt.Println("用法: analyze <文本>")
				continue
			}
			text := strings.Join(parts[1:], " ")
			analyzeText(analyzer, text)
		case "show-unicode":
			if len(parts) < 2 {
				showUnicodeInfo(analyzer, true)
			} else {
				enable := parts[1] == "on" || parts[1] == "true" || parts[1] == "1"
				showUnicodeInfo(analyzer, enable)
			}
		default:
			fmt.Printf("未知命令: %s，输入 'help' 查看帮助\n", command)
		}
	}
}

func parseUnicodeEscape(s string) rune {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "\\u") {
		return 0
	}

	hexStr := strings.TrimPrefix(s, "\\u")
	var codePoint uint32
	_, err := fmt.Sscanf(hexStr, "%X", &codePoint)
	if err != nil {
		return 0
	}

	return rune(codePoint)
}

func parseNumber(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
