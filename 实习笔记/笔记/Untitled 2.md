# 需求

1、针对敏感字段加密，腾讯企业敏感字段使用一套密钥，非腾讯企业敏感字段使用另一套密钥

2、此次梳理的敏感字段，均包含在 event_value 里边

3、针对腾讯企业的敏感字段加密，需要确认以下两个信息

- (a) 腾讯企业用户创建会 (event_code = e#create_meeting and corp_id = 1400115281)  **uid，meeting_code和19个敏感字段需要加密**
- (b) 非 e#create_meeting 事件，腾讯企业用户 (corp_id = 1400115281) 且creator_corp_id为腾讯企业， **uid，meeting_code和19个敏感字段需要加密**
- (c) 非 e#create_meeting 事件，腾讯企业用户 (corp_id = 1400115281) 且creator_corp_id**不为腾讯企业**， **uid，meeting_code和19个敏感字段需要加密**
- (d) 非腾讯企业用户 (corp_id ！= 1400115281) ， **uid，meeting_code和19个敏感字段需要加密**

确保所有敏感字段，只做一次加密，避免二次加密

详细敏感字段内容：https://doc.weixin.qq.com/sheet/e3_ABkAXAaEACc0zcSan89RPSM0ca0GC?scode=AJEAIQdfAAozvNWp0CABkAXAaEACc&tab=lih15v



# 原有方案

![image-20240814172443795](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240814172443795.png)



# 现有方案

![image-20240814180018125](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240814180018125.png)

## 主要改造点

#### 一。是否是腾讯企业判断的改造

**将原来写死在代码中的是否是腾讯企业的判断改造为通过七彩石拉取腾讯企业配置进行判断**

**原有代码**

<img src="/Users/giffinhao/Library/Application Support/typora-user-images/image-20240814174239288.png" alt="image-20240814174239288" style="zoom:50%;" />

**现在代码**

<img src="/Users/giffinhao/Library/Application Support/typora-user-images/image-20240814174339756.png" alt="image-20240814174339756" style="zoom:50%;" />

**七彩石中配置**

```json
{
  // 客户端判断是否是腾讯企业的配置，通过corp_id进行判断
  "client": {
    "corp_id": [
      "1400115281",
      "1400143280",
      "1400212767"
    ]
  },
  // 后端判断是否是腾讯企业的配置，通过corp_id和creatorCorpId进行判断，只要一项符合即返回
  "server": {
    "corp_id": [
      "1400115281",
      "1400143280",
      "1400212767"
    ],
    "creatorCorpId": [
      "1400212767"
    ]
  }
}
```

#### 

### 二。根据不同的加密规则加密制定字段的改造

**将原来写死在代码中的加密规则改造为根据govaluate匹配json中的condition后，加密action中的public公参，和private私参**

#### 详细json规则如下

```json
{
  "client": [
    {
      "id": 1,
      "comment": "腾讯企业客户端：私参使用TXGCbcKey加密",
      "isTencent": true,
      "condition": {
        "fields": [],
        "rule": "true"
      },
      "action": {
        "public": {
          "fields": []
        },
        "private": {
          "fields": [
            "meeting_subject", "subject", "nick_name", "nickname", "src_user_name", "tar_user_name",
            "nick_name_0", "old_nick_name_0", "creator_nick_name", "phone", "phone_number", "phone-number", "tar_phone",
            "src_phone", "phone_id_0", "email", "tar_email", "src_email", "rooms_report_issue_user_emails"
          ]
        }
      }
    },
    {
      "id": 2,
      "comment": "非腾讯企业客户端：私参使用GCGbcKey加密",
      "isTencent": false,
      "condition": {
        "fields": [],
        "rule": "true"
      },
      "action": {
        "public": {
          "fields": []
        },
        "private": {
          "fields": [
            "meeting_subject", "subject", "nick_name", "nickname", "src_user_name", "tar_user_name",
            "nick_name_0", "old_nick_name_0", "creator_nick_name", "phone", "phone_number", "phone-number", "tar_phone",
            "src_phone", "phone_id_0", "email", "tar_email", "src_email", "rooms_report_issue_user_emails"
          ]
        }
      }
    }
  ],
  "server": [
    {
      "id": 1,
      "comment": "腾讯企业后台创建会议：公参和私参使用TXGCbcKey加密",
      "isTencent": true,
      "condition": {
        "fields": [
          {
            "name": "event_code",
            "type": "string"
          }
        ],
        "rule": "event_code == 'create_meeting'"
      },
      "action": {
        "public": {
          "fields": ["uid", "meeting_code"]
        },
        "private": {
          "fields": [
            "meeting_subject", "subject", "nick_name", "nickname", "src_user_name", "tar_user_name",
            "nick_name_0", "old_nick_name_0", "creator_nick_name", "phone", "phone_number", "phone-number", "tar_phone",
            "src_phone", "phone_id_0", "email", "tar_email", "src_email", "rooms_report_issue_user_emails"
          ]
        }
      }
    },
    {
      "id": 2,
      "comment": "腾讯企业后台其他事件且creator_corp_id等于指定值：公参和私参使用TXGCbcKey加密",
      "isTencent": true,
      "condition": {
        "fields": [
          {
            "name": "creator_corp_id",
            "type": "string"
          }
        ],
        "rule": "creator_corp_id == '1400115281'"
      },
      "action": {
        "public": {
          "fields": ["uid", "meeting_code"]
        },
        "private": {
          "fields": [
            "creator_uid",
            "meeting_subject", "subject", "nick_name", "nickname", "src_user_name", "tar_user_name",
            "nick_name_0", "old_nick_name_0", "creator_nick_name", "phone", "phone_number", "phone-number", "tar_phone",
            "src_phone", "phone_id_0", "email", "tar_email", "src_email", "rooms_report_issue_user_emails"
          ]
        }
      }
    },
    {
      "id": 3,
      "comment": "腾讯企业后台其他事件且creator_corp_id不等于指定值：公参和私参使用TXGCbcKey加密",
      "isTencent": true,
      "condition": {
        "fields": [
          {
            "name": "creator_corp_id",
            "type": "string"
          }
        ],
        "rule": "creator_corp_id != '1400115281'"
      },
      "action": {
        "public": {
          "fields": ["uid", "meeting_code"]
        },
        "private": {
          "fields": [
            "meeting_subject", "subject", "nick_name", "nickname", "src_user_name", "tar_user_name",
            "nick_name_0", "old_nick_name_0", "creator_nick_name", "phone", "phone_number", "phone-number", "tar_phone",
            "src_phone", "phone_id_0", "email", "tar_email", "src_email", "rooms_report_issue_user_emails"
          ]
        }
      }
    },
    {
      "id": 4,
      "comment": "非腾讯企业后台：公参和私参使用GCGbcKey加密",
      "isTencent": false,
      "condition": {
        "fields": [],
        "rule": "true"
      },
      "action": {
        "public": {
          "fields": ["uid", "meeting_code"]
        },
        "private": {
          "fields": [
            "meeting_subject", "subject", "nick_name", "nickname", "src_user_name", "tar_user_name",
            "nick_name_0", "old_nick_name_0", "creator_nick_name", "phone", "phone_number", "phone-number", "tar_phone",
            "src_phone", "phone_id_0", "email", "tar_email", "src_email", "rooms_report_issue_user_emails"
          ]
        }
      }
    }
  ]
}

```



#### 现在代码如下

```go
// 后台加密方法
func EncryptBackendData(c *TMtxContext, reqMap map[string]string, config *rpc.SecretRuleConfig, isTencent bool) int {
	log.InfoContextf(c.Ctx, "hzp EncryptBackendData before: %v", reqMap)
	for _, ruleSet := range config.Server {
		// 先匹配 isTencent
		if ruleSet.IsTencent != isTencent {
			continue
		}

		// 解析规则字符串，创建可评估表达式
		expression, err := govaluate.NewEvaluableExpression(ruleSet.Condition.Rule)
		if err != nil {
			log.ErrorContextf(c.Ctx, "Error creating expression: %v", err)
			continue
		}

		// 创建参数映射，将 reqMap 中的字段值填入参数映射中
		parameters := make(map[string]interface{})
		for _, field := range ruleSet.Condition.Fields {
			parameters[field.Name] = reqMap[field.Name]
		}

		// 评估表达式，判断是否满足条件
		result, err := expression.Evaluate(parameters)
		if err != nil {
			log.ErrorContextf(c.Ctx, "Error evaluating expression: %v", err)
			continue
		}

		if !result.(bool) {
			log.InfoContextf(c.Ctx, "Rule ID: %d not matched", ruleSet.ID)
			continue
		}

		// 打印匹配的规则集ID
		log.InfoContextf(c.Ctx, "Backend Matched rule ID: %d", ruleSet.ID)

		// 加密公有字段
		if ruleSet.Action.Public.Key != "" {
			for _, field := range ruleSet.Action.Public.Fields {
				EncryptField(c, reqMap, field, ruleSet.Action.Public.Key)
			}
		}

		// 加密私有字段
		if ruleSet.Action.Private.Key != "" {
			for _, field := range ruleSet.Action.Private.Fields {
				EncryptField(c, reqMap, field, ruleSet.Action.Private.Key)
			}
		}
		log.InfoContextf(c.Ctx, "hzp EncryptBackendData after: %v", reqMap)
		// 返回匹配的规则集 ID
		return ruleSet.ID
	}
	return -1
}

// 客户端加密方法
func EncryptClientData(c *TMtxContext, reqMap map[string]string, config *rpc.SecretRuleConfig, isTencent bool) int {
	log.InfoContextf(c.Ctx, "hzp EncryptClientData before: %v", reqMap)
	for _, ruleSet := range config.Client {
		// 先匹配 isTencent
		if ruleSet.IsTencent != isTencent {
			continue
		}

		// 解析规则字符串，创建可评估表达式
		expression, err := govaluate.NewEvaluableExpression(ruleSet.Condition.Rule)
		if err != nil {
			log.ErrorContextf(c.Ctx, "Error creating expression: %v", err)
			continue
		}

		// 创建参数映射，将 reqMap 中的字段值填入参数映射中
		parameters := make(map[string]interface{})
		for _, field := range ruleSet.Condition.Fields {
			parameters[field.Name] = reqMap[field.Name]
		}

		// 评估表达式，判断是否满足条件
		result, err := expression.Evaluate(parameters)
		if err != nil {
			log.ErrorContextf(c.Ctx, "Error evaluating expression: %v", err)
			continue
		}

		if !result.(bool) {
			log.ErrorContextf(c.Ctx, "Rule ID: %d not matched", ruleSet.ID)
			continue
		}

		// 打印匹配的规则集ID
		log.InfoContextf(c.Ctx, "Client Matched rule ID: %d", ruleSet.ID)

		// 加密公有字段
		if ruleSet.Action.Public.Key != "" {
			for _, field := range ruleSet.Action.Public.Fields {
				EncryptField(c, reqMap, field, ruleSet.Action.Public.Key)
			}
		}

		// 加密私有字段
		if ruleSet.Action.Private.Key != "" {
			for _, field := range ruleSet.Action.Private.Fields {
				EncryptField(c, reqMap, field, ruleSet.Action.Private.Key)
			}
		}
		log.InfoContextf(c.Ctx, "hzp EncryptClientData after: %v", reqMap)
		// 返回匹配的规则集 ID
		return ruleSet.ID
	}
	return -1
}
```