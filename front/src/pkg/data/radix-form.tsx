import React from "react";
import {
  Card,
  Flex,
  Box,
  Text,
  TextField,
  Checkbox,
  Button,
} from "@radix-ui/themes";
import type { DynamicFormUI } from "./render";
import { TrashIcon } from "@radix-ui/react-icons";

export function makeRadixUI(maxWidth: number): DynamicFormUI {
  return {
    Form: ({ label, children }) => (
      <Card>
        <Text as="div" size="2" mb="2" weight="bold" color="indigo">
          {label}
        </Text>
        <Box
          data-testid="form"
          style={{
            maxWidth: maxWidth,
          }}
        >
          {children}
          <Flex gap="2" justify="between">
            <Button tabIndex={1} size="2" variant="outline" color="crimson">
              <TrashIcon />
              Delete
            </Button>
            <Button tabIndex={1} size="2">
              Save
            </Button>
          </Flex>
        </Box>
      </Card>
    ),
    FormObject: ({ label, children }) => (
      <Card mb="2">
        <Text as="div" size="2" mb="2" weight="bold">
          {label}
        </Text>
        <Box mb="2" data-testid={`object-${label}`}>
          {children}
        </Box>
      </Card>
    ),
    FormInput: ({ label, name, inputType, value }) => (
      <TextField.Root
        mb="4"
        data-testid={`input-${name}`}
        defaultValue={value}
        type={inputType as "text" | "number"}
      >
        <TextField.Slot>{label}</TextField.Slot>
      </TextField.Root>
    ),
    FormCheckBox: ({ label, name, checked }) => (
      <Text as="label" size="2" data-testid={`checkbox-${name}`}>
        <Flex gap="2">
          <Checkbox checked={checked} />
          {label}
        </Flex>
      </Text>
    ),
  };
}
