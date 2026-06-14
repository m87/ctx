import { HttpErrorResponse } from "@angular/common/http";
import { toast } from "ngx-sonner";

export interface ErrorResponse {
  code: string;
  description: string;
}

export function toastError(error: unknown): void {
  if (error instanceof Error) {
    toast.error(error.message);
  } else if (typeof error === 'string') {
    toast.error(error);
  } else if (error instanceof HttpErrorResponse) {
    toast.error(`${error.error.description}`);
    return;
  } else {
    toast.error('An unknown error occurred');
  }
}
