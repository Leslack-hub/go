package main

import (
	"reflect"
	"testing"
)

func TestFindBestIndices(t *testing.T) {
	// 启用 debug 模式以查看日志
	debugMode = true

	tests := []struct {
		name        string
		targetIdx   int
		segments    []*FieldSegment
		expected    []int
		description string
	}{
		{
			name:      "目标位置和下一个位置都可用",
			targetIdx: 2,
			segments: []*FieldSegment{
				{State: "1"}, // 0
				{State: "1"}, // 1
				{State: "0"}, // 2 - 目标
				{State: "0"}, // 3
				{State: "2"}, // 4
			},
			expected:    []int{2, 3},
			description: "应返回目标位置和下一个位置",
		},
		{
			name:      "目标位置的上一个位置和目标位置都可用",
			targetIdx: 3,
			segments: []*FieldSegment{
				{State: "1"}, // 0
				{State: "2"}, // 1
				{State: "0"}, // 2
				{State: "0"}, // 3 - 目标
				{State: "2"}, // 4
			},
			expected:    []int{2, 3},
			description: "应返回上一个位置和目标位置",
		},
		{
			name:      "目标位置不可用_向后找到两个连续",
			targetIdx: 2,
			segments: []*FieldSegment{
				{State: "1"}, // 0
				{State: "1"}, // 1
				{State: "2"}, // 2 - 目标(不可用)
				{State: "0"}, // 3
				{State: "0"}, // 4
				{State: "2"}, // 5
			},
			expected:    []int{3, 4},
			description: "目标不可用,应向后找到连续的3和4",
		},
		{
			name:      "目标位置不可用_向前找到两个连续",
			targetIdx: 5,
			segments: []*FieldSegment{
				{State: "0"}, // 0
				{State: "0"}, // 1
				{State: "2"}, // 2
				{State: "0"}, // 3
				{State: "0"}, // 4
				{State: "2"}, // 5 - 目标(不可用)
				{State: "2"}, // 6
			},
			expected:    []int{3, 4},
			description: "目标不可用,应向前找到连续的3和4",
		},
		{
			name:      "找不到两个连续_返回单个可用位置",
			targetIdx: 3,
			segments: []*FieldSegment{
				{State: "2"}, // 0
				{State: "0"}, // 1
				{State: "2"}, // 2
				{State: "2"}, // 3 - 目标(不可用)
				{State: "2"}, // 4
				{State: "0"}, // 5
				{State: "2"}, // 6
			},
			expected:    []int{1},
			description: "没有两个连续的,应返回最近的单个位置1",
		},
		{
			name:      "目标位置可用但没有连续的_返回单个",
			targetIdx: 3,
			segments: []*FieldSegment{
				{State: "2"}, // 0
				{State: "2"}, // 1
				{State: "2"}, // 2
				{State: "0"}, // 3 - 目标(可用但不连续)
				{State: "2"}, // 4
				{State: "2"}, // 5
			},
			expected:    []int{3},
			description: "目标可用但没有连续的,应返回目标位置",
		},
		{
			name:      "全部不可用_返回nil",
			targetIdx: 2,
			segments: []*FieldSegment{
				{State: "2"}, // 0
				{State: "2"}, // 1
				{State: "2"}, // 2 - 目标
				{State: "2"}, // 3
				{State: "2"}, // 4
			},
			expected:    nil,
			description: "全部不可用,应返回nil",
		},
		{
			name:      "索引越界_返回nil",
			targetIdx: 10,
			segments: []*FieldSegment{
				{State: "0"}, // 0
				{State: "0"}, // 1
			},
			expected:    nil,
			description: "索引越界,应返回nil",
		},
		{
			name:      "偏移超过5_应该找不到",
			targetIdx: 0,
			segments: []*FieldSegment{
				{State: "2"}, // 0 - 目标
				{State: "2"}, // 1
				{State: "2"}, // 2
				{State: "2"}, // 3
				{State: "2"}, // 4
				{State: "2"}, // 5
				{State: "2"}, // 6
				{State: "0"}, // 7 - 超过±5范围
				{State: "0"}, // 8
			},
			expected:    nil,
			description: "可用位置超过±5范围,应返回nil",
		},
		{
			name:      "边界情况_最后两个位置可用",
			targetIdx: 3,
			segments: []*FieldSegment{
				{State: "2"}, // 0
				{State: "2"}, // 1
				{State: "2"}, // 2
				{State: "2"}, // 3 - 目标
				{State: "0"}, // 4
				{State: "0"}, // 5
			},
			expected:    []int{4, 5},
			description: "最后两个位置可用,应返回4和5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findBestIndices(tt.targetIdx, tt.segments)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("%s\n期望: %v, 实际: %v", tt.description, tt.expected, result)
			}
		})
	}
}
