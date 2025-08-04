package welcome

import (
	"fmt"
)

const (
	WELCOME_SWITCH_ON         = 1
	WELCOME_SWITCH_OFF        = 0
	WELCOME_DEFAULT_STATE_ON  = 1
	WELCOME_DEFAULT_STATE_OFF = 0
	WELCOME_SAFE_PASS         = 0
	WELCOME_SAFE_FAIL         = 1
	DEFAULT_WELCOME_TITLE     = "欢迎进入腾讯会议网络研讨会！"
	DEFAULT_WELCOME_CONTENT   = "请大家共同维护大会秩序。会议过程中禁止出现违法违规、低俗色情、吸烟酗酒等内容，若发现违规行为请及时向我们反馈。如主办方在会议过程中引导交易，请谨慎判断，注意财产安全，谨防诈骗！"
	WELCOME_TITLE_NAME        = "welcome_title"
	WELCOME_CONTENT_NAME      = "welcome_content"

	ENTERPRISE_EDITION = 2  // 企业版
	EDUCATE_EDITION    = 16 // 教育版

	MEETING_STATE_VALID_ERR_TIPS   = "会议正在进行中，无法修改入会欢迎语信息"
	MEETING_STATE_INVALID_ERR_TIPS = "会议无效，无法修改入会欢迎语信息"

	UPDATE_FAIL_DEFAULT_TIPS = "未保存成功，请修改或稍后再试"

	WelComeDefaultExpire = 30 * 24 * 3600
)

// IsDefaultWelcome 检查是否默认欢迎语
func IsDefaultWelcome(title string, content string) bool {
	return title == DEFAULT_WELCOME_TITLE && content == DEFAULT_WELCOME_CONTENT
}

// IsSwitchOpen 检查开关是否打开
func IsSwitchOpen(switch_state uint32) bool {
	return switch_state == WELCOME_SWITCH_ON
}

// MakeWelcomeInfoKey ...
func MakeWelcomeInfoKey(meetID uint64) string {
	return fmt.Sprintf("webinar_welcome_info_%v", meetID)
}
