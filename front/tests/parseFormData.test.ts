import { describe, it, expect } from "vitest";
import { parseFormData } from "../src/pkg/data/parseFormData";

describe("parseFormData", () => {
  const tests: Record<
    string,
    {
      FormValues: Record<string, string[]>;
      Definition: Record<string, unknown>;
      Expected: Record<string, unknown>;
    }
  > = {
    "flat object": {
      FormValues: {
        firstname: ["John"],
        email: ["john@example.com"],
        newsletter: ["true"],
      },
      Definition: {
        firstname: "Alice",
        email: "Alice@example.com",
        newsletter: false,
      },
      Expected: {
        firstname: "John",
        email: "john@example.com",
        newsletter: true,
      },
    },

    "nested object": {
      FormValues: {
        "profile.name.first": ["Alice"],
        "profile.name.last": ["Smith"],
        "profile.age": ["28"],
      },
      Definition: {
        profile: {
          name: { first: "Bob", last: "Morane" },
          age: 10,
        },
      },
      Expected: {
        profile: {
          name: { first: "Alice", last: "Smith" },
          age: 28,
        },
      },
    },

    "missing field default to old value": {
      FormValues: {
        "profile.name.first": ["Alice"],
        "profile.name.last": ["Smith"],
        "profile.age": ["28"],
        "profile.gender": ["Mme"],
      },
      Definition: {
        profile: {
          name: { first: "Bob", last: "Morane" },
          age: 10,
        },
      },
      Expected: {
        profile: {
          name: { first: "Alice", last: "Smith" },
          age: 28,
        },
      },
    },

    "extra field ignored": {
      FormValues: {
        "profile.name.first": ["Alice"],
        "profile.name.last": ["Smith"],
      },
      Definition: {
        profile: {
          name: { first: "Bob", last: "Morane" },
          age: 10,
        },
      },
      Expected: {
        profile: {
          name: { first: "Alice", last: "Smith" },
          age: 10,
        },
      },
    },
  };

  for (const [name, tt] of Object.entries(tests)) {
    it(name, async () => {
      const formBody = new URLSearchParams();
      for (const key in tt.FormValues) {
        for (const v of tt.FormValues[key]) {
          formBody.append(key, v);
        }
      }

      const req = new Request("http://localhost/submit", {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: formBody.toString(),
      });

      const result = await parseFormData(req, tt.Definition);

      expect(result).toEqual(tt.Expected);
    });
  }
});
