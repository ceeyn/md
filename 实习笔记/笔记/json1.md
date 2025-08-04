



<img src="/Users/giffinhao/Library/Application Support/typora-user-images/image-20240802151136806.png" alt="image-20240802151136806" style="zoom:150%;" />

<img src="/Users/giffinhao/Library/Application Support/typora-user-images/image-20240802151305443.png" alt="image-20240802151305443" style="zoom:50%;" />

综合考虑后，因为原有的逻辑是在sendAction方法中根据在七彩石上配置规则的不同的type做分发，使用sourceId会和原来的代码逻辑有所冲突，因此去掉sourceId，在不同的方法中已经在sendaction中确定了是后台还是客户端，因此将client和server作为不同的组，针对client和server根据不同的条件进行加密，因此现在打算编写一个客户端加密方法，根据配置读入的需要加密的参数做加密，另一个是后台加密方法，根据配置读入的需要加密的参数做加密，

**好处**：针对客户端和后台的不同场景，解决了原来进行后台加密时需要先遍历一遍客户端规则消耗的cpu和时间



```
{
  "corp_id": ["1400143280", "1400212767"]
}
```



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
                    "key": "",
                    "fields": []
                },
                "private": {
                    "key": "TXGCbcKey",
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
                    "key": "",
                    "fields": []
                },
                "private": {
                    "key": "GCGbcKey",
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
                    "key": "TXGCbcKey",
                    "fields": ["uid", "meeting_code"]
                },
                "private": {
                    "key": "TXGCbcKey",
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
                "rule": "creator_corp_id in ['1400115281']"
            },
            "action": {
                "public": {
                    "key": "TXGCbcKey",
                    "fields": ["uid", "meeting_code"]
                },
                "private": {
                    "key": "TXGCbcKey",
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
                "rule": "creator_corp_id not in ['1400115281']"
            },
            "action": {
                "public": {
                    "key": "TXGCbcKey",
                    "fields": ["uid", "meeting_code"]
                },
                "private": {
                    "key": "TXGCbcKey",
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
                    "key": "GCGbcKey",
                    "fields": ["uid", "meeting_code"]
                },
                "private": {
                    "key": "GCGbcKey",
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





