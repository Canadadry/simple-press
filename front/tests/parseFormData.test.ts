import { describe, it, expect } from "vitest";
import { extractData, Dict, FlatDict } from "../src/pkg/data/parseFormData";

describe("extractData", () => {
  type Test = {
    FormValues: FlatDict;
    Definition: Dict;
    Expected: Dict;
  };
  const tests: Record<string, Test> = {
    "flat object": {
      FormValues: {
        firstname: "John",
        email: "john@example.com",
        newsletter: "true",
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
        "profile.name.first": "Alice",
        "profile.name.last": "Smith",
        "profile.age": "28",
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
        "profile.name.first": "Alice",
        "profile.name.last": "Smith",
        "profile.age": "28",
        "profile.gender": "Mme",
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
        "profile.name.first": "Alice",
        "profile.name.last": "Smith",
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
      const result = extractData(tt.FormValues, tt.Definition, "");
      expect(result).toEqual(tt.Expected);
    });
  }
});
