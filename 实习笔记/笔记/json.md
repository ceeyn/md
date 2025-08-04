# 一。整体概况

<img src="/Users/giffinhao/Library/Application Support/typora-user-images/image-20240802095538443.png" alt="image-20240802095538443" style="zoom:50%;" />



# 二。策略规则

**// corp_id 企业号数组**
**// "SourceID"`   //来源ID来源id#枚举值：1：客户端 2：后台**
**// event_code 事件**
**// public 公参部分，key，加密密钥，fields：指定字段**
**// private 私参部分，key，加密密钥，fields：指定字段**

```json
[
    // 腾讯企业客户端
    {
        "condition": {
            "SourceID": 1,
            "event_code": "",
            "fields": [
                {
                    "field_name": "corp_id",
                    "field_value": [1400115281]
                }
            ]
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
    // 非腾讯企业客户端
    {
        "condition": {
            "SourceID": 1,
            "event_code": "",
            "fields": [
                {
                    "field_name": "corp_id",
                    "field_value": []
                }
            ]
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
    },
    // 腾讯企业后台创建会议
    {
        "condition": {
            "SourceID": 2,
            "event_code": "create_meeting",
            "fields": [
                {
                    "field_name": "corp_id",
                    "field_value": [1400115281]
                }
            ]
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
    // 腾讯企业后台其他事件且creator_corp_id 等于1400115281
    {
        "condition": {
            "SourceID": 2,
            "event_code": "",
            "fields": [
                {
                    "field_name": "corp_id",
                    "field_value": [1400115281]
                },
                {
                    "field_name": "creator_corp_id",
                    "field_value": [1400115281]
                }
            ]
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
  // 腾讯企业后台其他事件且creator_corp_id 不等于1400115281
    {
        "condition": {
            "SourceID": 2,
            "event_code": "",
            "fields": [
                {
                    "field_name": "corp_id",
                    "field_value": [1400115281]
                },
                {
                    "field_name": "creator_corp_id",
                    "field_value": []
                }
            ]
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
    // 非腾讯企业后台
    {
        "condition": {
            "SourceID": 2,
            "event_code": "",
            "fields": [
                {
                    "field_name": "corp_id",
                    "field_value": []
                }
            ]
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

```

# 三。具体例子

### 测试腾讯企业 创建会议的私参

curl -X POST http://11.160.91.152:39004/api/server-report-standard \
    -H "Content-Type: application/json" \
    -d '{
        "request_time": 1690214725000,
        **"event_code": "e#create_meeting",**
        "comm_params": {
            **"corp_id": "1400115281",**
            "meeting_code": "meeting456",
            "product": "meeting",
            **"source_id": "2",**
            "user": "forestfu",
            **"uid": "123",**
            "ip": "10.64.35.55",
            "creator_corp_id": "1400143280",
            "creator_uid": "1400143280"
        },
        "event_value": {
            "getParams": "{\"corp_ids\":\"1400115280\"}",
            "postParams": "[]",
            "response": "[{\"corp_id\":\"1400115281\",\"status\":1}]",
            "creator_corp_id": "1400143280",
            "creator_uid": "1400143280",
            "uid": "123",
            **"meeting_subject": "aaa",**
            "subject": "test_subject",
            "nick_name": "test_nick_name",
            "nickname": "test_nickname",
            "src_user_name": "test_src_user_name",
            "tar_user_name": "test_tar_user_name",
            "nick_name_0": "test_nick_name_0",
            "old_nick_name_0": "test_old_nick_name_0",
            "creator_nick_name": "test_creator_nick_name",
            "phone": "1234567890",
            "phone_number": "0987654321",
            "phone-number": "1122334455",
            "tar_phone": "2233445566",
            "src_phone": "3344556677",
            "phone_id_0": "4455667788",
            "email": "test@example.com",
            "tar_email": "tar_test@example.com",
            "src_email": "src_test@example.com",
            "rooms_report_issue_user_emails": "issue_report@example.com"
        },
        "client_ip": "192.168.1.1",
        "client_params": {
            "client_param1": "value1"
        },
        "server_params": {
            "server_param1": "value1"
        },
        "version": "1.0"
    }'







### 测试腾讯企业 不创建会议的私参

curl -X POST http://11.160.91.152:39004/api/server-report-standard \
    -H "Content-Type: application/json" \
    -d '{
        "request_time": 1690214725000,
        **"event_code": "security_access_log",**
        "comm_params": {
            **"corp_id": "1400115284",**
            "meeting_code": "meeting456",
            "product": "meeting",
            "ip": "10.64.35.99",
            **"source_id": "2",**
            "user": "forestfu",
            "uid": "123",
            "creator_corp_id": "1400143284",
            "creator_uid": "1400143284"
        },
        "event_value": {
            "getParams": "{\"corp_ids\":\"1400115280\"}",
            "postParams": "[]",
            "response": "[{\"corp_id\":\"1400115281\",\"status\":1}]",
            "creator_corp_id": "1400143284",
            **"creator_uid": "1400143284",**
            "uid": "123",
            "meeting_subject": "aaa",
            "subject": "test_subject",
            "nick_name": "test_nick_name",
            "nickname": "test_nickname",
            "src_user_name": "test_src_user_name",
            "tar_user_name": "test_tar_user_name",
            "nick_name_0": "test_nick_name_0",
            "old_nick_name_0": "test_old_nick_name_0",
            "creator_nick_name": "test_creator_nick_name",
            "phone": "1234567890",
            "phone_number": "0987654321",
            "phone-number": "1122334455",
            "tar_phone": "2233445566",
            "src_phone": "3344556677",
            "phone_id_0": "4455667788",
            "email": "test@example.com",
            "tar_email": "tar_test@example.com",
            "src_email": "src_test@example.com",
            "rooms_report_issue_user_emails": "issue_report@example.com"
        },
        "client_ip": "192.168.1.1",
        "client_params": {
            "client_param1": "value1"
        },
        "server_params": {
            "server_param1": "value1"
        },
        "version": "1.0"
    }'



### 测试非腾讯会议后台

curl -X POST http://11.160.91.152:39004/api/server-report-standard \
     -H "Content-Type: application/json" \
     -d '{
          "request_time": 1721873603000,
          **"event_code": "security_access_log",**
          "comm_params": {
            **"corp_id": "1400115984",**
            "access_obj": ".meet.dashboard.real-time-mra-connect-count|",
            "access_time": "2024-07-25 10:13:23",
            "access_type": "未知",
            "ip": "10.64.35.33",
            "log_msg": "",

​	   **"source_id": "2",**

​            "module_name": "https://cowork-test.console.woa.com/meeting-connector",
​            "platform_name": "cowork",
​            "port": "56424",
​            "security_level": "3",
​            "server_ip": "11.160.91.152",
​            "server_time": "1721873603",
​            "user": "forestfu",
​            "meeting_code": "meeting789",
​            "uid": "123"
​          },
​          "event_value": {
​            "getParams": "[]",
​            "postParams": "{\"corp_id\":\"1400143280\"}",
​            "response": "{\"mraConnectCount\":\"0\",\"nowSipNum\":0,\"mraVideoConnectCount\":\"0\",\"mraAudioConnectCount\":\"0\"}"
​          },
​          "client_ip": "11.141.178.92",
​          "client_params": {},
​          "server_params": {},
​          "version": "1.0"
​        }'



### 测试非腾讯企业客户端私参

curl -X POST http://11.160.91.152:39004/api/server-report-standard \
    -H "Content-Type: application/json" \
    -d '{
        "request_time": 1690214725000,
        **"event_code": "security_access_log",**
        "comm_params": {
            **"corp_id": "1400990999",**
            "meeting_code": "meeting456",
            "product": "meeting",
            "ip": "10.64.35.55",
            **"source_id": "1",**
            "user": "forestfu",
            "uid": "123",
            "creator_corp_id": "1400143280",
            "creator_uid": "1400143280"
        },
        "event_value": {
            "getParams": "{\"corp_ids\":\"1400115280\"}",
            "postParams": "[]",
            "response": "[{\"corp_id\":\"1400115281\",\"status\":1}]",
            "creator_corp_id": "1400143280",
            "creator_uid": "1400143280",
            "uid": "123",
            "meeting_subject": "aaa",
            "subject": "test_subject",
            "nick_name": "test_nick_name",
            "nickname": "test_nickname",
            "src_user_name": "test_src_user_name",
            "tar_user_name": "test_tar_user_name",
            "nick_name_0": "test_nick_name_0",
            "old_nick_name_0": "test_old_nick_name_0",
            "creator_nick_name": "test_creator_nick_name",
            "phone": "1234567890",
            "phone_number": "0987654321",
            "phone-number": "1122334455",
            "tar_phone": "2233445566",
            "src_phone": "3344556677",
            "phone_id_0": "4455667788",
            "email": "test@example.com",
            "tar_email": "tar_test@example.com",
            "src_email": "src_test@example.com",
            "rooms_report_issue_user_emails": "issue_report@example.com"
        },
        "client_ip": "192.168.1.1",
        "client_params": {
            "client_param1": "value1"
        },
        "server_params": {
            "server_param1": "value1"
        },
        "version": "1.0"
    }'







### 测试腾讯企业 客户端私参

curl -X POST http://11.160.91.152:39004/api/server-report-standard \
    -H "Content-Type: application/json" \
    -d '{
        "request_time": 1690214725000,
        **"event_code": "security_access_log",**
        "comm_params": {
            **"corp_id": "1400115281",**
            "meeting_code": "meeting456",
            "product": "meeting",
            **"source_id": "1",**
            "user": "forestfu",
            "uid": "123",
            "creator_corp_id": "1400143280",
            "creator_uid": "1400143280"
        },
        "event_value": {
            "getParams": "{\"corp_ids\":\"1400115280\"}",
            "postParams": "[]",
            "response": "[{\"corp_id\":\"1400115281\",\"status\":1}]",
            "creator_corp_id": "1400143280",
            "creator_uid": "1400143280",
            "uid": "123",
            "meeting_subject": "aaa"
        },
        "client_ip": "192.168.1.1",
        "client_params": {
            "client_param1": "value1"
        },
        "server_params": {
            "server_param1": "value1"
        },
        "version": "1.0"
    }'



```
[
    {
        "id": 1,
        "condition": {
            "fields": [
                {
                    "name": "SourceID",
                    "type": "int"
                },
                {
                    "name": "corp_id",
                    "type": "string"
                }
            ],
            "rule": "SourceID == 1 && (corp_id == '1400115281' || corp_id == '1400115280')"
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
        "condition": {
            "fields": [
                {
                    "name": "SourceID",
                    "type": "int"
                },
                {
                    "name": "corp_id",
                    "type": "string"
                }
            ],
            "rule": "SourceID == 1 && corp_id != '1400115281' && corp_id != '1400115280'"
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
    },
    {
        "id": 3,
        "condition": {
            "fields": [
                {
                    "name": "SourceID",
                    "type": "int"
                },
                {
                    "name": "event_code",
                    "type": "string"
                },
                {
                    "name": "corp_id",
                    "type": "string"
                }
            ],
            "rule": "SourceID == 2 && event_code == 'create_meeting' && (corp_id == '1400115281' || corp_id == '1400115280')"
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
        "condition": {
            "fields": [
                {
                    "name": "SourceID",
                    "type": "int"
                },
                {
                    "name": "corp_id",
                    "type": "string"
                },
                {
                    "name": "creator_corp_id",
                    "type": "string"
                }
            ],
            "rule": "SourceID == 2 && (corp_id == '1400115281' || corp_id == '1400115280') && creator_corp_id == 'envCorpId'"
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
        "id": 5,
        "condition": {
            "fields": [
                {
                    "name": "SourceID",
                    "type": "int"
                },
                {
                    "name": "corp_id",
                    "type": "string"
                },
                {
                    "name": "creator_corp_id",
                    "type": "string"
                }
            ],
            "rule": "SourceID == 2 && (corp_id == '1400115281' || corp_id == '1400115280') && creator_corp_id != 'envCorpId'"
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
        "id": 6,
        "condition": {
            "fields": [
                {
                    "name": "SourceID",
                    "type": "int"
                },
                {
                    "name": "corp_id",
                    "type": "string"
                }
            ],
            "rule": "SourceID == 2 && corp_id != '1400115281' && corp_id != '1400115280'"
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

```





```json
[
    // 腾讯企业客户端
    {
        "id": 1,
        "condition": {
            "fields": [
                {
                    "name": "corp_id",
                    "type": "int"
                }
            ],
            "rule": "corp_id in [1400115281, 1400115282]"
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
    // 非腾讯企业客户端
    {
        "id": 2,
        "condition": {
            "fields": [
                {
                    "name": "corp_id",
                    "type": "int"
                }
            ],
            "rule": "corp_id not in [1400115281, 1400115282]"
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
]

```





```json
[// 腾讯企业后台创建会议
    {
        "id": 1,
        "condition": {
            "fields": [
                {
                    "name": "corp_id",
                    "type": "int"
                }
            ],
            "rule": "corp_id in [1400115281] && event_code == 'create_meeting'"
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
    // 腾讯企业后台其他事件且creator_corp_id 等于指定值
    {
        "id": 2,
        "condition": {
            "fields": [
                {
                    "name": "corp_id",
                    "type": "int"
                },
                {
                    "name": "creator_corp_id",
                    "type": "int"
                }
            ],
            "rule": "corp_id in [1400115281] && creator_corp_id in [1400115281]"
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
    // 腾讯企业后台其他事件且creator_corp_id 不等于指定值
    {
        "id": 3,
        "condition": {
            "fields": [
                {
                    "name": "corp_id",
                    "type": "int"
                },
                {
                    "name": "creator_corp_id",
                    "type": "int"
                }
            ],
            "rule": "corp_id in [1400115281] && creator_corp_id not in [1400115281]"
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
    // 非腾讯企业后台
    {
        "id": 4,
        "condition": {
            "fields": [
                {
                    "name": "corp_id",
                    "type": "int"
                }
            ],
            "rule": "corp_id not in [1400115281]"
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
```

