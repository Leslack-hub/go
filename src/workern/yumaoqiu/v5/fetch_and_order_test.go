package main

import (
	"encoding/json"
	"fmt"
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

// TestFetchAndOrderFlow 测试从拉取订单到下单的完整流程
func TestFetchAndOrderFlow(t *testing.T) {
	// 你提供的真实接口响应数据
	responseJSON := `{
  "sysdate": "2026-02-04 08:00:59",
  "cancelPaySendCodeTag": "0",
  "pageInfo": {
    "pageNum": 1,
    "pageSize": 10,
    "size": 2,
    "startRow": 1,
    "endRow": 2,
    "total": 2,
    "pages": 1,
    "list": [
      {
        "acceptDate": "2026-02-04 08:00:04",
        "tradeTicketList": [
          {
            "regularPrice": 0,
            "discount": 0,
            "fullTag": "0",
            "playerNum": 6,
            "couponAmount": 0,
            "ticketNo": "82893387022379",
            "state": "9",
            "fieldId": 20007729,
            "groupTag": "0",
            "priceItem": 2023080900125739,
            "ticketId": 2026020489338702,
            "fieldName": "全民健身中心羽毛球馆（收费日）-8号场",
            "ecardNo": "",
            "ticketSourceType": "0",
            "startSegment": 32,
            "payMoney": 0,
            "venueId": 5003000103,
            "expireDate": "2026-02-05 00:00:00",
            "startTime": "1600",
            "serviceId": 1002,
            "effectDate": "2026-02-05 00:00:00",
            "ticketDrawer": 1000060,
            "ticketType": 0,
            "custName": "15346503765",
            "endSegment": 35,
            "createTime": "2026-02-04 08:00:04",
            "fieldTypeName": "全民健身中心羽毛球馆（收费日）",
            "endTime": "1800",
            "tradeId": 2026020493646400
          }
        ],
        "payTfee": 0,
        "payState": "0",
        "tradeDesc": "02月05日 羽毛球数量：1\n全民健身中心羽毛球馆（收费日）-8号场 16:00-18:00",
        "title": "羽毛球场地预订",
        "tradeStaffId": 1000060,
        "productRes": "",
        "orderState": 2,
        "cancelAllTicketsTag": true,
        "venueId": 5003000103,
        "place": "全民健身中心综合馆",
        "serviceId": 1002,
        "tradeTypeCode": 11,
        "channelId": 11,
        "acceptMonth": "02",
        "centerId": 50030001,
        "netUserId": 2025101002519673,
        "subscribeState": "0",
        "cancelTag": "0",
        "subscribeId": 2026020493646406,
        "priority": 0,
        "serviceName": "羽毛球",
        "venueName": "全民健身中心综合馆",
        "tradeTypeCodeName": "场地预订",
        "expireTime": "2026-02-04 08:05:04",
        "createTime": "2026-02-04 08:00:04",
        "liveTag": "0",
        "tradeId": 2026020493646400
      },
      {
        "acceptDate": "2026-02-04 08:00:02",
        "tradeTicketList": [
          {
            "regularPrice": 0,
            "discount": 0,
            "fullTag": "0",
            "playerNum": 6,
            "couponAmount": 0,
            "ticketNo": "82893386859365",
            "state": "9",
            "fieldId": 20007599,
            "groupTag": "0",
            "priceItem": 2023080900125739,
            "ticketId": 2026020489338685,
            "fieldName": "全民健身中心羽毛球馆（收费日）-2号场",
            "ecardNo": "",
            "ticketSourceType": "0",
            "startSegment": 32,
            "payMoney": 0,
            "venueId": 5003000103,
            "expireDate": "2026-02-05 00:00:00",
            "startTime": "1600",
            "serviceId": 1002,
            "effectDate": "2026-02-05 00:00:00",
            "ticketDrawer": 1000060,
            "ticketType": 0,
            "custName": "15346503765",
            "endSegment": 35,
            "createTime": "2026-02-04 08:00:02",
            "fieldTypeName": "全民健身中心羽毛球馆（收费日）",
            "endTime": "1800",
            "tradeId": 2026020493646357
          }
        ],
        "payTfee": 0,
        "payState": "0",
        "tradeDesc": "02月05日 羽毛球数量：1\n全民健身中心羽毛球馆（收费日）-2号场 16:00-18:00",
        "title": "羽毛球场地预订",
        "tradeStaffId": 1000060,
        "productRes": "",
        "orderState": 2,
        "cancelAllTicketsTag": true,
        "venueId": 5003000103,
        "place": "全民健身中心综合馆",
        "serviceId": 1002,
        "tradeTypeCode": 11,
        "channelId": 11,
        "acceptMonth": "02",
        "centerId": 50030001,
        "netUserId": 2025101002519673,
        "subscribeState": "0",
        "cancelTag": "0",
        "subscribeId": 2026020493646360,
        "priority": 0,
        "serviceName": "羽毛球",
        "venueName": "全民健身中心综合馆",
        "tradeTypeCodeName": "场地预订",
        "expireTime": "2026-02-04 08:05:02",
        "createTime": "2026-02-04 08:00:02",
        "liveTag": "0",
        "tradeId": 2026020493646357
      }
    ],
    "firstPage": 1,
    "prePage": 0,
    "nextPage": 0,
    "lastPage": 1,
    "isFirstPage": true,
    "isLastPage": true,
    "hasPreviousPage": false,
    "hasNextPage": false,
    "navigatePages": 8,
    "navigatepageNums": [
      1
    ]
  },
  "error": 0,
  "message": "ok"
}`

	// 1. 解析订单响应
	var orderResp OrderResponse
	err := json.Unmarshal([]byte(responseJSON), &orderResp)
	if err != nil {
		t.Fatalf("解析订单响应失败: %v", err)
	}

	// 2. 验证响应基本信息
	if orderResp.Error != 0 {
		t.Errorf("订单接口返回错误: %s", orderResp.Message)
	}

	if orderResp.Message != "ok" {
		t.Errorf("期望 message 为 'ok', 实际为: %s", orderResp.Message)
	}

	// 3. 验证分页信息
	if orderResp.PageInfo == nil {
		t.Fatal("PageInfo 为空")
	}

	t.Logf("分页信息: pageNum=%d, pageSize=%d, total=%d",
		orderResp.PageInfo.PageNum,
		orderResp.PageInfo.PageSize,
		orderResp.PageInfo.Total)

	if orderResp.PageInfo.Total != 2 {
		t.Errorf("期望订单总数为 2, 实际为: %d", orderResp.PageInfo.Total)
	}

	// 4. 验证订单列表
	if len(orderResp.PageInfo.List) == 0 {
		t.Fatal("订单列表为空")
	}

	t.Logf("找到 %d 个订单", len(orderResp.PageInfo.List))

	// 5. 提取并验证 TradeId
	var tradeIds []string
	for i, order := range orderResp.PageInfo.List {
		if order.TradeId == nil {
			t.Errorf("订单 %d 的 TradeId 为空", i)
			continue
		}

		var tradeIdStr string
		switch v := order.TradeId.(type) {
		case string:
			tradeIdStr = v
		case float64:
			tradeIdStr = fmt.Sprintf("%.0f", v)
		case int, int32, int64, uint, uint32, uint64:
			tradeIdStr = fmt.Sprintf("%v", v)
		default:
			t.Errorf("订单 %d 的 TradeId 类型未知: %T, 值: %v", i, v, v)
			continue
		}

		if tradeIdStr == "" {
			t.Errorf("订单 %d 的 TradeId 转换后为空", i)
			continue
		}

		tradeIds = append(tradeIds, tradeIdStr)
		t.Logf("订单 %d: TradeId=%s", i+1, tradeIdStr)
	}

	// 6. 验证提取的 TradeId 数量
	expectedTradeIds := []string{"2026020493646400", "2026020493646357"}
	if len(tradeIds) != len(expectedTradeIds) {
		t.Errorf("期望提取 %d 个 TradeId, 实际提取 %d 个", len(expectedTradeIds), len(tradeIds))
	}

	// 7. 验证 TradeId 的值
	for i, expectedId := range expectedTradeIds {
		if i >= len(tradeIds) {
			t.Errorf("缺少第 %d 个 TradeId", i+1)
			continue
		}
		if tradeIds[i] != expectedId {
			t.Errorf("第 %d 个 TradeId 不匹配: 期望=%s, 实际=%s", i+1, expectedId, tradeIds[i])
		}
	}

	// 8. 模拟下单流程（验证订单是否可以被处理）
	t.Log("\n=== 模拟下单流程 ===")
	ecardNos = append(ecardNos, "E5003000100037657")
	openId = "o_9oO5YDtVXIBdWsrliOwuBEXi3M"
	apiSecret = "UzpnKLpYdQJgVBxQFeKyxw=="
	apiVersion = 31
	httpClient = createHTTPClient()
	for _, tradeId := range tradeIds {
		t.Logf("处理订单: TradeId=%s", tradeId)

		// 这里可以添加更多的业务逻辑验证
		payOrdersForUser("2025101002519673", 0, tradeIds)

		// 例如：验证订单状态、支付状态等
		if tradeId == "" {
			t.Error("TradeId 不能为空")
		}
	}

	t.Log("\n=== 测试完成 ===")
	t.Logf("成功解析并处理了 %d 个订单", len(tradeIds))
}
