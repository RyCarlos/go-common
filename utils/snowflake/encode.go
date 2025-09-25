// Package snowflake
//
//	@file		: token.go
//	@author		: Carlos
//	@contact	: 534994749@qq.com
//	@time		: 2025/5/6 12:17
//
// -------------------------------------------
package snowflake

import (
	"strings"
)

func IdToCode(Id int64) string {
	// 1. 使用改进的混淆算法
	mixed := betterShuffle(uint64(Id))

	// 2. 转换为base62
	code := toBase62(mixed)

	// 3. 添加基于完整Id的校验码
	checksum := calculateChecksum(Id)
	code += string(base62Chars[checksum])

	// 4. 格式化为固定长度
	return formatCode(code, 4, 4)
}

// 改进的混淆算法
func betterShuffle(num uint64) uint64 {
	// 乘法混淆（使用大质数）
	num *= 0x9E3779B97F4A7C15

	// 位旋转
	num = (num << 32) | (num >> 32)

	// XOR 混淆
	num ^= 0xAAAAAAAAAAAAAAAA

	return num
}

// 转换为base62
func toBase62(num uint64) string {
	const base = 62
	if num == 0 {
		return string(base62Chars[0])
	}

	var builder strings.Builder
	for num > 0 {
		remainder := num % base
		builder.WriteByte(base62Chars[remainder])
		num /= base
	}

	// 反转字符串
	runes := []rune(builder.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// 格式化代码
func formatCode(code string, part1Len, part2Len int) string {
	totalLen := part1Len + part2Len
	if len(code) < totalLen {
		// 左侧填充0
		code = strings.Repeat(string(base62Chars[0]), totalLen-len(code)) + code
	} else if len(code) > totalLen {
		// 保留最重要的部分
		code = code[:totalLen]
	}

	return code[:part1Len] + "-" + code[part1Len:part1Len+part2Len]
}

func calculateChecksum(id int64) int {
	sum := 0
	temp := id
	for temp != 0 {
		sum += int(temp % 62)
		temp /= 62
	}
	return sum % 62
}
