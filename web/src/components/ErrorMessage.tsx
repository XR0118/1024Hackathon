import React from 'react';
import { Card, CardBody, Button } from '@heroui/react';
import { X } from 'lucide-react';
import DOMPurify from 'dompurify';

interface ErrorMessageProps {
  message: string;
  onDismiss: () => void;
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({ message, onDismiss }) => {
  if (!message) {
    return null;
  }

  const sanitizedMessage = DOMPurify.sanitize(message);

  return (
    <Card className="bg-danger-50 border-danger mb-4">
      <CardBody className="flex flex-row items-start justify-between gap-3">
        <div 
          className="flex-1 text-danger" 
          dangerouslySetInnerHTML={{ __html: sanitizedMessage }} 
        />
        <Button
          isIconOnly
          size="sm"
          variant="light"
          color="danger"
          onPress={onDismiss}
          aria-label="Close"
        >
          <X size={18} />
        </Button>
      </CardBody>
    </Card>
  );
};

export default ErrorMessage;
