import { useState } from "react";
import type { FormEvent } from "react";
import { calculate } from "../api/calculatorApi";
import type { Operation } from "../api/calculatorApi";

const operations: { value: Operation; label: string; symbol: string; hint: string }[] = [
  { value: "add", label: "Addition", symbol: "+", hint: "Add two numbers" },
  { value: "subtract", label: "Subtraction", symbol: "−", hint: "Subtract the second number" },
  { value: "multiply", label: "Multiplication", symbol: "×", hint: "Multiply two numbers" },
  { value: "divide", label: "Division", symbol: "÷", hint: "Divide the first number" },
  { value: "power", label: "Exponentiation", symbol: "xʸ", hint: "Raise a number to a power" },
  { value: "sqrt", label: "Square root", symbol: "√", hint: "Find the square root" },
  { value: "percentage", label: "Percentage", symbol: "%", hint: "Find what percentage the value is of the base" },
];

function parseInput(value: string): number | null { if (value.trim() === "") return null; const number = Number(value); return Number.isFinite(number) ? number : null; }

export default function Calculator() {
  const [operation, setOperation] = useState<Operation>("add");
  const [a, setA] = useState("");
  const [b, setB] = useState("");
  const [result, setResult] = useState<number | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const needsSecond = operation !== "sqrt";
  const selected = operations.find((item) => item.value === operation)!;

  function selectOperation(value: Operation) { setOperation(value); setResult(null); setError(""); }

  async function onSubmit(event: FormEvent) {
    event.preventDefault(); setError(""); setResult(null);
    const first = parseInput(a); const second = parseInput(b);
    if (first === null) { setError("Enter a valid number for the first value."); return; }
    if (needsSecond && second === null) { setError("Enter a valid number for the second value."); return; }
    setLoading(true);
    try { const response = await calculate(needsSecond ? { operation, a: first, b: second! } : { operation, a: first }); setResult(response.result); }
    catch (caught) { setError(caught instanceof Error ? caught.message : "The calculation could not be completed."); }
    finally { setLoading(false); }
  }

  return <section className="calculator-layout" aria-label="Calculator workspace">
    <form className="calculator-panel" onSubmit={onSubmit}>
      <div className="panel-heading"><div><p className="eyebrow">Compute precisely</p><h2>Choose an operation</h2></div><span className="operation-symbol" aria-hidden="true">{selected.symbol}</span></div>
      <div className="operation-grid" role="group" aria-label="Operations">{operations.map((item) => <button type="button" className={`operation-button ${operation === item.value ? "selected" : ""}`} aria-pressed={operation === item.value} key={item.value} onClick={() => selectOperation(item.value)}><span className="operation-icon">{item.symbol}</span><span><strong>{item.label}</strong><small>{item.hint}</small></span></button>)}</div>
      <div className="input-grid">
        <label><span>{operation === "percentage" ? "Value" : "First value"}</span><input aria-label="First value" inputMode="decimal" type="number" step="any" value={a} onChange={(event) => setA(event.target.value)} placeholder="0" /></label>
        {needsSecond && <label><span>{operation === "percentage" ? "Base" : operation === "power" ? "Exponent" : "Second value"}</span><input aria-label={operation === "percentage" ? "Base" : operation === "power" ? "Exponent" : "Second value"} inputMode="decimal" type="number" step="any" value={b} onChange={(event) => setB(event.target.value)} placeholder="0" /></label>}
      </div>
      <button className="submit-button" type="submit" disabled={loading}>{loading ? "Calculating…" : "Calculate"}<span aria-hidden="true">→</span></button>
      {error && <p className="feedback error" role="alert">{error}</p>}
    </form>
    <aside className={`result-panel ${result !== null ? "has-result" : ""}`} aria-live="polite"><p className="eyebrow">Result</p>{result === null ? <div className="empty-result"><span aria-hidden="true">=</span><p>Your answer will appear here.</p></div> : <div className="result-value"><span>{selected.label}</span><strong>{result.toLocaleString(undefined, { maximumFractionDigits: 12 })}</strong></div>}<div className="result-rule" /><p className="result-note">Results are validated by the calculator API before they appear.</p></aside>
  </section>;
}
