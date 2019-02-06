# tyut-osc
理工助手 - 内网老系统爬虫系统

### 配置方式
*   配置文件统一置于 /tyuter/configs/ 目录下
*   CrawlerConfig.json - **Crawler 配置**
    *   BaseLocationURP: URP教务系统网站的基础URL,不需包含末尾 '*/*'
    *   BaseLocationGPA: GPA查询教务系统的基础URL,不包含末尾 '*/*'
    *   TempDir: 缓存目录,不影响功能,目前不需要使用
    *   UrpLoginAttempt: URP教务系统每次登录尝试次数, 考虑OCR引擎的正确率
    ```json
       {
         "BaseLocationURP" : ["http://202.207.247.44:8089", "http://202.207.247.44:8065", "http://202.207.247.44:8059", "http://202.207.247.44:8064" ,"http://202.207.247.51:8065","http://202.207.247.49"],
         "BaseLocationGPA" : ["http://202.207.247.60"],
         "TempDir": "/tyuter/tmp",
         "UrpLoginAttempt" : 10
       }
    ```

*   ZapConfig.json - **Uber Zap logging config** 
    *   disableCaller: false 输出caller方法名
    ```json
       {
         "level": "debug",
         "encoding": "json",
         "outputPaths": ["stdout"],
         "errorOutputPaths": ["stderr"],
         "disableCaller": false,
         "encoderConfig": {
           "messageKey": "message",
           "levelKey": "level",
           "levelEncoder": "lowercase",
           "callerKey": "caller",
           "callerEncoder": "short"
         }
       }
    ```
    
    
### 依赖
*   tesseract: `sudo apt install tesseract-ocr`
*   libtesseract: `sudo apt install libtesseract-dev`
*   环境变量: `TESSDATA_PREFIX` 指向`tessdata/`目录,包含 `rnd.traineddata`
*   网络环境: 需解决VPN连接问题

### 用法
```go
import "tyut-osc"
```

#### 1. GPACrawler

1. tyut-osc.NewGpaCrawler() *GpaCrawler 获取crawler实例(指针),线程安全的示例，通常使用单例模式
2. func (crawler *GpaCrawler) GetGpaDetail(stuid string, stuPassword string, targetStuid string) ( *DataModel.GpaDetail, error)
    获取tyut-osc/DataModel/GpaDetail 结构体实例，err包含errid，可在日志中查找具体内容；可以直接将err内容返回给最终用户
        
#### 2. URPCrawler

1. tyut-osc.NewUrpCrawler() *UrpCrawler 创建一个UrpCrawler实例，线程安全，通常使用单例模式
2. (urp *UrpCrawler) CreateClientAndLogin(stuid string, stuPassword string) (client *http.Client, activateUrlIdx int, err error)
    提供urp系统学号和密码，创建一个具有登录状态的http client，随后所有的获取数据方法都将依赖此http.Client实例, 相应获取数据为登录用户的对应数据
3. (urp *UrpCrawler) GetPassedCourses(client *http.Client, activateUrlIdx int) (terms []DataModel.Term, err error)
    获取通过成绩列表，按学期排列
    ```go
    type Term struct {
    	TermDescription string
    	TermYear        int
    	TermOrder       int // 0-春季学期, 1-秋季学期
    	PassedCourses   []PassedCourse
    }
    ```
4. (urp *UrpCrawler) GetFailedCourses(client *http.Client, activateUrlIdx int) (fcourses []DataModel.FailedCourse, err error)
    获取挂科列表，包括曾挂科和在挂科
    ```go
    type FailedCourse struct {
    	CourseId             string    //课程号
    	CourseSequenceNumber string    //课序号
    	CourseName           string    // 课程名
    	EnglishCourseName    string    //英文课程名
    	CourseCredit         float64   //课程学分
    	SelectionProperty    string    //课程属性 - 综合必修
    	Score                float64   //得分
    	ExamTime             time.Time //20170626
    	StillFail            bool
    	Reason               string
    }
    ```
    
5. (urp *UrpCrawler) GetCourseList(client *http.Client, activeUrlIdx int) (seletects []DataModel.SelectedCourse, err error)
     获取已选课表