package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

func LoadCfg(path string) *Config {
	if path == "" {
		panic("配置文件路径为空")
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("配置文件 %s 不存在，生成默认配置\n", path)
			fmt.Printf("请修改配置文件后重新运行程序\n")
			return loadDefaultCfg()
		}
		fmt.Printf("打开配置文件时出错: %v\n", err)
		return nil
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Printf("读取配置文件时出错: %v\n", err)
		return nil
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Printf("解析配置文件时出错: %v\n", err)
		return nil
	}

	// 验证配置是否包含空值
	if err := cfg.Validate(); err != nil {
		fmt.Printf("配置验证失败: %v\n", err)
		fmt.Println("请修改配置文件后重新运行程序")
		os.Exit(1)
		return nil
	}

	return &cfg
}

func loadDefaultCfg() *Config {
	defaultConfig := Config{
		Account: AccountConfig{
			PlayerID: "your_player_id_here",
		},
		Server: ServerConfig{
			Address: "",
		},
		Log: LogConfig{
			LogDir:     "logs",
			MaxSize:    10,
			Debug:      false,
			EnableFile: true,
		},
	}

	data, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		fmt.Printf("生成默认配置时出错: %v\n", err)
		return nil
	}

	// 获取程序运行目录
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("无法获取程序路径: %v\n", err)
		return nil
	}
	execDir := filepath.Dir(execPath)
	configPath := filepath.Join(execDir, "config.yaml")

	// 写入默认配置文件
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		fmt.Printf("写入默认配置文件时出错: %v\n", err)
		return nil
	}

	fmt.Printf("已在程序运行目录生成默认配置文件: %s\n", configPath)
	fmt.Println("请修改配置文件后重新运行程序")
	os.Exit(1)
	return nil
}

// Validate 检查配置是否有效（不包含关键空值）
func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("配置为空")
	}

	var empties []string

	var validate func(v reflect.Value, path string)
	validate = func(v reflect.Value, path string) {
		if !v.IsValid() {
			empties = append(empties, path)
			return
		}
		switch v.Kind() {
		case reflect.Ptr:
			if v.IsNil() {
				empties = append(empties, path)
				return
			}
			validate(v.Elem(), path)
		case reflect.Struct:
			t := v.Type()
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				// skip unexported fields
				if f.PkgPath != "" {
					continue
				}
				yamlTag := f.Tag.Get("yaml")
				name := strings.SplitN(yamlTag, ",", 2)[0]
				if name == "" || name == "-" {
					name = strings.ToLower(f.Name)
				}
				subPath := name
				if path != "" {
					subPath = path + "." + name
				}
				validate(v.Field(i), subPath)
			}
		case reflect.String:
			if strings.TrimSpace(v.String()) == "" {
				empties = append(empties, path)
			}
		case reflect.Slice, reflect.Array, reflect.Map:
			if v.Len() == 0 {
				empties = append(empties, path)
			}
		case reflect.Bool:
			// bool 默认为 false，通常不视为空，跳过检查
		default:
			// 对于数字等其他类型，使用 IsZero 判定是否为空
			if v.IsZero() {
				empties = append(empties, path)
			}
		}
	}

	validate(reflect.ValueOf(c).Elem(), "")

	if len(empties) > 0 {
		return fmt.Errorf("以下配置项不能为空: %s", strings.Join(empties, ", "))
	}
	return nil
}
