import React from 'react';

interface ErrorMessageProps {
  message: string;
  onDismiss: () => void;
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({ message, onDismiss }) => {
  if (!message) {
    return null;
  }

  return (
    <div className="alert alert-danger alert-dismissible" role="alert">
      {message}
      <button type="button" className="btn-close" onClick={onDismiss} aria-label="Close"></button>
    </div>
  );
};

export default ErrorMessage;
