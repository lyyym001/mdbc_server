package pb



//同步客户端玩家ID
type SyncPID struct {
 PID int32
}


//登录回调
type SyncLoginB struct {
    CtrlFlag string
    Code string
}

//同步客户端玩家ID
type SyncLogin struct {
  TID string //老师账号
  CID string //学生账号

  //如果是学生 需要都传，如果是老师 两个账号填一样的
}


//玩家聊天数据
type Talk struct{
 Content string    //聊天内容
}

//本地年级分类信息列表
type Sync_GetLocalGradeTypeList struct {
    LocalGradeTypeArr []LocalGradeTypeData
}

//本地年级分类信息
type LocalGradeTypeData struct {
    ID						int
    TypeName				string
    Visible					int
}


//本地课程类型信息列表
type Sync_GetLocalCourseTypeList struct {
    LocalCourseTypeArr []LocalCourseTypeData
}

//本地课程类型信息
type LocalCourseTypeData struct {
    ID						int
    TypeName				string
    BEdit					int
    BNotClassified			int
    InClassType 			int
    InClassTypeSort 			int
}

//本地课程信息列表
type Sync_GetLocalCourseList struct {
    LocalCourseArr []CoursewareData
}

//本地课程信息
type CoursewareData struct
{
    ID    int
    Name string                //名称
    IconName string             //图标名称
    CourseID string            //课程ID
    CourseType int            //课程类型(1-安全课程 3-视频课程,  这个参数不是真正的左边的课程类型)
    CourseOwner int           //课程拥有者（1-飞蝶  2-老师自主添加）
    InCourseType string       //属于哪个课程类型(本地数据使用) 会属于多个课程类型下
    InCourseTypeSort string       //课程类型下的排序ID
    ThirdType    int			//第三方课程的类型(0-视频 1-APK)
    ThirdMsg string             //第三方课程信息(APK包名)
    Md5 string
    GameUrl string
    ResVersion string
}

//老师登录成功后  跟老师客户端通信消息结构
type Sync_LoginTeacher_Send struct {
    Result     string
    CID        int
    AllCourses []Sync_LoginTeacher_Info     //老师发过来的课程
}


type Sync_LoginTeacher_Info struct {
    Name          string //名称
    IconName      string //图标名称
    CourseID      string
    CourseType    string //课程类型
    CourseSubType string //课程子类型
    Extras		  string //附加信息
    Md5			  string
    GameUrl       string
    ResVersion    string
}

//删除第三方课程协议
type DeleteCourse struct {
    ID						int
}

type UpdateCourse struct
{
    CourseID string
    Md5 string
    GameUrl string
    ResVersion string
    GameIcoPath string
    GameName string
    VersionNo string
}


//学生信息
type StudentInfo struct {
    StuUserName string
    Flag        int         //0-未登录 1-已登录
}


//老师控制信息
type Sync_TeacherControlData struct {
    IsTeacherControl string	//是否控制中
}

type Sync_SelectCourseware_Send struct {
    CourseID    string
    CourseMode  string
}

type Sync_GetCourse struct {
    CourseID string
    Mode     string
}

type StuName struct {
    StuUserName string //设备账号
}

type Sync_LeaveData struct {
    LeaveType int //离开类型 0-正常退出 1-被踢
}

//老师端通知结束课程
type Sync_CloseCourse struct {
    StuAccountID string //账号
}

type Sync_SendCourse struct {
    CourseID string
    Md5      string
    Mode     string
}

//老师端通知输入学号
type Sync_TInputSnum struct {
    StuAccountID string //账号
}


//学号输入不正常返回消息给学生
type Sync_SNumReturn struct {
    Result string //结果 "srep" 重复输入学号 "snull" 学号不存在/输入有误
}

type Sync_SNumToTeacher struct {
    SNum        string //学号
    PName       string //姓名
    StuUserName string //账号
}

type Sync_SInputSnum struct {
    SNum string //学号
}

//所有学生信息
type AllStudentData struct {
    AllData []StudentData
}

//一个学生的信息
type StudentData struct {
    Snum	string			//学号
    Sname 	string			//学生姓名
    Sclass	string			//学生年级班级
}

type GetConditionData struct {
    Condition string
}

//所有学生学习记录
type AllStudentRecordData struct {
    AllRecordData []StudentRecordData
}

type StudentRecordData struct{
    StudentSnum		string	//学生学号
    StudentName		string	//学生姓名
    StudentClass	string	//学生班级
    CourseID    	string	//课程ID
    StudyTime 		string	//学习时间
    CourseName		string	//课程名称
    StudyMode		string	//学习课程类别
    StudyTotalTime	int		//学习总用时
    StudyScore		int		//考试得分
    StudyAbility	string	//考试知识点
}

//所有学习记录信息
type AllStudyInfoData struct {
    StudyInfoData []SingleStudyInfoData
}

//单次学习记录信息
type SingleStudyInfoData struct {
    StudyTime 			string
    StudyCourseMode		string
    StudentRecordData 	[]StudentRecordData
}


type Sync_SCourseId struct {
    CourseID string
}

type UpdatePasswordData struct {
    Password 	string		//密码
}

//电量和存储空间
type Sync_BatteryAndSpaceData struct {
    BatteryLevel string
    AvailableSpace string
    TotalSpace string
}

//电量和存储空间
type Sync_BatteryAndSpaceData_Send struct {
    StuName	string
    BatteryLevel string
    AvailableSpace string
    TotalSpace string
}

//课程进度节点信息
type Sync_CourseNodeData struct {
    NodeName string		//节点名称
    NodeIndex int		//节点索引
    NodeTotal int		//节点总数
}

type Sync_CourseNodeData_Send struct {
    StuName	string		//学生姓名
    NodeName string		//节点名称
    NodeIndex int		//节点索引
    NodeTotal int		//节点总数
}

type Sync_GetReportBySnumData struct {
    Snum 					string
}


type Sync_GetReportBySnumData_Send struct {
    Snum 					string
    ReportBySnumData		[]StudentReportData
}

//学生报表数据
type StudentReportData struct {
    CourseName				string		//课程名称
    CourseType				string		//课程类型
    StudyCount				int			//学习次数
    AverageTime				float32		//平均时长
    PersonalAverageScore	float32		//个人平均分
    PersonalHighestScore	float32		//个人最高分
    PersonalLowestScore		float32		//个人最低分
    TotalAverageScore		float32		//总平均分
    TotalHighestScore		float32		//总最高分
}

type Sync_GetReportBySnameData_Send struct {
    Sname 					string
    ReportBySnumData		[]StudentReportData
}

type Sync_GetReportBySnameData struct {
    Sname 					string
}

type Sync_GetReportByCnameData_Send struct {
    Cname 					string
    ReportByCnameData		[]CourseReportData
}

//课程报表数据
type CourseReportData struct {
    CourseName				string		//课程名称
    StudyCount				int			//学习次数
    AverageTime				float32		//平均时长
    AverageScore			float32		//平均分
}

type Sync_GetReportByCnameData struct {
    Cname 					string
}

type Sync_GetReportByCtypeData struct {
    Ctype 					string
}

type Sync_GetReportByCtypeData_Send struct {
    Ctype 					string
    ReportByCtypeData		[]CourseReportData
}

type Sync_EnterRoom struct {
    Students []StudentInfo
}

//打开投影
type Sync_OpentheprojectionData struct {
    StuAccountID			string
    OpenTheProjection		string
}

type Sync_ClosetheprojectionData struct {
    StuAccountID			string
}

type CreateCourseTypeData_Send struct {
    ID						int64  //ID
}

//删除一个课程类型
type DeleteCourseTypeData struct {
    ID						int  //ID
}

//通知学生更新单一课程
type Sync_UpdataCourse struct {
    StuAccountID string     //学生账号
}

type Sync_GetCreatexCourse struct {
    Pid string
    Uid	string
    Index int
}

type Sync_SetCreatexCourse struct {
    Pid string
    Uid	string
}

//老师端通学生重启课程
type Sync_GetRestarCreatex struct {
    StuAccountID string //账号
    Pid string //PID
    Uid string //UID
    Index int
}


type Sync_ScoreGet struct {
    SNum     string  //学号
    CourseID string  //课程ID
    Score    float32 //成绩
    Ability  string  //知识点
}

type Sync_ScoreSend struct {
    AccountName string  //账号
    SNum        string  //学号
    CourseID    string  //课程ID
    Score       float32 //成绩
    Ability		string	//知识点对错
}

type DB_SaveData struct {
    Tid         string
    CourseId    string
    CourseName  string
    CourseMode  string
    CourseType	string
    StuTimeLong int64
    SNum        string
    Score       float32
    SetupDate   int64
    LeaveType   int 		//离开类型 0-正常退出 1-被踢
    Ability		string
}


//==============================new==================================
type Sync_Hello struct {
    Ip string           //ip
    Port        int     //端口
}

//同步客户端玩家ID
type Tcp_Login struct {
    UserName string     //账号
    AccountType int     //账号类型，0-学生 1-老师
    Code    string      //机器码

}

//学生信息
type Tcp_Info struct {
    Flag        int         //0-重复连接 -1-授权失败 1-成功
    UName       string      //用户名
    WorkRuning int          //活动是否开启中（对于客户端来说，如果不在活动中，需要到大厅，如果在活动中则不处理）
}

//玩家广播数据
type BroadCast struct {
    PID int32
    Pos  Position
}

type Position struct {
    X                    float32
    Y                    float32
    Z                    float32
    V                    float32
}

type SyncPlayers struct {
    Ps []BroadCast
}

type Tcp_RequestScene struct {
    SceneId string
    //Multiple byte
    Mode byte
}

type Tcp_ResponseScene struct {
    Code int
    Mode byte
    UNumber int    //本次参与教学人数
    WorkId string     //作品ID
    MainCtroller string  //主控用户
}

type TCP_RegisterObj struct {
    ObjId int
    X float32
    Y float32
    Z float32
    RX float32
    RY float32
    RZ float32
    InteractiveType int
    Tb int
    Visiable int
}

type TCP_Grab_Call struct {
    Uid int
    ObjId int
    UName string
}

type Tcp_Object struct {
    ObjId int
}


type TCP_DetachObj struct {
    ObjId int
    X float32
    Y float32
    Z float32
    RX float32
    RY float32
    RZ float32
    RIX float32
    RIY float32
    RIZ float32
    AX float32
    AY float32
    AZ float32
}

type TCP_TbObj struct {
    ObjId int
    X float32
    Y float32
    Z float32
    RX float32
    RY float32
    RZ float32
}

type Tcp_ObjectStatus struct {
    ObjId int
    Status int
}

type Tcp_Step struct {
    StepId int
    StepDate string
    StepState int
    UName string
}

type Tcp_Leave struct {
    State int
}

type Tcp_Shixun struct {
    Users [] string
    Mode byte
    WorkId string
}

type Tcp_QuestionInfo struct {
    Qid int
    Code int
    StepDate string
    StepState int
    UName string
}

type Tcp_WorkRecord struct {
    Workid string
    Workname string
    Date string
    Mode int
    Partnumber int
    MaxScore float32
}

type Tcp_ResponseRecord struct {
    Code int
    UniqueId int32
}

type Tcp_WorkInfoRecord struct {
    Username string
    Uname string
    Date string
    Type int
    SetId int
    State int
    Score float32
    Content string
    UniqueId int32
}

type Tcp_ResponseRecordInfo struct {
    Code int
    RecordId int32
}

type Tcp_Tj1 struct {
    Cid int
    Mode2 int
    Mode3 int
    Page int
}

type Tcp_Tj1Data struct {
    MaxPage int
    Data [] Tcp_Tj1Info
}

type Tcp_Tj1Info struct {
    Date string
    Mode int
    Number int
    MaxScore float32
    Score float32
    UniqueId int
}

type Tcp_Tj2 struct {
    UniqueId int
    Page int
}

type Tcp_Tj2Data struct {
    MaxPage int
    MaxNumber int
    Data [] Tcp_Tj2Info
}

type Tcp_Tj2Info struct {
    UName string
    Score float32
}

type Tcp_Tj3 struct {
    UniqueId int
    UName string
    WorkName string
    Mode string
    MaxScore float32
    Number int
}

type Tcp_Tj3Data struct {
    WorkName string
    Mode string
    MaxScore float32
    Number int
    Data [] Tcp_Tj3Info
}

type Tcp_Tj3Info struct {
    UName string
    Date string
    State int
    Type int
    Score float32
    Content string
}