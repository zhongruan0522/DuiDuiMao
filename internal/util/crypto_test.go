package util

import (
	"testing"
)

func TestDoubleEncode(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "普通字符串",
			input: "hello world",
		},
		{
			name:  "中文字符串",
			input: "你好世界",
		},
		{
			name:  "空字符串",
			input: "",
		},
		{
			name:  "特殊字符",
			input: "!@#$%^&*()",
		},
		{
			name:  "JSON字符串",
			input: `{"name":"test","value":123}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := DoubleEncode(tt.input)
			if encoded == "" && tt.input != "" {
				t.Errorf("DoubleEncode() 返回空字符串，输入: %s", tt.input)
			}
		})
	}
}

func TestDoubleDecode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "普通字符串",
			input:   "hello world",
			wantErr: false,
		},
		{
			name:    "中文字符串",
			input:   "你好世界",
			wantErr: false,
		},
		{
			name:    "空字符串",
			input:   "",
			wantErr: false,
		},
		{
			name:    "特殊字符",
			input:   "!@#$%^&*()",
			wantErr: false,
		},
		{
			name:    "JSON字符串",
			input:   `{"name":"test","value":123}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := DoubleEncode(tt.input)
			decoded, err := DoubleDecode(encoded)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoubleDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if decoded != tt.input {
				t.Errorf("DoubleDecode() = %v, want %v", decoded, tt.input)
			}
		})
	}
}

func TestDoubleDecodeInvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "无效的Base64字符串",
			input:   "这不是一个有效的Base64字符串!!!",
			wantErr: true,
		},
		{
			name:    "只有一层编码的字符串",
			input:   "aGVsbG8gd29ybGQ=", // "hello world" 的单次 Base64 编码
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DoubleDecode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoubleDecode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// 基准测试
func BenchmarkDoubleEncode(b *testing.B) {
	testData := "这是一个测试字符串用于基准测试"
	for i := 0; i < b.N; i++ {
		DoubleEncode(testData)
	}
}

func BenchmarkDoubleDecode(b *testing.B) {
	testData := "这是一个测试字符串用于基准测试"
	encoded := DoubleEncode(testData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DoubleDecode(encoded)
	}
}
