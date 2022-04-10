
# 筛选vqa数据集

> 此项目目的是为了生成正确数据：跑通[TRAR-VQA](https://github.com/zput/TRAR-VQA)深度学习代码

## 使用

```makefile
build: xos win

xos:
	CGO_ENABLED=0 GOOS=darwin go build -o parse-json main.go

win:
	CGO_ENABLED=0 GOOS=windows go build -o parse-json.exe main.go

train: xos
	./parse-json \
          -fa "/Users/edz/Desktop/DL/test-back/vqa/raw/v2_mscoco_train2014_annotations.json" \
          -fq "/Users/edz/Desktop/DL/test-back/vqa/raw/v2_OpenEnded_mscoco_train2014_questions.json" \
          -dele "/Users/edz/code/github/DL/TRAR-VQA/data/vqa/feats" \
          > result.json

val: xos
	./parse-json \
          -fa "/Users/edz/Desktop/DL/test-back/vqa/raw/v2_mscoco_val2014_annotations.json" \
          -fq "/Users/edz/Desktop/DL/test-back/vqa/raw/v2_OpenEnded_mscoco_val2014_questions.json" \
          -dele "/Users/edz/code/github/DL/TRAR-VQA/data/vqa/feats" \
          > result.json

test: train val
```


