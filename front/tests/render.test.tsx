import { describe, it, expect } from "vitest";
import { render, screen, within } from "@testing-library/react";
import { DynamicForm } from "../src/pkg/data/render";
import type { DynamicFormUI } from "../src/pkg/data/render";
import React from "react";
import { Dict } from "../src/pkg/data/parseFormData";

function makeTestUI(): DynamicFormUI {
  return {
    Form: ({ children }) => <form data-testid="form">{children}</form>,
    FormObject: ({ label, children }) => (
      <fieldset data-testid={`object-${label}`}>{children}</fieldset>
    ),
    FormInput: ({ label, name, inputType, value }) => (
      <input
        aria-label={label}
        data-testid={`input-${name}`}
        type={inputType}
        defaultValue={value}
      />
    ),
    FormCheckBox: ({ label, name, checked }) => (
      <input
        aria-label={label}
        data-testid={`checkbox-${name}`}
        type="checkbox"
        defaultChecked={checked}
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
    render(<DynamicForm data={meta} ui={ui} />);

    const form = screen.getByTestId("form");
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
});
