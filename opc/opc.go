package opc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

// PersonInfo struct
type PersonInfo struct {
	Name    string
	age     int32
	Sex     bool
	Hobbies []string
}

var personInfo = []PersonInfo{
	{"David", 30, true, []string{"跑步", "读书", "看电影"}},
	{"Lee", 27, false, []string{"工作", "读书", "看电影"}},
}

// ReadFile json
func ReadFile() {

	filePtr, err := os.Open("./conf/person_info.json")
	if err != nil {
		fmt.Printf("Open file failed [Err:%s]", err.Error())
		return
	}
	defer filePtr.Close()

	var person []PersonInfo

	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&person)
	if err != nil {
		fmt.Println("Decoder failed", err.Error())

	} else {
		fmt.Println("Decoder success")
		fmt.Println(person)
	}
}

// WriteFile json
func WriteFile(v interface{}) {

	// 创建文件
	filePtr, err := os.Create("./conf/d_opcda_client.json")
	if err != nil {
		fmt.Println("创建文件失败！", err.Error())
		return
	}
	defer filePtr.Close()

	// 创建Json编码器
	encoder := json.NewEncoder(filePtr)

	if err := encoder.Encode(v); err != nil {
		fmt.Println("写配置文件失败！", err.Error())

	} else {
		fmt.Println("写配置文件成功！")
	}
}

//CopyFile 文件复制
func CopyFile(src, des string) (err error) {
	b, err1 := ioutil.ReadFile(src)
	if err1 != nil {
		return err1
	}
	err2 := ioutil.WriteFile(des, b, 0666)
	if err2 != nil {
		return err1
	}
	fmt.Println("读取成功！")
	return nil

}

//OpcdaForm struct
type OpcdaForm struct {
	Module             string `form:"module" binding:"required" json:"Module"`
	MainServerIP       string `form:"main_server_ip" binding:"required" json:"main_server_ip" `
	MainServerPrgid    string `form:"main_server_prgid" binding:"required" json:"main_server_prgid"`
	MainServerClsid    string `form:"main_server_clsid" binding:"required" json:"main_server_clsid"`
	MainServerDomain   string `form:"main_server_domain" binding:"required" json:"main_server_domain"`
	MainServerUser     string `form:"main_server_user" binding:"required" json:"main_server_user"`
	MainServerPassword string `form:"main_server_password" binding:"required"  json:"main_server_password"`
	BakServerIP        string `form:"bak_server_ip" binding:"required" json:"bak_server_ip"`
	BakServerPrgid     string `form:"bak_server_prgid" binding:"required" json:"bak_server_prgid"`
	BakServerClsid     string `form:"bak_server_clsid" binding:"required" json:"bak_server_clsid"`
	BakServerDomain    string `form:"bak_server_domain" binding:"required" json:"bak_server_domain"`
	BakServerUser      string `form:"bak_server_user" binding:"required" json:"bak_server_user"`
	BakServerPassword  string `form:"bak_server_password" binding:"required" json:"bak_server_password"`
}

//Opcdaget return opc da config page
func Opcdaget(c *gin.Context) {

	c.HTML(http.StatusOK, "da.html", gin.H{
		"title": "opc da",
	})

}

//Opcdapost handle post
func Opcdapost(c *gin.Context) {
	var form OpcdaForm
	if c.ShouldBind(&form) == nil {

		file, err := c.FormFile("file")
		if err != nil {
			fmt.Println(err)
		}

		// 上传文件到指定的路径
		// dst := filepath.Base(`D:\Web\go\src\hello\upload\` + file.Filename)
		dst := fmt.Sprintf(`./upload/` + file.Filename)
		if e := c.SaveUploadedFile(file, dst); e != nil {
			fmt.Println(e)
		}

		c.HTML(http.StatusOK, "opc_show.html", gin.H{
			"title":      "opc show",
			"now":        (time.Now()).Format("2006-01-02 15:04:05"),
			"opc_config": form,
		})

		outputFile, err := os.OpenFile("./conf/da.json", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Printf("An error occurred with file opening or creation\n")
			return
		}
		defer outputFile.Close()

		outputWriter := bufio.NewWriter(outputFile)
		outputString := "hello golang!\n"

		for i := 0; i < 1; i++ {
			outputWriter.WriteString(outputString)
		}
		outputWriter.Flush()
		// WriteExcel()
		ReadExcel(form)
	}
}

type tags struct {
	TagID          int    `json:"tag_id"`
	PublishTagName string `json:"publish_tag_name"`
	SourceTagName  string `json:"source_tag_name"`
	DataType       string `json:"data_type"`
}

type groups struct {
	GroupID      int    `json:"group_id"`
	GroupName    string `json:"group_name"`
	CollectCycle int    `json:"collect_cycle"`
	Tags         []tags `json:"tags"`
}

//Gslist GROUPS
type Gslist struct {
	GroupList []groups `json:"groups"`
}

//Da opc da config struct
type Da struct {
	OpcdaForm
	Gslist
}

var ts tags
var gs groups
var gl Gslist

// ReadExcel read tag from excel
func ReadExcel(form OpcdaForm) {
	f, err := excelize.OpenFile("./upload/opc.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get value from cell by given worksheet name and axis.
	// cell := f.GetCellValue("group1-10", "B2")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(cell)
	// Get all the rows in the Sheet.
	// rows := f.GetRows("group1-10")

	// f.SetActiveSheet(2)
	// index := f.GetActiveSheetIndex()
	// name := f.GetSheetName(index)
	// fmt.Println(name)
	m := f.GetSheetMap()
	var s []int
	for key := range m {
		s = append(s, key)
	}
	sort.Ints(s)
	var i = 0
	for _, v := range s {
		rows := f.GetRows(m[v])
		s := strings.Split(m[v], "-")
		// fmt.Println(s, s[0][len(s[0])-1:], s[1])
		gs.GroupName = s[0]
		gs.GroupID, _ = strconv.Atoi(s[0][len(s[0])-1:])
		gs.CollectCycle, _ = strconv.Atoi(s[1])
		tmp := []tags{}
		for _, row := range rows[1:] {
			for _, colCell := range row[:1] {
				// fmt.Println(colCell)
				ts.PublishTagName = colCell
				ts.SourceTagName = colCell
				ts.DataType = "ENUM_INT32"
				ts.TagID = i
			}
			tmp = append(tmp, ts)
			// gs.Tags=tmp
			gs.Tags = make([]tags, len(tmp), cap(tmp))
			copy(gs.Tags, tmp)
			i++
		}

		gl.GroupList = append(gl.GroupList, gs)

	}
	// fmt.Println(form, gl)
	// b, _ := json.Marshal(GroupList)
	// fmt.Println(string(b))
	dajson := Da{
		form,
		gl,
	}
	WriteFile(dajson)
}

// WriteExcel write data to excel
func WriteExcel() {
	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet2")
	// Set value of a cell.
	f.SetCellValue("Sheet2", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "B2", 100)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save xlsx file by the given path.
	if err := f.SaveAs("./upload/Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
