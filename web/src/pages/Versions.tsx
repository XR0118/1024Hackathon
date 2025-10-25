import React, { useEffect, useState, useCallback } from 'react'
import { versionApi } from '@/services/api'
import { formatDate } from '@/utils'
import type { Version } from '@/types'
import { Search, RefreshCw } from 'lucide-react'
import { useErrorStore } from '@/store/error'
import {
  Button,
  Input,
  Card,
  CardBody,
  Table,
  TableHeader,
  TableColumn,
  TableBody,
  TableRow,
  TableCell,
  Chip,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  useDisclosure,
} from '@heroui/react'

const Versions: React.FC = () => {
  const { setError } = useErrorStore();
  const [versions, setVersions] = useState<Version[]>([])
  const [loading, setLoading] = useState(false)
  const [searchText, setSearchText] = useState('')
  const [selectedVersion, setSelectedVersion] = useState<Version | null>(null)
  const { isOpen, onOpen, onClose } = useDisclosure()

  const loadVersions = useCallback(async () => {
    setLoading(true)
    try {
      const data = await versionApi.list({
        search: searchText || undefined,
      })
      setVersions(data)
    } catch (error) {
      setError('Failed to load versions.')
    } finally {
      setLoading(false)
    }
  }, [searchText, setError])

  useEffect(() => {
    const timer = setTimeout(() => {
      loadVersions()
    }, 500)
    return () => clearTimeout(timer)
  }, [loadVersions])

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchText(e.target.value)
  }

  const openVersionDetail = (version: Version) => {
    setSelectedVersion(version)
    onOpen()
  }

  const getHealthColor = (health: number) => {
    if (health >= 80) return 'success'
    if (health >= 50) return 'warning'
    return 'danger'
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-3xl font-bold">版本管理</h2>
      </div>

      <Card>
        <CardBody>
          <div className="flex gap-2 mb-4">
            <Input
              type="text"
              placeholder="搜索版本号或标签..."
              value={searchText}
              onChange={handleSearchChange}
              startContent={<Search className="w-4 h-4" />}
              className="flex-1"
            />
            <Button
              color="primary"
              onPress={loadVersions}
              isDisabled={loading}
              startContent={<RefreshCw className="w-4 h-4" />}
            >
              刷新
            </Button>
          </div>

          <Table aria-label="版本列表">
            <TableHeader>
              <TableColumn>版本号</TableColumn>
              <TableColumn>Git Tag</TableColumn>
              <TableColumn>应用信息</TableColumn>
              <TableColumn>创建时间</TableColumn>
              <TableColumn>操作</TableColumn>
            </TableHeader>
            <TableBody>
              {versions.map((version) => (
                <TableRow key={version.version}>
                  <TableCell>{version.version}</TableCell>
                  <TableCell>
                    <a
                      href={`https://github.com/XR0118/1024Hackathon/releases/tag/${version.git.tag}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-primary hover:underline"
                    >
                      {version.git.tag}
                    </a>
                  </TableCell>
                  <TableCell>
                    <div className="flex flex-col gap-1">
                      {version.applications.map((app) => (
                        <div key={app.name} className="flex items-center gap-2">
                          <Chip color="secondary" variant="flat" size="sm">{app.name}</Chip>
                          <small className="text-gray-500">
                            覆盖度: {app.coverage}% | 健康度: {app.health}% | 
                            更新: {formatDate(app.lastUpdatedAt)}
                          </small>
                        </div>
                      ))}
                    </div>
                  </TableCell>
                  <TableCell>{formatDate(version.createdAt)}</TableCell>
                  <TableCell>
                    <Button
                      size="sm"
                      variant="light"
                      color="primary"
                      onPress={() => openVersionDetail(version)}
                    >
                      详情
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardBody>
      </Card>

      <Modal isOpen={isOpen} onClose={onClose} size="2xl">
        <ModalContent>
          <ModalHeader>版本详情</ModalHeader>
          <ModalBody className="pb-6">
            {selectedVersion && (
              <div className="space-y-4">
                <div>
                  <h3 className="text-lg font-semibold mb-2">基本信息</h3>
                  <div className="space-y-2">
                    <p><strong>版本号:</strong> {selectedVersion.version}</p>
                    <p><strong>Git Tag:</strong> {selectedVersion.git.tag}</p>
                    <p><strong>创建时间:</strong> {formatDate(selectedVersion.createdAt)}</p>
                  </div>
                </div>

                <div>
                  <h3 className="text-lg font-semibold mb-2">应用信息</h3>
                  <div className="space-y-2">
                    {selectedVersion.applications.map((app) => (
                      <Card key={app.name} shadow="sm">
                        <CardBody>
                          <div className="flex justify-between items-center">
                            <strong>{app.name}</strong>
                            <div className="flex flex-col gap-1 text-sm">
                              <div>
                                <span className="text-gray-500">覆盖度:</span>{' '}
                                <Chip color="primary" variant="flat" size="sm">{app.coverage}%</Chip>
                              </div>
                              <div>
                                <span className="text-gray-500">健康度:</span>{' '}
                                <Chip color={getHealthColor(app.health)} variant="flat" size="sm">
                                  {app.health}%
                                </Chip>
                              </div>
                              <div>
                                <span className="text-gray-500">最后更新:</span>{' '}
                                {formatDate(app.lastUpdatedAt)}
                              </div>
                            </div>
                          </div>
                        </CardBody>
                      </Card>
                    ))}
                  </div>
                </div>
              </div>
            )}
          </ModalBody>
        </ModalContent>
      </Modal>
    </div>
  )
}

export default Versions
