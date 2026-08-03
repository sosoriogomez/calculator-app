import { describe, expect, it, vi, beforeEach } from "vitest";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Calculator from "./Calculator";
import * as api from "../api/calculatorApi";

vi.mock("../api/calculatorApi");
const mockedCalculate = vi.mocked(api.calculate);

describe("Calculator", () => {
  beforeEach(() => { vi.resetAllMocks(); });
  it("validates missing values before calling the API", async () => { render(<Calculator />); await userEvent.click(screen.getByRole("button", { name: "Calculate" })); expect(screen.getByRole("alert")).toHaveTextContent("first value"); expect(mockedCalculate).not.toHaveBeenCalled(); });
  it("calculates and displays an answer", async () => { mockedCalculate.mockResolvedValue({ operation: "add", result: 5 }); render(<Calculator />); await userEvent.type(screen.getByLabelText("First value"), "2"); await userEvent.type(screen.getByLabelText("Second value"), "3"); await userEvent.click(screen.getByRole("button", { name: "Calculate" })); await waitFor(() => expect(screen.getByText("5")).toBeInTheDocument()); expect(mockedCalculate).toHaveBeenCalledWith({ operation: "add", a: 2, b: 3 }); });
  it("supports square root with one operand", async () => { mockedCalculate.mockResolvedValue({ operation: "sqrt", result: 4 }); render(<Calculator />); await userEvent.click(screen.getByRole("button", { name: /Square root/ })); expect(screen.queryByLabelText("Second value")).not.toBeInTheDocument(); await userEvent.type(screen.getByLabelText("First value"), "16"); await userEvent.click(screen.getByRole("button", { name: "Calculate" })); await waitFor(() => expect(screen.getByText("4")).toBeInTheDocument()); });
  it("shows backend errors", async () => { mockedCalculate.mockRejectedValue(new Error("division by zero")); render(<Calculator />); await userEvent.click(screen.getByRole("button", { name: /Division/ })); await userEvent.type(screen.getByLabelText("First value"), "1"); await userEvent.type(screen.getByLabelText("Second value"), "0"); await userEvent.click(screen.getByRole("button", { name: "Calculate" })); await waitFor(() => expect(screen.getByRole("alert")).toHaveTextContent("division by zero")); });
  it("validates the second operand and clears old state on operation change", async () => { render(<Calculator />); await userEvent.type(screen.getByLabelText("First value"), "2"); await userEvent.click(screen.getByRole("button", { name: /Percentage/ })); await userEvent.click(screen.getByRole("button", { name: "Calculate" })); expect(screen.getByRole("alert")).toHaveTextContent("second value"); fireEvent.click(screen.getByRole("button", { name: /Square root/ })); expect(screen.queryByRole("alert")).not.toBeInTheDocument(); });
});
