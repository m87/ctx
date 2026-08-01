import { HttpErrorResponse } from '@angular/common/http';
import { toast } from 'ngx-sonner';

export interface ApiErrorResponse {
  code: string;
  description: string;
}

export interface AppError {
  code?: string;
  status?: number;
  message: string;
}

const unknownErrorMessage = 'An unknown error occurred';
const connectionErrorMessage =
  'Unable to connect to the server. Check your connection and try again.';
const maxQueryRetries = 2;

export function normalizeError(error: unknown): AppError {
  if (error instanceof HttpErrorResponse) {
    if (isConnectionError(error)) {
      return {
        status: error.status,
        message: connectionErrorMessage,
      };
    }

    if (isApiErrorResponse(error.error)) {
      return {
        code: error.error.code,
        status: error.status,
        message: error.error.description,
      };
    }

    const responseMessage = getNonEmptyString(error.error);
    return {
      status: error.status,
      message: responseMessage ?? getNonEmptyString(error.message) ?? unknownErrorMessage,
    };
  }

  if (error instanceof Error) {
    return { message: getNonEmptyString(error.message) ?? unknownErrorMessage };
  }

  return { message: getNonEmptyString(error) ?? unknownErrorMessage };
}

export function hasApiErrorCode(error: unknown, code: string): boolean {
  return error instanceof HttpErrorResponse && isApiErrorResponse(error.error)
    ? error.error.code === code
    : false;
}

export function shouldRetryQuery(failureCount: number, error: unknown): boolean {
  if (!(error instanceof HttpErrorResponse) || failureCount >= maxQueryRetries) {
    return false;
  }

  return error.status === 0 || error.status === 408 || error.status === 429 || error.status >= 500;
}

export function isConnectionError(error: unknown): boolean {
  if (!(error instanceof HttpErrorResponse)) {
    return false;
  }

  if (error.status === 0) {
    return true;
  }

  if (isApiErrorResponse(error.error)) {
    return false;
  }

  if (error.status === 502 || error.status === 503 || error.status === 504) {
    return true;
  }

  return error.status === 500 && getNonEmptyString(error.error) === undefined;
}

export function toastError(error: unknown): void {
  const normalizedError = normalizeError(error);
  const toastId =
    normalizedError.code ?? (isConnectionError(error) ? 'server-connection-error' : undefined);

  toast.error(normalizedError.message, toastId ? { id: toastId } : undefined);
}

function isApiErrorResponse(value: unknown): value is ApiErrorResponse {
  if (typeof value !== 'object' || value === null) {
    return false;
  }

  const response = value as Record<string, unknown>;
  return (
    getNonEmptyString(response['code']) !== undefined &&
    getNonEmptyString(response['description']) !== undefined
  );
}

function getNonEmptyString(value: unknown): string | undefined {
  if (typeof value !== 'string') {
    return undefined;
  }

  const normalized = value.trim();
  return normalized.length > 0 ? normalized : undefined;
}
