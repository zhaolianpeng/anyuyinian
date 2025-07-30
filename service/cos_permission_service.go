package service

import (
	"context"
	"fmt"
	"wxcloudrun-golang/config"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// COSPermissionService COS权限管理服务
type COSPermissionService struct {
	client *cos.Client
}

// NewCOSPermissionService 创建COS权限管理服务
func NewCOSPermissionService() *COSPermissionService {
	return &COSPermissionService{
		client: config.GetCOSClient(),
	}
}

// SetObjectPublicRead 设置对象为公共读取权限（所有用户可读，仅创建者可读写）
func (s *COSPermissionService) SetObjectPublicRead(objectKey string) error {
	LogStep("设置对象为公共读取权限", map[string]string{"objectKey": objectKey})

	// 设置ACL为public-read
	_, err := s.client.Object.PutACL(context.Background(), objectKey, &cos.ObjectPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosACL: "public-read",
		},
	})

	if err != nil {
		LogError("设置对象公共读取权限失败", err)
		return fmt.Errorf("设置对象公共读取权限失败: %v", err)
	}

	LogStep("对象公共读取权限设置成功", map[string]string{"objectKey": objectKey})
	return nil
}

// SetObjectPrivate 设置对象为私有权限（仅创建者可读写）
func (s *COSPermissionService) SetObjectPrivate(objectKey string) error {
	LogStep("设置对象为私有权限", map[string]string{"objectKey": objectKey})

	// 设置ACL为private
	_, err := s.client.Object.PutACL(context.Background(), objectKey, &cos.ObjectPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosACL: "private",
		},
	})

	if err != nil {
		LogError("设置对象私有权限失败", err)
		return fmt.Errorf("设置对象私有权限失败: %v", err)
	}

	LogStep("对象私有权限设置成功", map[string]string{"objectKey": objectKey})
	return nil
}

// GetObjectACL 获取对象ACL权限
func (s *COSPermissionService) GetObjectACL(objectKey string) (string, error) {
	LogStep("获取对象ACL权限", map[string]string{"objectKey": objectKey})

	result, _, err := s.client.Object.GetACL(context.Background(), objectKey)
	if err != nil {
		LogError("获取对象ACL权限失败", err)
		return "", fmt.Errorf("获取对象ACL权限失败: %v", err)
	}

	acl := "private"
	if len(result.AccessControlList) > 0 {
		for _, grant := range result.AccessControlList {
			if grant.Grantee.URI == "http://cam.qcloud.com/groups/global/AllUsers" {
				acl = "public-read"
				break
			}
		}
	}

	LogStep("获取对象ACL权限成功", map[string]string{"objectKey": objectKey, "acl": acl})
	return acl, nil
}

// DeleteObject 删除对象
func (s *COSPermissionService) DeleteObject(objectKey string) error {
	LogStep("删除对象", map[string]string{"objectKey": objectKey})

	_, err := s.client.Object.Delete(context.Background(), objectKey)
	if err != nil {
		LogError("删除对象失败", err)
		return fmt.Errorf("删除对象失败: %v", err)
	}

	LogStep("对象删除成功", map[string]string{"objectKey": objectKey})
	return nil
}

// CheckObjectExists 检查对象是否存在
func (s *COSPermissionService) CheckObjectExists(objectKey string) (bool, error) {
	LogStep("检查对象是否存在", map[string]string{"objectKey": objectKey})

	_, err := s.client.Object.Head(context.Background(), objectKey, nil)
	if err != nil {
		// 如果对象不存在，会返回404错误
		if cos.IsNotFoundError(err) {
			LogStep("对象不存在", map[string]string{"objectKey": objectKey})
			return false, nil
		}
		LogError("检查对象存在性失败", err)
		return false, fmt.Errorf("检查对象存在性失败: %v", err)
	}

	LogStep("对象存在", map[string]string{"objectKey": objectKey})
	return true, nil
}
