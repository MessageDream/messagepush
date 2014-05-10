package main

import (
	"github.com/golang/glog"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
	"time"
)

var (
	db *mgo.Database
)

func InitDb() {
	conn := Conf.DbConn
	if conn == "" {
		glog.Error("The database connection string has not been configured.")
		os.Exit(1)
	}

	session, err := mgo.Dial(conn)
	if err != nil {
		glog.Error("MongoDB connected error:", err.Error())
		os.Exit(1)
	}

	session.SetMode(mgo.Monotonic, true)

	db = session.DB("messagepush")
}

type Device struct {
	Id_    bson.ObjectId `bson:"_id"`
	AppKey bson.ObjectId `bson:"appkey"`
	Token  string
	Tags   []string
	Alias  string
}

type App struct {
	AppKey         bson.ObjectId `bson:"_id"`
	AppName        string
	AppIcon        string
	AndroidPkgName string
	ApnsEnv        int
	BundleID       string
	DevCerUrl      string
	DevCerPwd      string
	ProCerUrl      string
	ProCerPwd      string
	UserID         bson.ObjectId `bson:"uid"`
}

type Statistics struct {
	AllCount     int
	SuccessCount int
}

type PushRecord struct {
	Id_               bson.ObjectId `bson:"_id"`
	AppKey            bson.ObjectId `bson:"appkey"`
	Content           string
	IosStatistics     Statistics
	AndroidStatistics Statistics
	WpStatistics      Statistics
	Type              int
	PushTime          time.Time
}

type UserContactInfo struct {
	MobilePhone string
	Tell        string
	Company     string
	Addr        string
}

type User struct {
	Id_      bson.ObjectId `bson:"_id"`
	UserName string
	Password string
	Portrait string
	Email    string
	Contact  UserContactInfo
}

func (u *User) GetApps() ([]App, error) {
	c := db.C("apps")
	var apps []App
	err := c.Find(&bson.M{"uid": u.Id_}).All(&apps)
	return apps, err
}
