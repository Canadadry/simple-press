import { describe, it, expect } from "vitest";
import { render, WriterRenderer, Dict } from "../src/pkg/data/render";

class MockWriterRenderer implements WriterRenderer {
  private lines: string[] = [];

  private log(line: string): void {
    this.lines.push(line);
  }

  beginForm(): void {
    this.log("BeginForm");
  }

  endForm(): void {
    this.log("EndForm");
  }

  beginObject(label: string): void {
    this.log(`BeginFieldset label=${label}`);
  }

  endObject(): void {
    this.log("EndFieldset");
  }

  input(label: string, name: string, inputType: string, value: string): void {
    this.log(
      `Input label=${label} name=${name} type=${inputType} value=${value}`,
    );
  }

  checkbox(label: string, name: string, checked: boolean): void {
    this.log(
      `Checkbox label=${label} name=${name} checked=${checked.toString()}`,
    );
  }

  select(label: string, name: string, options: string[], value: string): void {
    this.log(
      `Select label=${label} name=${name} options=${JSON.stringify(
        options,
      )} value=${value}`,
    );
  }

  getOutput(): string[] {
    return this.lines.slice();
  }
}

function compareStacks(got: string[], want: string[]) {
  expect(got.length).toBe(want.length);
  for (let i = 0; i < want.length; i++) {
    expect(got[i]).toBe(want[i]);
  }
}

describe("render", () => {
  it("nested object", () => {
    const input: Dict = {
      profile: {
        name: {
          first: "Jane",
          last: "Doe",
        },
        age: 42,
      },
    };

    const expectLines = [
      "BeginForm",
      "BeginFieldset label=profile",
      "Input label=age name=profile.age type=number value=42",
      "BeginFieldset label=name",
      "Input label=first name=profile.name.first type=text value=Jane",
      "Input label=last name=profile.name.last type=text value=Doe",
      "EndFieldset",
      "EndFieldset",
      "EndForm",
    ];

    const mock = new MockWriterRenderer();
    const err = render(input, mock);

    expect(err).toBeNull();
    compareStacks(mock.getOutput(), expectLines);
  });

  it("throws on array", () => {
    const input: Dict = {
      profile: [1, 2, 3],
    };

    const mock = new MockWriterRenderer();
    const err = render(input, mock);

    expect(err).not.toBeNull();
  });
});
