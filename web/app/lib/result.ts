export type ValidationIssue = { field: string; tag: string }

export type Result<T> =
  | { ok: true; data: T }
  | { ok: false; code: "VALIDATION_ERROR"; issues: ValidationIssue[] }
  | { ok: false; code: string }

export function ok<T>(data: T): Result<T> {
  return { ok: true, data }
}

export function err(code: string): Result<never> {
  return { ok: false, code }
}

export function validationErr(issues: ValidationIssue[]): Result<never> {
  return { ok: false, code: "VALIDATION_ERROR", issues }
}
