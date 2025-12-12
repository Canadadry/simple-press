import { describe, it, expect } from "vitest";
import { updateData } from "../src/pkg/data/updateData";
import { Dict } from "../src/api/api";

describe("updateData", () => {
  const base: Dict = {
    firstname: "Alice",
    age: 42,
    profile: {
      name: { first: "Bob", last: "Morane" },
      active: false,
    },
  };

  it("updates a flat string field", () => {
    const updated = updateData(base, "firstname", "John");

    expect(updated).toEqual({
      firstname: "John",
      age: 42,
      profile: {
        name: { first: "Bob", last: "Morane" },
        active: false,
      },
    });
  });

  it("updates a flat number field", () => {
    const updated = updateData(base, "age", "50");

    expect(updated).toEqual({
      firstname: "Alice",
      age: 50,
      profile: {
        name: { first: "Bob", last: "Morane" },
        active: false,
      },
    });
  });

  it("updates a nested string field", () => {
    const updated = updateData(base, "profile.name.first", "Alice");

    expect(updated).toEqual({
      firstname: "Alice",
      age: 42,
      profile: {
        name: { first: "Alice", last: "Morane" },
        active: false,
      },
    });
  });

  it("updates a boolean field", () => {
    const updated = updateData(base, "profile.active", "true");

    expect(updated).toEqual({
      firstname: "Alice",
      age: 42,
      profile: {
        name: { first: "Bob", last: "Morane" },
        active: true,
      },
    });
  });

  it("creates a correct deep copy without modifying originals", () => {
    const original = JSON.parse(JSON.stringify(base));
    const updated = updateData(base, "profile.name.first", "Z");

    expect(updated).toEqual({
      firstname: "Alice",
      age: 42,
      profile: {
        name: { first: "Z", last: "Morane" },
        active: false,
      },
    });

    expect(base).toEqual(original);
  });
  it("does nothing if the key does not exist in the original data", () => {
    const updated = updateData(base, "profile.unknown", "XXX");

    expect(updated).toEqual({
      firstname: "Alice",
      age: 42,
      profile: {
        name: { first: "Bob", last: "Morane" },
        active: false,
      },
    });
  });
});
