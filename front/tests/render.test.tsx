import { describe, it, expect, vi } from "vitest";
import { render, screen, within, fireEvent } from "@testing-library/react";
import { DynamicForm } from "../src/pkg/data/render";
import type { DynamicFormUI } from "../src/pkg/data/render";
import React from "react";
import { Dict } from "../src/api/api";

function makeTestUI(): DynamicFormUI {
  return {
    Form: ({ label, children }) => (
      <form data-testid={`form-${label}`}>{children}</form>
    ),
    FormObject: ({ label, children }) => (
      <fieldset data-testid={`object-${label}`}>{children}</fieldset>
    ),
    FormInput: ({ label, name, inputType, value, setData }) => (
      <input
        aria-label={label}
        data-testid={`input-${name}`}
        type={inputType}
        defaultValue={value}
        onChange={(e) => setData(name, e.target.value)}
      />
    ),
    FormCheckBox: ({ label, name, checked, setData }) => (
      <input
        aria-label={label}
        data-testid={`checkbox-${name}`}
        type="checkbox"
        defaultChecked={checked}
        onChange={(e) => setData(name, String(e.target.checked))}
      />
    ),
  };
}

describe("DynamicForm rendering", () => {
  it("renders nested object structure correctly", () => {
    const meta: Dict = {
      profile: {
        name: {
          first: "Jane",
          last: "Doe",
        },
        age: 42,
      },
    };

    const ui = makeTestUI();
    const setData = vi.fn();

    render(<DynamicForm name="test" data={meta} setData={setData} ui={ui} />);

    const form = screen.getByTestId("form-test");
    expect(form).toBeInTheDocument();

    const profileObj = screen.getByTestId("object-profile");
    expect(profileObj).toBeInTheDocument();

    const nameObj = within(profileObj).getByTestId("object-name");
    expect(nameObj).toBeInTheDocument();

    const first = within(nameObj).getByTestId("input-profile.name.first");
    const last = within(nameObj).getByTestId("input-profile.name.last");
    expect(first).toBeInTheDocument();
    expect(last).toBeInTheDocument();

    const age = within(profileObj).getByTestId("input-profile.age");
    expect(age).toBeInTheDocument();

    const ageInName = within(nameObj).queryByTestId("input-profile.age");
    expect(ageInName).toBeNull();
  });

  it("passes setData to each FormInput", () => {
    const meta: Dict = {
      profile: { age: 42 },
    };
    const ui = makeTestUI();
    const setData = vi.fn();

    render(<DynamicForm name="test" data={meta} setData={setData} ui={ui} />);

    const input = screen.getByTestId("input-profile.age");

    // NEW: vérifie qu'un changement déclenche l'appel de setData via l'UI
    fireEvent.change(input, { target: { value: "43" } });

    expect(setData).toHaveBeenCalledTimes(1);
    expect(setData.mock.calls[0][0]).toEqual({
      profile: { age: 43 },
    });
  });

  it("updates deeply nested values through setData", () => {
    const meta: Dict = {
      profile: {
        name: { first: "Jane" },
      },
    };
    const ui = makeTestUI();
    const setData = vi.fn();

    render(<DynamicForm name="test" data={meta} setData={setData} ui={ui} />);

    const first = screen.getByTestId("input-profile.name.first");

    fireEvent.change(first, { target: { value: "Alice" } });

    expect(setData).toHaveBeenCalledTimes(1);
    expect(setData.mock.calls[0][0]).toEqual({
      profile: {
        name: { first: "Alice" },
      },
    });
  });

  it("updates checkbox values through setData", () => {
    const meta: Dict = {
      flags: { active: false },
    };
    const ui = makeTestUI();
    const setData = vi.fn();

    render(<DynamicForm name="test" data={meta} setData={setData} ui={ui} />);

    const checkbox = screen.getByTestId("checkbox-flags.active");

    fireEvent.click(checkbox);

    expect(setData).toHaveBeenCalledTimes(1);
    expect(setData.mock.calls[0][0]).toEqual({
      flags: { active: true },
    });
  });
});
