package kdb

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"

	_ "github.com/kweaver-ai/proton-rds-sdk-go/driver"
)

type blobType int32

const (
	creatorIDLength     = 36 // 创建者ID长度，固定36字节
	versionDataPathSize = 64 // 版本数据存在位置，前32位为版本ID，后32位为对象存储ID
	uuidLen             = 32 // uuid 长度，不带中划线

	versionBlob blobType = 3 // customMetadata 中的类型，3：版本对象 0：回收站对象 -1：未知
)

type VersionCustomMetadata struct {
	CreatorID          string // 创建者ID 36字节
	FileName           string // 文件名
	FileSize           int64  // 文件大小，8字节
	ClientModifiedTime int64  // 客户端本地修改时间，8字节

	// 对象存储结构：https://127.0.0.1:443/test/8990cff5-2987-4eef-8839-0e898a2b488a/E91D39890BD7445B9D88BD09AE1637A2/22B710E9EFF24C91B992573EE6D4B066
	// 第一级：ossID: 对于 metadata 表，代表存储id (8990cff5-2987-4eef-8839-0e898a2b488a)
	// 第二级：cid 文档库ID，对象存储的具体位置之一 (E91D39890BD7445B9D88BD09AE1637A2)
	// 第三级：versionID 版本ID (22B710E9EFF24C91B992573EE6D4B066)
	// 即对象存储中，用户名是 bucket，文档库是文件夹，版本id是文件名
	// 注意：因为文件可以移动，其CID可能会发生变化，但存储中其位置仍然未变更，故下载等相关操作始终需要从本字段中获取第二级路径
	VersionID string // 版本ID，32字节
	DoclibID  string // 文档库ID，32字节
}

func serializeVersionBlob(info VersionCustomMetadata) ([]byte, error) {
	var buf bytes.Buffer

	if len(info.CreatorID) != creatorIDLength || len(info.VersionID) != uuidLen || len(info.DoclibID) != uuidLen {
		return nil, errors.New("invalid CreatorID or VersionID or DoclibID")
	}

	// 写入类型：4字节
	if err := binary.Write(&buf, binary.BigEndian, int32(versionBlob)); err != nil {
		return nil, err
	}

	// 写入创建者ID 4 + 36字节
	// creatorID := []byte("123e4567-e89b-12d3-a456-426614174000") // 示例创建者ID
	if err := writeString(&buf, info.CreatorID); err != nil {
		return nil, err
	}

	// 文件名，4字节 + N 字节
	if err := writeString(&buf, info.FileName); err != nil {
		return nil, err
	}

	// 版本数据大小（文件大小），8字节
	if err := binary.Write(&buf, binary.BigEndian, info.FileSize); err != nil {
		return nil, err
	}

	// 客户端本地修改时间，8字节
	if err := binary.Write(&buf, binary.BigEndian, info.ClientModifiedTime); err != nil {
		return nil, err
	}

	// 版本数据存放位置，4 + 64字节，前32字节存放版本ID，后32字节存放所在文档库ID
	if err := writeString(&buf, info.VersionID+info.DoclibID); err != nil {
		return nil, err
	}

	// 最后，在缓冲区的前8个字节处，写入 body 长度，即：
	// 存放内容总大小，8 字节（！！！！注意，不包括保留位的4字节）
	// 版本号(值为3，占4字节) + "创建者ID"长度(值为36，占4字节) + 创建者ID(占36字节) + 文件名长度(值为X，占4字节) +
	// 文件名(占X字节) + 版本数据大小(占8字节) + 客户端本地修改时间(占8字节) + "版本数据存放位置"长度(值为64，占4字节) +
	// 版本数据存放位置(占64字节，前32字节存放版本ID，后32字节存放所在文档库ID) + 保留字段个数(占4字节)
	var full bytes.Buffer
	body := buf.Bytes()

	// 写入总大小，8字节
	if err := binary.Write(&full, binary.BigEndian, uint64(len(body))); err != nil {
		return nil, err
	}
	// 写入 body
	if err := binary.Write(&full, binary.BigEndian, body); err != nil {
		return nil, err
	}
	// 最后，写入保留字段个数，4字节
	if err := binary.Write(&full, binary.BigEndian, int32(0)); err != nil {
		return nil, err
	}
	return full.Bytes(), nil
}

// writeString 写入 string
func writeString(out *bytes.Buffer, content string) error {
	if err := binary.Write(out, binary.BigEndian, uint32(len(content))); err != nil {
		return err
	}
	if _, err := out.WriteString(content); err != nil {
		return err
	}
	return nil
}

const (
	v0 = 0
	v1 = 1
	v2 = 2
	v3 = 3
)

func deserializeVersionBlob(data []byte) (*VersionCustomMetadata, error) {
	buf := bytes.NewReader(data)
	vo := &VersionCustomMetadata{}

	// 读取内容总大小，8字节
	contentSize := uint64(0)
	if err := binary.Read(buf, binary.BigEndian, &contentSize); err != nil {
		return nil, err
	}

	// 读取内容
	bodyBuffer := make([]byte, contentSize)
	if err := binary.Read(buf, binary.BigEndian, &bodyBuffer); err != nil {
		return nil, err
	}
	body := bytes.NewReader(bodyBuffer)

	// 从 buf 中读取4字节保留位，这里没用到，直接省略
	// reserved := int32(0)
	// binary.Read(buf, binary.BigEndian, &reserved)

	// 解析内容
	{
		// 读取4字节，得到协议版本号
		protocolVersion := int32(0)
		if err := binary.Read(body, binary.BigEndian, &protocolVersion); err != nil {
			return nil, err
		}

		// 解决没有默认值导致大小错误的问题
		vo.FileSize = -2
		vo.ClientModifiedTime = 0

		// 读取版本0
		if protocolVersion >= v0 {
			// 读取创建者ID
			creatorID, err := readString(body)
			if err != nil {
				return nil, err
			}
			if len(creatorID) != creatorIDLength {
				return nil, errors.New("bad data")
			}

			// 读取文件名
			fileName, err := readString(body)
			if err != nil {
				return nil, err
			}

			vo.CreatorID = creatorID
			vo.FileName = fileName
		}

		// 读取版本1
		if protocolVersion >= v1 {
			// 读取8字节，文件大小
			if err := binary.Read(body, binary.BigEndian, &vo.FileSize); err != nil {
				return nil, err
			}
		}

		// 读取版本2
		if protocolVersion >= v2 {
			// 读取8字节，客户端本地修改时间
			if err := binary.Read(body, binary.BigEndian, &vo.ClientModifiedTime); err != nil {
				return nil, err
			}
		}

		// 读取版本3
		if protocolVersion >= v3 {
			// 读取4字节，值固定为64，versionID 和 ossID 分别占 32 字节
			versionAndOss, err := readString(body)
			if err != nil {
				return nil, err
			}

			if len(versionAndOss) != versionDataPathSize {
				return nil, errors.New("bad data")
			}

			vo.VersionID = versionAndOss[:32]
			vo.DoclibID = versionAndOss[32:]
		}
	}
	return vo, nil
}

// readString 读取 string
func readString(in io.Reader) (string, error) {
	size := uint32(0)
	if err := binary.Read(in, binary.BigEndian, &size); err != nil {
		return "", err
	}

	strBuffer := make([]byte, size)
	if _, err := in.Read(strBuffer); err != nil {
		return "", err
	}

	return string(strBuffer), nil
}

func TestBlob(op *sql.DB) {
	tableSql := "CREATE TABLE IF NOT EXISTS `test_blob`(" +
		"`id` INT," +
		"`blob` BYTEA" +
		")"
	_, err := op.Exec(tableSql)
	if err != nil {
		fmt.Println(err)
		return
	}

	orgMetadata := VersionCustomMetadata{
		CreatorID:          "f4e6adaa-169b-11f0-9549-a2cac61ed55c",
		FileName:           "a.text",
		FileSize:           1024,
		ClientModifiedTime: 1744351958210069,
		VersionID:          "50D96DA9494C4B4281FE60A65439082C",
		DoclibID:           "CAC97EBD4CA1402194F6525410B5DE75",
	}

	blob, err := serializeVersionBlob(orgMetadata)
	if err != nil {
		fmt.Println(err)
		return
	}

	type st struct {
		id int64
		b  []byte
	}

	os := st{
		id: 1,
		b:  blob,
	}
	_, err = op.Exec("INSERT INTO `test_blob` VALUES(?, ?)", os.id, os.b)
	if err != nil {
		fmt.Println(err)
		return
	}

	row := op.QueryRow("SELECT `id`, `blob` FROM `test_blob` WHERE id=?", 1)
	ns := st{}
	err = row.Scan(
		&ns.id,
		&ns.b,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	newMetadata, _ := deserializeVersionBlob(ns.b)
	if ns.id != os.id ||
		orgMetadata.CreatorID != newMetadata.CreatorID ||
		orgMetadata.FileName != newMetadata.FileName ||
		orgMetadata.FileSize != newMetadata.FileSize ||
		orgMetadata.ClientModifiedTime != newMetadata.ClientModifiedTime ||
		orgMetadata.VersionID != newMetadata.VersionID ||
		orgMetadata.DoclibID != newMetadata.DoclibID {
		log.Fatalf("data not match: new: %v, org: %v", ns, os)
	}

	fmt.Println("success")
}
