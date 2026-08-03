import { beforeEach, describe, expect, it, vi } from "vitest";
import { calculate } from "./calculatorApi";

describe("calculate API client", () => {
  beforeEach(() => { vi.restoreAllMocks(); });
  it("returns a successful calculation", async () => { vi.stubGlobal("fetch", vi.fn().mockResolvedValue({ ok: true, json: async () => ({ operation: "add", result: 5 }) })); await expect(calculate({ operation: "add", a: 2, b: 3 })).resolves.toEqual({ operation: "add", result: 5 }); expect(fetch).toHaveBeenCalledWith("/api/v1/calculations", expect.objectContaining({ method: "POST", body: JSON.stringify({ operation: "add", a: 2, b: 3 }) })); });
  it("surfaces the backend error", async () => { vi.stubGlobal("fetch", vi.fn().mockResolvedValue({ ok: false, json: async () => ({ error: { message: "division by zero" } }) })); await expect(calculate({ operation: "divide", a: 1, b: 0 })).rejects.toThrow("division by zero"); });
  it("uses a fallback for an empty backend error", async () => { vi.stubGlobal("fetch", vi.fn().mockResolvedValue({ ok: false, json: async () => ({}) })); await expect(calculate({ operation: "add", a: 1, b: 2 })).rejects.toThrow("The calculation could not be completed."); });
});
