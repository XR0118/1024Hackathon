import React from 'react';
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
    <div className="alert alert-danger alert-dismissible" role="alert">
      <div dangerouslySetInnerHTML={{ __html: sanitizedMessage }} />
      <button type="button" className="btn-close" onClick={onDismiss} aria-label="Close"></button>
    </div>
  );
};

export default ErrorMessage;
