## sql注入可能语句

```go
querySql := fmt.Sprintf(
    "select count(*) from label_load_record where label_id = '%s'  and ts >= %d", labelId, ts)
```



`fmt.Sprintf` 是格式化字符串的函数，它**会将传入的参数直接插入到字符串中，而不对这些参数进行任何转义或安全检查**。因此，使用 `fmt.Sprintf` 生成 SQL 查询时，如果传入的参数包含恶意的 SQL 代码，会直接成为最终 SQL 查询的一部分，导致 SQL 注入漏洞。

 `fmt.Sprintf` 防范不了 SQL 注入的详细解释：



### 1. 基本 SQL 注入

攻击者在 `labelId` 参数中注入恶意代码，如下所示：

#### 恶意输入：

```sql
labelId = "1' OR '1'='1"
```

#### 拼接后的 SQL 查询：

```sql
SELECT count(*) FROM label_load_record WHERE label_id = '1' OR '1'='1' AND ts >= 1625097600
```

这样会导致查询返回所有记录，而不是特定的 `labelId`。



### 2. 联合查询注入

攻击者在 `labelId` 参数中注入联合查询，试图获取其他表的信息：

#### 恶意输入：

```sql
labelId = "1'; SELECT * FROM users; --"
```

#### 拼接后的 SQL 查询：

```sql
SELECT count(*) FROM label_load_record WHERE label_id = '1'; SELECT * FROM users; --' AND ts >= 1625097600

```

这样会导致执行两个查询，第一个查询正常执行，第二个查询获取 `users` 表的信息。







### 为什么参数化查询可以防范 SQL 注入？

参数化查询的关键在于**数据库驱动程序会将参数与 SQL 语句分开处理。它不会将参数视为 SQL 代码的一部分，而是作为数据进行处理。这意味着任何参数中的特殊字符（如单引号 `'`）都不会被解释为 SQL 语法的一部分，而是作为纯文本进行处理。数据库驱动程序会自动对参数进行转义处理，确保它们不会被解释为 SQL 代码的一部分。**例如：

### 防范示例

输入

```sql
labelId = "1' OR '1'='1"
```

在使用参数化查询时，SQL 查询会被处理成如下形式：

```sql
SELECT count(*) FROM label_load_record WHERE label_id = '1'' OR ''1''=''1' AND ts >= 1625112000;
```





1. **检查当前分支**：首先检查你当前所在的分支。

   ```
   sh
   复制代码
   git branch
   ```

2. **创建新分支**：在本地仓库中创建一个新分支。

   ```
   sh
   复制代码
   git checkout -b <new-branch-name>
   ```

3. **添加文件**：将你修改或新增的文件添加到暂存区。

   ```
   sh
   复制代码
   git add .
   ```

4. **提交更改**：将添加到暂存区的文件提交到本地仓库。

   ```
   sh
   复制代码
   git commit -m "Your commit message"
   ```

5. **推送新分支到远程仓库**：将新创建的分支推送到远程仓库，并且在远程仓库中创建该分支。

   ```
   sh
   复制代码
   git push origin <new-branch-name>
   ```

6. **验证远程分支**：确保新分支已经推送到远程仓库。

   ```
   sh
   复制代码
   git branch -r
   ```



```
// QueryTaskTodayRecord ...
func (s *LabelLoadSvrServiceImpl) QueryTaskTodayRecord(ctx context.Context, taskId, labelId, cosUrl string) error {
	// 获取今天凌晨的时间戳
	t := time.Now()
	ts := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()

	// 查询某个标签id今天的任务
	var recordCnt int
	querySql := "select count(*) from label_load_record where label_id = ? and ts >= ?"
	log.InfoContextf(ctx, "Executing query: %s with params: labelId=%s, ts=%d", querySql, labelId, ts)
	err := s.DB.QueryRow(ctx, []interface{}{&recordCnt}, querySql, labelId, ts)
	if err != nil {
		log.ErrorContextf(ctx, "query db failed, labelId:%s, err:%+v", labelId, err)
		return errors.New("query db failed")
	}
	log.InfoContextf(ctx, "Query result: %d", recordCnt)

	if recordCnt >= 3 {
		return errors.New(fmt.Sprintf("labelId request too frequency, id:%s", labelId))
	} else {
		// 没有就插入记录
		insertSql := "insert into label_load_record (task_id, label_id, cos_url, ts) values (?,?,?,?)"
		log.InfoContextf(ctx, "Executing insert: %s with params: taskId=%s, labelId=%s, cosUrl=%s, ts=%d", insertSql, taskId, labelId, cosUrl, t.Unix())
		_, insertErr := s.DB.Exec(ctx, insertSql, taskId, labelId, cosUrl, t.Unix())
		if insertErr != nil {
			log.ErrorContextf(ctx, "insert record failed,labelId:%s, err:%+v", labelId, insertErr)
			return errors.New("insert db failed")
		}
	}
	log.InfoContextf(ctx, "QueryTaskTodayRecord ok, labelId: %s", labelId)
	return nil
}

```



### . 确保本地 `master` 分支是最新的

首先，确保你的本地 `master` 分支是最新的。可以通过以下命令拉取最新的 `master` 分支代码：

```
sh
复制代码
git checkout master  # 切换到 master 分支
git pull origin master  # 拉取远程 master 分支的最新代码
```

### 2. 创建并切换到新分支

创建一个新的本地分支，并在该分支上进行修改和提交：

```
sh
复制代码
git checkout -b <新分支名>
```

这将创建一个名为 `<新分支名>` 的新分支，并立即切换到这个分支上。例如：

```
sh
复制代码
git checkout -b my-feature-branch
```

### 3. 将新分支推送到远程仓库

一旦在新分支上进行了修改并且做了一些提交，你可以将这个新分支推送到远程仓库：

```
sh
复制代码
git push origin <新分支名>
```

如果是第一次推送该分支，需要设置远程跟踪：

```
sh
复制代码
git push -u origin <新分支名>
```

例如：

```
sh
复制代码
git push origin my-feature-branch
```







### 区别总结

- **`git fetch`**: 只获取远程仓库的更新，但不合并。这使得你可以在本地查看远程分支的最新状态，并选择何时进行合并操作。
- **`git merge`**: 将一个分支的更改合并到当前分支。这通常是在你使用 `git fetch` 后进行的操作，以便将远程更改合并到本地分支。
- **`git pull`**: 相当于先执行 `git fetch`，然后再执行 `git merge`，自动将远程更改合并到当前分支。

### 使用场景

- **使用 `git fetch` + `git merge`**: 适用于需要在合并之前检查远程分支变化的情况。这样可以确保你在合并之前了解所有的更改，并可以在合并之前处理任何潜在的冲突。

  ```
  bash
  复制代码
  git fetch origin
  git merge origin/dev
  ```

- **使用 `git pull`**: 适用于想要快速将远程更改合并到本地分支的情况。当你对远程分支的更改有信心，且希望自动处理合并时，可以使用 `git pull`。

  ```
  bash
  复制代码
  git pull origin dev
  ```

### 例子

```
bash
复制代码
# fetch 例子
git fetch origin

# merge 例子
git merge origin/dev

# pull 例子
git pull origin dev
```

### pull

当前在哪个本地分支，远程代码就合并到这个分支



### push

要将本地的 `master` 分支更改提交到远程的 `dev/encrypt_sensitive_info` 分支，可以按照以下步骤进行：

1. **确保当前在 `master` 分支**： 确保你在本地的 `master` 分支上，并且已经完成了所有的更改和提交。

   ```
   bash
   复制代码
   git checkout master
   ```

2. **将本地 `master` 分支的更改推送到远程的 `dev/encrypt_sensitive_info` 分支**： 使用 `git push` 命令将本地 `master` 分支的更改推送到远程的 `dev/encrypt_sensitive_info` 分支。

   ```
   bash
   复制代码
   git push origin master:dev/encrypt_sensitive_info
   ```





// QueryTaskTodayRecord ...
func (s *LabelLoadSvrServiceImpl) QueryTaskTodayRecord(ctx context.Context, taskId, labelId, cosUrl string) error {
	// 获取今天凌晨的时间戳
	t := time.Now()
	ts := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()

	// 构建 SQL 查询语句
	querySql := "select count(*) from label_load_record where label_id = ? and ts >= ?"
	log.InfoContextf(ctx, "Executing SQL query: %s, labelId: %s, ts: %d", querySql, labelId, ts)
	
	// 查询某个标签id今天的任务
	var recordCnt int
	dest := []interface{}{&recordCnt}
	err := s.DB.QueryRow(ctx, dest, querySql, labelId, ts)
	if err != nil {
		log.ErrorContextf(ctx, "query db failed, labelId:%s, err:%+v", labelId, err)
		return errors.New("query db failed")
	}
	
	// 检查记录数量
	if recordCnt > 10 {
		log.WarnContextf(ctx, "recordCnt greater than 10, labelId: %s", labelId)
		s.Monitor.Report(monitor.ImportLabelCntGreaterThan10)
	}
	if recordCnt > 60 { // 修改限制次数
		log.ErrorContextf(ctx, "limit request for labelId, labelId: %s, cosUrl: %s", labelId, cosUrl)
		return fmt.Errorf("labelId request too frequent, id: %s", labelId)
	}
	
	// 没有记录则插入
	_, insertErr := s.DB.Exec(ctx,
		"insert into label_load_record (task_id, label_id, cos_url, ts) values (?,?,?,?)",
		taskId, labelId, cosUrl, t.Unix())
	if insertErr != nil {
		log.ErrorContextf(ctx, "insert record failed, labelId: %s, err: %+v", labelId, insertErr)
		return errors.New("insert db failed")
	}
	
	log.InfoContextf(ctx, "QueryTaskTodayRecord ok, labelId: %s, cosUrl: %+v", labelId, cosUrl)
	return nil
}







```go
// QueryTaskTodayRecord ...
func (s *LabelLoadSvrServiceImpl) QueryTaskTodayRecord(ctx context.Context, taskId, labelId, cosUrl string) error {
	// 获取今天凌晨的时间戳
	t := time.Now()
	ts := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()

	// 构建 SQL 查询语句
	querySql := "select count(*) from label_load_record where label_id = ? and ts >= ?"
	log.InfoContextf(ctx, "Executing SQL query: %s, labelId: %s, ts: %d", querySql, labelId, ts)

	// 查询某个标签id今天的任务
	var recordCnt int
	dest := []interface{}{&recordCnt}
	err := s.DB.QueryRow(ctx, dest, querySql, labelId, ts)
	// 打印 querySql 的执行结果 dest
	log.InfoContextf(ctx, "QueryTaskTodayRecord result:  dest=%+v", dest)
	if err != nil {
		log.ErrorContextf(ctx, "query db failed, labelId:%s, err:%+v", labelId, err)
		return errors.New("query db failed")
	}

	// 检查记录数量
	if recordCnt > 10 {
		log.WarnContextf(ctx, "recordCnt greater than 10, labelId: %s", labelId)
		s.Monitor.Report(monitor.ImportLabelCntGreaterThan10)
	}
	if recordCnt > 60 { // 修改限制次数
		log.ErrorContextf(ctx, "limit request for labelId, labelId: %s, cosUrl: %s", labelId, cosUrl)
		return fmt.Errorf("labelId request too frequent, id: %s", labelId)
	} else {
		// 没有就插入记录
		_, insertErr := s.DB.Exec(ctx,
			"insert into label_load_record (task_id, label_id, cos_url, ts) values (?,?,?,?)",
			taskId, labelId, cosUrl, t.Unix())
		if insertErr != nil {
			log.ErrorContextf(ctx, "insert record failed,labelId:%+s, err:%+v", labelId, insertErr)
			return errors.New("insert db failed")
		}
	}
	log.InfoContextf(ctx, "QueryTaskTodayRecord ok ,labelId：%s, cosUrl:%+v", labelId, cosUrl)
	return nil

}
```





```
func EncryptField(c *TMtxContext, reqMap map[string]string, field string, enKey string) {
	iv := enKey
	tmpField, ok := reqMap[field]
	if ok {
		err, enField := util.CBCEncrypter([]byte(tmpField), []byte(enKey), []byte(iv))
		if err != nil {
			log.ErrorContextf(c.Ctx, "EncryptField CBCEncrypter error:%+v,filed:%+v", err, field)
		} else {
			reqMap[field] = string(enField)
		}

		// 获取密钥的前四位和后四位
		displayKey := fmt.Sprintf("%s****%s", enKey[:4], enKey[len(enKey)-4:])

		if field == "uid" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [uid]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "meeting_code" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [meeting_code]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "nick_name" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [nick_name]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "nickname" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [nickname]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "nick_name_0" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [nick_name_0]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "old_nick_name_0" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [old_nick_name_0]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "creator_nick_name" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [creator_nick_name]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "src_user_name" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [src_user_name]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "tar_user_name" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [tar_user_name]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "phone" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [phone]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "phone_number" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [phone_number]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "phone-number" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [phone-number]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "tar_phone" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [tar_phone]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "src_phone" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [src_phone]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "phone_id_0" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [phone_id_0]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "email" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [email]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "tar_email" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [tar_email]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "src_email" {
			log.InfoContextf(c.Ctx, "EncryptData-EncryptField [src_email]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		} else if field == "rooms_report_issue_user_emails" {
			log.InfoContextf(c.Ctx,
				"EncryptData-EncryptField [rooms_report_issue_user_emails]:%+v,enField:%+v,GCbcKey:%+v",
				tmpField, string(enField), displayKey)
		}
	}
}

```

