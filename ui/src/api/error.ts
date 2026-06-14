import { HttpErrorResponse } from '@angular/common/http';
import { toast } from 'ngx-sonner';

export interface ErrorResponse {
  code: string;
  description: string;
}

export function toastError(error: unknown): void {
  if (error instanceof HttpErrorResponse) {
    const response = error.error as Partial<ErrorResponse> | null;
    toast.error(response?.description ?? error.message);
  } else if (error instanceof Error) {
    toast.error(error.message);
  } else if (typeof error === 'string') {
    toast.error(error);
  } else {
    toast.error('An unknown error occurred');
  }
}
