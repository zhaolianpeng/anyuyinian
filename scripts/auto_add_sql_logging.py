#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
自动为DAO文件添加SQL日志的Python脚本
这个脚本会扫描所有DAO文件并为SQL操作添加日志记录
"""

import os
import re
import sys
from pathlib import Path

class DAOLoggingAdder:
    def __init__(self, dao_dir="db/dao"):
        self.dao_dir = Path(dao_dir)
        self.processed_files = []
        self.skipped_files = []
        
    def should_skip_file(self, filename):
        """判断是否应该跳过某个文件"""
        skip_patterns = [
            "interface.go",  # 接口文件
            "dao.go",        # 主DAO文件
            "sql_logger.go", # 日志工具文件
        ]
        
        for pattern in skip_patterns:
            if pattern in filename:
                return True
        return False
    
    def add_sql_logging_to_file(self, filepath):
        """为单个DAO文件添加SQL日志"""
        print(f"处理文件: {filepath}")
        
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                content = f.read()
            
            # 检查是否已经包含SQL日志
            if "NewSQLLogger" in content:
                print(f"  - 文件已包含SQL日志，跳过")
                self.skipped_files.append(filepath.name)
                return
            
            # 检查是否包含sql_logger导入
            if "sql_logger" not in content:
                print(f"  - 需要添加sql_logger导入")
                # 这里可以添加导入逻辑，但为了安全起见，我们手动处理
            
            # 为常见的SQL操作模式添加日志
            modified_content = self.add_logging_patterns(content, filepath.name)
            
            if modified_content != content:
                # 备份原文件
                backup_path = filepath.with_suffix(filepath.suffix + '.backup')
                with open(backup_path, 'w', encoding='utf-8') as f:
                    f.write(content)
                print(f"  - 已创建备份文件: {backup_path}")
                
                # 写入修改后的内容
                with open(filepath, 'w', encoding='utf-8') as f:
                    f.write(modified_content)
                print(f"  - 已添加SQL日志")
                self.processed_files.append(filepath.name)
            else:
                print(f"  - 未找到需要添加日志的SQL操作")
                self.skipped_files.append(filepath.name)
                
        except Exception as e:
            print(f"  - 处理文件时出错: {e}")
            self.skipped_files.append(filepath.name)
    
    def add_logging_patterns(self, content, filename):
        """为内容添加日志模式"""
        # 这里可以添加更复杂的模式匹配和替换逻辑
        # 为了安全起见，我们只做基本的检查
        
        # 检查是否包含常见的SQL操作
        sql_patterns = [
            r'\.Table\([^)]+\)\.Where\([^)]+\)\.First\([^)]+\)\.Error',
            r'\.Table\([^)]+\)\.Where\([^)]+\)\.Find\([^)]+\)\.Error',
            r'\.Table\([^)]+\)\.Create\([^)]+\)\.Error',
            r'\.Table\([^)]+\)\.Where\([^)]+\)\.Updates\([^)]+\)\.Error',
            r'\.Table\([^)]+\)\.Where\([^)]+\)\.Count\([^)]+\)\.Error',
        ]
        
        has_sql_operations = any(re.search(pattern, content) for pattern in sql_patterns)
        
        if has_sql_operations:
            print(f"  - 发现SQL操作，需要手动添加日志")
            return content  # 返回原内容，需要手动处理
        
        return content
    
    def process_all_dao_files(self):
        """处理所有DAO文件"""
        if not self.dao_dir.exists():
            print(f"错误: DAO目录不存在: {self.dao_dir}")
            return
        
        dao_files = list(self.dao_dir.glob("*.go"))
        
        if not dao_files:
            print(f"在 {self.dao_dir} 中未找到Go文件")
            return
        
        print(f"找到 {len(dao_files)} 个Go文件")
        print("=" * 50)
        
        for filepath in dao_files:
            if self.should_skip_file(filepath.name):
                print(f"跳过文件: {filepath.name} (接口文件或工具文件)")
                continue
            
            self.add_sql_logging_to_file(filepath)
            print()
        
        self.print_summary()
    
    def print_summary(self):
        """打印处理摘要"""
        print("=" * 50)
        print("处理摘要:")
        print(f"已处理文件: {len(self.processed_files)}")
        for filename in self.processed_files:
            print(f"  ✓ {filename}")
        
        print(f"跳过文件: {len(self.skipped_files)}")
        for filename in self.skipped_files:
            print(f"  - {filename}")
        
        print("\n注意:")
        print("1. 此脚本只做了基本的检查，大部分文件需要手动添加SQL日志")
        print("2. 请参考 scripts/dao_logging_template.go 中的模板")
        print("3. 确保所有DAO文件都导入了sql_logger包")
        print("4. 运行测试确保功能正常")

def main():
    if len(sys.argv) > 1:
        dao_dir = sys.argv[1]
    else:
        dao_dir = "db/dao"
    
    adder = DAOLoggingAdder(dao_dir)
    adder.process_all_dao_files()

if __name__ == "__main__":
    main()
