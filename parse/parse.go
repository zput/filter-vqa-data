package parse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
)

type JoinAQ struct {
	Questions   *Question   `json:"questions,omitempty"`
	Annotations *Annotation `json:"annotations,omitempty"`
}

type Annotation struct {
	AnswerType string `json:"answer_type"`
	Answers    []struct {
		Answer           string `json:"answer"`
		AnswerConfidence string `json:"answer_confidence"`
		AnswerID         int    `json:"answer_id"`
	} `json:"answers"`
	ImageID              int    `json:"image_id"`
	MultipleChoiceAnswer string `json:"multiple_choice_answer"`
	QuestionID           int    `json:"question_id"`
	QuestionType         string `json:"question_type"`
}

type Question struct {
	ImageID    int    `json:"image_id"`
	Question   string `json:"question"`
	QuestionID int    `json:"question_id"`
}

type T struct {
	TaskType  string      `json:"task_type,omitempty"`
	Questions []*Question `json:"questions,omitempty"`

	Annotations []*Annotation `json:"annotations,omitempty"`

	DataSubtype string `json:"data_subtype"`
	DataType    string `json:"data_type"`
	Info        struct {
		Contributor string `json:"contributor"`
		DateCreated string `json:"date_created"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Version     string `json:"version"`
		Year        int    `json:"year"`
	} `json:"info"`
	License struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"license"`
}

func Parse(pathAnnotations, pathQuestions string, imageIds []int) []int {

	var (
		val T
		err error
	)
	if err = ReadDataFromFile(pathAnnotations, &val); err != nil {
		panic(err)
	}

	var tmp T
	if err = ReadDataFromFile(pathQuestions, &tmp); err != nil {
		panic(err)
	}

	val.Questions = tmp.Questions

	recordInvalidImageIds, joinData := JointImagIdPointed(&val, imageIds)

	printImagIdPointed(path.Base(pathAnnotations), path.Base(pathQuestions), &val, &tmp, joinData)

	return recordInvalidImageIds
}

func printImagIdPointed(fa, fq string, val, tmp *T, m map[int]*JoinAQ) {
	var (
		ids   []int
		anno  T
		quest T
	)

	// 只保留证书/信息
	anno = *val
	anno.TaskType = ""
	anno.Questions = nil
	anno.Annotations = nil

	quest = *tmp
	//quest.TaskType = ""
	quest.Questions = nil
	quest.Annotations = nil

	var keys []int
	for k, _ := range m {
		keys = append(keys, k)
	}
	// 排序
	sort.Ints(keys)
	for _, k := range keys {
		ids = append(ids, k)
		anno.Annotations = append(anno.Annotations, m[k].Annotations)
		quest.Questions = append(quest.Questions, m[k].Questions)
	}

	fmt.Println("===============打印图片IDs")
	Write("v2_imgIds", printIds(ids))

	fmt.Println("===============Anno")
	v, _ := json.Marshal(anno)
	Write(fa, v)

	fmt.Println("===============quest")
	vv, _ := json.Marshal(quest)
	Write(fq, vv)
}

func printIds(ids []int) []byte {
	// [1,2,3]
	var (
		buf bytes.Buffer
		i   int
	)

	buf.WriteString("[ ")

	for _, e := range ids {
		buf.WriteString(fmt.Sprintf("%d", e))

		i++

		if i < len(ids)-1 {
			buf.WriteString(", ")
		}
	}

	buf.WriteString(" ]")

	return buf.Bytes()
}

func Write(fileName string, data []byte) {
	fp, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	_, err = fp.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func JointImagIdPointed(val *T, imageIds []int) ([]int, map[int]*JoinAQ) {

	fmt.Println()
	var (
		img2Anno = make(map[int]*Annotation)
		img2Ques = make(map[int]*Question)

		valImg = make(map[int]*JoinAQ)
	)

	for _, e := range val.Annotations {
		img2Anno[e.ImageID] = e
	}

	for _, e := range val.Questions {
		img2Ques[e.ImageID] = e
	}

	var recordInvalidImageIds []int
	for _, e := range imageIds {
		// 不存在直接返回
		if _, ok := img2Ques[e]; !ok {
			recordInvalidImageIds = append(recordInvalidImageIds, e)
			continue
		}
		if _, ok := img2Anno[e]; !ok {
			recordInvalidImageIds = append(recordInvalidImageIds, e)
			continue
		}

		if _, ok := valImg[e]; !ok {
			valImg[e] = new(JoinAQ)
		}
		valImg[e].Questions = img2Ques[e]
		valImg[e].Annotations = img2Anno[e]
	}

	fmt.Println("-------value lenght:", len(valImg))

	fmt.Println()

	return recordInvalidImageIds, valImg
}

func printTenElement(val *T) {

	fmt.Println()

	if len(val.Annotations) > 0 {
		waitP, _ := json.Marshal(val.Annotations[:10])
		fmt.Println(string(waitP))
	}

	fmt.Println()

	if len(val.Questions) > 0 {
		waitPQues, _ := json.Marshal(val.Questions[:10])
		fmt.Println(string(waitPQues))
	}
}

func ReadDataFromFile(path string, val interface{}) (err error) {

	// 打开json文件
	var jsonFile *os.File

	jsonFile, err = os.Open(path)

	// 最好要处理以下错误
	if err != nil {
		fmt.Println(err)
		return
	}

	// 要记得关闭
	defer jsonFile.Close()

	var byteValue []byte

	byteValue, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(byteValue, val)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
