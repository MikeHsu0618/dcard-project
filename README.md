# dcard-project

2022 Dcard 後端實習作業題目：

請使用 Golang 或是 Node.js 設計以及實現出一套短網址系統並且附上測試：

- 短網址共有兩支 APIs, 請依照下列要求實現：
  - 一個可以上傳網址以及過期時間並返回短網址的 Restful Api
  - 一個輸入短網址就可以轉導到原網址的 API, 如過期則返回 404
- 請對這兩支 APIs, 做出合理的約束以及錯誤處理
- 不需考慮權限驗證問題
- 許多客戶端可能會同時發出請求或連接到不存在的短網址, 請將效能問題納入考量 
- 請安心使用任何第三方套件

## 專案講解

### Demo 連結

- [http://www.d-project.link](http://www.d-project.link)

### 目的

本專案以『盡可能的做出一個好的以及完整的短網址服務體驗』為目標

所以建構了一個前後端環境來運行起整個完整的短網址服務

### 如何實現短網址服務

* [實現介紹以及細節](./docs/README.md)


### 短網址好處

短網址的好處非常多, 在於社群軟體上也有廣泛的使用, 在商業應用上還可以站在開發者角度針對短網址進行流量與點擊等統計, 
挖掘出各種有價值的訊息, 對未來的決策有相當大的幫助
- 在 Twitter、微博等每條訊息有限制字數, 可以有效解決原網址字數過多的問題, 尤其在傳統文字短信上錙銖必較的字數, 不小心超過的話很容易多付出不少冤枉錢
- 跟 QRCode 有絕妙的搭配, 縮短過後的網址可以根本上的解決 QRCode 資料量越大或文字數愈多, `產生出來的圖片顆粒密度越大導致辨識度低落的痛點`
- 除了上述所說的, 短網址對於 SEO 應用是非常靈活, 除了收集統計點擊數, 甚至可以因此知道使用者的使用裝置或是瀏覽器, 以及來自哪個網頁到轉導到哪個網頁,
大大的幫助我們建立出使用者的足跡以及輪廓

### 需求分析

小弟對此專案的需求了解分成下列幾點 :
- 需要一個前端輸入介面同時來呈現返回的短網址相關資訊結果
- 一個正常情況下`讀大於寫`，且沒有`更新`需求的服務
- 有高流量高併發的可能
- 需要預防惡意多次請求攻擊
- 需要長期運行且穩定, 對資料量擴張以及持久性需要把持
- 在顧及上述條件時需兼顧性能以及一致性

### 使用環境及工具

![image](./resources/asset/img/stucture.png)

- 前端: Javascript(Vue) 
- 後端: Golang(Gin, Gorm)
- 資料庫: PostgreSql
- 快取: Redis
- CI/CD: Github Action
- Test: Testify

### 如何運行該專案(使用docker-compose)

可利用本專案的docker-compose.yaml會一次啟動Backend、PostgreSql、Redis，方便直接運行測試。請確保主機有docker環境，如果是Linux環境則需要另外安裝docker-compose套件。而如果是Windows、Mac則只需要安裝Docker Desktop即可。

#### Clone 專案

```bash
# 透過 git clone 專案到主機任意路徑下
git clone https://github.com/MikeHsu0618/dcard-project.git
```
#### 運行專案

````bash
# 在本專案的根目錄下執行以下指令即可
cp .env.example .env
docker compose up -d pg-master pg-slave dcard-project redis
````
#### 初始化資料庫
```bash
# migrate/init schema
docker compose up -d migrate
```


#### Schema
```
// migrate up 

DROP TABLE IF EXISTS urls CASCADE;

CREATE TABLE urls
(
    id         BIGSERIAL,
    org_url    varchar(255) NOT NULL UNIQUE,
    
    // 如果有產生時間的排序需求, 非常建議在 created_at 加上索引
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

COMMENT ON COLUMN urls.org_url IS '原網址'

----------------------------------------------------

// migrate down
DROP table IF EXISTS urls;
```

#### 運行測試
```bash
# run test
go test -v ./test
```

#### Swagger Document
```
// 運行成功起後, 可於以下路徑顯示 API 文件
http://localhost:8080/swagger/index.html
```

#### Available Services

* starts a server on localhost port 8080 (by default).
* http://localhost:8080

| Method | Path        | Usage                       |
|--------|-------------|-----------------------------|
| POST   | /           | get short url and meta info |
| GET    | /{shortUrl} | redirect to origin url      |

## 總結

這裡需要特別感謝 [KennyChenFight 肯尼攻城獅](https://github.com/KennyChenFight) 分享的各種資源,
無論是 Blog Youtube 或者是 Github 資源, 幾乎都拜讀了無數遍, 才有這份專案

本專案為小弟從零學習 golang 的第一個 demo 作品, 結合 Redis, PostgreSql, Docker, 練習系統設計與工具使用, 嘗試做出一個完整的服務

希望持續優化的方向為學習使用 
- 學習使用 k8s 建立強健以及彈性的架構
- 將溝通方式從 Restful Api 轉成 gRPC, 把本專案轉成一個微服務