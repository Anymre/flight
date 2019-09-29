package service

import (
	"encoding/json"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

const url = "http://www.juneyaoair.com/PriceCache"
const changzhi = "CIH"
const nanjing = "NKG"

type Forward struct {
	Id    int64     `orm:"auto"`
	Date  time.Time `orm:"type(datetime)"`
	Now   time.Time `orm:"type(datetime)"`
	Price float64
}
type Back struct {
	Id    int64     `orm:"auto"`
	Date  time.Time `orm:"type(datetime)"`
	Now   time.Time `orm:"type(datetime)"`
	Price float64
}

func init() {
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	_ = orm.RegisterDataBase("default", "mysql", "flight:123456@tcp(server.anymre.top:3306)/flight?charset=utf8&parseTime=true&charset=utf8&loc=Asia%2FShanghai")
	orm.RegisterModel(new(Forward), new(Back))
	_ = orm.RunSyncdb("default", false, true)
}

func Get(s, t, n string) (string, error) {
	req := httplib.Get(url)
	req.Param("flightType", "OW")
	req.Param("departureDate", n)
	req.Param("returnDate", n)
	req.Param("sendCode", "CIH")
	req.Param("arrCode", "NKG")
	req.Param("periodType", "Line")
	req.Param("_", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	return req.String()
}

func Search(n string) ([]Forward, []Back) {
	str0, _ := Get(changzhi, nanjing, n)
func Search() (OW, error) {
	str0, err := Get(changzhi, nanjing)
	if err != nil {
		return OW{nil, nil}, err
	}
	var sth1 []interface{}
	_ = json.Unmarshal([]byte(str0), &sth1)

	arr0 := sth1[0].(map[string]interface{})
	var forward []Forward
	var back []Back

	for i := range arr0 {
		v := arr0[i].(float64)
		flight := new(Forward)
		flight.Date = timeFormat(i)
		flight.Now = time.Now()
		flight.Price = v
		forward = append(forward, *flight)
	}
	//如果有返程
	var arr1 map[string]interface{}

	if len(sth1) > 1 {
		arr1 = sth1[1].(map[string]interface{})
	} else {
		str2, _ := Get(nanjing, changzhi, n)
		var sth2 []interface{}
		_ = json.Unmarshal([]byte(str2), &sth2)
		arr1 = sth2[0].(map[string]interface{})

	}
	for i := range arr1 {
		v := arr0[i].(float64)
		flight := new(Back)
		flight.Date = timeFormat(i)
		flight.Now = time.Now()
		flight.Price = v
		back = append(back, *flight)
	}

	return forward, back
	return OW{forward, back}, nil
}

func timeFormat(r string) time.Time {
	HourTemplate := "15:04:05"
	DayTemplate := "2006-01-02 15:04:05" //常规类型

	now := time.Now().Format(HourTemplate)
	result, _ := time.ParseInLocation(DayTemplate, r+" "+now, time.Local)
	return result
}

func Perform() {
	o := orm.NewOrm()
	now0 := time.Now().AddDate(0, 0, 14)
	now1 := now0.AddDate(0, 0, 44)

	f0, b0 := Search(getTimeStr(now0))
	f1, b1 := Search(getTimeStr(now1))

	f := append(f0, f1...)
	b := append(b0, b1...)

	_, _ = o.InsertMulti(len(f), f)
	_, _ = o.InsertMulti(len(b), b)
}

func getTimeStr(t time.Time) string {
	return t.Format("2006-01-02")
}

type OW struct {
	Forward []Forward
	Back    []Back
	r, err := Search()
	if err != nil {
		fmt.Println(err)
	} else {
		for e := range r.forward {
			_, _ = o.Insert(&r.forward[e])
		}
		for e := range r.back {
			_, _ = o.Insert(&r.back[e])
		}
	}
}

func GetAll() OW {
	o := orm.NewOrm()
	var f []Forward
	_, _ = o.QueryTable("forward").All(&f)
	var b []Back
	_, _ = o.QueryTable("back").All(&b)
	var ow = OW{f, b}
	return ow
}
