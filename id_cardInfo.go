package utils

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type IDCardInfo struct {
	IDCardNo string
	Year     string
	Month    string
	Day      string
	BirthDay string
	Area struct{
		Status bool
		Result string
		Province string
		City string
		County string
		Using int
	}
	province string
	sufixProvince uint  //获取省级行政区划代码
	sufixCity uint   //获取地市级行政区划代码
	county uint   //获取县级行政区划代码
	Constellation string //星座
	Zodiac string   //属相
	Sex      uint
	Age      uint
}
//星座
var aConstellations =[]string{
	"水瓶座",
	"双鱼座",  // 2.20-3.20 [Pisces]
	"白羊座",  // 3.21-4.19 [Aries]
	"金牛座",  // 4.20-5.20 [Taurus]
	"双子座",  // 5.21-6.21 [Gemini]
	"巨蟹座",  // 6.22-7.22 [Cancer]
	"狮子座",  // 7.23-8.22 [Leo]
	"处女座",  // 8.23-9.22 [Virgo]
	"天秤座",  // 9.23-10.23 [Libra]
	"天蝎座",  // 10.24-11.21 [Scorpio]
	"射手座",  // 11.22-12.20 [Sagittarius]
	"魔羯座",  // 12.21-1.20 [Capricorn]
}
var aProvinces = map[string]string{
	"11": "北京",
	"12": "天津",
	"13": "河北",
	"14": "山西",
	"15": "内蒙古",
	"21": "辽宁",
	"22": "吉林",
	"23": "黑龙江",
	"31": "上海",
	"32": "江苏",
	"33": "浙江",
	"34": "安徽",
	"35": "福建",
	"36": "江西",
	"37": "山东",
	"41": "河南",
	"42": "湖北",
	"43": "湖南",
	"44": "广东",
	"45": "广西",
	"46": "海南",
	"50": "重庆",
	"51": "四川",
	"52": "贵州",
	"53": "云南",
	"54": "西藏",
	"61": "陕西",
	"62": "甘肃",
	"63": "青海",
	"64": "宁夏",
	"65": "新疆",
}
//星座边缘日切数据
var aConstellationEdgeDays = []int{21, 20, 21, 20, 21, 22, 23, 23, 23, 24, 22, 21,}
//实例化居民身份证结构体
func NewIDCard(IDCardNo string) *IDCardInfo {
	if IDCardNo == "" || len(IDCardNo) != 18 {
		return nil
	}

	ins := IDCardInfo{
		IDCardNo: IDCardNo,
	}
	ins.province = ins.IDCardNo[:2]
	ins.sufixProvince = ins.stringCtvInt(ins.IDCardNo[:2]+"0000")
	ins.sufixCity = ins.stringCtvInt(ins.IDCardNo[:4]+"00")
	ins.county = ins.stringCtvInt(ins.IDCardNo[:6])
	ins.GetArea()
	ins.Year = ins.GetYear()
	ins.Month = ins.GetMonth()
	ins.Day = ins.GetDay()
	ins.Sex = ins.GetSex()
	ins.BirthDay = ins.GetBirthDayStr()
	ins.Age = ins.GetAge()
	ins.Constellation = ins.GetConstellation()
	ins.Zodiac = ins.GetZodiac()
	return &ins
}
func (s *IDCardInfo)stringCtvInt(pra string) uint{
	r,err :=strconv.Atoi(pra)
	if err!=nil{
		log.Fatal(err)
		return 0
	}
	return uint(r)
}

//根据身份证号获取生日（时间格式）
func (s *IDCardInfo) GetBirthDay() *time.Time {
	if s == nil {
		return nil
	}

	dayStr := s.IDCardNo[6:14] + "000001"
	birthDay, err := time.Parse("20060102150405", dayStr)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &birthDay
}

//根据身份证号获取生日（字符串格式 yyyy-MM-dd HH:mm:ss）
func (s *IDCardInfo) GetBirthDayStr() string {
	defaultDate := "1999-01-01 00:00:01"
	if s == nil {
		return defaultDate
	}

	birthDay := s.GetBirthDay()
	if birthDay == nil {
		return defaultDate
	}

	return birthDay.Format("2006-01-02 15:04:05")
}

//根据身份证号获取生日的年份
func (s *IDCardInfo) GetYear() string {
	if s == nil {
		return ""
	}

	return s.IDCardNo[6:10]
}

//根据身份证号获取生日的月份
func (s *IDCardInfo) GetMonth() string {
	if s == nil {
		return ""
	}

	return s.IDCardNo[10:12]
}

//根据身份证号获取生日的日份
func (s *IDCardInfo) GetDay() string {
	if s == nil {
		return ""
	}

	return s.IDCardNo[12:14]
}

//根据身份证号获取性别
func (s *IDCardInfo) GetSex() uint {
	var unknown uint = 3
	if s == nil {
		return unknown
	}

	sexStr := s.IDCardNo[16:17]
	if sexStr == "" {
		return unknown
	}

	i, err := strconv.Atoi(sexStr)
	if err != nil {
		return unknown
	}

	if i%2 != 0 {
		return 1
	}

	return 0
}

//根据身份证号获取年龄
func (s *IDCardInfo) GetAge() uint {
	if s == nil {
		return 19
	}

	birthDay := s.GetBirthDay()
	if birthDay == nil {
		return 19
	}

	now := time.Now()

	age := now.Year() - birthDay.Year()
	if now.Month() > birthDay.Month() {
		age = age - 1
	}

	if age <= 0 {
		return 19
	}

	if age <= 0 || age >= 150 {
		return 19
	}

	return uint(age)
}
//根据身份证证号获取所在地区.
func (s *IDCardInfo)GetArea() {
	var result,_provinceName,_cityName string
	if _,ok :=aProvinces[s.province];ok==true{
		_province,err :=getDivision(s.sufixProvince)
		if err!=nil{
			_provinceName =""
		}else{
			_provinceName =string(_province["name"])
		}

		_city,err := getDivision(s.sufixCity)
		if err!=nil{
			_cityName =""
		}else{
			_cityName =string(_city["name"])
		}
		_county,err := getDivision(s.county)
		if _county!=nil{
			result = _provinceName+_cityName+string(_county["name"])
		}else{
			result = _provinceName+_cityName
		}
		s.Area.Status=true
		s.Area.Result=result
		s.Area.Province=_provinceName
		s.Area.City=_cityName
		s.Area.County=string(_county["name"])
		s.Area.Using = 0
	}else{
		s.Area.Status =false
	}
}
func getDivision(id uint) (map[string][]byte,error){
	orm, err := xorm.NewEngine("sqlite3", workDir()+"/db/database.sqlite")
	if err != nil {
		log.Printf("orm failed to initialized: %v", err)
		return nil,err
	}
	rows, err := orm.Query("SELECT divisions.id,divisions.name,divisions.status,divisions.year FROM divisions WHERE divisions.id = ?", id)
	if err != nil {
		log.Printf("orm failed to initialized: %v", err)
		return nil,err
	}
	return rows[0],nil
}
//获取生肖
func(s *IDCardInfo) GetZodiac() string{
	start :=1901
	end,_ :=strconv.Atoi(s.IDCardNo[6:10])
	x :=(start-end)%12
	if x==1 || x==-11{
		return "鼠"
	}
	if x==0 {
		return "牛"
	}
	if x==2 || x==-10{
		return "猪"
	}
	if x==3 || x==-9{
		return "狗"
	}
	if x==4 || x==-8{
		return "鸡"
	}
	if x==5 || x==-7{
		return "猴"
	}
	if x==6 || x==-6{
		return "羊"
	}
	if x==7 || x==-5{
		return "马"
	}
	if x==8 || x==-4{
		return "蛇"
	}
	if x==9 || x==-3{
		return "龙"
	}
	if x==10 || x==-2{
		return "兔"
	}
	if x==11 || x==-1{
		return "虎"
	}
	return ""
}
//获取星座
func(s *IDCardInfo) GetConstellation() string{
	monthInt,_ :=strconv.Atoi(s.Month)
	month :=monthInt-1
	dayInt,_ :=strconv.Atoi(s.Day)
	if dayInt<aConstellationEdgeDays[monthInt]{
		month= month-1
	}
	if month>=0{
		return aConstellations[month]
	}
	return aConstellations[11]
}
func workDir() string {
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	wd := filepath.Dir(execPath)
	if filepath.Base(wd) == "bin" {
		wd = filepath.Dir(wd)
	}

	return wd
}
