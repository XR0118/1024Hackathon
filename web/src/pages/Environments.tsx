import React, { useEffect, useState, useCallback } from 'react'
import { Card, CardBody, Button, Table, TableHeader, TableColumn, TableBody, TableRow, TableCell, Chip } from '@heroui/react'
import { environmentApi } from '@/services/api'
import type { Environment } from '@/types'
import { Cloud, Plus } from 'lucide-react'
import { useErrorStore } from '@/store/error'

const Environments: React.FC = () => {
  const { setError } = useErrorStore();
  const [environments, setEnvironments] = useState<Environment[]>([])

  const loadEnvironments = useCallback(async () => {
    try {
      const data = await environmentApi.list()
      setEnvironments(data)
    } catch (error) {
      setError('Failed to load environments.')
    }
  }, [setError])

  useEffect(() => {
    loadEnvironments()
  }, [loadEnvironments])

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-3xl font-bold">环境管理</h2>
        <Button color="primary" startContent={<Plus size={16} />}>
          添加环境
        </Button>
      </div>

      <Card>
        <CardBody>
          <Table aria-label="环境列表">
            <TableHeader>
              <TableColumn>环境名称</TableColumn>
              <TableColumn>类型</TableColumn>
              <TableColumn>状态</TableColumn>
              <TableColumn>应用数量</TableColumn>
              <TableColumn>操作</TableColumn>
            </TableHeader>
            <TableBody>
              {environments.map((env) => (
                <TableRow key={env.id}>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <Cloud size={16} />
                      {env.name}
                    </div>
                  </TableCell>
                  <TableCell>
                    <Chip color={env.type === 'k8s' ? 'primary' : 'success'} variant="flat">
                      {env.type === 'k8s' ? 'Kubernetes' : '物理机'}
                    </Chip>
                  </TableCell>
                  <TableCell>
                    <Chip color={env.status === 'active' ? 'success' : 'default'} variant="flat">
                      {env.status === 'active' ? '运行中' : '已停止'}
                    </Chip>
                  </TableCell>
                  <TableCell>{env.applicationCount}</TableCell>
                  <TableCell>
                    <Button size="sm" color="primary" variant="light">
                      查看详情
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardBody>
      </Card>
    </div>
  )
}

export default Environments
