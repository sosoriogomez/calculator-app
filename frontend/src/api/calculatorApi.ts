export type Operation = "add" | "subtract" | "multiply" | "divide" | "power" | "sqrt" | "percentage";

export interface CalculationRequest { operation: Operation; a: number; b?: number }
export interface CalculationResponse { operation: string; result: number }
export interface ApiError { error?: { code?: string; message?: string } }

export async function calculate(request: CalculationRequest): Promise<CalculationResponse> {
  const response = await fetch(`${import.meta.env.VITE_API_URL ?? ""}/api/v1/calculations`, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(request) });
  const payload = await response.json() as CalculationResponse & ApiError;
  if (!response.ok) throw new Error(payload.error?.message ?? "The calculation could not be completed.");
  return payload;
}
