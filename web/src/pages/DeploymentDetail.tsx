import React, { useEffect, useState, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { deploymentApi } from '@/services/api'
import { formatDate, formatDuration, getStatusColor, getStatusText } from '@/utils'
import type { DeploymentDetail as DeploymentDetailType } from '@/types'
import { ArrowLeft, Check, Undo2 } from 'lucide-react'
import { useErrorStore } from '@/store/error'
import DOMPurify from 'dompurify'
import {
  Button,
  Card,
  CardBody,
  CardHeader,
  Chip,
  Textarea,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  useDisclosure,
  Slider,
} from '@heroui/react'

const DeploymentDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [deployment, setDeployment] = useState<DeploymentDetailType | null>(null)
  const [loading, setLoading] = useState(false)
  const [note, setNote] = useState('')
  const [reason, setReason] = useState('')
  const { isOpen: isConfirmOpen, onOpen: onConfirmOpen, onClose: onConfirmClose } = useDisclosure()
  const { isOpen: isRollbackOpen, onOpen: onRollbackOpen, onClose: onRollbackClose } = useDisclosure()

  const loadDeployment = useCallback(async () => {
    if (!id) return
    setLoading(true)
    try {
      const data = await deploymentApi.get(id)
      setDeployment(data)
    } catch (error) {
      useErrorStore.getState().setError('Failed to load deployment details.')
    } finally {
      setLoading(false)
    }
  }, [id])

  useEffect(() => {
    loadDeployment()
    const interval = setInterval(loadDeployment, 3000)
    return () => clearInterval(interval)
  }, [id, loadDeployment])

  const handleConfirm = async () => {
    if (!id) return
    try {
      await deploymentApi.confirm(id, note)
      setNote('')
      loadDeployment()
      onConfirmClose()
    } catch (error) {
      useErrorStore.getState().setError('Failed to confirm deployment.')
    }
  }

  const handleRollback = async () => {
    if (!id) return
    try {
      await deploymentApi.rollback(id, reason)
      setReason('')
      loadDeployment()
      onRollbackClose()
    } catch (error) {
      useErrorStore.getState().setError('Failed to rollback deployment.')
    }
  }

  if (loading && !deployment) {
    return <div className="flex justify-center items-center min-h-screen">Loading...</div>
  }

  if (!deployment) {
    return <div className="flex justify-center items-center min-h-screen">Deployment not found</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Button
          isIconOnly
          variant="light"
          onPress={() => navigate('/deployments')}
        >
          <ArrowLeft className="w-5 h-5" />
        </Button>
        <h2 className="text-3xl font-bold">部署详情</h2>
      </div>

      <Card>
        <CardBody>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <p><strong>部署ID:</strong> {deployment.id}</p>
              <p><strong>版本:</strong> {deployment.version}</p>
              <p><strong>应用:</strong> {deployment.applications.join(', ')}</p>
            </div>
            <div className="space-y-2">
              <p>
                <strong>状态:</strong>{' '}
                <Chip color={getStatusColor(deployment.status) as any} variant="flat">
                  {getStatusText(deployment.status)}
                </Chip>
              </p>
              <p><strong>环境:</strong> {deployment.environments.join(', ')}</p>
              <p><strong>创建时间:</strong> {formatDate(deployment.createdAt)}</p>
            </div>
            {deployment.duration && (
              <div className="col-span-full">
                <p><strong>执行时长:</strong> {formatDuration(deployment.duration)}</p>
              </div>
            )}
          </div>
        </CardBody>
      </Card>

      {deployment.status === 'waiting_confirm' && (
        <Card className="bg-warning-50">
          <CardBody>
            <div className="flex justify-between items-center">
              <div>
                <strong>需要人工确认:</strong> 此部署需要人工确认后才能继续。
              </div>
              <div className="flex gap-2">
                <Button
                  color="success"
                  onPress={onConfirmOpen}
                  startContent={<Check className="w-4 h-4" />}
                >
                  确认继续
                </Button>
                <Button
                  color="danger"
                  onPress={onRollbackOpen}
                  startContent={<Undo2 className="w-4 h-4" />}
                >
                  回滚
                </Button>
              </div>
            </div>
          </CardBody>
        </Card>
      )}

      {deployment.grayscaleEnabled && (
        <Card>
          <CardHeader>
            <h3 className="text-xl font-semibold">灰度发布</h3>
          </CardHeader>
          <CardBody>
            <Slider
              label="当前灰度比例"
              value={deployment.grayscaleRatio}
              isDisabled
              maxValue={100}
              minValue={0}
              className="max-w-md"
              formatOptions={{ style: 'percent', maximumFractionDigits: 0 }}
            />
          </CardBody>
        </Card>
      )}

      <Card>
        <CardHeader>
          <h3 className="text-xl font-semibold">部署流程</h3>
        </CardHeader>
        <CardBody>
          <div className="flex gap-2 flex-wrap">
            {deployment.steps.map((step, index) => (
              <Chip
                key={index}
                color={step.status === 'success' ? 'success' : 'default'}
                variant={step.status === 'success' ? 'flat' : 'bordered'}
              >
                {step.name}
              </Chip>
            ))}
          </div>
        </CardBody>
      </Card>

      <Card>
        <CardHeader>
          <h3 className="text-xl font-semibold">实时日志</h3>
        </CardHeader>
        <CardBody>
          <div className="bg-black text-green-500 font-mono text-xs p-4 rounded max-h-96 overflow-auto">
            {deployment.logs.map((log, index) => (
              <div key={index}>
                <span className="text-gray-600">[{log.timestamp}]</span>{' '}
                <span className={log.level === 'error' ? 'text-red-500' : log.level === 'warn' ? 'text-orange-500' : 'text-green-500'}>
                  [{log.level.toUpperCase()}]
                </span>{' '}
                <span dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(log.message) }} />
              </div>
            ))}
          </div>
        </CardBody>
      </Card>

      <Modal isOpen={isConfirmOpen} onClose={onConfirmClose}>
        <ModalContent>
          <ModalHeader>确认部署</ModalHeader>
          <ModalBody>
            <p className="mb-4">确认继续此部署吗?</p>
            <Textarea
              placeholder="备注(可选)"
              value={note}
              onValueChange={setNote}
              minRows={4}
            />
          </ModalBody>
          <ModalFooter>
            <Button variant="light" onPress={onConfirmClose}>
              取消
            </Button>
            <Button color="primary" onPress={handleConfirm}>
              确认
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>

      <Modal isOpen={isRollbackOpen} onClose={onRollbackClose}>
        <ModalContent>
          <ModalHeader>回滚部署</ModalHeader>
          <ModalBody>
            <p className="mb-4">确认回滚此部署吗?</p>
            <Textarea
              placeholder="回滚原因(可选)"
              value={reason}
              onValueChange={setReason}
              minRows={4}
            />
          </ModalBody>
          <ModalFooter>
            <Button variant="light" onPress={onRollbackClose}>
              取消
            </Button>
            <Button color="danger" onPress={handleRollback}>
              回滚
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </div>
  )
}

export default DeploymentDetailPage
