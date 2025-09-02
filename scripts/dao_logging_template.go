// DAO SQL日志添加模板
// 这个文件展示了如何为DAO方法添加SQL日志的标准模式

package dao

import (
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

// 示例：查询操作
func (imp *ExampleInterfaceImp) GetExampleById(id int32) (*model.ExampleModel, error) {
	var example = new(model.ExampleModel)
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", "ExampleTable", map[string]interface{}{
		"id": id,
	})

	err := cli.Table("ExampleTable").Where("id = ?", id).First(example).Error
	logger.LogQuery(example, err)

	return example, err
}

// 示例：插入操作
func (imp *ExampleInterfaceImp) CreateExample(example *model.ExampleModel) error {
	cli := db.Get()
	example.CreatedAt = time.Now()
	example.UpdatedAt = time.Now()

	// 记录SQL操作日志
	logger := NewSQLLogger("插入", "ExampleTable", map[string]interface{}{
		"name": example.Name,
		"type": example.Type,
	})

	err := cli.Table("ExampleTable").Create(example).Error
	logger.LogInsert(example, err)

	return err
}

// 示例：更新操作
func (imp *ExampleInterfaceImp) UpdateExample(example *model.ExampleModel) error {
	cli := db.Get()
	example.UpdatedAt = time.Now()

	// 记录SQL操作日志
	logger := NewSQLLogger("更新", "ExampleTable", map[string]interface{}{
		"id":   example.Id,
		"name": example.Name,
	})

	result := cli.Table("ExampleTable").Where("id = ?", example.Id).Updates(example)
	logger.LogUpdate(result.RowsAffected, result.Error)

	return result.Error
}

// 示例：删除操作
func (imp *ExampleInterfaceImp) DeleteExample(id int32) error {
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("删除", "ExampleTable", map[string]interface{}{
		"id": id,
	})

	result := cli.Table("ExampleTable").Where("id = ?", id).Delete(&model.ExampleModel{})
	logger.LogDelete(result.RowsAffected, result.Error)

	return result.Error
}

// 示例：计数操作
func (imp *ExampleInterfaceImp) CountExamplesByType(typeName string) (int64, error) {
	var count int64
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("计数", "ExampleTable", map[string]interface{}{
		"type": typeName,
	})

	err := cli.Table("ExampleTable").Where("type = ?", typeName).Count(&count).Error
	logger.LogCount(count, err)

	return count, err
}

// 示例：分页查询操作
func (imp *ExampleInterfaceImp) GetExamplesByType(typeName string, page, pageSize int) ([]*model.ExampleModel, int64, error) {
	var examples []*model.ExampleModel
	var total int64
	cli := db.Get()

	// 记录SQL操作日志
	logger := NewSQLLogger("查询", "ExampleTable", map[string]interface{}{
		"type":     typeName,
		"page":     page,
		"pageSize": pageSize,
	})

	// 获取总数
	err := cli.Table("ExampleTable").Where("type = ?", typeName).Count(&total).Error
	if err != nil {
		logger.LogCount(0, err)
		return nil, 0, err
	}
	logger.LogCount(total, nil)

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = cli.Table("ExampleTable").
		Where("type = ?", typeName).
		Order("createdAt DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&examples).Error

	logger.LogQuery(examples, err)
	return examples, total, err
}
